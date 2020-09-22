package smtp_client

import (
	"fmt"
	"log"

	"github.com/davecourtois/Globular/smtp/smtp_client"
	"github.com/davecourtois/Globular/smtp/smtppb"

	"testing"
)

var (
	client *smtp_client.SMTP_Client
)

// First test create a fresh new connection...
func TestCreateConnection(t *testing.T) {
	var err error
	client, err = smtp_client.NewSmtp_Client("globular.io", "smtp_server")
	if err != nil {
		log.Panicln(err)
	}

	fmt.Println("Connection creation test.")
	err = client.CreateConnection("test_smtp", "username", "password", 25, "localhost")
	if err != nil {
		log.Panicln(err)
	}
}

/**
 * Test send email whitout attachements.
 */
func TestSendEmail(t *testing.T) {

	from := "dave.courtois60@localhost"
	to := []string{"dave.courtois60@gmail.com"}
	cc := []*smtppb.CarbonCopy{&smtppb.CarbonCopy{Name: "Dave Courtois", Address: "dave.courtois60@gmail.com"}}
	subject := "Smtp Test"
	body := `<meta http-equiv="Content-Type" content="text/html; charset=utf-8"><div dir="ltr">Message test.</div>`
	bodyType := int32(smtppb.BodyType_HTML)

	err := client.SendEmail("test_smtp", from, to, cc, subject, body, bodyType)

	if err != nil {
		log.Panicln(err)
	}
}

/**
 * Test send email whit attachements.
 */
/*func TestSendEmailWithAttachements(t *testing.T) {

	from := "dave.courtois@safrangroup.com"
	to := []string{"dave.courtois@safrangroup.com"}
	cc := []*smtppb.CarbonCopy{&smtppb.CarbonCopy{Name: "Dave Courtois", Address: "dave.courtois60@gmail.com"}}
	subject := "Smtp Test"
	body := `<meta http-equiv="Content-Type" content="text/html; charset=utf-8"><div dir="ltr">Message test.</div>`
	bodyType := int32(smtppb.BodyType_HTML)
	attachments := []string{"attachements/Daft Punk - Get Lucky (Official Audio) ft. Pharrell Williams, Nile Rodgers.mp3", "attachements/NGEN3549.JPG", "attachements/NGEN3550.JPG"}

	err := client.SendEmailWithAttachements("test_smtp", from, to, cc, subject, body, bodyType, attachments)

	if err != nil {
		log.Panicln(err)
	}

}*/
