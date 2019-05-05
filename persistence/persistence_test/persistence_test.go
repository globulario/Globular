package Globular

import (
	"context"
	"fmt"
	"log"

	"github.com/davecourtois/Globular/persistence/persistencepb"
	"google.golang.org/grpc"

	"testing"

	"time"

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

// Persist an employe to the store.
func persistEmploye(c persistencepb.PersistenceServiceClient, emp_no float64, birth_date float64, first_name string, last_name string, gender string, hire_date float64) error {
	// Here I will persist entity...
	rqst := &persistencepb.PersistEntityRqst{
		Entity: &persistencepb.Entity{
			UUID:     Utility.RandomUUID(),
			Typename: "Employee",
			Attibutes: []*persistencepb.Attribute{
				&persistencepb.Attribute{
					Name: "emp_no",
					Value: &persistencepb.Attribute_NumericVal{
						NumericVal: emp_no,
					},
				},
				&persistencepb.Attribute{
					Name: "birth_date",
					Value: &persistencepb.Attribute_NumericVal{
						NumericVal: birth_date,
					},
				},
				&persistencepb.Attribute{
					Name: "first_name",
					Value: &persistencepb.Attribute_StrVal{
						StrVal: first_name,
					},
				},
				&persistencepb.Attribute{
					Name: "last_name",
					Value: &persistencepb.Attribute_StrVal{
						StrVal: last_name,
					},
				},
				&persistencepb.Attribute{
					Name: "gender",
					Value: &persistencepb.Attribute_StrVal{
						StrVal: gender,
					},
				}, &persistencepb.Attribute{
					Name: "hire_date",
					Value: &persistencepb.Attribute_NumericVal{
						NumericVal: hire_date,
					},
				},
			},
		},
	}

	rsp, err := c.PersistEntity(context.Background(), rqst)
	if err != nil {
		return err
	}

	log.Println(rsp.Result)

	return nil
}

// First test create a fresh new connection...
func TestPersist(t *testing.T) {
	fmt.Println("Connection creation test.")

	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := persistencepb.NewPersistenceServiceClient(cc)

	// Here I will convert the values to basic type.
	emp_no := float64(10906)                                // int to numeric
	hire_date, _ := time.Parse("2006-01-02", "2011-01-19")  // date string to unix time
	birth_date, _ := time.Parse("1979-01-28", "2011-01-19") // date string to unix time

	// TODO make the test with employees from sql table to see performance.
	err := persistEmploye(c, emp_no, float64(hire_date.Unix()), "Dave", "Courtois", "M", float64(birth_date.Unix()))

	if err != nil {
		log.Panicln("fail test persist!", err)
	}
}
