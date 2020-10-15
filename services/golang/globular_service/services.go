package globular_service

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"

	//	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"

	"log"
	"os"

	"fmt"
	"os/signal"
	"time"

	"github.com/davecourtois/Globular/services/golang/admin/admin_client"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	//"google.golang.org/grpc/grpclog"

	"errors"
)

// The client service interface.
type Service interface {

	/** Getter/Setter **/

	// The id of a particular service instance.
	GetId() string
	SetId(string)

	// The name of a service, must be the gRpc Service name.
	GetName() string
	SetName(string)

	// The path of the executable.
	GetPath() string
	SetPath(string)

	// The path of the .proto file.
	GetProto() string
	SetProto(string)

	// The gRpc port.
	GetPort() int
	SetPort(int)

	// The reverse proxy port (use by gRpc Web)
	GetProxy() int
	SetProxy(int)

	// Can be one of http/https/tls
	GetProtocol() string
	SetProtocol(string)

	// Return true if all Origins are allowed to access the mircoservice.
	GetAllowAllOrigins() bool
	SetAllowAllOrigins(bool)

	// If AllowAllOrigins is false then AllowedOrigins will contain the
	// list of address that can reach the services.
	GetAllowedOrigins() string // comma separated string.
	SetAllowedOrigins(string)

	// Can be a ip address or domain name.
	GetDomain() string
	SetDomain(string)

	// TLS section

	// If true the service run with TLS. The
	GetTls() bool
	SetTls(bool)

	// The certificate authority file
	GetCertAuthorityTrust() string
	SetCertAuthorityTrust(string)

	// The certificate file.
	GetCertFile() string
	SetCertFile(string)

	// The key file.
	GetKeyFile() string
	SetKeyFile(string)

	// The service version
	GetVersion() string
	SetVersion(string)

	// The publisher id.
	GetPublisherId() string
	SetPublisherId(string)

	GetKeepUpToDate() bool
	SetKeepUptoDate(bool)

	GetKeepAlive() bool
	SetKeepAlive(bool)

	GetPermissions() []interface{} // contains the action permission for the services.
	SetPermissions([]interface{})

	/** Initialyse the service configuration **/
	Init() error

	/** Save the service configuration **/
	Save() error

	/** Stop the service **/
	StopService() error

	/** Start the service **/
	StartService() error
}

/**
 * Initialise a globular service from it configuration file.
 */
func InitService(path string, s Service) error {

	// Here I will retreive the list of connections from file if there are some...
	file, err := ioutil.ReadFile(path)
	if err == nil {
		return json.Unmarshal([]byte(file), s)
	} else {
		// Generate an id if none exist in the given configuration.
		if len(s.GetId()) == 0 {
			// Generate random id for the server instance.
			s.SetId(Utility.RandomUUID())
		}

		execPath, _ := os.Executable()
		execPath = strings.ReplaceAll(execPath, "\\", "/")
		s.SetPath(execPath)

		package_ := strings.Split(s.GetProto(), "/")[0]
		path_ := execPath[0:strings.Index(execPath, package_)]
		s.SetProto(path_ + s.GetProto())

		// save the service configuation.
		return SaveService(path, s)
	}
}

func GetTLSConfig(key string, cert string, ca string) *tls.Config {
	tlsCer, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		log.Fatalf("Failed to generate credentials: %v", err)
	}

	certPool := x509.NewCertPool()
	clientCA, err := ioutil.ReadFile(ca)
	if err != nil {
		log.Fatalf("failed to read client ca cert: %s", err)
	}
	ok := certPool.AppendCertsFromPEM(clientCA)
	if !ok {
		log.Fatal("failed to append client certs")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{tlsCer},
		ClientAuth:   tls.RequireAnyClientCert,
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			opts := x509.VerifyOptions{
				Roots:         certPool,
				CurrentTime:   time.Now(),
				Intermediates: x509.NewCertPool(),
				KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
			}

			for _, cert := range rawCerts[1:] {
				opts.Intermediates.AppendCertsFromPEM(cert)
			}

			c, err := x509.ParseCertificate(rawCerts[0])
			if err != nil {
				return errors.New("tls: failed to verify client certificate: " + err.Error())
			}
			_, err = c.Verify(opts)
			if err != nil {
				return errors.New("tls: failed to verify client certificate: " + err.Error())
			}
			return nil
		},
	}
}

/**
 * Initilalsyse the grpc server that will run the service.
 */
func InitGrpcServer(s Service, unaryInterceptor grpc.UnaryServerInterceptor, streamInterceptor grpc.StreamServerInterceptor) (*grpc.Server, error) {
	var server *grpc.Server

	if s.GetTls() {

		// Create the TLS credentials
		creds := credentials.NewTLS(GetTLSConfig(s.GetKeyFile(), s.GetCertFile(), s.GetCertAuthorityTrust()))

		// Create the gRPC server with the credentials

		opts := []grpc.ServerOption{grpc.Creds(creds)}
		if unaryInterceptor != nil {
			opts = append(opts, grpc.UnaryInterceptor(unaryInterceptor))
		}
		if streamInterceptor != nil {
			opts = append(opts, grpc.StreamInterceptor(streamInterceptor))
		}

		server = grpc.NewServer(opts...)

	} else {
		if unaryInterceptor != nil && streamInterceptor != nil {
			server = grpc.NewServer(
				grpc.UnaryInterceptor(unaryInterceptor),
				grpc.StreamInterceptor(streamInterceptor))
		} else {
			server = grpc.NewServer()
		}
	}

	return server, nil
}

var (
	admin_client_ *admin_client.Admin_Client
)

/**
 * Get a the local ressource client.
 */
func getAdminClient(domain string) (*admin_client.Admin_Client, error) {
	var err error
	if admin_client_ == nil {
		admin_client_, err = admin_client.NewAdminService_Client(domain, "admin.AdminService")
		if err != nil {
			return nil, err
		}
	}

	return admin_client_, nil
}

/**
 * Save a globular service.
 */
func UpdateServiceConfig(s Service) error {

	// Create the file...
	str, err := Utility.ToJson(s)
	if err != nil {
		return err
	}
	admin_client_, err := getAdminClient(s.GetDomain())
	if err != nil {
		return err
	}

	return admin_client_.SaveConfig(str)
}

/**
 * Save a globular service.
 */
func SaveService(path string, s Service) error {

	// Create the file...
	str, err := Utility.ToJson(s)
	if err != nil {
		return err
	}

	ioutil.WriteFile(path, []byte(str), 0644)
	return nil
}

func StartService(s Service, server *grpc.Server) error {
	// First of all I will creat a listener.
	// Create the channel to listen on
	lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(s.GetPort()))
	if err != nil {
		return errors.New("could not list at domain " + s.GetDomain() + err.Error())
	}

	// Here I will make a signal hook to interrupt to exit cleanly.
	go func() {
		// no web-rpc server.
		fmt.Println(s.GetName() + " grpc service is starting")
		if err := server.Serve(lis); err != nil {
			if err.Error() == "signal: killed" {
				fmt.Println("service ", s.GetId(), s.GetName(), " was stop!")
			}
		}
	}()

	// Wait for signal to stop.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	fmt.Println(s.GetId() + " is now stopped!")
	return nil
}

func StopService(s Service, server *grpc.Server) error {
	// Stop the service.
	server.GracefulStop()
	return nil
}
