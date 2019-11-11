package services

//"encoding/json"
import (
	"io/ioutil"
	"log"
	"testing"
)

var (
	// Connect to the services client.
	services_discovery  = NewServicesDiscovery_Client("localhost", 10005, false, "", "", "", "")
	services_repository = NewServicesRepository_Client("localhost", 10007, false, "", "", "", "")
)

// Test publish a service.
func TestPublishServiceDescriptor(t *testing.T) {
	s := &ServiceDescriptor{
		Id:          "echo_server",
		PublisherId: "globular.app",
		Version:     "1.0.4",
		Description: "Simple service with one function named Echo. It's mostly a test service.",
		Keywords:    []string{"Test", "Echo"},
	}
	err := services_discovery.PublishService(s)
	if err != nil {
		log.Panic(err)
	}

	log.Print("Service was publish with success!!!")
}

func TestGetServiceDescriptor(t *testing.T) {

	values, err := services_discovery.GetServiceDescriptor("echo_server", "localhost")

	if err != nil {
		log.Panic(err)
	}

	log.Print("Service was retreived with success!!!", values)
}

func TestUploadServiceBundle(t *testing.T) {

	// The service bundle...
	err := services_repository.UploadBundle("localhost:10005", "echo_server", "localhost", 0, "C:\\temp\\globular\\echo_server.7z")
	if err != nil {
		log.Panicln(err)
	}
}

func TestDownloadServiceBundle(t *testing.T) {
	bundle, err := services_repository.DownloadLastVersionBundle("localhost:10005", "echo_server", "localhost", 0)

	if err != nil {
		log.Panicln(err)
	}

	ioutil.WriteFile("C:\\temp\\echo_server.7z", bundle.Binairies, 777)
}
