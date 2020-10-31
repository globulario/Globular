package file_client

import (
	//"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/globulario/Globular/file/filepb"
	"google.golang.org/grpc"

	"testing"
)

// Set the correct addresse here as needed.
var (
	addresse = "localhost:10011"
)

/**
 * Get the client connection.
 */
func getClientConnection() *grpc.ClientConn {
	// So here I will read the server configuration to see if the connection
	// is secure...
	cc, err := grpc.Dial(addresse, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}

	return cc
}

// First test create a fresh new connection...
func _TestReadDir(t *testing.T) {
	fmt.Println("Read dir test")

	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := filepb.NewFileServiceClient(cc)

	rqst := &filepb.ReadDirRequest{
		Path:           "C:\\Temp\\Cargo\\WebApp\\Cargo\\Apps\\BrisOutil",
		Recursive:      true,
		ThumnailHeight: 256,
		ThumnailWidth:  256,
	}

	stream, err := c.ReadDir(api.GetClientContext(self), rqst)
	if err != nil {
		log.Fatalf("Query error %v", err)
	}

	// Here I will create the final array
	data := make([]byte, 0)
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// end of stream...
			log.Println("----> ", string(data))
			break
		}

		data = append(data, msg.Data...)
		if err != nil {
			log.Fatalf("error while TestReadDir: %v", err)
		}
	}

	log.Println("TestReadDir successed!")
}

func TestGetThumbnails(t *testing.T) {
	fmt.Println("Get Thumbnails")

	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := filepb.NewFileServiceClient(cc)

	rqst := &filepb.GetThumbnailsRequest{
		Path:           "/test/filePane", //"C:\\Temp\\Cargo\\WebApp\\Cargo\\Apps\\BrisOutil",
		Recursive:      true,
		ThumnailHeight: 256,
		ThumnailWidth:  256,
	}

	stream, err := c.GetThumbnails(api.GetClientContext(self), rqst)
	if err != nil {
		log.Fatalf("Query error %v", err)
	}

	// Here I will create the final array
	data := make([]byte, 0)
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// end of stream...
			break
		}

		data = append(data, msg.Data...)
		if err != nil {
			log.Fatalf("error while TestReadDir: %v", err)
		}
	}

	log.Println("TestReadDir successed!")
}

/**
 * Create a new directory.
 */
func TestCreateDir(t *testing.T) {
	fmt.Println("Create dir test")
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := filepb.NewFileServiceClient(cc)

	rqst := &filepb.CreateDirRequest{
		Path: "C:\\Temp",
		Name: "TestDir",
	}

	rsp, err := c.CreateDir(api.GetClientContext(self), rqst)
	if err != nil {
		log.Fatalf("error while TestCreateDir: %v", err)
	}

	log.Println("Response from TestCreateDir:", rsp.Result)
}

/**
 * Rename a directory
 */
func TestRenameDir(t *testing.T) {
	fmt.Println("Rename dir test")
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := filepb.NewFileServiceClient(cc)

	rqst := &filepb.RenameRequest{
		Path:    "C:\\Temp",
		OldName: "TestDir",
		NewName: "TestTestTestDir",
	}

	rsp, err := c.Rename(api.GetClientContext(self), rqst)
	if err != nil {
		log.Fatalf("error while TestRenameDir: %v", err)
	}

	log.Println("Response from TestRenameDir:", rsp.Result)
}

/**
 * Rename a directory
 */
func TestDeleteDir(t *testing.T) {
	fmt.Println("Delete dir test")
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := filepb.NewFileServiceClient(cc)

	rqst := &filepb.DeleteDirRequest{
		Path: "C:\\Temp\\TestTestTestDir",
	}

	rsp, err := c.DeleteDir(api.GetClientContext(self), rqst)
	if err != nil {
		log.Fatalf("error while TestDeleteDir: %v", err)
	}

	log.Println("Response from TestDeleteDir:", rsp.Result)
}

////////////////////////////////////////////////////////////////////////////////
// File test
////////////////////////////////////////////////////////////////////////////////
func TestGetFileInof(t *testing.T) {
	fmt.Println("Get File info test")
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := filepb.NewFileServiceClient(cc)

	rqst := &filepb.GetFileInfoRequest{
		Path:           "C:\\Temp\\Cargo\\WebApp\\Cargo\\Apps\\BrisOutil\\Upload\\515\\NGEN3603.JPG",
		ThumnailHeight: 256,
		ThumnailWidth:  256,
	}

	rsp, err := c.GetFileInfo(api.GetClientContext(self), rqst)
	if err != nil {
		log.Fatalf("error while testing get file info: %v", err)
	}

	log.Println("Response form Get file info response :", string(rsp.Data))
}

// Read file test.
func TestReadFile(t *testing.T) {
	fmt.Println("Read file test")

	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := filepb.NewFileServiceClient(cc)

	rqst := &filepb.ReadFileRequest{
		Path: "C:\\Temp\\Cargo\\WebApp\\Cargo\\Apps\\BrisOutil\\Upload\\515\\NGEN3603.JPG",
	}

	stream, err := c.ReadFile(api.GetClientContext(self), rqst)
	if err != nil {
		log.Fatalf("ReadFile error %v", err)
	}

	// Here I will create the final array
	data := make([]byte, 0)
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// end of stream...
			log.Println("Read file ----> ", len(data))
			break
		}

		data = append(data, msg.Data...)
		if err != nil {
			log.Fatalf("error while read File: %v", err)
		}
	}

}

/**
 * Test send email whit attachements.
 */
func TestSaveFile(t *testing.T) {
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := filepb.NewFileServiceClient(cc)

	// Open the stream...
	stream, err := c.SaveFile(api.GetClientContext(self))
	if err != nil {
		log.Fatalf("error while TestSendEmailWithAttachements: %v", err)
	}

	err = stream.Send(&filepb.SaveFileRequest{
		File: &filepb.SaveFileRequest_Path{
			Path: "C:\\Temp\\toto.bmp", // Where the file will be save...
		},
	})
	if err != nil {
		log.Println("--> error: ")
	}

	// Where the file is read from.
	path := "C:\\Temp\\Cargo\\WebApp\\Cargo\\Apps\\BrisOutil\\Upload\\515\\NGEN3603.JPG"
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Fail to open file "+path+" with error: %v", err)
	}

	// close the file when done.
	defer file.Close()

	const BufferSize = 1024 * 5 // the chunck size.
	buffer := make([]byte, BufferSize)
	for {
		bytesread, err := file.Read(buffer)
		if bytesread > 0 {
			rqst := &filepb.SaveFileRequest{
				File: &filepb.SaveFileRequest_Data{
					Data: buffer[:bytesread],
				},
			}
			err = stream.Send(rqst)
		}

		if err != nil {
			if err != io.EOF {
				fmt.Println(err)
			}
			break
		}
	}
	rsp, err := stream.CloseAndRecv()
	if err != nil {
		log.Println(err)
	}

	log.Println("save file succeed ", rsp.Result)
}

// Test delete file on the server
func TestDeleteFile(t *testing.T) {
	fmt.Println("Get File info test")
	cc := getClientConnection()

	// when done the connection will be close.
	defer cc.Close()

	// Create a new client service...
	c := filepb.NewFileServiceClient(cc)

	rqst := &filepb.DeleteFileRequest{
		Path: "C:\\Temp\\toto.bmp",
	}

	rsp, err := c.DeleteFile(api.GetClientContext(self), rqst)
	if err != nil {
		log.Fatalf("error while testing get file info: %v", err)
	}

	log.Println("Delete file succeed:", rsp.Result)
}
