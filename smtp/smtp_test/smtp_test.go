package Globular

import (
	"context"
	"fmt"
	"log"

	"github.com/davecourtois/Globular/smtp/smtppb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

// Set the correct addresse here as needed.
var (
	addresse = "localhost:10007"
)

/**
 *
 */
func getClientConnection() *grpc.ClientConn {
	// So here I will read the server configuration to see if the connection
	// is secure...
	config := make(map[string]interface{})
	data, err := ioutil.ReadFile("../smtp_server/config.json")
	if err != nil {
		log.Fatal("fail to read configuration")
	}

	// Read the config file.
	json.Unmarshal(data, &config)

	var cc *grpc.ClientConn
	if cc == nil {
		if config["TLS"].(bool) {
			creds, sslErr := credentials.NewClientTLSFromFile(config["CertAuthorityTrust"].(string), "")
			if err != nil {
				log.Fatalf("Error while loading CA trust certificate: %v", sslErr)
			}
			opts := grpc.WithTransportCredentials(creds)
			cc, err = grpc.Dial(addresse, opts)
		} else {
			cc, err = grpc.Dial(addresse, grpc.WithInsecure())
			if err != nil {
				log.Fatalf("could not connect: %v", err)
			}
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

/**
 * Here I will make use of a stream
 */
func sendFile(id string, path string, stream smtppb.SmtpService_SendEmailWithAttachementsClient) {

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Fail to open file "+path+" with error: %v", err)
	}

	// close the file when done.
	defer file.Close()

	const BufferSize = 1024 * 5 // the chunck size.

	buffer := make([]byte, BufferSize)
	for {
		bytesread, err := file.Read(buffer)
		if bytesread > 0 {
			rqst := &smtppb.SendEmailWithAttachementsRqst{
				Id: id,
				Data: &smtppb.SendEmailWithAttachementsRqst_Attachements{
					Attachements: &smtppb.Attachement{
						FileName: path,
						FileData: buffer[:bytesread],
					},
				},
			}
			err = stream.Send(rqst)
		}

		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
	}
}

/**
 * Test send email whit attachements.
 */
func TestSendEmailWithAttachements(t *testing.T) {
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := smtppb.NewSmtpServiceClient(cc)

	// Open the stream...
	stream, err := c.SendEmailWithAttachements(context.Background())
	if err != nil {
		log.Fatalf("error while TestSendEmailWithAttachements: %v", err)
	}

	// Send file attachment as a stream, not need to be send first.
	sendFile("test_smtp", "attachements/Daft Punk - Get Lucky (Official Audio) ft. Pharrell Williams, Nile Rodgers.mp3", stream)
	sendFile("test_smtp", "attachements/NGEN3549.JPG", stream)
	sendFile("test_smtp", "attachements/NGEN3550.JPG", stream)

	// Send the email message...
	rqst := &smtppb.SendEmailWithAttachementsRqst{
		Id: "test_smtp",
		Data: &smtppb.SendEmailWithAttachementsRqst_Email{
			Email: &smtppb.Email{
				From:     "dave.courtois@safrangroup.com",
				To:       []string{"dave.courtois@safrangroup.com"},
				Cc:       []*smtppb.CarbonCopy{&smtppb.CarbonCopy{Name: "Dave Courtois", Address: "dave.courtois60@gmail.com"}},
				Subject:  "Smtp Test",
				Body:     "This is a simple mail test!",
				BodyType: smtppb.BodyType_TEXT,
			},
		},
	}

	err = stream.Send(rqst)
	if err != nil {
		log.Fatalf("error while TestSendEmailWithAttachements: %v", err)
	}

	rsp, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Fail Send email with attachement %v", err)
	}

	log.Println("Response form SendEmailWithAttachements:", rsp.Result)

}
