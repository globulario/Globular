// dns.go
package main

import (
	"log"
	"net"
	"strconv"

	"github.com/miekg/dns"
)

var domainsToAddresses map[string]string = map[string]string{
	"steve.clients.omniscient.app.": "135.19.213.242",
	"dave.clients.omniscient.app.":  "70.54.171.141",
}

var certificate_validation_key string

type handler struct{}

func (this *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {

	msg := dns.Msg{}
	msg.SetReply(r)
	log.Println("-----> dns resquest receive... ", msg)
	switch r.Question[0].Qtype {
	case dns.TypeA:
		msg.Authoritative = true
		domain := msg.Question[0].Name
		address, ok := domainsToAddresses[domain]
		log.Println("---> ask for domain: ", domain, " address to redirect is ", address)
		if ok {
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: 60},
				A:   net.ParseIP(address),
			})
		}
	case dns.TypeTXT:
		log.Println("--------> dig for text value ", msg.Question[0].Name)
		log.Println("--------> certificate_validation_key", certificate_validation_key)

		msg.Answer = append(msg.Answer, &dns.TXT{
			Hdr: dns.RR_Header{Name: "", Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 60},
			Txt: []string{certificate_validation_key},
		})
	}
	w.WriteMsg(&msg)
}

func ServeDns(port int, key string) {
	certificate_validation_key = key
	// Now I will start the dns server.
	srv := &dns.Server{Addr: "127.0.0.1:" + strconv.Itoa(port), Net: "udp"}
	srv.Handler = &handler{}
	if err := srv.ListenAndServe(); err != nil {
		log.Println("Failed to set udp listener %s\n", err.Error())
	}
}
