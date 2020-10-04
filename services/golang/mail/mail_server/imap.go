package main

import (
	"log"

	"crypto/tls"

	//	"github.com/emersion/go-imap"
	//	"github.com/emersion/go-imap/backend"
	"github.com/davecourtois/Utility"
	"github.com/emersion/go-imap/backend/memory"
	imap_server "github.com/emersion/go-imap/server"
)

/**
 * So here I will implement Globular/MongoDB backend for IMAP.
 */

/*type Backend struct {
}

func (self *Backend) Login(connInfo *imap.ConnInfo, username, password string) (backend.User, error) {
	log.Println("----------> try to login to IMAP backend.", username, password)
	return nil, nil
}*/

func startImap(pwd string, keyFile string, certFile string, port int) {

	// Create a memory backend
	be := memory.New() //new(Backend)

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
