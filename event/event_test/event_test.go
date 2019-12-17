package Globular

import (
	"log"
	"strconv"
	"testing"

	"time"

	"github.com/davecourtois/Globular/event/event_client"
)

func subscribeTo(client *event_client.Event_Client, subject string) string {
	data_chan := make(chan []byte)
	uuid, err := client.Subscribe(subject, data_chan)

	if err != nil {
		log.Println(err)
	}

	go func() {
		for {
			select {
			case msg := <-data_chan:
				log.Println(string(msg))
			}
		}
	}()
	return uuid
}

/**
 * Test event
 */
func TestEventService(t *testing.T) {
	log.Println("Test event service")
	domain := "localhost"

	// The topic.
	subject := "my topic"
	size := 1000 // test with 500 client...
	clients := make([]*event_client.Event_Client, size)
	uuids := make([]string, size)
	for i := 0; i < size; i++ {
		c := event_client.NewEvent_Client(domain, "event_server")
		uuids[i] = subscribeTo(c, subject)
		log.Println("client ", i)
		clients[i] = c
	}
	for i := 0; i < 50; i++ {
		clients[0].Publish(subject, []byte("--->"+strconv.Itoa(i)+" this is a message!"))
	}

	// Here I will simply suspend this thread to give time to publish message
	time.Sleep(time.Second * 1)

	for i := 0; i < size; i++ {
		log.Println("---> close the client")
		clients[i].UnSubscribe(subject, uuids[i])
	}

}
