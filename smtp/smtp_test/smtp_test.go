package Globular

import (
	"context"
	"fmt"
	"log"

	"github.com/davecourtois/Globular/smtp/smtppb"
	"google.golang.org/grpc"

	// "encoding/json"
	"testing"
)

/**
 *
 */
func getClientConnection() *grpc.ClientConn {
	var err error
	var cc *grpc.ClientConn
	if cc == nil {
		cc, err = grpc.Dial("localhost:50051", grpc.WithInsecure())
		if err != nil {
			log.Fatalf("could not connect: %v", err)
		}

	}
	return cc
}

// First test create a fresh new connection...
func TestCreateConnection(t *testing.T) {
	fmt.Println("Connection creation test.")

	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := smtppb.NewSmtpServiceClient(cc)

	rqst := &smtppb.CreateConnectionRqst{
		Connection: &smtppb.Connection{
			Id:       "test_smtp",
			User:     "mrmfct037@UD6.UF6",
			Password: "Dowty123",
			Port:     25,
			Host:     "mon-print-01",
		},
	}

	rsp, err := c.CreateConnection(context.Background(), rqst)
	if err != nil {
		log.Fatalf("error while CreateConnection: %v", err)
	}

	log.Println("Response form CreateConnection:", rsp.Result)
}

/**
 * Test send email whitout attachements.
 */
func TestSendEmail(t *testing.T) {
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := smtppb.NewSmtpServiceClient(cc)

	rqst := &smtppb.SendEmailRqst{
		Id: "test_smtp",
		Email: &smtppb.Email{
			From:     "dave.courtois@safrangroup.com",
			To:       []string{"dave.courtois@safrangroup.com"},
			Cc:       []*smtppb.CarbonCopy{&smtppb.CarbonCopy{Name: "Dave Courtois", Address: "dave.courtois60@gmail.com"}},
			Subject:  "Smtp Test",
			Body:     "This is a simple mail test!",
			BodyType: smtppb.BodyType_TEXT,
		},
	}

	rsp, err := c.SendEmail(context.Background(), rqst)
	if err != nil {
		log.Fatalf("error while CreateConnection: %v", err)
	}

	log.Println("Response form CreateConnection:", rsp.Result)
}
