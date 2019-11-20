package dns_client

import (
	"strconv"

	"github.com/davecourtois/Globular/api"
	"github.com/davecourtois/Globular/dns/dnspb"
	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// echo Client Service
////////////////////////////////////////////////////////////////////////////////

type DNS_Client struct {
	cc *grpc.ClientConn
	c  dnspb.DnsServiceClient

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
func NewDns_Client(domain string, port int, hasTLS bool, keyFile string, certFile string, caFile string) *DNS_Client {
	client := new(DNS_Client)
	client.domain = domain
	client.port = port
	client.name = "dns"
	client.hasTLS = hasTLS
	client.keyFile = keyFile
	client.certFile = certFile
	client.caFile = caFile
	client.cc = api.GetClientConnection(client)
	client.c = dnspb.NewDnsServiceClient(client.cc)
	return client
}

// Return the domain
func (self *DNS_Client) GetDomain() string {
	return self.domain
}

// Return the address
func (self *DNS_Client) GetAddress() string {
	return self.domain + ":" + strconv.Itoa(self.port)
}

// Return the name of the service
func (self *DNS_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *DNS_Client) Close() {
	self.cc.Close()
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

func (self *DNS_Client) GetA(domain string) (string, error) {

	// I will execute a simple ldap search here...
	rqst := &dnspb.GetARequest{
		Domain: domain,
	}

	rsp, err := self.c.GetA(api.GetClientContext(self), rqst)
	if err != nil {
		return "", err
	}
	return rsp.A, nil
}

func (self *DNS_Client) SetA(name string, ipv4 string, ttl uint32) (string, error) {

	// I will execute a simple ldap search here...
	rqst := &dnspb.SetARequest{
		Name: name,
		A:    ipv4,
		Ttl:  ttl,
	}

	rsp, err := self.c.SetA(api.GetClientContext(self), rqst)
	if err != nil {
		return "", err
	}
	return rsp.Message, nil
}

func (self *DNS_Client) RemoveA(name string) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.RemoveARequest{
		Name: name,
	}

	_, err := self.c.RemoveA(api.GetClientContext(self), rqst)
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

	rsp, err := self.c.GetAAAA(api.GetClientContext(self), rqst)
	if err != nil {
		return "", err
	}
	return rsp.Aaaa, nil
}

func (self *DNS_Client) SetAAAA(name string, ipv6 string, ttl uint32) (string, error) {

	// I will execute a simple ldap search here...
	rqst := &dnspb.SetAAAARequest{
		Name: name,
		Aaaa: ipv6,
		Ttl:  ttl,
	}

	rsp, err := self.c.SetAAAA(api.GetClientContext(self), rqst)
	if err != nil {
		return "", err
	}
	return rsp.Message, nil
}

func (self *DNS_Client) RemoveAAAA(name string) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.RemoveAAAARequest{
		Name: name,
	}

	_, err := self.c.RemoveAAAA(api.GetClientContext(self), rqst)
	if err != nil {
		return err
	}
	return nil
}

func (self *DNS_Client) GetText(id string) ([]string, error) {

	// I will execute a simple ldap search here...
	rqst := &dnspb.GetTextRequest{
		Id: id,
	}

	rsp, err := self.c.GetText(api.GetClientContext(self), rqst)
	if err != nil {
		return nil, err
	}

	return rsp.GetValues(), nil
}

func (self *DNS_Client) SetText(id string, values []string, ttl uint32) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.SetTextRequest{
		Id:     id,
		Values: values,
		Ttl:    ttl,
	}

	_, err := self.c.SetText(api.GetClientContext(self), rqst)
	return err
}

func (self *DNS_Client) RemoveText(id string) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.RemoveTextRequest{
		Id: id,
	}

	_, err := self.c.RemoveText(api.GetClientContext(self), rqst)
	if err != nil {
		return err
	}
	return nil
}

func (self *DNS_Client) GetNs(id string) (string, error) {

	// I will execute a simple ldap search here...
	rqst := &dnspb.GetNsRequest{
		Id: id,
	}

	rsp, err := self.c.GetNs(api.GetClientContext(self), rqst)
	if err != nil {
		return "", err
	}

	return rsp.GetNs(), nil
}

func (self *DNS_Client) SetNs(id string, ns string, ttl uint32) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.SetNsRequest{
		Id:  id,
		Ns:  ns,
		Ttl: ttl,
	}

	_, err := self.c.SetNs(api.GetClientContext(self), rqst)
	return err
}

func (self *DNS_Client) RemoveNs(id string) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.RemoveNsRequest{
		Id: id,
	}

	_, err := self.c.RemoveNs(api.GetClientContext(self), rqst)
	if err != nil {
		return err
	}
	return nil
}

func (self *DNS_Client) GetCName(id string) (string, error) {

	// I will execute a simple ldap search here...
	rqst := &dnspb.GetCNameRequest{
		Id: id,
	}

	rsp, err := self.c.GetCName(api.GetClientContext(self), rqst)
	if err != nil {
		return "", err
	}

	return rsp.GetCname(), nil
}

func (self *DNS_Client) SetCName(id string, cname string, ttl uint32) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.SetCNameRequest{
		Id:    id,
		Cname: cname,
		Ttl:   ttl,
	}

	_, err := self.c.SetCName(api.GetClientContext(self), rqst)
	return err
}

func (self *DNS_Client) RemoveCName(id string) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.RemoveCNameRequest{
		Id: id,
	}

	_, err := self.c.RemoveCName(api.GetClientContext(self), rqst)
	if err != nil {
		return err
	}
	return nil
}

func (self *DNS_Client) GetMx(id string) (map[string]interface{}, error) {

	// I will execute a simple ldap search here...
	rqst := &dnspb.GetMxRequest{
		Id: id,
	}

	rsp, err := self.c.GetMx(api.GetClientContext(self), rqst)
	if err != nil {
		return nil, err
	}

	mx := make(map[string]interface{}, 0)
	mx["Preference"] = uint16(rsp.GetResult().Preference)
	mx["Mx"] = rsp.GetResult().Mx

	return mx, nil
}

func (self *DNS_Client) SetMx(id string, preference uint16, mx string, ttl uint32) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.SetMxRequest{
		Id: id,
		Mx: &dnspb.MX{
			Preference: int32(preference),
			Mx:         mx,
		},
		Ttl: ttl,
	}

	_, err := self.c.SetMx(api.GetClientContext(self), rqst)
	return err
}

func (self *DNS_Client) RemoveMx(id string) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.RemoveMxRequest{
		Id: id,
	}

	_, err := self.c.RemoveMx(api.GetClientContext(self), rqst)
	if err != nil {
		return err
	}
	return nil
}

func (self *DNS_Client) GetSoa(id string) (map[string]interface{}, error) {

	// I will execute a simple ldap search here...
	rqst := &dnspb.GetSoaRequest{
		Id: id,
	}

	rsp, err := self.c.GetSoa(api.GetClientContext(self), rqst)
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

func (self *DNS_Client) SetSoa(id string, ns string, mbox string, serial uint32, refresh uint32, retry uint32, expire uint32, minttl uint32, ttl uint32) error {

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

	_, err := self.c.SetSoa(api.GetClientContext(self), rqst)
	return err
}

func (self *DNS_Client) RemoveSoa(id string) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.RemoveSoaRequest{
		Id: id,
	}

	_, err := self.c.RemoveSoa(api.GetClientContext(self), rqst)
	if err != nil {
		return err
	}
	return nil
}

func (self *DNS_Client) GetUri(id string) (map[string]interface{}, error) {

	// I will execute a simple ldap search here...
	rqst := &dnspb.GetUriRequest{
		Id: id,
	}

	rsp, err := self.c.GetUri(api.GetClientContext(self), rqst)
	if err != nil {
		return nil, err
	}

	uri := make(map[string]interface{}, 0)
	uri["Priority"] = rsp.GetResult().Priority
	uri["Weight"] = rsp.GetResult().Weight
	uri["Target"] = rsp.GetResult().Target

	return uri, nil
}

func (self *DNS_Client) SetUri(id string, priority uint32, weight uint32, target string, ttl uint32) error {

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

	_, err := self.c.SetUri(api.GetClientContext(self), rqst)
	return err
}

func (self *DNS_Client) RemoveUri(id string) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.RemoveUriRequest{
		Id: id,
	}

	_, err := self.c.RemoveUri(api.GetClientContext(self), rqst)
	if err != nil {
		return err
	}
	return nil
}

func (self *DNS_Client) GetCaa(id string) (map[string]interface{}, error) {

	// I will execute a simple ldap search here...
	rqst := &dnspb.GetCaaRequest{
		Id: id,
	}

	rsp, err := self.c.GetCaa(api.GetClientContext(self), rqst)
	if err != nil {
		return nil, err
	}

	caa := make(map[string]interface{}, 0)
	caa["Flag"] = rsp.GetResult().Flag
	caa["Tag"] = rsp.GetResult().Tag
	caa["Value"] = rsp.GetResult().Value

	return caa, nil
}

func (self *DNS_Client) SetCaa(id string, flag uint32, tag string, value string, ttl uint32) error {

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

	_, err := self.c.SetCaa(api.GetClientContext(self), rqst)
	return err
}

func (self *DNS_Client) RemoveCaa(id string) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.RemoveCaaRequest{
		Id: id,
	}

	_, err := self.c.RemoveCaa(api.GetClientContext(self), rqst)
	if err != nil {
		return err
	}
	return nil
}

func (self *DNS_Client) GetAfsdb(id string) (map[string]interface{}, error) {

	// I will execute a simple ldap search here...
	rqst := &dnspb.GetAfsdbRequest{
		Id: id,
	}

	rsp, err := self.c.GetAfsdb(api.GetClientContext(self), rqst)
	if err != nil {
		return nil, err
	}

	afsdb := make(map[string]interface{}, 0)
	afsdb["Subtype"] = rsp.GetResult().Subtype
	afsdb["Hostname"] = rsp.GetResult().Hostname

	return afsdb, nil
}

func (self *DNS_Client) SetAfsdb(id string, subtype uint32, hostname string, ttl uint32) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.SetAfsdbRequest{
		Id: id,
		Afsdb: &dnspb.AFSDB{
			Subtype:  subtype,
			Hostname: hostname,
		},
		Ttl: ttl,
	}

	_, err := self.c.SetAfsdb(api.GetClientContext(self), rqst)
	return err
}

func (self *DNS_Client) RemoveAfsdb(id string) error {

	// I will execute a simple ldap search here...
	rqst := &dnspb.RemoveAfsdbRequest{
		Id: id,
	}

	_, err := self.c.RemoveAfsdb(api.GetClientContext(self), rqst)
	if err != nil {
		return err
	}
	return nil
}
