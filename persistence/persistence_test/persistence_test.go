package Globular

import (
	"context"
	"fmt"
	"log"

	"github.com/davecourtois/Globular/persistence/persistencepb"
	"google.golang.org/grpc"

	"testing"

	"github.com/davecourtois/Utility"
)

// Set the correct addresse here as needed.
var (
	addresse = "localhost:50051"
)

/**
 * Get the client connection.
 */
func getClientConnection() *grpc.ClientConn {
	var err error
	var cc *grpc.ClientConn
	if cc == nil {
		cc, err = grpc.Dial(addresse, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("could not connect: %v", err)
		}

	}
	return cc
}

// First test create a fresh new connection...
func TestPersist(t *testing.T) {
	fmt.Println("Connection creation test.")

	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := persistencepb.NewPersistenceServiceClient(cc)

	// Here I will persist entity...
	rqst := &persistencepb.PersistEntityRqst{
		Entity: &persistencepb.Entity{
			UUID:     Utility.RandomUUID(),
			Typename: "Employee",
			Attibutes: []*persistencepb.Attribute{
				// String value test
				&persistencepb.Attribute{
					Name: "id",
					Value: &persistencepb.Attribute_StrVal{
						StrVal: "mm006819",
					},
				},
				// Numric value test
				&persistencepb.Attribute{
					Name: "age",
					Value: &persistencepb.Attribute_NumericVal{
						NumericVal: 43,
					},
				},
				// Bool value test
				&persistencepb.Attribute{
					Name: "isProgrammer",
					Value: &persistencepb.Attribute_BoolVal{
						BoolVal: true,
					},
				},
				// Array test (Numeric array and bool array will be use the same way)
				&persistencepb.Attribute{
					Name: "programmingLanguage",
					Value: &persistencepb.Attribute_StrArr{
						StrArr: &persistencepb.StringArray{
							Values: []string{"c++", "go"},
						},
					},
				},
			},
		},
	}

	rsp, err := c.PersistEntity(context.Background(), rqst)
	if err != nil {
		log.Fatalf("error while CreateConnection: %v", err)
	}

	log.Println("Response form CreateConnection:", rsp.Result)
}
