package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc/metadata"

	"github.com/globulario/Globular/Interceptors"
	"github.com/globulario/Globular/services/golang/file/file_client"
	"github.com/globulario/Globular/services/golang/file/filepb"
	globular "github.com/globulario/Globular/services/golang/globular_service"
	"github.com/davecourtois/Utility"
	"github.com/nfnt/resize"
	"github.com/polds/imgbase64"
	"github.com/tealeg/xlsx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// TODO take care of TLS/https
var (
	defaultPort  = 10043
	defaultProxy = 10044

	// By default all origins are allowed.
	allow_all_origins = true

	// comma separeated values.
	allowed_origins string = ""

	// The default domain.
	domain string = "localhost"

	//s *server
)

// Value need by Globular to start the services...
type server struct {
	// The global attribute of the services.
	Id                 string
	Name               string
	Path               string
	Proto              string
	Port               int
	Proxy              int
	AllowAllOrigins    bool
	AllowedOrigins     string // comma separated string.
	Protocol           string
	Domain             string
	Description        string
	Keywords           []string
	Repositories       []string
	Discoveries        []string
	CertFile           string
	CertAuthorityTrust string
	KeyFile            string
	TLS                bool
	Version            string
	PublisherId        string
	KeepUpToDate       bool
	KeepAlive          bool
	Permissions        []interface{} // contains the action permission for the services.

	// The grpc server.
	grpcServer *grpc.Server

	// Specific to file server.
	Root string
}

// Globular services implementation...
// The id of a particular service instance.
func (self *server) GetId() string {
	return self.Id
}
func (self *server) SetId(id string) {
	self.Id = id
}

// The name of a service, must be the gRpc Service name.
func (self *server) GetName() string {
	return self.Name
}
func (self *server) SetName(name string) {
	self.Name = name
}

// The description of the service
func (self *server) GetDescription() string {
	return self.Description
}
func (self *server) SetDescription(description string) {
	self.Description = description
}

// The list of keywords of the services.
func (self *server) GetKeywords() []string {
	return self.Keywords
}
func (self *server) SetKeywords(keywords []string) {
	self.Keywords = keywords
}

func (self *server) GetRepositories() []string {
	return self.Repositories
}
func (self *server) SetRepositories(repositories []string) {
	self.Repositories = repositories
}

func (self *server) GetDiscoveries() []string {
	return self.Discoveries
}
func (self *server) SetDiscoveries(discoveries []string) {
	self.Discoveries = discoveries
}

// Dist
func (self *server) Dist(path string) error {

	return globular.Dist(path, self)
}

func (self *server) GetPlatform() string {
	return globular.GetPlatform()
}

func (self *server) PublishService(address string, user string, password string) error {
	return globular.PublishService(address, user, password, self)
}

// The path of the executable.
func (self *server) GetPath() string {
	return self.Path
}
func (self *server) SetPath(path string) {
	self.Path = path
}

// The path of the .proto file.
func (self *server) GetProto() string {
	return self.Proto
}
func (self *server) SetProto(proto string) {
	self.Proto = proto
}

// The gRpc port.
func (self *server) GetPort() int {
	return self.Port
}
func (self *server) SetPort(port int) {
	self.Port = port
}

// The reverse proxy port (use by gRpc Web)
func (self *server) GetProxy() int {
	return self.Proxy
}
func (self *server) SetProxy(proxy int) {
	self.Proxy = proxy
}

// Can be one of http/https/tls
func (self *server) GetProtocol() string {
	return self.Protocol
}
func (self *server) SetProtocol(protocol string) {
	self.Protocol = protocol
}

// Return true if all Origins are allowed to access the mircoservice.
func (self *server) GetAllowAllOrigins() bool {
	return self.AllowAllOrigins
}
func (self *server) SetAllowAllOrigins(allowAllOrigins bool) {
	self.AllowAllOrigins = allowAllOrigins
}

// If AllowAllOrigins is false then AllowedOrigins will contain the
// list of address that can reach the services.
func (self *server) GetAllowedOrigins() string {
	return self.AllowedOrigins
}

func (self *server) SetAllowedOrigins(allowedOrigins string) {
	self.AllowedOrigins = allowedOrigins
}

// Can be a ip address or domain name.
func (self *server) GetDomain() string {
	return self.Domain
}
func (self *server) SetDomain(domain string) {
	self.Domain = domain
}

// TLS section

// If true the service run with TLS. The
func (self *server) GetTls() bool {
	return self.TLS
}
func (self *server) SetTls(hasTls bool) {
	self.TLS = hasTls
}

// The certificate authority file
func (self *server) GetCertAuthorityTrust() string {
	return self.CertAuthorityTrust
}
func (self *server) SetCertAuthorityTrust(ca string) {
	self.CertAuthorityTrust = ca
}

// The certificate file.
func (self *server) GetCertFile() string {
	return self.CertFile
}
func (self *server) SetCertFile(certFile string) {
	self.CertFile = certFile
}

// The key file.
func (self *server) GetKeyFile() string {
	return self.KeyFile
}
func (self *server) SetKeyFile(keyFile string) {
	self.KeyFile = keyFile
}

// The service version
func (self *server) GetVersion() string {
	return self.Version
}
func (self *server) SetVersion(version string) {
	self.Version = version
}

// The publisher id.
func (self *server) GetPublisherId() string {
	return self.PublisherId
}
func (self *server) SetPublisherId(publisherId string) {
	self.PublisherId = publisherId
}

func (self *server) GetKeepUpToDate() bool {
	return self.KeepUpToDate
}
func (self *server) SetKeepUptoDate(val bool) {
	self.KeepUpToDate = val
}

func (self *server) GetKeepAlive() bool {
	return self.KeepAlive
}
func (self *server) SetKeepAlive(val bool) {
	self.KeepAlive = val
}

func (self *server) GetPermissions() []interface{} {
	return self.Permissions
}
func (self *server) SetPermissions(permissions []interface{}) {
	self.Permissions = permissions
}

// Create the configuration file if is not already exist.
func (self *server) Init() error {

	// That function is use to get access to other server.
	Utility.RegisterFunction("NewFileService_Client", file_client.NewFileService_Client)

	// Get the configuration path.
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	err := globular.InitService(dir+"/config.json", self)
	if err != nil {
		return err
	}

	// Initialyse GRPC server.
	self.grpcServer, err = globular.InitGrpcServer(self, Interceptors.ServerUnaryInterceptor, Interceptors.ServerStreamInterceptor)
	if err != nil {
		return err
	}

	return nil

}

// Save the configuration values.
func (self *server) Save() error {
	// Create the file...
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return globular.SaveService(dir+"/config.json", self)
}

func (self *server) StartService() error {
	return globular.StartService(self, self.grpcServer)
}

func (self *server) StopService() error {
	return globular.StopService(self, self.grpcServer)
}

func (self *server) Stop(context.Context, *filepb.StopRequest) (*filepb.StopResponse, error) {
	return &filepb.StopResponse{}, self.StopService()
}

/**
 * Create a thumbnail...
 */
func createThumbnail(file *os.File, thumbnailMaxHeight int, thumbnailMaxWidth int) string {
	// Set the buffer pointer back to the begening of the file...
	file.Seek(0, 0)
	var originalImg image.Image
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

func getFileInfo(s *server, path string) (*fileInfo, error) {
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
func readDir(s *server, path string, recursive bool, thumbnailMaxWidth int32, thumbnailMaxHeight int32) (*fileInfo, error) {

	// get the file info
	info, err := getFileInfo(s, path)
	if err != nil {
		return nil, err
	}
	if info.IsDir == false {
		return nil, errors.New(path + " is a directory")
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, f := range files {

		if f.IsDir() {
			if recursive {
				info_, err := readDir(s, path+string(os.PathSeparator)+f.Name(), recursive, thumbnailMaxWidth, thumbnailMaxHeight)
				if err != nil {
					return nil, err
				}
				info.Files = append(info.Files, info_)
			}
		} else {
			info_, err := getFileInfo(s, path+string(os.PathSeparator)+f.Name())

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

	info, err := readDir(self, path, rqst.GetRecursive(), rqst.GetThumnailWidth(), rqst.GetThumnailHeight())
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

	// Set Create dir Permission.
	ressource_client, err := Interceptors.GetRessourceClient(self.GetDomain())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// The token and the application id.
	var token string
	var clientId string

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		token = strings.Join(md["token"], "")
	}

	var expiredAt int64
	clientId, _, expiredAt, err = Interceptors.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	if expiredAt < time.Now().Unix() {
		return nil, errors.New("The token is expired!")
	}

	err = ressource_client.CreateDirPermissions(token, rqst.GetPath(), rqst.GetName())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Here I will set the ressource owner for the directory.
	if strings.HasSuffix(rqst.GetPath(), "/") {
		ressource_client.SetRessourceOwner(clientId, rqst.GetPath()+rqst.GetName(), "")
	} else {
		ressource_client.SetRessourceOwner(clientId, rqst.GetPath()+"/"+rqst.GetName(), "")
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

	// Rename permissions
	ressource_client, err := Interceptors.GetRessourceClient(self.GetDomain())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err = ressource_client.RenameFilePermission(rqst.GetPath(), rqst.GetOldName(), rqst.GetNewName())
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

	ressource_client, err := Interceptors.GetRessourceClient(self.GetDomain())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	// Set permission on the ressource.
	err = ressource_client.DeleteDirPermissions(rqst.GetPath())
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

	info, err := getFileInfo(self, path)

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

	ressource_client, err := Interceptors.GetRessourceClient(self.GetDomain())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			Utility.JsonErrorStr(Utility.FunctionName(), Utility.FileLine(), err))
	}

	err = ressource_client.DeleteFilePermissions(rqst.GetPath())
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

	// The actual server implementation.
	s_impl := new(server)

	// The name must the same as the grpc service name.
	s_impl.Name = string(filepb.File_services_proto_file_proto.Services().Get(0).FullName())
	s_impl.Proto = filepb.File_services_proto_file_proto.Path()
	s_impl.Port = defaultPort
	s_impl.Proxy = defaultProxy
	s_impl.Protocol = "grpc"
	s_impl.Domain = domain
	s_impl.Version = "0.0.1"
	s_impl.AllowAllOrigins = allow_all_origins
	s_impl.AllowedOrigins = allowed_origins
	s_impl.PublisherId = "globulario"
	s_impl.Permissions = make([]interface{}, 12)
	s_impl.Keywords = make([]string, 0)
	s_impl.Repositories = make([]string, 0)
	s_impl.Discoveries = make([]string, 0)

	// So here I will set the default permissions for services actions.
	// Permission are use in conjonctions of ressource.
	s_impl.Permissions[0] = map[string]interface{}{"action": "/file.FileService/ReadDir", "actionParameterRessourcePermissions": []interface{}{map[string]interface{}{"Index": 0, "Permission": 4}}}
	s_impl.Permissions[1] = map[string]interface{}{"action": "/file.FileService/CreateDir", "actionParameterRessourcePermissions": []interface{}{map[string]interface{}{"Index": 0, "Permission": 2}}}
	s_impl.Permissions[2] = map[string]interface{}{"action": "/file.FileService/DeleteDir", "actionParameterRessourcePermissions": []interface{}{map[string]interface{}{"Index": 0, "Permission": 1}}}
	s_impl.Permissions[3] = map[string]interface{}{"action": "/file.FileService/Rename", "actionParameterRessourcePermissions": []interface{}{map[string]interface{}{"Index": 0, "Permission": 2}}}
	s_impl.Permissions[4] = map[string]interface{}{"action": "/file.FileService/GetFileInfo", "actionParameterRessourcePermissions": []interface{}{map[string]interface{}{"Index": 0, "Permission": 4}}}
	s_impl.Permissions[5] = map[string]interface{}{"action": "/file.FileService/ReadFile", "actionParameterRessourcePermissions": []interface{}{map[string]interface{}{"Index": 0, "Permission": 4}}}
	s_impl.Permissions[6] = map[string]interface{}{"action": "/file.FileService/SaveFile", "actionParameterRessourcePermissions": []interface{}{map[string]interface{}{"Index": 0, "Permission": 2}}}
	s_impl.Permissions[7] = map[string]interface{}{"action": "/file.FileService/DeleteFile", "actionParameterRessourcePermissions": []interface{}{map[string]interface{}{"Index": 0, "Permission": 1}}}
	s_impl.Permissions[8] = map[string]interface{}{"action": "/file.FileService/GetThumbnails", "actionParameterRessourcePermissions": []interface{}{map[string]interface{}{"Index": 0, "Permission": 4}}}
	s_impl.Permissions[9] = map[string]interface{}{"action": "/file.FileService/WriteExcelFile", "actionParameterRessourcePermissions": []interface{}{map[string]interface{}{"Index": 0, "Permission": 2}}}
	s_impl.Permissions[10] = map[string]interface{}{"action": "/file.FileService/CreateAchive", "actionParameterRessourcePermissions": []interface{}{map[string]interface{}{"Index": 0, "Permission": 2}}}
	s_impl.Permissions[11] = map[string]interface{}{"action": "/file.FileService/FileUploadHandler", "actionParameterRessourcePermissions": []interface{}{map[string]interface{}{"Index": 0, "Permission": 1}}}

	// Set the root path if is pass as argument.
	if len(s_impl.Root) == 0 {
		s_impl.Root = os.TempDir()
	}

	// Here I will retreive the list of connections from file if there are some...
	err := s_impl.Init()
	if err != nil {
		log.Fatalf("Fail to initialyse service %s: %s", s_impl.Name, s_impl.Id, err)
	}

	if len(os.Args) == 2 {
		s_impl.Port, _ = strconv.Atoi(os.Args[1]) // The second argument must be the port number
	}

	if len(os.Args) > 2 {
		publishCommand := flag.NewFlagSet("publish", flag.ExitOnError)
		publishCommand_domain := publishCommand.String("a", "", "The address(domain ex. my.domain.com:8080) of your backend (Required)")
		publishCommand_user := publishCommand.String("u", "", "The user (Required)")
		publishCommand_password := publishCommand.String("p", "", "The password (Required)")

		switch os.Args[1] {
		case "publish":
			publishCommand.Parse(os.Args[2:])
		default:
			flag.PrintDefaults()
			os.Exit(1)
		}

		if publishCommand.Parsed() {
			// Required Flags
			if *publishCommand_domain == "" {
				publishCommand.PrintDefaults()
				os.Exit(1)
			}

			if *publishCommand_user == "" {
				publishCommand.PrintDefaults()
				os.Exit(1)
			}

			if *publishCommand_password == "" {
				publishCommand.PrintDefaults()
				os.Exit(1)
			}

			err := s_impl.PublishService(*publishCommand_domain, *publishCommand_user, *publishCommand_password)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("Your service was publish successfuly!")
			}
		}
	} else {
		// Register the echo services
		filepb.RegisterFileServiceServer(s_impl.grpcServer, s_impl)
		reflection.Register(s_impl.grpcServer)

		// Start the service.
		s_impl.StartService()
	}
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

	info, err := readDir(self, path, rqst.GetRecursive(), rqst.GetThumnailHeight(), rqst.GetThumnailWidth())
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
