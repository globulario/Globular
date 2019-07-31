package file_client

import (
	"context"
	"log"

	"io"
	"os"

	"github.com/davecourtois/Globular/api"
	"github.com/davecourtois/Globular/file/filepb"

	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////////////////////////
// File Client Service
////////////////////////////////////////////////////////////////////////////////

type File_Client struct {
	cc *grpc.ClientConn
	c  filepb.FileServiceClient

	// The name of the service
	name string

	// The ipv4 address
	addresse string

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
func NewFile_Client(domain string, addresse string, hasTLS bool, keyFile string, certFile string, caFile string) *File_Client {
	client := new(File_Client)

	client.addresse = addresse
	client.name = "file"
	client.hasTLS = hasTLS
	client.keyFile = keyFile
	client.certFile = certFile
	client.caFile = caFile
	client.domain = domain
	client.cc = api.GetClientConnection(client)
	client.c = filepb.NewFileServiceClient(client.cc)

	return client
}

// Return the ipv4 address
func (self *File_Client) GetAddress() string {
	return self.addresse
}

// Return the domain
func (self *File_Client) GetDomain() string {
	return self.domain
}

// Return the name of the service
func (self *File_Client) GetName() string {
	return self.name
}

// must be close when no more needed.
func (self *File_Client) Close() {
	self.cc.Close()
}

////////////////// TLS ///////////////////

// Get if the client is secure.
func (self *File_Client) HasTLS() bool {
	return self.hasTLS
}

// Get the TLS certificate file path
func (self *File_Client) GetCertFile() string {
	return self.certFile
}

// Get the TLS key file path
func (self *File_Client) GetKeyFile() string {
	return self.keyFile
}

// Get the TLS key file path
func (self *File_Client) GetCaFile() string {
	return self.caFile
}

// Read the content of a dir and return it info.
func (self *File_Client) ReadDir(path interface{}, recursive interface{}, thumbnailHeight interface{}, thumbnailWidth interface{}) (string, error) {

	// Create a new client service...
	rqst := &filepb.ReadDirRequest{
		Path:           Utility.ToString(path),
		Recursive:      Utility.ToBool(recursive),
		ThumnailHeight: int32(Utility.ToInt(thumbnailHeight)),
		ThumnailWidth:  int32(Utility.ToInt(thumbnailWidth)),
	}

	stream, err := self.c.ReadDir(context.Background(), rqst)
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

/**
 * Create a new directory on the server.
 */
func (self *File_Client) CreateDir(path interface{}, name interface{}) error {

	rqst := &filepb.CreateDirRequest{
		Path: Utility.ToString(path),
		Name: Utility.ToString(name),
	}

	_, err := self.c.CreateDir(context.Background(), rqst)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Read file data
 */
func (self *File_Client) ReadFile(path interface{}) ([]byte, error) {

	rqst := &filepb.ReadFileRequest{
		Path: Utility.ToString(path),
	}

	stream, err := self.c.ReadFile(context.Background(), rqst)
	if err != nil {
		return nil, err
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
			return nil, err
		}
	}

	return data, err
}

/**
 * Rename a directory.
 */
func (self *File_Client) RenameDir(path interface{}, oldname interface{}, newname interface{}) error {

	rqst := &filepb.RenameRequest{
		Path:    Utility.ToString(path),
		OldName: Utility.ToString(oldname),
		NewName: Utility.ToString(newname),
	}

	_, err := self.c.Rename(context.Background(), rqst)

	return err
}

/**
 * Delete a directory
 */
func (self *File_Client) DeleteDir(path string) error {
	rqst := &filepb.DeleteDirRequest{
		Path: Utility.ToString(path),
	}

	_, err := self.c.DeleteDir(context.Background(), rqst)
	return err
}

/**
 * Get a single file info.
 */
func (self *File_Client) GetFileInfo(path interface{}, recursive interface{}, thumbnailHeight interface{}, thumbnailWidth interface{}) (string, error) {

	rqst := &filepb.GetFileInfoRequest{
		Path:           Utility.ToString(path),
		ThumnailHeight: int32(Utility.ToInt(thumbnailHeight)),
		ThumnailWidth:  int32(Utility.ToInt(thumbnailWidth)),
	}

	rsp, err := self.c.GetFileInfo(context.Background(), rqst)
	if err != nil {
		return "", err
	}

	return string(rsp.Data), nil
}

/**
 * That function move a file from a directory to another... (mv) in unix.
 */
func (self *File_Client) MoveFile(path interface{}, dest interface{}) error {

	// Open the stream...
	stream, err := self.c.SaveFile(context.Background())
	if err != nil {
		return err
	}

	err = stream.Send(&filepb.SaveFileRequest{
		File: &filepb.SaveFileRequest_Path{
			Path: Utility.ToString(dest), // Where the file will be save...
		},
	})

	if err != nil {
		return err
	}

	// Where the file is read from.
	file, err := os.Open(Utility.ToString(path))
	if err != nil {
		return err
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
				return err
			}
			break
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		return err
	}

	return nil
}

/**
 * Delete a file whit a given path.
 */
func (self *File_Client) DeleteFile(path string) error {

	rqst := &filepb.DeleteFileRequest{
		Path: Utility.ToString(path),
	}

	_, err := self.c.DeleteFile(context.Background(), rqst)
	if err != nil {
		log.Fatalf("error while testing get file info: %v", err)
	}

	return err
}

// Read the content of a dir and return all images as thumbnails.
func (self *File_Client) GetThumbnails(path interface{}, recursive interface{}, thumbnailHeight interface{}, thumbnailWidth interface{}) (string, error) {

	// Create a new client service...
	rqst := &filepb.GetThumbnailsRequest{
		Path:           Utility.ToString(path),
		Recursive:      Utility.ToBool(recursive),
		ThumnailHeight: int32(Utility.ToInt(thumbnailHeight)),
		ThumnailWidth:  int32(Utility.ToInt(thumbnailWidth)),
	}

	stream, err := self.c.GetThumbnails(context.Background(), rqst)
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
