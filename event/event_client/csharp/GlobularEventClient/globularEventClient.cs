using System.Runtime.InteropServices.WindowsRuntime;
using System.Diagnostics;
using System;
using System.Threading;
using System.Threading.Channels;

namespace Globular
{
    
    public class GlobularEventClient : Client
    {
        private Event.EventService.EventServiceClient client;

        public struct Action {
            string name;


        }

        /// <summary>
        /// gRPC client for event service.
        /// </summary>
        /// <param name="address">Can be a domain or a IP address ex: localhos or 127.0.0.1</param>
        /// <param name="name">The name of the service on the server. ex: persistence_server</param>
        /// <returns>Return the instance of the client with it connection ready to be use.</returns>
        public GlobularEventClient(string address, string name) : base(address, name)
        {
            // Here I will create grpc connection with the service...
            this.client = new Event.EventService.EventServiceClient(this.channel);

            // process...
            Thread t0 = new Thread(this.processEvent);
            t0.Start();
        }

        private void processEvent(){
            // Here I will start on event processing.
            var data_channel = Channel.CreateUnbounded<Event.Event>();

            var evt = data_channel.Reader.WaitToReadAsync();

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

        public void Subscribe(string name, string uuid, Func<string> callback)
        {
            var rqst = new Event.SubscribeRequest();
            rqst.Name = name;
            rqst.Uuid = uuid;

            // Subscribe to a given event.
            this.client.Subscribe(rqst, this.GetClientContext());

            // Register the uuid in local subscriber 
        }

        public void UnSubscribe(string name, string uuid)
        {
            var rqst = new Event.UnSubscribeRequest();
            rqst.Name = name;
            rqst.Uuid = uuid;

            // Subscribe to a given event.
            this.client.UnSubscribe(rqst, this.GetClientContext());

            // Register the uuid in local subscriber 
        }

        public void OnEvent(string uuid, Channel<Event.Event> channel){
            // Here I will open a channel write.
            var writer = channel.Writer;
            
        }

    }
}
