package event_client

import (
	"io"
	"log"
	"strconv"

	"github.com/davecourtois/Globular/api"
	"github.com/davecourtois/Globular/event/eventpb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// echo Client Service
////////////////////////////////////////////////////////////////////////////////

type Event_Client struct {
	cc *grpc.ClientConn
	c  eventpb.EventServiceClient

	// The name of the service
	name string

	// The client domain
	domain string

	// The port
	port int

	// is the connection is secure?
	hasTLS bool

	// Link to client key file
	keyFile string

	// Link to client certificate file.
	certFile string

	// certificate authority file
	caFile string

	// the client uuid.
	uuid string

	// The event channel.
	actions chan map[string]interface{}
}

// Create a connection to the service.
func NewEvent_Client(address string, name string) *Event_Client {
	client := new(Event_Client)
	api.InitClient(client, address, name)
	client.cc = api.GetClientConnection(client)
	client.c = eventpb.NewEventServiceClient(client.cc)
	client.uuid = Utility.RandomUUID()

	// The channel where data will be exchange.
	client.actions = make(chan map[string]interface{})

	// Here I will also open the event stream to receive event from server.
	go func() {
		err := client.run()
		if err != nil {
			log.Println("Fail to create event client: ", err)
		}
	}()

	return client
}

/**
 * Process event from the server. Only one stream is needed between the server
 * and the client. Local handler are kept in a map with a unique uuid, so many
 * handler can exist for a single event.
 */
func (self *Event_Client) run() error {

	// Create the channel.
	data_channel := make(chan *eventpb.Event, 0)

	// start listenting to events from the server...
	err := self.onEvent(self.uuid, data_channel)
	if err != nil {
		return err
	}

	// the map that will contain the event handler.
	handlers := make(map[string]map[string]func(*eventpb.Event))

	for {
		select {
		case evt := <-data_channel:
			// So here I received an event, I will dispatch it to it function.
			handlers_ := handlers[evt.Name]
			if handlers_ != nil {
				for _, fct := range handlers_ {
					// Call the handler.
					fct(evt)
				}
			}
		case action := <-self.actions:
			if action["action"].(string) == "subscribe" {
				if handlers[action["name"].(string)] == nil {
					handlers[action["name"].(string)] = make(map[string]func(*eventpb.Event))
				}
				// Set it handler.
				handlers[action["name"].(string)][action["uuid"].(string)] = action["fct"].(func(*eventpb.Event))
			} else if action["action"].(string) == "unsubscribe" {
				// Now I will remove the handler...
				for _, handler := range handlers {
					if handler[action["uuid"].(string)] != nil {
						delete(handler, action["uuid"].(string))
					}
				}
			} else if action["action"].(string) == "stop" {
				return nil
			}
		}
	}
}

// Return the domain
func (self *Event_Client) GetDomain() string {
	return self.domain
}

// Return the address
func (self *Event_Client) GetAddress() string {
	return self.domain + ":" + strconv.Itoa(self.port)
}

// Return the name of the service
func (self *Event_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *Event_Client) Close() {
	self.cc.Close()
	action := make(map[string]interface{})
	action["action"] = "stop"

	// set the action.
	self.actions <- action
}

// Set grpc_service port.
func (self *Event_Client) SetPort(port int) {
	self.port = port
}

// Set the client name.
func (self *Event_Client) SetName(name string) {
	self.name = name
}

// Set the domain.
func (self *Event_Client) SetDomain(domain string) {
	self.domain = domain
}

////////////////// TLS ///////////////////

// Get if the client is secure.
func (self *Event_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *Event_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *Event_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *Event_Client) GetCaFile() string {
	return self.caFile
}

// Set the client is a secure client.
func (self *Event_Client) SetTLS(hasTls bool) {
	self.hasTLS = hasTls
}

// Set TLS certificate file path
func (self *Event_Client) SetCertFile(certFile string) {
	self.certFile = certFile
}

// Set TLS key file path
func (self *Event_Client) SetKeyFile(keyFile string) {
	self.keyFile = keyFile
}

// Set TLS authority trust certificate file path
func (self *Event_Client) SetCaFile(caFile string) {
	self.caFile = caFile
}

///////////////////// API ///////////////////////
// Publish and event over the network
func (self *Event_Client) Publish(name string, data interface{}) error {
	//log.Println("publish event ", name)
	rqst := &eventpb.PublishRequest{
		Evt: &eventpb.Event{
			Name: name,
			Data: data.([]byte),
		},
	}

	_, err := self.c.Publish(api.GetClientContext(self), rqst)
	if err != nil {
		return err
	}

	return nil
}

func (self *Event_Client) onEvent(uuid string, data_channel chan *eventpb.Event) error {

	rqst := &eventpb.OnEventRequest{
		Uuid: uuid,
	}

	stream, err := self.c.OnEvent(api.GetClientContext(self), rqst)
	if err != nil {
		log.Println("----> fail to connect to event server ", err)
		return err
	}

	// Run in it own goroutine.
	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				// end of stream...
				break
			}
			if err != nil {
				// oher error stop processing stream.
				log.Println("event stream error: ", err)
				break
			}

			// Get the result...
			data_channel <- msg.Evt

		}
	}()

	// Wait for subscriber uuid and return it to the function caller.
	return nil
}

// Subscribe to an event it return it subscriber uuid. The uuid must be use
// to unsubscribe from the channel. data_channel is use to get event data.
func (self *Event_Client) Subscribe(name string, uuid string, fct func(evt *eventpb.Event)) error {
	//log.Println("Subscribe to event ", name)
	rqst := &eventpb.SubscribeRequest{
		Name: name,
		Uuid: self.uuid,
	}

	_, err := self.c.Subscribe(api.GetClientContext(self), rqst)
	if err != nil {
		return err
	}

	action := make(map[string]interface{})
	action["action"] = "subscribe"
	action["uuid"] = uuid
	action["name"] = name
	action["fct"] = fct

	// set the action.
	self.actions <- action
	return nil
}

// Exit event channel.
func (self *Event_Client) UnSubscribe(name string, uuid string) error {
	//log.Println("UnSubscribe event ", name)

	// Unsubscribe from the event channel.
	rqst := &eventpb.UnSubscribeRequest{
		Name: name,
		Uuid: self.uuid,
	}

	_, err := self.c.UnSubscribe(api.GetClientContext(self), rqst)
	if err != nil {
		return err
	}

	action := make(map[string]interface{})
	action["action"] = "unsubscribe"
	action["uuid"] = uuid
	action["name"] = name

	// set the action.
	self.actions <- action
	return nil
}
