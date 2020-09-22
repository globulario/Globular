package dns_client

import (
	"context"
	"strconv"

	"github.com/davecourtois/Globular/services/golang/dns/dnspb"
	globular "github.com/davecourtois/Globular/services/golang/globular_client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

////////////////////////////////////////////////////////////////////////////////
// echo Client Service
////////////////////////////////////////////////////////////////////////////////

type DNS_Client struct {
	cc *grpc.ClientConn
	c  dnspb.DnsServiceClient

	// The id of the service
	id string

	// The name of the service
	name string

	// The client domain
	domain string

	// The port
	port int

	// is the connection is secure?
	hasTLS bool

	// Link to client key file
	keyFile string

	// Link to client certificate file.
	certFile string

	// certificate authority file
	caFile string
}

// Create a connection to the service.
func NewDns_Client(address string, id string) (*DNS_Client, error) {
	client := new(DNS_Client)
	err := globular.InitClient(client, address, id)

	if err != nil {

		return nil, err
	}

	client.cc, err = globular.GetClientConnection(client)
	if err != nil {
		return nil, err
	}
	client.c = dnspb.NewDnsServiceClient(client.cc)
	return client, nil
}

func (self *DNS_Client) Invoke(method string, rqst interface{}, ctx context.Context) (interface{}, error) {
	if ctx == nil {
		ctx = globular.GetClientContext(self)
	}
	return globular.InvokeClientRequest(self.c, ctx, method, rqst)
}

// Return the domain
func (self *DNS_Client) GetDomain() string {
	return self.domain
}

// Return the address
func (self *DNS_Client) GetAddress() string {
	return self.domain + ":" + strconv.Itoa(self.port)
}

// Return the id of the service instance
func (self *DNS_Client) GetId() string {
	return self.id
}

// Return the name of the service
func (self *DNS_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *DNS_Client) Close() {
	self.cc.Close()
}

// Set grpc_service port.
func (self *DNS_Client) SetPort(port int) {
	self.port = port
}

// Set the client instance id.
func (self *DNS_Client) SetId(id string) {
	self.id = id
}

// Set the client name.
func (self *DNS_Client) SetName(name string) {
	self.name = name
}

// Set the domain.
func (self *DNS_Client) SetDomain(domain string) {
	self.domain = domain
}

////////////////// TLS ///////////////////

// Get if the client is secure.
func (self *DNS_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *DNS_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *DNS_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *DNS_Client) GetCaFile() string {
	return self.caFile
}

// Set the client is a secure client.
func (self *DNS_Client) SetTLS(hasTls bool) {
	self.hasTLS = hasTls
}

// Set TLS certificate file path
func (self *DNS_Client) SetCertFile(certFile string) {
	self.certFile = certFile
}

// Set TLS key file path
func (self *DNS_Client) SetKeyFile(keyFile string) {
	self.keyFile = keyFile
}

// Set TLS authority trust certificate file path
func (self *DNS_Client) SetCaFile(caFile string) {
	self.caFile = caFile
}

// The domain of the globule responsible to do ressource validation.
// That domain will be use by the interceptor and access validation will
// be evaluated by the ressource manager at the domain address.
func (self *DNS_Client) getDomainContext(domain string) context.Context {
	// Here I will set the targeted domain as domain in the context.
	md := metadata.New(map[string]string{"domain": domain})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	return ctx
}

///////////////// API ////////////////////
func (self *DNS_Client) GetA(domain string) (string, error) {

	// I will execute a simple ldap search here...
	rqst := &dnspb.GetARequest{
		Domain: domain,
	}

	rsp, err := self.c.GetA(self.getDomainContext(domain), rqst)
	if err != nil {
		return "", err
	}

	return rsp.A, nil
}

// Register a subdomain to a domain.
// ex: toto.globular.io is the subdomain to globular.io, so here
// domain will be globular.io and subdomain toto.globular.io. The validation will
// be done by globular.io and not the dns itself.
func (self *DNS_Client) SetA(domain string, subdomain string, ipv4 string, ttl uint32) (string, error) {

	// I will execute a simple ldap search here...
	rqst := &dnspb.SetARequest{
		Domain: subdomain,
		A:      ipv4,
		Ttl:    ttl,
	}

	rsp, err := self.c.SetA(self.getDomainContext(domain), rqst)
	if err != nil {
		return "", err
	}
	return rsp.Message, nil
}

func (self *DNS_Client) RemoveA(domain string) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.RemoveARequest{
		Domain: domain,
	}

	_, err := self.c.RemoveA(self.getDomainContext(domain), rqst)
	if err != nil {
		return err
	}
	return nil
}

func (self *DNS_Client) GetAAAA(domain string) (string, error) {

	// I will execute a simple ldap search here...
	rqst := &dnspb.GetAAAARequest{
		Domain: domain,
	}

	rsp, err := self.c.GetAAAA(self.getDomainContext(domain), rqst)
	if err != nil {
		return "", err
	}
	return rsp.Aaaa, nil
}

func (self *DNS_Client) SetAAAA(domain string, ipv6 string, ttl uint32) (string, error) {

	// I will execute a simple ldap search here...
	rqst := &dnspb.SetAAAARequest{
		Domain: domain,
		Aaaa:   ipv6,
		Ttl:    ttl,
	}

	rsp, err := self.c.SetAAAA(self.getDomainContext(domain), rqst)
	if err != nil {
		return "", err
	}
	return rsp.Message, nil
}

func (self *DNS_Client) RemoveAAAA(domain string) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.RemoveAAAARequest{
		Domain: domain,
	}

	_, err := self.c.RemoveAAAA(self.getDomainContext(domain), rqst)
	if err != nil {
		return err
	}
	return nil
}

func (self *DNS_Client) GetText(domain string, id string) ([]string, error) {

	// I will execute a simple ldap search here...
	rqst := &dnspb.GetTextRequest{
		Id: id,
	}

	rsp, err := self.c.GetText(self.getDomainContext(domain), rqst)
	if err != nil {
		return nil, err
	}

	return rsp.GetValues(), nil
}

func (self *DNS_Client) SetText(domain string, id string, values []string, ttl uint32) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.SetTextRequest{
		Id:     id,
		Values: values,
		Ttl:    ttl,
	}

	_, err := self.c.SetText(self.getDomainContext(domain), rqst)
	return err
}

func (self *DNS_Client) RemoveText(domain string, id string) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.RemoveTextRequest{
		Id: id,
	}

	_, err := self.c.RemoveText(self.getDomainContext(domain), rqst)
	if err != nil {
		return err
	}
	return nil
}

func (self *DNS_Client) GetNs(domain string, id string) (string, error) {

	// I will execute a simple ldap search here...
	rqst := &dnspb.GetNsRequest{
		Id: id,
	}

	rsp, err := self.c.GetNs(self.getDomainContext(domain), rqst)
	if err != nil {
		return "", err
	}

	return rsp.GetNs(), nil
}

func (self *DNS_Client) SetNs(domain string, id string, ns string, ttl uint32) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.SetNsRequest{
		Id:  id,
		Ns:  ns,
		Ttl: ttl,
	}

	_, err := self.c.SetNs(self.getDomainContext(domain), rqst)
	return err
}

func (self *DNS_Client) RemoveNs(domain string, id string) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.RemoveNsRequest{
		Id: id,
	}

	_, err := self.c.RemoveNs(self.getDomainContext(domain), rqst)
	if err != nil {
		return err
	}
	return nil
}

func (self *DNS_Client) GetCName(domain string, id string) (string, error) {

	// I will execute a simple ldap search here...
	rqst := &dnspb.GetCNameRequest{
		Id: id,
	}

	rsp, err := self.c.GetCName(self.getDomainContext(domain), rqst)
	if err != nil {
		return "", err
	}

	return rsp.GetCname(), nil
}

func (self *DNS_Client) SetCName(domain string, id string, cname string, ttl uint32) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.SetCNameRequest{
		Id:    id,
		Cname: cname,
		Ttl:   ttl,
	}

	_, err := self.c.SetCName(self.getDomainContext(domain), rqst)
	return err
}

func (self *DNS_Client) RemoveCName(domain string, id string) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.RemoveCNameRequest{
		Id: id,
	}

	_, err := self.c.RemoveCName(self.getDomainContext(domain), rqst)
	if err != nil {
		return err
	}
	return nil
}

func (self *DNS_Client) GetMx(domain string, id string) (map[string]interface{}, error) {

	// I will execute a simple ldap search here...
	rqst := &dnspb.GetMxRequest{
		Id: id,
	}

	rsp, err := self.c.GetMx(self.getDomainContext(domain), rqst)
	if err != nil {
		return nil, err
	}

	mx := make(map[string]interface{}, 0)
	mx["Preference"] = uint16(rsp.GetResult().Preference)
	mx["Mx"] = rsp.GetResult().Mx

	return mx, nil
}

func (self *DNS_Client) SetMx(domain string, id string, preference uint16, mx string, ttl uint32) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.SetMxRequest{
		Id: id,
		Mx: &dnspb.MX{
			Preference: int32(preference),
			Mx:         mx,
		},
		Ttl: ttl,
	}

	_, err := self.c.SetMx(self.getDomainContext(domain), rqst)
	return err
}

func (self *DNS_Client) RemoveMx(domain string, id string) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.RemoveMxRequest{
		Id: id,
	}

	_, err := self.c.RemoveMx(self.getDomainContext(domain), rqst)
	if err != nil {
		return err
	}
	return nil
}

func (self *DNS_Client) GetSoa(domain string, id string) (map[string]interface{}, error) {

	// I will execute a simple ldap search here...
	rqst := &dnspb.GetSoaRequest{
		Id: id,
	}

	rsp, err := self.c.GetSoa(self.getDomainContext(domain), rqst)
	if err != nil {
		return nil, err
	}

	soa := make(map[string]interface{}, 0)
	soa["Ns"] = rsp.GetResult().Ns
	soa["Mbox"] = rsp.GetResult().Mbox
	soa["Serial"] = rsp.GetResult().Serial
	soa["Refresh"] = rsp.GetResult().Refresh
	soa["Retry"] = rsp.GetResult().Retry
	soa["Expire"] = rsp.GetResult().Expire
	soa["Minttl"] = rsp.GetResult().Minttl

	return soa, nil
}

func (self *DNS_Client) SetSoa(domain string, id string, ns string, mbox string, serial uint32, refresh uint32, retry uint32, expire uint32, minttl uint32, ttl uint32) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.SetSoaRequest{
		Id: id,
		Soa: &dnspb.SOA{
			Ns:      ns,
			Mbox:    mbox,
			Serial:  serial,
			Refresh: refresh,
			Retry:   retry,
			Expire:  expire,
			Minttl:  minttl,
		},
		Ttl: ttl,
	}

	_, err := self.c.SetSoa(self.getDomainContext(domain), rqst)
	return err
}

func (self *DNS_Client) RemoveSoa(domain string, id string) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.RemoveSoaRequest{
		Id: id,
	}

	_, err := self.c.RemoveSoa(self.getDomainContext(domain), rqst)
	if err != nil {
		return err
	}
	return nil
}

func (self *DNS_Client) GetUri(domain string, id string) (map[string]interface{}, error) {

	// I will execute a simple ldap search here...
	rqst := &dnspb.GetUriRequest{
		Id: id,
	}

	rsp, err := self.c.GetUri(self.getDomainContext(domain), rqst)
	if err != nil {
		return nil, err
	}

	uri := make(map[string]interface{}, 0)
	uri["Priority"] = rsp.GetResult().Priority
	uri["Weight"] = rsp.GetResult().Weight
	uri["Target"] = rsp.GetResult().Target

	return uri, nil
}

func (self *DNS_Client) SetUri(domain string, id string, priority uint32, weight uint32, target string, ttl uint32) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.SetUriRequest{
		Id: id,
		Uri: &dnspb.URI{
			Priority: priority,
			Weight:   weight,
			Target:   target,
		},
		Ttl: ttl,
	}

	_, err := self.c.SetUri(self.getDomainContext(domain), rqst)
	return err
}

func (self *DNS_Client) RemoveUri(domain string, id string) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.RemoveUriRequest{
		Id: id,
	}

	_, err := self.c.RemoveUri(self.getDomainContext(domain), rqst)
	if err != nil {
		return err
	}
	return nil
}

func (self *DNS_Client) GetCaa(domain string, id string) (map[string]interface{}, error) {

	// I will execute a simple ldap search here...
	rqst := &dnspb.GetCaaRequest{
		Id: id,
	}

	rsp, err := self.c.GetCaa(self.getDomainContext(domain), rqst)
	if err != nil {
		return nil, err
	}

	caa := make(map[string]interface{}, 0)
	caa["Flag"] = rsp.GetResult().Flag
	caa["Tag"] = rsp.GetResult().Tag
	caa["Value"] = rsp.GetResult().Value

	return caa, nil
}

func (self *DNS_Client) SetCaa(domain string, id string, flag uint32, tag string, value string, ttl uint32) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.SetCaaRequest{
		Id: id,
		Caa: &dnspb.CAA{
			Flag:  flag,
			Tag:   tag,
			Value: value,
		},
		Ttl: ttl,
	}

	_, err := self.c.SetCaa(self.getDomainContext(domain), rqst)
	return err
}

func (self *DNS_Client) RemoveCaa(domain string, id string) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.RemoveCaaRequest{
		Id: id,
	}

	_, err := self.c.RemoveCaa(self.getDomainContext(domain), rqst)
	if err != nil {
		return err
	}
	return nil
}

func (self *DNS_Client) GetAfsdb(domain string, id string) (map[string]interface{}, error) {

	// I will execute a simple ldap search here...
	rqst := &dnspb.GetAfsdbRequest{
		Id: id,
	}

	rsp, err := self.c.GetAfsdb(self.getDomainContext(domain), rqst)
	if err != nil {
		return nil, err
	}

	afsdb := make(map[string]interface{}, 0)
	afsdb["Subtype"] = rsp.GetResult().Subtype
	afsdb["Hostname"] = rsp.GetResult().Hostname

	return afsdb, nil
}

func (self *DNS_Client) SetAfsdb(domain string, id string, subtype uint32, hostname string, ttl uint32) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.SetAfsdbRequest{
		Id: id,
		Afsdb: &dnspb.AFSDB{
			Subtype:  subtype,
			Hostname: hostname,
		},
		Ttl: ttl,
	}

	_, err := self.c.SetAfsdb(self.getDomainContext(domain), rqst)
	return err
}

func (self *DNS_Client) RemoveAfsdb(domain string, id string) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.RemoveAfsdbRequest{
		Id: id,
	}

	_, err := self.c.RemoveAfsdb(self.getDomainContext(domain), rqst)
	if err != nil {
		return err
	}
	return nil
}
