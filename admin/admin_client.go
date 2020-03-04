package admin

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strconv"

	"github.com/davecourtois/Globular/api"

	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

////////////////////////////////////////////////////////////////////////////////
// admin Client Service
////////////////////////////////////////////////////////////////////////////////

type Admin_Client struct {
	cc *grpc.ClientConn
	c  AdminServiceClient

	// The name of the service
	name string

	// The port
	port int

	// The client domain
	domain string

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
func NewAdmin_Client(address string, name string) *Admin_Client {
	client := new(Admin_Client)
	api.InitClient(client, address, name)
	client.cc = api.GetClientConnection(client)
	client.c = NewAdminServiceClient(client.cc)

	return client
}

// Return the address
func (self *Admin_Client) GetAddress() string {
	return self.domain + ":" + strconv.Itoa(self.port)
}

// Return the domain
func (self *Admin_Client) GetDomain() string {
	return self.domain
}

// Return the name of the service
func (self *Admin_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *Admin_Client) Close() {
	self.cc.Close()
}

// Set grpc_service port.
func (self *Admin_Client) SetPort(port int) {
	self.port = port
}

// Set the client name.
func (self *Admin_Client) SetName(name string) {
	self.name = name
}

// Set the domain.
func (self *Admin_Client) SetDomain(domain string) {
	self.domain = domain
}

////////////////// TLS ///////////////////

// Get if the client is secure.
func (self *Admin_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *Admin_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *Admin_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *Admin_Client) GetCaFile() string {
	return self.caFile
}

// Set the client is a secure client.
func (self *Admin_Client) SetTLS(hasTls bool) {
	self.hasTLS = hasTls
}

// Set TLS certificate file path
func (self *Admin_Client) SetCertFile(certFile string) {
	self.certFile = certFile
}

// Set TLS key file path
func (self *Admin_Client) SetKeyFile(keyFile string) {
	self.keyFile = keyFile
}

// Set TLS authority trust certificate file path
func (self *Admin_Client) SetCaFile(caFile string) {
	self.caFile = caFile
}

/////////////////////// API /////////////////////

// Get server configuration.
func (self *Admin_Client) GetConfig() (string, error) {
	rqst := new(GetConfigRequest)

	rsp, err := self.c.GetConfig(api.GetClientContext(self), rqst)
	if err != nil {
		return "", err
	}

	return rsp.GetResult(), nil
}

// Get the server configuration with all detail must be secured.
func (self *Admin_Client) GetFullConfig() (string, error) {
	rqst := new(GetConfigRequest)

	rsp, err := self.c.GetFullConfig(api.GetClientContext(self), rqst)
	if err != nil {
		return "", err
	}

	return rsp.GetResult(), nil
}

func (self *Admin_Client) SaveConfig(config string) error {
	rqst := &SaveConfigRequest{
		Config: config,
	}

	_, err := self.c.SaveConfig(api.GetClientContext(self), rqst)
	if err != nil {
		return err
	}

	return nil
}

func (self *Admin_Client) StartService(id string) (int, int, error) {
	rqst := new(StartServiceRequest)
	rqst.ServiceId = id
	rsp, err := self.c.StartService(api.GetClientContext(self), rqst)
	if err != nil {
		return -1, -1, err
	}

	return int(rsp.ServicePid), int(rsp.ProxyPid), nil
}

func (self *Admin_Client) StopService(id string) error {
	rqst := new(StopServiceRequest)
	rqst.ServiceId = id
	_, err := self.c.StopService(api.GetClientContext(self), rqst)
	if err != nil {
		return err
	}

	return nil
}

// Register and start an application.
func (self *Admin_Client) RegisterExternalApplication(id string, path string, args []string) (int, error) {
	rqst := &RegisterExternalApplicationRequest{
		ServiceId: id,
		Path:      path,
		Args:      args,
	}

	rsp, err := self.c.RegisterExternalApplication(api.GetClientContext(self), rqst)

	if err != nil {
		return -1, err
	}

	return int(rsp.ServicePid), nil
}

/////////////////////////// Services management functions ////////////////////////

/** Create a service package **/
func (self *Admin_Client) createServicePackage(publisherId string, serviceId string, version string, platform Platform, servicePath string) (string, error) {

	// Take the information from the configuration...
	id := publisherId + "%" + serviceId + "%" + version
	if platform == Platform_LINUX32 {
		id += "%LINUX32"
	} else if platform == Platform_LINUX64 {
		id += "%LINUX64"
	} else if platform == Platform_WIN32 {
		id += "%WIN32"
	} else if platform == Platform_WIN64 {
		id += "%WIN64"
	}

	// So here I will create a directory and put file in it...
	path := id
	Utility.CreateDirIfNotExist(path)

	// copy all the data.
	Utility.CopyDirContent(servicePath, path)

	// tar + gzip
	var buf bytes.Buffer
	Utility.CompressDir("", path, &buf)

	// write the .tar.gzip
	fileToWrite, err := os.OpenFile(os.TempDir()+string(os.PathSeparator)+id+".tar.gz", os.O_CREATE|os.O_RDWR, os.FileMode(0755))
	if err != nil {
		panic(err)
	}

	if _, err := io.Copy(fileToWrite, &buf); err != nil {
		panic(err)
	}

	// close the file.
	fileToWrite.Close()

	// Remove the dir when the archive is created.
	err = os.RemoveAll(path)

	if err != nil {
		return "", err
	}

	return os.TempDir() + string(os.PathSeparator) + id + ".tar.gz", nil
}

/**
 * Create and Upload the service archive on the server.
 */
func (self *Admin_Client) UploadServicePackage(path string, publisherId string, serviceId string, version string, token string) (string, error) {
	// Here I will try to read the service configuation from the path.
	configs, _ := Utility.FindFileByName(path, "config.json")
	if len(configs) == 0 {
		return "", errors.New("No configuration file was found")
	}

	_, err := Utility.FindFileByName(path, ".proto")
	if len(configs) == 0 {
		return "", errors.New("No prototype file was found")
	}

	s := make(map[string]interface{}, 0)
	data, err := ioutil.ReadFile(configs[0])
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(data, &s)
	if err != nil {
		return "", err
	}

	// set the correct information inside the configuration
	s["PublisherId"] = publisherId
	s["Version"] = version
	s["Name"] = serviceId

	jsonStr, _ := Utility.ToJson(&s)
	ioutil.WriteFile(configs[0], []byte(jsonStr), 0644)

	var platform Platform

	// The first step will be to create the archive.
	if runtime.GOOS == "windows" {
		if runtime.GOARCH == "amd64" {
			platform = Platform_WIN64
		} else if runtime.GOARCH == "386" {
			platform = Platform_WIN32
		}
	} else if runtime.GOOS == "linux" { // also can be specified to FreeBSD
		if runtime.GOARCH == "amd64" {
			platform = Platform_LINUX64
		} else if runtime.GOARCH == "386" {
			platform = Platform_LINUX32
		}
	} else if runtime.GOOS == "darwin" {
		/** TODO Deploy services on other platforme here... **/
	}

	md := metadata.New(map[string]string{"token": token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// First of all I will create the archive for the service.
	// If a path is given I will take it entire content. If not
	// the proto, the config and the executable only will be taken.
	packagePath, err := self.createServicePackage(publisherId, serviceId, version, platform, path)
	if err != nil {
		return "", err
	}

	// Remove the file when it's transfer on the server...
	defer os.Remove(packagePath)

	// Read the package data.
	packageFile, err := os.Open(packagePath)
	if err != nil {
		log.Fatal(err)
	}
	defer packageFile.Close()

	// Now I will create the request to upload the package on the server.
	// Open the stream...
	stream, err := self.c.UploadServicePackage(ctx)
	if err != nil {
		log.Fatalf("error while TestSendEmailWithAttachements: %v", err)
	}

	const chunksize = 1024 * 5 // the chunck size.
	var count int
	reader := bufio.NewReader(packageFile)
	part := make([]byte, chunksize)

	for {
		if count, err = reader.Read(part); err != nil {
			break
		}
		rqst := &UploadServicePackageRequest{
			Data: part[:count],
		}
		// send the data to the server.
		err = stream.Send(rqst)
	}
	if err != io.EOF {
		return "", err
	} else {
		err = nil
	}

	// get the file path on the server where the package is store before being
	// publish.
	rsp, err := stream.CloseAndRecv()
	if err != nil {
		return "", err
	}

	return rsp.Path, nil
}

/**
 * Publish a service from a runing globular server.
 */
func (self *Admin_Client) PublishService(path string, serviceId string, publisherId string, discoveryAddress string, repositoryAddress string, description string, version string, platform int32, keywords []string, token string) error {

	rqst := new(PublishServiceRequest)
	rqst.Path = path
	rqst.PublisherId = publisherId
	rqst.Description = description
	rqst.DicorveryId = discoveryAddress
	rqst.RepositoryId = repositoryAddress
	rqst.Keywords = keywords
	rqst.Version = version
	rqst.ServiceId = serviceId
	rqst.Platform = Platform(platform)

	// Set the token into the context and send the request.
	md := metadata.New(map[string]string{"token": token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := self.c.PublishService(ctx, rqst)

	return err
}

/**
 * Intall a new service or update an existing one.
 */
func (self *Admin_Client) InstallService(discoveryId string, publisherId string, serviceId string) error {

	rqst := new(InstallServiceRequest)
	rqst.DicorveryId = discoveryId
	rqst.PublisherId = publisherId
	rqst.ServiceId = serviceId

	_, err := self.c.InstallService(api.GetClientContext(self), rqst)

	return err
}

/**
 * Intall a new service or update an existing one.
 */
func (self *Admin_Client) UninstallService(publisherId string, serviceId string, version string) error {

	rqst := new(UninstallServiceRequest)
	rqst.PublisherId = publisherId
	rqst.ServiceId = serviceId
	rqst.Version = version

	_, err := self.c.UninstallService(api.GetClientContext(self), rqst)

	return err
}

/**
 * Deploy the content of an application with a given name to the server.
 */
func (self *Admin_Client) DeployApplication(name string, path string, token string) error {

	rqst := new(DeployApplicationRequest)
	rqst.Name = name

	Utility.CreateDirIfNotExist(Utility.GenerateUUID(name))
	Utility.CopyDirContent(path, Utility.GenerateUUID(name))

	// Now I will open the data and create a archive from it.
	var buffer bytes.Buffer
	err := Utility.CompressDir("", Utility.GenerateUUID(name), &buffer)
	if err != nil {
		return err
	}
	// remove the dir and keep the archive in memory
	os.RemoveAll(Utility.GenerateUUID(name))

	// Set the token into the context and send the request.
	md := metadata.New(map[string]string{"token": string(token), "application": name})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Open the stream...
	stream, err := self.c.DeployApplication(ctx)
	if err != nil {
		log.Panicln(err)
		return err
	}

	const BufferSize = 1024 * 5 // the chunck size.
	for {
		var data [BufferSize]byte
		bytesread, err := buffer.Read(data[0:BufferSize])
		if bytesread > 0 {
			rqst := &DeployApplicationRequest{
				Data: data[0:bytesread],
				Name: name,
			}
			// send the data to the server.
			err = stream.Send(rqst)
		}

		if err == io.EOF {
			err = nil
			break
		} else if err != nil {
			return err
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		os.RemoveAll(Utility.GenerateUUID(name))
		return err
	}

	return nil

}
