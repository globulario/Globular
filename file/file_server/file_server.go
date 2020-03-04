package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/davecourtois/Globular/Interceptors/server"
	"github.com/davecourtois/Globular/file/filepb"
	"github.com/davecourtois/Utility"
	"github.com/nfnt/resize"
	"github.com/polds/imgbase64"
	"github.com/tealeg/xlsx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// TODO take care of TLS/https
var (
	defaultPort  = 10011
	defaultProxy = 10012

	// By default all origins are allowed.
	allow_all_origins = true

	// comma separeated values.
	allowed_origins string = ""

	// The default domain.
	domain string = "localhost"

	s *server
)

// Value need by Globular to start the services...
type server struct {
	// The global attribute of the services.
	Name               string
	Port               int
	Proxy              int
	AllowAllOrigins    bool
	AllowedOrigins     string // comma separated string.
	Protocol           string
	Domain             string
	CertFile           string
	CertAuthorityTrust string
	KeyFile            string
	TLS                bool
	Version            string
	PublisherId        string
	KeepUpToDate       bool
	KeepAlive          bool

	// Specific to file server.
	Root string
}

// Create the configuration file if is not already exist.
func (self *server) init() {
	// Here I will retreive the list of connections from file if there are some...
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	file, err := ioutil.ReadFile(dir + "/config.json")
	if err == nil {
		json.Unmarshal([]byte(file), self)
	} else {
		self.save()
	}
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

/**
 * Create a thumbnail...
 */
func createThumbnail(file *os.File, thumbnailMaxHeight int, thumbnailMaxWidth int) string {
	// Set the buffer pointer back to the begening of the file...
	file.Seek(0, 0)
	var originalImg image.Image
	var format string
	var err error

	if strings.HasSuffix(file.Name(), ".png") || strings.HasSuffix(file.Name(), ".PNG") {
		originalImg, err = png.Decode(file)
	} else if strings.HasSuffix(file.Name(), ".jpeg") || strings.HasSuffix(file.Name(), ".jpg") || strings.HasSuffix(file.Name(), ".JPEG") || strings.HasSuffix(file.Name(), ".JPG") {
		originalImg, err = jpeg.Decode(file)
	} else if strings.HasSuffix(file.Name(), ".gif") || strings.HasSuffix(file.Name(), ".GIF") {
		originalImg, err = gif.Decode(file)
	} else {
		return ""
	}

	if err != nil {
		log.Println("File ", file.Name(), " Format is ", format, " error: ", err)
		return ""
	}

	// I will get the ratio for the new image size to respect the scale.
	hRatio := thumbnailMaxHeight / originalImg.Bounds().Size().Y
	wRatio := thumbnailMaxWidth / originalImg.Bounds().Size().X

	var h int
	var w int

	// First I will try with the height
	if hRatio*originalImg.Bounds().Size().Y < thumbnailMaxWidth {
		h = thumbnailMaxHeight
		w = hRatio * originalImg.Bounds().Size().Y
	} else {
		// So here i will use it width
		h = wRatio * thumbnailMaxHeight
		w = thumbnailMaxWidth
	}

	// do not zoom...
	if hRatio > 1 {
		h = originalImg.Bounds().Size().Y
	}

	if wRatio > 1 {
		w = originalImg.Bounds().Size().X
	}

	// Now I will calculate the image size...
	img := resize.Resize(uint(h), uint(w), originalImg, resize.Lanczos3)
	var buf bytes.Buffer
	jpeg.Encode(&buf, img, &jpeg.Options{jpeg.DefaultQuality})

	// Now I will save the buffer containt to the thumbnail...
	thumbnail := imgbase64.FromBuffer(buf)
	file.Seek(0, 0) // Set the reader back to the begenin of the file...
	return thumbnail
}

type fileInfo struct {
	Name    string      // base name of the file
	Size    int64       // length in bytes for regular files; system-dependent for others
	Mode    os.FileMode // file mode bits
	ModTime time.Time   // modification time
	IsDir   bool        // abbreviation for Mode().IsDir()
	Path    string      // The path on the server.

	Mime      string
	Thumbnail string
	Files     []*fileInfo
}

func getFileInfo(path string) (*fileInfo, error) {
	fileStat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	info := new(fileInfo)

	info.IsDir = fileStat.IsDir()
	info.Size = fileStat.Size()
	info.Name = fileStat.Name()
	info.ModTime = fileStat.ModTime()
	info.Path = path

	// Cut the Root part of the part.
	if len(s.Root) > 0 {
		startIndex := strings.Index(info.Path, s.Root)
		if startIndex == 0 {
			info.Path = info.Path[len(s.Root):]
			info.Path = strings.Replace(info.Path, "\\", "/", -1) // Set the slash instead of back slash.
		}
	}

	return info, nil
}

func getThumbnails(info *fileInfo) []interface{} {
	// The array of thumbnail
	thumbnails := make([]interface{}, 0)

	// Now from the info i will extract the thumbnail
	for i := 0; i < len(info.Files); i++ {
		if !info.Files[i].IsDir {
			thumbnail := make(map[string]string)
			thumbnail["path"] = info.Files[i].Path
			thumbnail["thumbnail"] = info.Files[i].Thumbnail
			thumbnails = append(thumbnails, thumbnail)
		} else {
			thumbnails = append(thumbnails, getThumbnails(info.Files[i])...)
		}
	}

	return thumbnails
}

/**
 * Read the directory and return the file info.
 */
func readDir(path string, recursive bool, thumbnailMaxWidth int32, thumbnailMaxHeight int32) (*fileInfo, error) {

	// get the file info
	info, err := getFileInfo(path)
	if err != nil {
		log.Println("fail to read file ", path)
		return nil, err
	}
	if info.IsDir == false {
		return nil, errors.New(path + " is a directory")
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println("fail to read directory ", path)
		return nil, err
	}

	for _, f := range files {

		if f.IsDir() {
			if recursive {
				info_, err := readDir(path+string(os.PathSeparator)+f.Name(), recursive, thumbnailMaxWidth, thumbnailMaxHeight)
				if err != nil {
					return nil, err
				}
				info.Files = append(info.Files, info_)
			}
		} else {
			info_, err := getFileInfo(path + string(os.PathSeparator) + f.Name())

			f_, err := os.Open(path + string(os.PathSeparator) + f.Name())
			defer f_.Close()

			if err != nil {
				return nil, err
			}

			info_.Mime, err = Utility.GetFileContentType(f_)

			// in case of image...
			if strings.HasPrefix(info_.Mime, "image/") {
				if thumbnailMaxHeight > 0 && thumbnailMaxWidth > 0 {
					info_.Thumbnail = createThumbnail(f_, int(thumbnailMaxHeight), int(thumbnailMaxWidth))
				}
			}

			if err != nil {
				return nil, err
			}
			info.Files = append(info.Files, info_)
		}

	}

	return info, err
}

////////////////////////////////////////////////////////////////////////////////
// Directory operations
////////////////////////////////////////////////////////////////////////////////
func (self *server) ReadDir(rqst *filepb.ReadDirRequest, stream filepb.FileService_ReadDirServer) error {
	path := rqst.GetPath()

	// set the correct os path separator.
	path = strings.ReplaceAll(strings.ReplaceAll(path, "\\", string(os.PathSeparator)), "/", string(os.PathSeparator))

	if strings.HasPrefix(path, string(os.PathSeparator)) {
		if len(path) > 1 {
			if strings.HasPrefix(path, string(os.PathSeparator)) {
				path = self.Root + path
			} else {
				path = self.Root + string(os.PathSeparator) + path
			}
		} else {
			path = self.Root
		}
	}

	log.Println("--> read dir: ", path)
	info, err := readDir(path, rqst.GetRecursive(), rqst.GetThumnailWidth(), rqst.GetThumnailHeight())
	if err != nil {
		return err
	}

	// Here I will serialyse the data into JSON.
	jsonStr, err := json.Marshal(info)
	if err != nil {
		return err
	}

	maxSize := 1024 * 5
	size := int(math.Ceil(float64(len(jsonStr)) / float64(maxSize)))
	for i := 0; i < size; i++ {
		start := i * maxSize
		end := start + maxSize
		var data []byte
		if end > len(jsonStr) {
			data = jsonStr[start:]
		} else {
			data = jsonStr[start:end]
		}
		stream.Send(&filepb.ReadDirResponse{
			Data: data,
		})
	}

	return nil
}

// Create a new directory
func (self *server) CreateDir(ctx context.Context, rqst *filepb.CreateDirRequest) (*filepb.CreateDirResponse, error) {
	path := rqst.GetPath()

	// set the correct os path separator.
	path = strings.ReplaceAll(strings.ReplaceAll(path, "\\", string(os.PathSeparator)), "/", string(os.PathSeparator))

	if strings.HasPrefix(path, string(os.PathSeparator)) {
		if len(path) > 1 {
			if strings.HasPrefix(path, string(os.PathSeparator)) {
				path = self.Root + path
			} else {
				path = self.Root + string(os.PathSeparator) + path
			}
		} else {
			path = self.Root
		}
	}

	err := Utility.CreateDirIfNotExist(path + string(os.PathSeparator) + rqst.GetName())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// The directory was successfuly created.
	return &filepb.CreateDirResponse{
		Result: true,
	}, nil
}

// Create an archive from a given dir and set it with name.
func (self *server) CreateAchive(ctx context.Context, rqst *filepb.CreateArchiveRequest) (*filepb.CreateArchiveResponse, error) {
	path := rqst.GetPath()

	// set the correct os path separator.
	path = strings.ReplaceAll(strings.ReplaceAll(path, "\\", string(os.PathSeparator)), "/", string(os.PathSeparator))

	if strings.HasPrefix(path, string(os.PathSeparator)) {
		if len(path) > 1 {
			if strings.HasPrefix(path, string(os.PathSeparator)) {
				path = self.Root + path
			} else {
				path = self.Root + string(os.PathSeparator) + path
			}
		} else {
			path = self.Root
		}
	}

	if !Utility.Exists(path) {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No file found with path '"+path+"'")))
	}

	var buf bytes.Buffer
	Utility.CompressDir(self.Root, path, &buf)

	dest := path[0:strings.LastIndex(path, string(os.PathSeparator))] + string(os.PathSeparator) + rqst.GetName() + ".tgz"

	// Now I will save the file to the destination.
	err := ioutil.WriteFile(dest, buf.Bytes(), 0644)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Remove the
	path = strings.Replace(rqst.Path, "\\", "/", -1)
	dest = path[0:strings.LastIndex(path, "/")] + "/" + rqst.GetName() + ".tgz"

	return &filepb.CreateArchiveResponse{
		Result: dest,
	}, nil

}

// Rename a file or a directory.
func (self *server) Rename(ctx context.Context, rqst *filepb.RenameRequest) (*filepb.RenameResponse, error) {

	path := rqst.GetPath()

	// set the correct os path separator.
	path = strings.ReplaceAll(strings.ReplaceAll(path, "\\", string(os.PathSeparator)), "/", string(os.PathSeparator))

	if strings.HasPrefix(path, string(os.PathSeparator)) {
		if len(path) > 1 {
			if strings.HasPrefix(path, string(os.PathSeparator)) {
				path = self.Root + path
			} else {
				path = self.Root + string(os.PathSeparator) + path
			}
		} else {
			path = self.Root
		}
	}

	err := os.Rename(path+string(os.PathSeparator)+rqst.OldName, path+string(os.PathSeparator)+rqst.NewName)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &filepb.RenameResponse{
		Result: true,
	}, nil
}

// Delete a directory
func (self *server) DeleteDir(ctx context.Context, rqst *filepb.DeleteDirRequest) (*filepb.DeleteDirResponse, error) {
	path := rqst.GetPath()

	// set the correct os path separator.
	path = strings.ReplaceAll(strings.ReplaceAll(path, "\\", string(os.PathSeparator)), "/", string(os.PathSeparator))

	if strings.HasPrefix(path, string(os.PathSeparator)) {
		if len(path) > 1 {
			if strings.HasPrefix(path, string(os.PathSeparator)) {
				path = self.Root + path
			} else {
				path = self.Root + string(os.PathSeparator) + path
			}
		} else {
			path = self.Root
		}
	}

	if !Utility.Exists(path) {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), errors.New("No directory with path "+path+" was found!")))
	}
	err := os.RemoveAll(path)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &filepb.DeleteDirResponse{
		Result: true,
	}, nil
}

////////////////////////////////////////////////////////////////////////////////
// File Operation
////////////////////////////////////////////////////////////////////////////////

// Get file info, can be use to get file thumbnail or knowing that a file exist
// or not.
func (self *server) GetFileInfo(ctx context.Context, rqst *filepb.GetFileInfoRequest) (*filepb.GetFileInfoResponse, error) {
	path := rqst.GetPath()

	// set the correct os path separator.
	path = strings.ReplaceAll(strings.ReplaceAll(path, "\\", string(os.PathSeparator)), "/", string(os.PathSeparator))

	if strings.HasPrefix(path, string(os.PathSeparator)) {
		if len(path) > 1 {
			if strings.HasPrefix(path, string(os.PathSeparator)) {
				path = self.Root + path
			} else {
				path = self.Root + string(os.PathSeparator) + path
			}
		} else {
			path = self.Root
		}
	}

	info, err := getFileInfo(path)

	// the file
	f_, err := os.Open(path)
	defer f_.Close()

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	info.Mime, err = Utility.GetFileContentType(f_)
	thumbnailMaxHeight := rqst.GetThumnailHeight()
	thumbnailMaxWidth := rqst.GetThumnailWidth()

	// in case of image...
	if strings.HasPrefix(info.Mime, "image/") {
		if thumbnailMaxHeight > 0 && thumbnailMaxWidth > 0 {
			info.Thumbnail = createThumbnail(f_, int(thumbnailMaxHeight), int(thumbnailMaxWidth))
		}
	}

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	var jsonStr string
	jsonStr, err = Utility.ToJson(info)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &filepb.GetFileInfoResponse{
		Data: jsonStr,
	}, nil
}

// Read file, can be use for small to medium file...
func (self *server) ReadFile(rqst *filepb.ReadFileRequest, stream filepb.FileService_ReadFileServer) error {
	path := rqst.GetPath()

	// set the correct os path separator.
	path = strings.ReplaceAll(strings.ReplaceAll(path, "\\", string(os.PathSeparator)), "/", string(os.PathSeparator))

	if strings.HasPrefix(path, string(os.PathSeparator)) {
		if len(path) > 1 {
			if strings.HasPrefix(path, string(os.PathSeparator)) {
				path = self.Root + path
			} else {
				path = self.Root + string(os.PathSeparator) + path
			}

		} else {
			path = self.Root
		}
	}

	file, err := os.Open(path)
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
			stream.Send(&filepb.ReadFileResponse{
				Data: buffer[:bytesread],
			})
		}
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
	}
	return nil
}

// Save a file on the server...
func (self *server) SaveFile(stream filepb.FileService_SaveFileServer) error {
	// Here I will receive the file
	data := make([]byte, 0)
	var path string
	for {
		rqst, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// Here all data is read...
				err := ioutil.WriteFile(path, data, 0644)

				if err != nil {
					return status.Errorf(
						codes.Internal,
						Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
				}

				// Close the stream...
				stream.SendAndClose(&filepb.SaveFileResponse{
					Result: true,
				})

				return nil
			} else {
				return err
			}
		}

		// Receive message informations.
		switch msg := rqst.File.(type) {
		case *filepb.SaveFileRequest_Path:
			// The roo will be the Root specefied by the server.
			path = msg.Path
			if strings.HasPrefix(path, "/") {
				path = self.Root + string(os.PathSeparator) + path
			}

		case *filepb.SaveFileRequest_Data:
			data = append(data, msg.Data...)
		}
	}

	return nil
}

// Delete file
func (self *server) DeleteFile(ctx context.Context, rqst *filepb.DeleteFileRequest) (*filepb.DeleteFileResponse, error) {
	path := rqst.GetPath()

	// set the correct os path separator.
	path = strings.ReplaceAll(strings.ReplaceAll(path, "\\", string(os.PathSeparator)), "/", string(os.PathSeparator))

	if strings.HasPrefix(path, string(os.PathSeparator)) {
		if len(path) > 1 {
			if strings.HasPrefix(path, string(os.PathSeparator)) {
				path = self.Root + path
			} else {
				path = self.Root + string(os.PathSeparator) + path
			}
		} else {
			path = self.Root
		}
	}

	err := os.Remove(path)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &filepb.DeleteFileResponse{
		Result: true,
	}, nil

}

// That service is use to give access to SQL.
// port number must be pass as argument.
func main() {

	// set the logger.
	grpclog.SetLogger(log.New(os.Stdout, "file_service: ", log.LstdFlags))

	// Set the log information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// The first argument must be the port number to listen to.
	port := defaultPort // the default value.

	if len(os.Args) > 1 {
		port, _ = strconv.Atoi(os.Args[1]) // The second argument must be the port number
	}

	// The actual server implementation.
	s_impl := new(server)
	s_impl.Name = Utility.GetExecName(os.Args[0])
	s_impl.Port = port
	s_impl.Proxy = defaultProxy
	s_impl.Protocol = "grpc"
	s_impl.Domain = domain
	s_impl.Version = "0.0.1"
	s_impl.AllowAllOrigins = allow_all_origins
	s_impl.AllowedOrigins = allowed_origins
	s_impl.PublisherId = "localhost"

	// Here I will retreive the list of connections from file if there are some...
	s_impl.init()

	// First of all I will creat a listener.
	// Create the channel to listen on
	lis, err := net.Listen("tcp", "0.0.0.0:"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("could not list on %s: %s", s_impl.Domain, err)
		return
	}

	// Set the root path if is pass as argument.
	if len(os.Args) > 2 {
		s_impl.Root = os.Args[2]
	}

	if len(s_impl.Root) == 0 {
		log.Fatalln("No root path was given!")
		return
	}

	s = s_impl // keep ref...
	var grpcServer *grpc.Server
	if s_impl.TLS {
		// Load the certificates from disk
		certificate, err := tls.LoadX509KeyPair(s_impl.CertFile, s_impl.KeyFile)
		if err != nil {
			log.Fatalf("could not load server key pair: %s", err)
			return
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
		opts := []grpc.ServerOption{grpc.Creds(creds), grpc.UnaryInterceptor(Interceptors.UnaryAuthInterceptor), grpc.StreamInterceptor(Interceptors.StreamAuthInterceptor)}
		grpcServer = grpc.NewServer(opts...)

	} else {
		grpcServer = grpc.NewServer([]grpc.ServerOption{grpc.UnaryInterceptor(Interceptors.UnaryAuthInterceptor), grpc.StreamInterceptor(Interceptors.StreamAuthInterceptor)}...)
	}

	filepb.RegisterFileServiceServer(grpcServer, s_impl)
	reflection.Register(grpcServer)

	// Here I will make a signal hook to interrupt to exit cleanly.
	go func() {
		log.Println(s_impl.Name + " grpc service is starting")
		// no web-rpc server.
		if err := grpcServer.Serve(lis); err != nil {
			f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
			if err != nil {
				log.Fatalf("error opening file: %v", err)
			}
			defer f.Close()
		}
		log.Println(s_impl.Name + " grpc service is closed")
	}()

	// Wait for signal to stop.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

}

////////////////////////////////////////////////////////////////////////////////
// Utility functions
////////////////////////////////////////////////////////////////////////////////

// Return the list of thumbnail for a given directory...
func (self *server) GetThumbnails(rqst *filepb.GetThumbnailsRequest, stream filepb.FileService_GetThumbnailsServer) error {
	path := rqst.GetPath()

	// The roo will be the Root specefied by the server.
	if strings.HasPrefix(path, "/") {
		path = self.Root + path
		// Set the path separator...
		path = strings.Replace(path, "/", string(os.PathSeparator), -1)
	}

	info, err := readDir(path, rqst.GetRecursive(), rqst.GetThumnailHeight(), rqst.GetThumnailWidth())
	if err != nil {
		return err
	}

	thumbnails := getThumbnails(info)

	// Here I will serialyse the data into JSON.
	jsonStr, err := json.Marshal(thumbnails)
	if err != nil {
		return err
	}

	maxSize := 1024 * 5
	size := int(math.Ceil(float64(len(jsonStr)) / float64(maxSize)))
	for i := 0; i < size; i++ {
		start := i * maxSize
		end := start + maxSize
		var data []byte
		if end > len(jsonStr) {
			data = jsonStr[start:]
		} else {
			data = jsonStr[start:end]
		}
		stream.Send(&filepb.GetThumbnailsResponse{
			Data: data,
		})
	}

	return nil
}

func (self *server) WriteExcelFile(ctx context.Context, rqst *filepb.WriteExcelFileRequest) (*filepb.WriteExcelFileResponse, error) {
	path := rqst.GetPath()

	// The roo will be the Root specefied by the server.
	if strings.HasPrefix(path, "/") {
		path = self.Root + path
		// Set the path separator...
		path = strings.Replace(path, "/", string(os.PathSeparator), -1)
	}

	sheets := make(map[string]interface{}, 0)

	err := json.Unmarshal([]byte(rqst.Data), &sheets)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err = self.writeExcelFile(path, sheets)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	return &filepb.WriteExcelFileResponse{
		Result: true,
	}, nil
}

/**
 * Save excel file to a given destination.
 * The sheets must contain a with values map[pageName] [[], [], []] // 2D array.
 */
func (self *server) writeExcelFile(path string, sheets map[string]interface{}) error {

	xlFile, err_ := xlsx.OpenFile(path)
	var xlSheet *xlsx.Sheet
	if err_ != nil {
		xlFile = xlsx.NewFile()
	}

	for name, data := range sheets {
		xlSheet, _ = xlFile.AddSheet(name)
		values := data.([]interface{})
		// So here I got the xl file open and sheet ready to write into.
		for i := 0; i < len(values); i++ {
			row := xlSheet.AddRow()
			for j := 0; j < len(values[i].([]interface{})); j++ {
				if values[i].([]interface{})[j] != nil {
					cell := row.AddCell()
					if reflect.TypeOf(values[i].([]interface{})[j]).String() == "string" {
						str := values[i].([]interface{})[j].(string)
						// here I will try to format the date time if it can be...
						dateTime, err := Utility.DateTimeFromString(str, "2006-01-02 15:04:05")
						if err != nil {
							cell.SetString(str)
						} else {
							cell.SetDateTime(dateTime)
						}
					} else {
						if values[i].([]interface{})[j] != nil {
							cell.SetValue(values[i].([]interface{})[j])
						}
					}
				}
			}
		}

	}

	// Here I will save the file at the given path...
	err := xlFile.Save(path)

	if err != nil {
		return nil
	}

	return nil
}
