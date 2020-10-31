package imap

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/globulario/Globular/services/golang/persistence/persistence_store"
	"github.com/davecourtois/Utility"

	// "github.com/emersion/go-imap/backend/memory"
	imap_server "github.com/emersion/go-imap/server"
)

////////////////////////////////////////////////////////////////////////////////
// The Backend implementation
////////////////////////////////////////////////////////////////////////////////
var (
	// The backend.
	store *persistence_store.MongoStore

	pwd string
)

/**
 * Save a message in the backend.
 */
func saveMessage(user string, mailBox string, body []byte, flags []string, date time.Time) error {

	data := make(map[string]interface{})

	data["Date"] = date
	data["Flags"] = flags
	data["Size"] = uint32(len(body))
	data["Body"] = body
	data["Uid"] = date.Unix() // I will use the unix time as Uid

	// Now I will insert the message into the inbox of the user.
	_, err := store.InsertOne(context.Background(), "local_ressource", user+"_db", mailBox, data, "")
	if err != nil {
		fmt.Println(err)
	}

	return err
}

/**
 * Rename the connection.
 */
func renameCollection(database string, name string, rename string) error {
	script := `db=db.getSiblingDB('admin');db.adminCommand({renameCollection:'` + database + `.` + name + `', to:'` + database + `.` + rename + `'})`
	err := store.RunAdminCmd(context.Background(), "local_ressource", "sa", pwd, script)
	return err
}

func startImap(port int, keyFile string, certFile string) {

	// Create backend instance.
	be := new(Backend_impl)

	// Create a new server
	s := imap_server.New(be)
	s.Addr = "0.0.0.0:" + Utility.ToString(port)

	go func() {
		if len(certFile) > 0 {
			cer, err := tls.LoadX509KeyPair(certFile, keyFile)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("start imap server at address ", s.Addr)
			s.TLSConfig = &tls.Config{Certificates: []tls.Certificate{cer}}
			if err := s.ListenAndServeTLS(); err != nil {
				log.Fatal(err)
			}
		} else {
			// Since we will use this server for testing only, we can allow plain text
			// authentication over unencrypted connections
			s.AllowInsecureAuth = true
			log.Println("start imap server at address ", s.Addr)
			if err := s.ListenAndServe(); err != nil {
				log.Fatal(err)
			}
		}
	}()
}

func StartImap(password string, keyFile string, certFile string, port int, tls_port int, alt_port int) {

	// Keep the password...
	pwd = password

	// The backend (mongodb)
	store = new(persistence_store.MongoStore)

	// The admin connection. The local ressource contain necessay information.
	err := store.Connect("local_ressource", "0.0.0.0", 27017, "sa", password, "local_ressource", 1000, "")
	if err != nil {
		log.Println("Fail to connect to the backend!", err)
	}

	// Create a memory backend
	startImap(port, "", "")
	startImap(tls_port, keyFile, certFile)
	startImap(alt_port, keyFile, certFile)
}
