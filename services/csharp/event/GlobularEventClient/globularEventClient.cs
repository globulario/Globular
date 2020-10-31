using System.Runtime.InteropServices.WindowsRuntime;
using System.Diagnostics;
using System;
using System.Collections.Generic;
using System.Threading.Channels;
using System.Threading.Tasks;
using System.Collections.Concurrent;

namespace Globular
{
    struct Subscriber
    {
        // The event to subscribe to.
        public string Name;

        // The subscriber uuid.
        public string Uuid;

        // The fct to run when event is received.
        public Action<Event.Event> Fct;
    }

    public class GlobularEventClient : Client
    {
        private Event.EventService.EventServiceClient client;
        private string uuid;

        private Channel<Subscriber> subscribe_channel;
        private Channel<Subscriber> unsubscribe_channel;

        /// <summary>
        /// gRPC client for event service.
        /// </summary>
        /// <param name="id"></param> The name or the id of the services.
        /// <param name="domain"></param> The domain of the services
        /// <param name="configurationPort"></param> The domain of the services
        /// <returns>Return the instance of the client with it connection ready to be use.</returns>
        public GlobularEventClient( string id, string domain, int configurationPort) : base(id, domain, configurationPort)
        {
            // Here I will create grpc connection with the service...
            this.client = new Event.EventService.EventServiceClient(this.channel);
            this.uuid = System.Guid.NewGuid().ToString();

            // Start run the 
            Task.Run(() =>
            {
                this.run();
            });

        }

        private void run()
        {
            // Here I will start on event processing.
            var data_channel = Channel.CreateUnbounded<Event.Event>();

            // start listenting to events from the server...
            this.OnEvent(data_channel);

            // Here I will keep handler's
            var handlers = new ConcurrentDictionary<string, Dictionary<string, Action<Event.Event>>>();

            // OnEvent processing loop.
            Task.Run(async () =>
            {
                while (await data_channel.Reader.WaitToReadAsync())
                {
                    var evt = await data_channel.Reader.ReadAsync();
                    if (handlers.ContainsKey(evt.Name))
                    {
                        foreach (var handler in handlers[evt.Name])
                        {
                            handler.Value(evt);
                        }
                    }
                }
            });

            // Now On subscribe processing loop.
            this.subscribe_channel = Channel.CreateUnbounded<Subscriber>();
            Task.Run(async () =>
            {
                while (await subscribe_channel.Reader.WaitToReadAsync())
                {
                    var subscriber = await subscribe_channel.Reader.ReadAsync();
                    if (!handlers.ContainsKey(subscriber.Name))
                    {
                        handlers[subscriber.Name] = new Dictionary<string, Action<Event.Event>>();
                    }
                    handlers[subscriber.Name][subscriber.Uuid] = subscriber.Fct;
                }
            });

            // Unsbuscribe from an event.
            this.unsubscribe_channel = Channel.CreateUnbounded<Subscriber>();
            Task.Run(async () =>
            {
                while (await unsubscribe_channel.Reader.WaitToReadAsync())
                {
                    var subscriber = await unsubscribe_channel.Reader.ReadAsync();
                    if (handlers.ContainsKey(subscriber.Name))
                    {
                        if(handlers[subscriber.Name].ContainsKey(subscriber.Uuid)){
                            handlers[subscriber.Name].Remove(subscriber.Uuid);
                        }
                    }
                }
            });
        }

        public void OnEvent(Channel<Event.Event> channel)
        {
            // Here I will open a channel write.
            var rqst = new Event.OnEventRequest();
            rqst.Uuid = this.uuid;
            
            // Run it in it own tread...
            Task.Run(() =>
            {
                var call = this.client.OnEvent(rqst, this.GetClientContext());
                bool hasNext = true;
                // read until no more values found...
                while (hasNext)
                {
                    var task = Task.Run(() => call.ResponseStream.MoveNext(default(global::System.Threading.CancellationToken)));
                    task.Wait(); // wait for the next value...
                    hasNext = task.Result;
                    if (hasNext)
                    {
                        // write event received from the server on the channel.
                        channel.Writer.WriteAsync(call.ResponseStream.Current.Evt);
                    }
                }
            });

            // The stream was close and no more event will be process.

        }

        /// <summary>
        /// Publish an event on the network.
        /// </summary>
        /// <param name="name">The name of the channel where the event will be publish, can be anything.</param>
        /// <param name="data">The data to be print on the channel as a bytes array.</param>
        public void Publish(string name, byte[] data)
        {
            var rqst = new Event.PublishRequest();
            var evt = new Event.Event();
            evt.Name = name;

            evt.Data = Google.Protobuf.ByteString.CopyFrom(data);
            this.client.Publish(rqst, this.GetClientContext());
        }

        public void Subscribe(string name, string uuid, Action<Event.Event> callback)
        {
            var rqst = new Event.SubscribeRequest();
            rqst.Name = name;
            rqst.Uuid = uuid;

            // Subscribe to a given event.
            this.client.Subscribe(rqst, this.GetClientContext());

            var subscriber = new Subscriber();
            subscriber.Fct = callback;
            subscriber.Uuid = uuid;
            subscriber.Name = name;

            // Register the uuid in local subscribers
            this.subscribe_channel.Writer.WriteAsync(subscriber);
        }

        public void UnSubscribe(string name, string uuid)
        {
            var rqst = new Event.UnSubscribeRequest();
            rqst.Name = name;
            rqst.Uuid = uuid;

            // Subscribe to a given event.
            this.client.UnSubscribe(rqst, this.GetClientContext());

            // Register the uuid in local subscriber
            var subscriber = new Subscriber();
            subscriber.Uuid = uuid;
            subscriber.Name = name;

            // Uregister the uuid in local subscribers
            this.subscribe_channel.Writer.WriteAsync(subscriber);
        }
    }
}
