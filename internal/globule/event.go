package globule

import (
	"fmt"

	"github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/event/event_client"
	"github.com/globulario/services/golang/event/eventpb"
	"github.com/globulario/services/golang/globular_client"
	Utility "github.com/globulario/utility"
)

// getEventClient returns an Event service client bound to this node's address.
func (g *Globule) getEventClient() (*event_client.Event_Client, error) {
	Utility.RegisterFunction("NewEventService_Client", event_client.NewEventService_Client)
	addr, _ := config.GetAddress()
	c, err := globular_client.GetClient(addr, "event.EventService", "NewEventService_Client")
	if err != nil {
		return nil, fmt.Errorf("getEventClient: %w", err)
	}
	return c.(*event_client.Event_Client), nil
}

// subscribe wraps the Event service's Subscribe so call-sites stay simple.
func (g *Globule) subscribe(topic string, handler func(evt *eventpb.Event)) error {
	ev, err := g.getEventClient()
	if err != nil {
		return fmt.Errorf("subscribe(%s): %w", topic, err)
	}
	// Use the node name as the subscriber id (same behavior as before).
	if err := ev.Subscribe(topic, g.Name, handler); err != nil {
		return fmt.Errorf("subscribe(%s): %w", topic, err)
	}
	// Subscription keeps the client open; ownership stays with the caller.
	return nil
}

// publish wraps the Event service's Publish to keep call-sites tidy.
func (g *Globule) publish(topic string, data []byte) error {
	ev, err := g.getEventClient()
	if err != nil {
		return fmt.Errorf("publish(%s): %w", topic, err)
	}
	defer ev.Close()
	if err := ev.Publish(topic, data); err != nil {
		return fmt.Errorf("publish(%s): %w", topic, err)
	}
	return nil
}

// (Optional) unsubscribe helper if you need it later.
//
// func (g *Globule) unsubscribe(topic string) error {
// 	ev, err := g.getEventClient()
// 	if err != nil {
// 		return fmt.Errorf("unsubscribe(%s): %w", topic, err)
// 	}
// 	if err := ev.Unsubscribe(topic, g.Name); err != nil {
// 		return fmt.Errorf("unsubscribe(%s): %w", topic, err)
// 	}
// 	return nil
// }
