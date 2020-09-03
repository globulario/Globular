package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/davecourtois/Globular/Interceptors"
	"github.com/davecourtois/Globular/dns/dnspb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	//"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	//"google.golang.org/grpc/grpclog"
	"github.com/davecourtois/Globular/dns/dns_client"
	"google.golang.org/grpc/reflection"

	//"google.golang.org/grpc/status"
	"encoding/binary"

	"github.com/davecourtois/Globular/storage/storage_store"
	"github.com/miekg/dns"
)

// TODO take care of TLS/https
var (
	defaultPort  = 10033
	defaultProxy = 10034

	// By default all origins are allowed.
	allow_all_origins = true

	// comma separeated values.
	allowed_origins string = ""

	// Thr IPV4 address
	domain       string = "localhost"
	connectionId string = "dns_service"

	// pointer to the sever implementation.
	s *server
)

// Value need by Globular to start the services...
type server struct {
	// The global attribute of the services.
	Id              string
	Name            string
	Path            string
	Proto           string
	Port            int
	Proxy           int
	AllowAllOrigins bool
	AllowedOrigins  string // comma separated string.
	Protocol        string
	Domain          string
	// self-signed X.509 public keys for distribution
	CertFile string
	// a private RSA key to sign and authenticate the public key
	KeyFile string
	// a private RSA key to sign and authenticate the public key
	CertAuthorityTrust string
	TLS                bool
	Version            string
	PublisherId        string
	KeepUpToDate       bool
	KeepAlive          bool
	Permissions        []interface{} // contains the action permission for the services.

	// Contain the configuration of the storage service use to store
	// the actual values.
	DnsPort         int      // the dns port
	Domains         []string // The list of domains managed by this dns.
	StorageDataPath string

	store *storage_store.LevelDB_store

	// The dns records... https://en.wikipedia.org/wiki/List_of_DNS_record_types
	// see the wikipedia page to know exactly what are the values that can
	// be use here.
	Records map[string][]interface{}

	connection_is_open bool
}

// Create the configuration file if is not already exist.
func (self *server) init() {

	// That function is use to get access to other server.
	Utility.RegisterFunction("NewDns_Client", dns_client.NewDns_Client)

	// Here I will retreive the list of connections from file if there are some...
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	file, err := ioutil.ReadFile(dir + "/config.json")

	// default value.
	self.DnsPort = 5353 // 53 is the default dns port.
	self.Domains = make([]string, 0)
	self.StorageDataPath = os.TempDir()

	if err == nil {
		json.Unmarshal([]byte(file), self)
	} else {
		if len(self.Id) == 0 {
			// Generate random id for the server instance.
			self.Id = Utility.RandomUUID()
		}
		self.save()
	}

	if len(self.StorageDataPath) == 0 {
		fmt.Println("The value StorageDataPath in the configuration must be given. You can use /tmp (on linux) if you don't want to keep values indefilnely on the storage server.")
	}

}

// Open the connection if it's close.
func (self *server) openConnection() error {
	if self.connection_is_open == true {
		return nil
	}

	// Open store.
	self.store = storage_store.NewLevelDB_store()
	err := self.store.Open(`{"path":"` + self.StorageDataPath + `", "name":"dns_data_store"}`)
	if err != nil {
		return err
	}

	self.connection_is_open = true

	// Init the records with that connection.
	self.initRecords()

	return nil
}

// Save the configuration values.
func (self *server) save() error {

	// Create the file...
	str, err := Utility.ToJson(self)
	if err != nil {
		return err
	}

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	ioutil.WriteFile(dir+"/config.json", []byte(str), 0644)
	return nil
}

func (self *server) isManaged(domain string) bool {
	for i := 0; i < len(self.Domains); i++ {
		if strings.HasSuffix(domain, self.Domains[i]) {
			return true
		}
	}
	return false
}

// Set a dns entry.
func (self *server) SetA(ctx context.Context, rqst *dnspb.SetARequest) (*dnspb.SetAResponse, error) {

	fmt.Println("Try set dns entry ", rqst.Domain)
	if !self.isManaged(rqst.Domain) {
		err := errors.New("The domain " + rqst.Domain + " is not manage by this dns.")
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))

	}

	domain := strings.ToLower(rqst.Domain)

	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	uuid := Utility.GenerateUUID("A:" + domain)
	err = self.store.SetItem(uuid, []byte(rqst.A))
	if err != nil {

		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	fmt.Println("domain: ", "A:"+domain, " with uuid", uuid, "is set with ipv4 address", rqst.A)
	self.setTtl(uuid, rqst.Ttl)

	return &dnspb.SetAResponse{
		Message: domain, // return the full domain.
	}, nil
}

func (self *server) RemoveA(ctx context.Context, rqst *dnspb.RemoveARequest) (*dnspb.RemoveAResponse, error) {
	fmt.Println("Try remove dns entry ", rqst.Domain)
	if !self.isManaged(rqst.Domain) {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("The domain "+rqst.Domain+" is not manage by this dns.")))
	}

	domain := strings.ToLower(rqst.Domain)
	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	uuid := Utility.GenerateUUID("A:" + domain)
	err = self.store.RemoveItem(uuid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &dnspb.RemoveAResponse{
		Result: true, // return the full domain.
	}, nil
}

func (self *server) get_ipv4(domain string) (string, uint32, error) {
	domain = strings.ToLower(domain)
	if strings.HasSuffix(domain, ".") {
		domain = domain[:len(domain)-1]
	}
	err := self.openConnection()
	if err != nil {
		return "", 0, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	uuid := Utility.GenerateUUID("A:" + domain)
	ipv4, err := self.store.GetItem(uuid)
	if err != nil {
		return "", 0, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return string(ipv4), self.getTtl(uuid), nil
}

func (self *server) GetA(ctx context.Context, rqst *dnspb.GetARequest) (*dnspb.GetAResponse, error) {
	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	domain := strings.ToLower(rqst.Domain)
	uuid := Utility.GenerateUUID("A:" + domain)
	fmt.Println("GetA --> try to find value: ", "A:"+rqst.Domain)
	ipv4, err := self.store.GetItem(uuid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	fmt.Println("ipv4 for", domain, "is", string(ipv4))
	return &dnspb.GetAResponse{
		A: string(ipv4), // return the full domain.
	}, nil
}

// Set a dns entry.
func (self *server) SetAAAA(ctx context.Context, rqst *dnspb.SetAAAARequest) (*dnspb.SetAAAAResponse, error) {
	fmt.Println("Try set dns entry ", rqst.Domain)
	if !self.isManaged(rqst.Domain) {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("The domain "+rqst.Domain+" is not manage by this dns.")))
	}

	domain := strings.ToLower(rqst.Domain)

	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	uuid := Utility.GenerateUUID("AAAA:" + domain)
	err = self.store.SetItem(uuid, []byte(rqst.Aaaa))
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	self.setTtl(uuid, rqst.Ttl)

	return &dnspb.SetAAAAResponse{
		Message: domain, // return the full domain.
	}, nil
}

func (self *server) RemoveAAAA(ctx context.Context, rqst *dnspb.RemoveAAAARequest) (*dnspb.RemoveAAAAResponse, error) {
	domain := strings.ToLower(rqst.Domain)
	fmt.Println("Try remove dns entry ", domain)
	if !self.isManaged(rqst.Domain) {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("The domain "+rqst.Domain+" is not manage by this dns.")))
	}

	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	uuid := Utility.GenerateUUID("AAAA:" + domain)
	err = self.store.RemoveItem(uuid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &dnspb.RemoveAAAAResponse{
		Result: true, // return the full domain.
	}, nil
}

func (self *server) get_ipv6(domain string) (string, uint32, error) {
	domain = strings.ToLower(domain)
	if strings.HasSuffix(domain, ".") {
		domain = domain[:len(domain)-1]
	}
	fmt.Println("Try get dns entry ", domain)
	err := self.openConnection()
	if err != nil {
		return "", 0, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	uuid := Utility.GenerateUUID("AAAA:" + domain)
	address, err := self.store.GetItem(uuid)
	if err != nil {
		return "", 0, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	return string(address), self.getTtl(uuid), nil
}

func (self *server) GetAAAA(ctx context.Context, rqst *dnspb.GetAAAARequest) (*dnspb.GetAAAAResponse, error) {
	domain := strings.ToLower(rqst.Domain)
	fmt.Println("Try get dns entry ", domain)

	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	uuid := Utility.GenerateUUID("AAAA:" + domain)
	ipv6, err := self.store.GetItem(uuid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	fmt.Println("ipv6 for", domain, "is", string(ipv6))
	return &dnspb.GetAAAAResponse{
		Aaaa: string(ipv6), // return the full domain.
	}, nil
}

// Set a text entry.
func (self *server) SetText(ctx context.Context, rqst *dnspb.SetTextRequest) (*dnspb.SetTextResponse, error) {
	fmt.Println("Try set dns text ", rqst.Id)

	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	values, err := json.Marshal(rqst.Values)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	uuid := Utility.GenerateUUID("TXT:" + rqst.Id)
	err = self.store.SetItem(uuid, values)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	self.setTtl(uuid, rqst.Ttl)

	return &dnspb.SetTextResponse{
		Result: true, // return the full domain.
	}, nil
}

// return the text.
func (self *server) getText(id string) ([]string, uint32, error) {
	fmt.Println("Try get dns text ", id)
	err := self.openConnection()
	if err != nil {
		return nil, 0, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	uuid := Utility.GenerateUUID("TXT:" + id)
	data, err := self.store.GetItem(uuid)
	if err != nil {
	}
	values := make([]string, 0)

	err = json.Unmarshal(data, &values)
	if err != nil {
		return nil, 0, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	return values, self.getTtl(uuid), nil
}

// Retreive a text value
func (self *server) GetText(ctx context.Context, rqst *dnspb.GetTextRequest) (*dnspb.GetTextResponse, error) {
	fmt.Println("Try get dns text ", domain)
	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	uuid := Utility.GenerateUUID("TXT:" + rqst.Id)
	data, err := self.store.GetItem(uuid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	values := make([]string, 0)
	err = json.Unmarshal(data, &values)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &dnspb.GetTextResponse{
		Values: values, // return the full domain.
	}, nil
}

// Remove a text entry
func (self *server) RemoveText(ctx context.Context, rqst *dnspb.RemoveTextRequest) (*dnspb.RemoveTextResponse, error) {
	fmt.Println("Try remove dns text ", rqst.Id)

	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	uuid := Utility.GenerateUUID("TXT:" + rqst.Id)
	err = self.store.RemoveItem(uuid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &dnspb.RemoveTextResponse{
		Result: true, // return the full domain.
	}, nil
}

// Set a NS entry.
func (self *server) SetNs(ctx context.Context, rqst *dnspb.SetNsRequest) (*dnspb.SetNsResponse, error) {
	fmt.Println("Try set dns ns ", rqst.Id)

	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	uuid := Utility.GenerateUUID("NS:" + rqst.Id)
	err = self.store.SetItem(uuid, []byte(rqst.Ns))
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	self.setTtl(uuid, rqst.Ttl)

	return &dnspb.SetNsResponse{
		Result: true, // return the full domain.
	}, nil
}

// return the text.
func (self *server) getNs(id string) (string, uint32, error) {
	fmt.Println("Try get dns ns ", id)
	err := self.openConnection()
	if err != nil {
		return "", 0, err
	}
	uuid := Utility.GenerateUUID("NS:" + id)
	data, err := self.store.GetItem(uuid)
	if err != nil {
		return "", 0, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	return string(data), self.getTtl(uuid), err
}

// Retreive a text value
func (self *server) GetNs(ctx context.Context, rqst *dnspb.GetNsRequest) (*dnspb.GetNsResponse, error) {

	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	uuid := Utility.GenerateUUID("NS:" + rqst.Id)
	data, err := self.store.GetItem(uuid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &dnspb.GetNsResponse{
		Ns: string(data), // return the full domain.
	}, nil
}

// Remove a text entry
func (self *server) RemoveNs(ctx context.Context, rqst *dnspb.RemoveNsRequest) (*dnspb.RemoveNsResponse, error) {
	fmt.Println("Try remove dns text ", rqst.Id)

	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	uuid := Utility.GenerateUUID("NS:" + rqst.Id)
	err = self.store.RemoveItem(uuid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &dnspb.RemoveNsResponse{
		Result: true, // return the full domain.
	}, nil
}

// Set a CName entry.
func (self *server) SetCName(ctx context.Context, rqst *dnspb.SetCNameRequest) (*dnspb.SetCNameResponse, error) {
	fmt.Println("Try set dns CName ", rqst.Id)

	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	uuid := Utility.GenerateUUID("CName:" + rqst.Id)
	err = self.store.SetItem(uuid, []byte(rqst.Cname))
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	self.setTtl(uuid, rqst.Ttl)

	return &dnspb.SetCNameResponse{
		Result: true, // return the full domain.
	}, nil
}

// return the CName.
func (self *server) getCName(id string) (string, uint32, error) {
	fmt.Println("Try get CName ", id)
	err := self.openConnection()
	if err != nil {
		return "", 0, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	uuid := Utility.GenerateUUID("CName:" + id)
	data, err := self.store.GetItem(uuid)
	if err != nil {
		return "", 0, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	return string(data), self.getTtl(uuid), err
}

// Retreive a CName value
func (self *server) GetCName(ctx context.Context, rqst *dnspb.GetCNameRequest) (*dnspb.GetCNameResponse, error) {

	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	uuid := Utility.GenerateUUID("CName:" + rqst.Id)
	data, err := self.store.GetItem(uuid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &dnspb.GetCNameResponse{
		Cname: string(data), // return the full domain.
	}, nil
}

// Remove a text entry
func (self *server) RemoveCName(ctx context.Context, rqst *dnspb.RemoveCNameRequest) (*dnspb.RemoveCNameResponse, error) {
	fmt.Println("Try remove dns CName ", rqst.Id)

	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	uuid := Utility.GenerateUUID("CName:" + rqst.Id)
	err = self.store.RemoveItem(uuid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &dnspb.RemoveCNameResponse{
		Result: true, // return the full domain.
	}, nil
}

// Set a text entry.
func (self *server) SetMx(ctx context.Context, rqst *dnspb.SetMxRequest) (*dnspb.SetMxResponse, error) {
	fmt.Println("Try set dns mx ", rqst.Id)

	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	values, err := json.Marshal(rqst.Mx)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	uuid := Utility.GenerateUUID("MX:" + rqst.Id)
	err = self.store.SetItem(uuid, values)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	self.setTtl(uuid, rqst.Ttl)

	return &dnspb.SetMxResponse{
		Result: true, // return the full domain.
	}, nil
}

// return the text.
func (self *server) getMx(id string) (map[string]interface{}, uint32, error) {
	fmt.Println("Try get dns text ", id)
	err := self.openConnection()
	if err != nil {
		return nil, 0, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	uuid := Utility.GenerateUUID("MX:" + id)
	data, err := self.store.GetItem(uuid)

	values := make(map[string]interface{}, 0) // use a map instead of Mx struct.
	err = json.Unmarshal(data, &values)
	if err != nil {
		return nil, 0, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	return values, self.getTtl(uuid), nil
}

// Retreive a text value
func (self *server) GetMx(ctx context.Context, rqst *dnspb.GetMxRequest) (*dnspb.GetMxResponse, error) {
	fmt.Println("Try get dns mx ", domain)
	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	uuid := Utility.GenerateUUID("MX:" + rqst.Id)
	data, err := self.store.GetItem(uuid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	values := make(map[string]interface{}, 0)
	err = json.Unmarshal(data, &values)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &dnspb.GetMxResponse{
		Result: &dnspb.MX{
			Preference: values["Preference"].(int32),
			Mx:         values["Mx"].(string),
		}, // return the full domain.
	}, nil
}

// Remove a text entry
func (self *server) RemoveMx(ctx context.Context, rqst *dnspb.RemoveMxRequest) (*dnspb.RemoveMxResponse, error) {
	fmt.Println("Try remove dns text ", rqst.Id)

	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	uuid := Utility.GenerateUUID("MX:" + rqst.Id)
	err = self.store.RemoveItem(uuid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &dnspb.RemoveMxResponse{
		Result: true, // return the full domain.
	}, nil
}

// Set a text entry.
func (self *server) SetSoa(ctx context.Context, rqst *dnspb.SetSoaRequest) (*dnspb.SetSoaResponse, error) {
	fmt.Println("Try set dns Soa ", rqst.Id)

	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	values, err := json.Marshal(rqst.Soa)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	uuid := Utility.GenerateUUID("SOA:" + rqst.Id)
	err = self.store.SetItem(uuid, values)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	self.setTtl(uuid, rqst.Ttl)

	return &dnspb.SetSoaResponse{
		Result: true, // return the full domain.
	}, nil
}

// return the text.
func (self *server) getSoa(id string) (*dnspb.SOA, uint32, error) {
	fmt.Println("Try get dns soa ", id)
	err := self.openConnection()
	if err != nil {
		return nil, 0, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	uuid := Utility.GenerateUUID("SOA:" + id)
	data, err := self.store.GetItem(uuid)

	soa := new(dnspb.SOA) // use a map instead of Mx struct.
	err = json.Unmarshal(data, soa)
	if err != nil {
		return nil, 0, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	return soa, self.getTtl(uuid), nil
}

// Retreive a text value
func (self *server) GetSoa(ctx context.Context, rqst *dnspb.GetSoaRequest) (*dnspb.GetSoaResponse, error) {
	fmt.Println("Try get dns soa ", domain)
	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	uuid := Utility.GenerateUUID("SOA:" + rqst.Id)
	data, err := self.store.GetItem(uuid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	soa := new(dnspb.SOA)
	err = json.Unmarshal(data, soa)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &dnspb.GetSoaResponse{
		Result: soa, // return the full domain.
	}, nil
}

// Remove a text entry
func (self *server) RemoveSoa(ctx context.Context, rqst *dnspb.RemoveSoaRequest) (*dnspb.RemoveSoaResponse, error) {
	fmt.Println("Try remove dns text ", rqst.Id)

	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	uuid := Utility.GenerateUUID("SOA:" + rqst.Id)
	err = self.store.RemoveItem(uuid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &dnspb.RemoveSoaResponse{
		Result: true, // return the full domain.
	}, nil
}

// Set a text entry.
func (self *server) SetUri(ctx context.Context, rqst *dnspb.SetUriRequest) (*dnspb.SetUriResponse, error) {
	fmt.Println("Try set dns uri ", rqst.Id)
	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	values, err := json.Marshal(rqst.Uri)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	uuid := Utility.GenerateUUID("URI:" + rqst.Id)
	err = self.store.SetItem(uuid, values)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	self.setTtl(uuid, rqst.Ttl)

	return &dnspb.SetUriResponse{
		Result: true, // return the full domain.
	}, nil
}

// return the text.
func (self *server) getUri(id string) (*dnspb.URI, uint32, error) {
	fmt.Println("Try get dns uri ", id)
	err := self.openConnection()
	if err != nil {
		return nil, 0, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	uuid := Utility.GenerateUUID("URI:" + id)
	data, err := self.store.GetItem(uuid)

	uri := new(dnspb.URI) // use a map instead of Mx struct.
	err = json.Unmarshal(data, uri)
	if err != nil {
		return nil, 0, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	return uri, self.getTtl(uuid), nil
}

// Retreive a text value
func (self *server) GetUri(ctx context.Context, rqst *dnspb.GetUriRequest) (*dnspb.GetUriResponse, error) {
	fmt.Println("Try get dns uri ", domain)
	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	uuid := Utility.GenerateUUID("URI:" + rqst.Id)
	data, err := self.store.GetItem(uuid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	uri := new(dnspb.URI)
	err = json.Unmarshal(data, uri)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &dnspb.GetUriResponse{
		Result: uri, // return the full domain.
	}, nil
}

// Remove AFSDB
func (self *server) RemoveUri(ctx context.Context, rqst *dnspb.RemoveUriRequest) (*dnspb.RemoveUriResponse, error) {
	fmt.Println("Try remove dns uri ", rqst.Id)

	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	uuid := Utility.GenerateUUID("URI:" + rqst.Id)
	err = self.store.RemoveItem(uuid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &dnspb.RemoveUriResponse{
		Result: true, // return the full domain.
	}, nil
}

// Set a AFSDB entry.
func (self *server) SetAfsdb(ctx context.Context, rqst *dnspb.SetAfsdbRequest) (*dnspb.SetAfsdbResponse, error) {
	fmt.Println("Try set dns afsdb ", rqst.Id)
	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	values, err := json.Marshal(rqst.Afsdb)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	uuid := Utility.GenerateUUID("AFSDB:" + rqst.Id)
	err = self.store.SetItem(uuid, values)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	self.setTtl(uuid, rqst.Ttl)

	return &dnspb.SetAfsdbResponse{
		Result: true, // return the full domain.
	}, nil
}

// return the AFSDB.
func (self *server) getAfsdb(id string) (*dnspb.AFSDB, uint32, error) {
	fmt.Println("Try get dns AFSDB ", id)
	err := self.openConnection()
	if err != nil {
		return nil, 0, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	uuid := Utility.GenerateUUID("AFSDB:" + id)
	data, err := self.store.GetItem(uuid)

	afsdb := new(dnspb.AFSDB) // use a map instead of Mx struct.
	err = json.Unmarshal(data, afsdb)
	if err != nil {
		return nil, 0, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	return afsdb, self.getTtl(uuid), nil
}

// Retreive a AFSDB value
func (self *server) GetAfsdb(ctx context.Context, rqst *dnspb.GetAfsdbRequest) (*dnspb.GetAfsdbResponse, error) {
	fmt.Println("Try get dns AFSDB ", domain)
	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	uuid := Utility.GenerateUUID("AFSDB:" + rqst.Id)
	data, err := self.store.GetItem(uuid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	afsdb := new(dnspb.AFSDB)
	err = json.Unmarshal(data, afsdb)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &dnspb.GetAfsdbResponse{
		Result: afsdb, // return the full domain.
	}, nil
}

// Remove AFSDB
func (self *server) RemoveAfsdb(ctx context.Context, rqst *dnspb.RemoveAfsdbRequest) (*dnspb.RemoveAfsdbResponse, error) {
	fmt.Println("Try remove dns Afsdb ", rqst.Id)

	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	uuid := Utility.GenerateUUID("AFSDB:" + rqst.Id)
	err = self.store.RemoveItem(uuid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &dnspb.RemoveAfsdbResponse{
		Result: true, // return the full domain.
	}, nil
}

// Set a CAA entry.
func (self *server) SetCaa(ctx context.Context, rqst *dnspb.SetCaaRequest) (*dnspb.SetCaaResponse, error) {
	fmt.Println("Try set dns caa ", rqst.Id)
	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	values, err := json.Marshal(rqst.Caa)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	uuid := Utility.GenerateUUID("CAA:" + rqst.Id)
	err = self.store.SetItem(uuid, values)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	self.setTtl(uuid, rqst.Ttl)

	return &dnspb.SetCaaResponse{
		Result: true, // return the full domain.
	}, nil
}

// return the CAA.
func (self *server) getCaa(id string) (*dnspb.CAA, uint32, error) {
	fmt.Println("Try get dns CAA ", id)
	err := self.openConnection()
	if err != nil {
		return nil, 0, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	uuid := Utility.GenerateUUID("CAA:" + id)
	data, err := self.store.GetItem(uuid)

	caa := new(dnspb.CAA) // use a map instead of Mx struct.
	err = json.Unmarshal(data, caa)
	if err != nil {
		return nil, 0, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}
	return caa, self.getTtl(uuid), nil
}

// Retreive a AFSDB value
func (self *server) GetCaa(ctx context.Context, rqst *dnspb.GetCaaRequest) (*dnspb.GetCaaResponse, error) {
	fmt.Println("Try get dns CAA ", domain)
	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	uuid := Utility.GenerateUUID("CAA:" + rqst.Id)
	data, err := self.store.GetItem(uuid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	caa := new(dnspb.CAA)
	err = json.Unmarshal(data, caa)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &dnspb.GetCaaResponse{
		Result: caa, // return the full domain.
	}, nil
}

// Remove CAA
func (self *server) RemoveCaa(ctx context.Context, rqst *dnspb.RemoveCaaRequest) (*dnspb.RemoveCaaResponse, error) {
	fmt.Println("Try remove dns Afsdb ", rqst.Id)

	err := self.openConnection()
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	uuid := Utility.GenerateUUID("CAA:" + rqst.Id)
	err = self.store.RemoveItem(uuid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &dnspb.RemoveCaaResponse{
		Result: true, // return the full domain.
	}, nil
}

/////////////////////// DNS Specific service //////////////////////
type handler struct{}

func (this *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := dns.Msg{}
	msg.SetReply(r)
	fmt.Println("-----> dns resquest receive... ", msg)
	switch r.Question[0].Qtype {
	case dns.TypeA:
		domain := msg.Question[0].Name
		msg.Authoritative = true
		address, ttl, err := s.get_ipv4(domain) // get the address name from the
		fmt.Println("---> look for A ", domain)
		if err == nil {
			fmt.Println("---> ask for domain: ", domain, " address to redirect is ", address)
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: ttl},
				A:   net.ParseIP(address),
			})
		} else {
			fmt.Println("fail to retreive ipv4 address for "+domain, err)
		}

	case dns.TypeAAAA:
		msg.Authoritative = true
		domain := msg.Question[0].Name
		address, ttl, err := s.get_ipv6(domain) // get the address name from the
		fmt.Println("---> look for AAAA ", domain)
		if err == nil {
			fmt.Println("---> ask for domain: ", domain, " address to redirect is ", address)
			msg.Answer = append(msg.Answer, &dns.AAAA{
				Hdr:  dns.RR_Header{Name: domain, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: ttl},
				AAAA: net.ParseIP(address),
			})
		} else {
			fmt.Println(err)
		}

	case dns.TypeAFSDB:

		msg.Authoritative = true
		id := msg.Question[0].Name
		afsdb, ttl, err := s.getAfsdb(id)
		if err == nil {
			msg.Answer = append(msg.Answer, &dns.AFSDB{
				Hdr:      dns.RR_Header{Name: domain, Rrtype: dns.TypeAFSDB, Class: dns.ClassINET, Ttl: ttl},
				Subtype:  uint16(afsdb.Subtype), //
				Hostname: afsdb.Hostname,
			})
		} else {
			fmt.Println(err)
		}

	case dns.TypeCAA:
		msg.Authoritative = true
		name := msg.Question[0].Name
		fmt.Println("---> look for CAA ", name)
		caa, ttl, err := s.getCaa(name)

		if err == nil {
			msg.Answer = append(msg.Answer, &dns.CAA{
				Hdr:   dns.RR_Header{Name: name, Rrtype: dns.TypeCAA, Class: dns.ClassINET, Ttl: ttl},
				Flag:  uint8(caa.Flag),
				Tag:   caa.Tag,
				Value: caa.Value,
			})
		} else {
			fmt.Println(err)
		}

	case dns.TypeCNAME:
		id := msg.Question[0].Name
		cname, ttl, err := s.getCName(id)
		if err == nil {
			// in case of empty string I will return the certificate validation key.
			msg.Answer = append(msg.Answer, &dns.CNAME{
				// keep text values.
				Hdr:    dns.RR_Header{Name: id, Rrtype: dns.TypeCNAME, Class: dns.ClassINET, Ttl: ttl},
				Target: cname,
			})
		}

	case dns.TypeTXT:
		id := msg.Question[0].Name
		fmt.Println("---> look for txt ", id)
		values, ttl, err := s.getText(id)
		if err == nil {
			// in case of empty string I will return the certificate validation key.
			msg.Answer = append(msg.Answer, &dns.TXT{
				// keep text values.
				Hdr: dns.RR_Header{Name: id, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: ttl},
				Txt: values,
			})
		} else {
			fmt.Println("fail to retreive txt!", err)
		}

	case dns.TypeNS:
		id := msg.Question[0].Name
		ns, ttl, err := s.getNs(id)
		if err == nil {
			// in case of empty string I will return the certificate validation key.
			msg.Answer = append(msg.Answer, &dns.NS{
				// keep text values.
				Hdr: dns.RR_Header{Name: id, Rrtype: dns.TypeNS, Class: dns.ClassINET, Ttl: ttl},
				Ns:  ns,
			})
		}

	case dns.TypeMX:
		id := msg.Question[0].Name // the id where the values is store.
		mx, ttl, err := s.getMx(id)

		if err == nil {
			// in case of empty string I will return the certificate validation key.
			msg.Answer = append(msg.Answer, &dns.MX{
				// keep text values.
				Hdr:        dns.RR_Header{Name: id, Rrtype: dns.TypeMX, Class: dns.ClassINET, Ttl: ttl},
				Preference: uint16(mx["Preference"].(int32)),
				Mx:         mx["Mx"].(string),
			})
		}

	case dns.TypeSOA:
		id := msg.Question[0].Name
		fmt.Println("---> look for soa ", id)
		soa, ttl, err := s.getSoa(id)
		if err == nil {
			// in case of empty string I will return the certificate validation key.
			msg.Answer = append(msg.Answer, &dns.SOA{
				// keep text values.
				Hdr:     dns.RR_Header{Name: id, Rrtype: dns.TypeSOA, Class: dns.ClassINET, Ttl: ttl},
				Ns:      soa.Ns,
				Mbox:    soa.Mbox,
				Serial:  soa.Serial,
				Refresh: soa.Refresh,
				Retry:   soa.Retry,
				Expire:  soa.Expire,
				Minttl:  soa.Minttl,
			})
		}

	case dns.TypeURI:
		id := msg.Question[0].Name
		fmt.Println("---> look for uri ", id)
		uri, ttl, err := s.getUri(id)
		if err == nil {
			// in case of empty string I will return the certificate validation key.
			msg.Answer = append(msg.Answer, &dns.URI{
				// keep text values.
				Hdr:      dns.RR_Header{Name: id, Rrtype: dns.TypeURI, Class: dns.ClassINET, Ttl: ttl},
				Priority: uint16(uri.Priority),
				Weight:   uint16(uri.Weight),
				Target:   uri.Target,
			})
		}
	}
	w.WriteMsg(&msg)
}

func ServeDns(port int) {
	// Now I will start the dns server.
	srv := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: "udp"}
	srv.Handler = &handler{}
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println("Failed to set udp listener", err.Error())
	}
}

func (self *server) initStringRecords(recordType string, ttl uint32, record map[string]interface{}) error {
	uuid := Utility.GenerateUUID(recordType + ":" + record["Id"].(string))
	err := self.setTtl(uuid, ttl)
	if err != nil {
		return err
	}
	return self.store.SetItem(uuid, []byte(record["Value"].(string)))
}

func (self *server) initSructRecords(recordType string, ttl uint32, record map[string]interface{}) error {

	data, err := json.Marshal(record["Value"].(map[string]interface{}))
	if err != nil {
		return err
	}
	uuid := Utility.GenerateUUID(recordType + ":" + record["Id"].(string))
	err = self.store.SetItem(uuid, data)
	if err != nil {
		return err
	}

	return self.setTtl(uuid, ttl)
}

func (self *server) initArrayRecords(recordType string, ttl uint32, record map[string]interface{}) error {

	data, err := json.Marshal(record["Value"].([]interface{}))
	if err != nil {
		return err
	}

	uuid := Utility.GenerateUUID(recordType + ":" + record["Id"].(string))

	err = self.store.SetItem(uuid, data)
	if err != nil {
		return err
	}

	return self.setTtl(uuid, ttl)
}

func (self *server) setTtl(uuid string, ttl uint32) error {
	data := make([]byte, 4)
	binary.LittleEndian.PutUint32(data, ttl)
	uuid = Utility.GenerateUUID("TTL:" + uuid)
	err := self.store.SetItem(uuid, data)
	return err
}

func (self *server) getTtl(uuid string) uint32 {
	uuid = Utility.GenerateUUID("TTL:" + uuid)
	data, err := self.store.GetItem(uuid)
	if err != nil {
		return 60 // the default value
	}
	return binary.LittleEndian.Uint32(data)
}

// Initialyse all the records from the configuration.
func (self *server) initRecords() error {
	if self.Records == nil {
		return nil
	}

	for name, records := range self.Records {
		for i := 0; i < len(records); i++ {
			var record = records[i].(map[string]interface{})
			var ttl uint32
			if record["ttl"] != nil {
				ttl = uint32(record["ttl"].(float64))
			} else {
				ttl = 60 // default value of time to live.
			}
			var err error
			if name == "A" {
				err = self.initStringRecords("A", ttl, record)
			} else if name == "AAAA" {
				err = self.initSructRecords("AAAA", ttl, record)
			} else if name == "AFSDB" {
				err = self.initSructRecords("AFSDB", ttl, record)
			} else if name == "CAA" {
				err = self.initSructRecords("CAA", ttl, record)
			} else if name == "CNAME" {
				err = self.initStringRecords("CNAME", ttl, record)
			} else if name == "MX" {
				err = self.initSructRecords("MX", ttl, record)
			} else if name == "SOA" {
				err = self.initSructRecords("SOA", ttl, record)
			} else if name == "TXT" {
				err = self.initArrayRecords("TXT", ttl, record)
			} else if name == "URI" {
				err = self.initSructRecords("URI", ttl, record)
			} else if name == "NS" {
				err = self.initStringRecords("NS", ttl, record)
			} else {
				return errors.New("No ns record with type" + name + "exist!")
			}
			if err != nil {
				fmt.Println("---> ", err)
				return err
			}
		}
	}
	return nil
}

// That service is use to give access to SQL.
// port number must be pass as argument.
func main() {

	// set the logger.
	// grpclog.SetLogger(log.New(os.Stdout, "dns_service: ", log.LstdFlags))

	// Set the log information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// The first argument must be the port number to listen to.
	port := defaultPort // the default value.

	if len(os.Args) > 1 {
		port, _ = strconv.Atoi(os.Args[1]) // The second argument must be the port number
	}

	// The actual server implementation.
	s_impl := new(server)
	Utility.RegisterType(s_impl) // must be call dynamically
	s_impl.Name = string(dnspb.File_dns_dnspb_dns_proto.Services().Get(0).FullName())
	s_impl.Proto = dnspb.File_dns_dnspb_dns_proto.Path()
	s_impl.Path, _ = os.Executable()
	package_ := string(dnspb.File_dns_dnspb_dns_proto.Package().Name())
	s_impl.Path = s_impl.Path[strings.Index(s_impl.Path, package_):]
	s_impl.Port = port
	s_impl.Proxy = defaultProxy
	s_impl.Protocol = "grpc"
	s_impl.Domain = domain
	s_impl.Version = "0.0.1"
	s_impl.PublisherId = domain // value by default.
	s_impl.Permissions = make([]interface{}, 0)
	// TODO set it from the program arguments...
	s_impl.AllowAllOrigins = allow_all_origins
	s_impl.AllowedOrigins = allowed_origins

	// Here I will retreive the list of connections from file if there are some...
	s_impl.init()

	s = s_impl // set the pointer to the server.

	// First of all I will creat a listener.
	// Create the channel to listen on
	lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("could not list on %s: %s", s_impl.Domain, err)
		return
	}

	var grpcServer *grpc.Server
	if s_impl.TLS {
		// Load the certificates from disk
		certificate, err := tls.LoadX509KeyPair(s_impl.CertFile, s_impl.KeyFile)
		if err != nil {
			log.Fatalf("could not load server key pair: %s", err)
			return
		} else {
			fmt.Println("load certificate from ", s_impl.CertFile, s_impl.KeyFile)
		}

		// Create a certificate pool from the certificate authority
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(s_impl.CertAuthorityTrust)
		if err != nil {
			log.Fatalf("could not read ca certificate: %s", err)
			return
		}

		// Append the client certificates from the CA
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			log.Fatalf("failed to append client certs")
			return
		}

		// Create the TLS credentials
		creds := credentials.NewTLS(&tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{certificate},
			ClientCAs:    certPool,
		})

		// Create the gRPC server with the credentials
		opts := []grpc.ServerOption{grpc.Creds(creds), grpc.UnaryInterceptor(Interceptors.ServerUnaryInterceptor), grpc.StreamInterceptor(Interceptors.ServerStreamInterceptor)}
		grpcServer = grpc.NewServer(opts...)

	} else {
		grpcServer = grpc.NewServer([]grpc.ServerOption{grpc.UnaryInterceptor(Interceptors.ServerUnaryInterceptor), grpc.StreamInterceptor(Interceptors.ServerStreamInterceptor)}...)
	}

	dnspb.RegisterDnsServiceServer(grpcServer, s_impl)
	reflection.Register(grpcServer)

	// Here I will make a signal hook to interrupt to exit cleanly.
	go func() {
		fmt.Println(s_impl.Name + " grpc service is starting")
		// no web-rpc server.
		if err := grpcServer.Serve(lis); err != nil {

			if err.Error() == "signal: killed" {
				fmt.Println("service ", s_impl.Name, " was stop!")
			}
		}
		err := s_impl.store.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	// start lisen on the network for dns queries...
	go func() {
		fmt.Println("--> start listen for dns queries at port", s_impl.DnsPort)
		ServeDns(s_impl.DnsPort)
	}()

	// Wait for signal to stop.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

}
