package main

import (
	"context"
	"log"

	"io"

	"github.com/davecourtois/Globular/file/filepb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
)

// The client service interface.
type Client interface {

	// Close the client.
	Close()
}

/**
 * Get the client connection.
 */
func getClientConnection(addresse string) *grpc.ClientConn {
	var err error
	var cc *grpc.ClientConn
	if cc == nil {
		cc, err = grpc.Dial(addresse, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("could not connect: %v", err)
		}
	}
	return cc
}

////////////////////////////////////////////////////////////////////////////////
// File Client Service
////////////////////////////////////////////////////////////////////////////////

type File_Client struct {
	cc *grpc.ClientConn
}

// Create a connection to the service.
func NewFile_Client(addresse string) *File_Client {

	client := new(File_Client)
	client.cc = getClientConnection(addresse)

	return client
}

// must be close when no more needed.
func (self *File_Client) Close() {
	self.cc.Close()
}

// Read the content of a dir and return it info.
func (self *File_Client) ReadDir(path string, recursiveStr string, thumbnailHeightStr string, thumbnailWidtStr string) (string, error) {

	recursive := Utility.ToBool(recursiveStr)
	thumbnailHeight := Utility.ToInt(thumbnailHeightStr)
	thumbnailWidth := Utility.ToInt(thumbnailWidtStr)

	// Create a new client service...
	c := filepb.NewFileServiceClient(self.cc)

	rqst := &filepb.ReadDirRequest{
		Path:           path,
		Recursive:      recursive,
		ThumnailHeight: int32(thumbnailHeight),
		ThumnailWidth:  int32(thumbnailWidth),
	}

	stream, err := c.ReadDir(context.Background(), rqst)
	if err != nil {
		return "", err
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
			return "", err
		}
	}

	return string(data), nil
}
