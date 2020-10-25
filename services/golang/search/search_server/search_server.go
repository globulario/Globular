package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/davecourtois/Globular/Interceptors"
	"github.com/davecourtois/Globular/services/golang/search/search_client"
	"github.com/davecourtois/Globular/services/golang/search/searchpb"
	"github.com/davecourtois/Utility"
	"google.golang.org/grpc"

	//	"google.golang.org/grpc/codes"
	globular "github.com/davecourtois/Globular/services/golang/globular_service"

	//"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	//	"google.golang.org/grpc/status"
	"reflect"

	"errors"
	"os/exec"

	"github.com/davecourtois/GoXapian/xapian"
)

// TODO take care of TLS/https
var (
	defaultPort  = 10041
	defaultProxy = 10042

	// By default all origins are allowed.
	allow_all_origins = true

	// comma separeated values.
	allowed_origins string = ""

	domain string = "localhost"
)

// Value need by Globular to start the services...
type server struct {
	// The global attribute of the services.
	Id              string
	Name            string
	Proto           string
	Path            string
	Port            int
	Proxy           int
	AllowAllOrigins bool
	AllowedOrigins  string // comma separated string.
	Protocol        string
	Domain          string
	Description     string
	Keywords        []string
	Repositories    []string
	Discoveries     []string
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

// The list of keywords of the services.
func (self *server) GetKeywords() []string {
	return self.Keywords
}
func (self *server) SetKeywords(keywords []string) {
	self.Keywords = keywords
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
	Utility.RegisterFunction("NewSearchService_Client", search_client.NewSearchService_Client)

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

func (self *server) Stop(context.Context, *searchpb.StopRequest) (*searchpb.StopResponse, error) {
	return &searchpb.StopResponse{}, self.StopService()
}

/**
 * Set base type indexation.
 */
func (self *server) indexJsonObjectField(db xapian.WritableDatabase, termgenerator xapian.TermGenerator, k string, v interface{}, indexs []string) error {
	typeOf := reflect.TypeOf(v).Kind()
	field := strings.ToLower(k)
	field = strings.ToUpper(field[0:1]) + field[1:]
	if typeOf == reflect.String {
		// Index each field with a suitable prefix.
		termgenerator.Index_text(v, uint(1), "X"+field)
		// # Index fields without prefixes for general search.
		if Utility.Contains(indexs, k) {
			termgenerator.Index_text(v)
			termgenerator.Increase_termpos()
		}
	} else if typeOf == reflect.Bool {

	} else if typeOf == reflect.Int || typeOf == reflect.Int8 || typeOf == reflect.Int16 || typeOf == reflect.Int32 || typeOf == reflect.Int64 {

	} else if typeOf == reflect.Float32 || typeOf == reflect.Float64 {

	} else if typeOf == reflect.Struct {
		//v := reflect.ValueOf(v)
	}
	return nil
}

// Return the underlying engine version.
func (self *server) GetEngineVersion(ctx context.Context, rqst *searchpb.GetEngineVersionRequest) (*searchpb.GetEngineVersionResponse, error) {
	v := xapian.Version_string()
	return &searchpb.GetEngineVersionResponse{
		Message: v,
	}, nil
}

// Remove a document from the db
func (self *server) DeleteDocument(ctx context.Context, rqst *searchpb.DeleteDocumentRequest) (*searchpb.DeleteDocumentResponse, error) {

	db := xapian.NewWritableDatabase(rqst.Path, xapian.DB_CREATE_OR_OPEN)
	defer xapian.DeleteWritableDatabase(db)

	// Begin the transaction.
	db.Begin_transaction(true)

	id := "Q" + strings.ToUpper(rqst.Id[0:1]) + strings.ToLower(rqst.Id[1:])

	// Delete a document from the database.
	db.Delete_document(id)

	db.Commit_transaction()

	db.Close()

	return &searchpb.DeleteDocumentResponse{}, nil
}

/**
 * Index a json object.
 */
func (self *server) indexJsonObject(db xapian.WritableDatabase, obj map[string]interface{}, language string, id string, indexs []string, data string) error {

	doc := xapian.NewDocument()
	defer xapian.DeleteDocument(doc)

	termgenerator := xapian.NewTermGenerator()
	stemmer := xapian.NewStem(language)
	defer xapian.DeleteStem(stemmer)
	termgenerator.Set_stemmer(xapian.NewStem(language))
	defer xapian.DeleteTermGenerator(termgenerator)
	termgenerator.Set_document(doc)

	// Here I will iterate over the object and append fields to the document.
	// Here I will index each field's
	for k, v := range obj {

		if v != nil {
			typeOf := reflect.TypeOf(v).Kind()
			field := strings.ToLower(k)
			field = strings.ToUpper(field[0:1]) + field[1:]
			if typeOf == reflect.Map {
				// In case of recursive structure.
				self.indexJsonObject(db, v.(map[string]interface{}), language, id, indexs, data)
			} else if typeOf == reflect.Slice {
				s := reflect.ValueOf(v)
				for i := 0; i < s.Len(); i++ {
					_v := s.Index(i)

					typeOf := reflect.TypeOf(_v).Kind()
					if typeOf == reflect.Map {
						// Slice of object.
						self.indexJsonObject(db, _v.Interface().(map[string]interface{}), language, id, indexs, data)
					} else {
						// Slice of literal type.

						self.indexJsonObjectField(db, termgenerator, k, _v.Interface(), indexs)
					}
				}
			} else {
				self.indexJsonObjectField(db, termgenerator, k, v, indexs)
			}
		}

	}

	// Here I will set object metadata.
	var infos map[string]interface{}
	if len(data) > 0 {
		infos = make(map[string]interface{}, 0)
		json.Unmarshal([]byte(data), &infos)
	} else {
		infos = make(map[string]interface{}, 0)
	}

	if len(id) > 0 {
		infos["id"] = id
		infos["type"] = "object"
	}

	jsonStr, _ := Utility.ToJson(infos)
	doc.Set_data(jsonStr)

	// Here If the object contain an id I will add it as boolean term and
	// replace existing document or create it.
	if len(id) > 0 {
		_id := "Q" + Utility.ToString(obj[id])
		doc.Add_boolean_term(_id)
		db.Replace_document(_id, doc)
	} else {
		db.Add_document(doc)
	}

	return nil
}

// Return the number of document in a database.
func (self *server) Count(ctx context.Context, rqst *searchpb.CountRequest) (*searchpb.CountResponse, error) {

	db := xapian.NewDatabase(rqst.Path)
	//defer db.Close()
	defer xapian.DeleteDatabase(db)
	count := int32(db.Get_doccount())
	return &searchpb.CountResponse{
		Result: count,
	}, nil

}

// Here I will append the sub-database.
func (self *server) addSubDBs(db xapian.Database, path string) []xapian.Database {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	// keep track of sub database.
	subDbs := make([]xapian.Database, 0)

	// Add the database.
	for i := 0; i < len(files); i++ {
		if files[i].IsDir() {
			_db := xapian.NewDatabase(path + "/" + files[i].Name())
			db.Add_database(_db)
			subDbs = append(subDbs, _db)
			subDbs = append(subDbs, self.addSubDBs(db, path+"/"+files[i].Name())...)
		}
	}

	return subDbs
}

// Search documents...
func (self *server) searchDocuments(paths []string, language string, fields []string, queryStr string, offset int32, pageSize int32, snippetLength int32) ([]*searchpb.SearchResult, error) {

	if len(paths) == 0 {
		return nil, errors.New("No database was path given!")
	}

	db := xapian.NewDatabase(paths[0])
	defer xapian.DeleteDatabase(db)

	// Open the db for read...
	for i := 1; i < len(paths); i++ {
		path := paths[i]
		_db := xapian.NewDatabase(path)
		defer xapian.DeleteDatabase(_db)

		// Now I will recursively append data base is there is some subdirectory...
		subDbs := self.addSubDBs(_db, path)
		db.Add_database(_db)

		// clear pointer memory...
		for j := 0; j < len(subDbs); j++ {
			defer xapian.DeleteDatabase(subDbs[j])
		}
	}

	queryParser := xapian.NewQueryParser()
	defer xapian.DeleteQueryParser(queryParser)

	stemmer := xapian.NewStem(language)
	defer xapian.DeleteStem(stemmer)

	queryParser.Set_stemmer(stemmer)
	queryParser.Set_stemming_strategy(xapian.XapianQueryParserStem_strategy(xapian.QueryParserSTEM_SOME))

	// Append the list of field to search for
	for i := 0; i < len(fields); i++ {
		field := strings.ToUpper(fields[i][0:1]) + strings.ToLower(fields[i][1:])
		queryParser.Add_prefix(field, "X")
	}

	// Generate the query from the given string.
	query := queryParser.Parse_query(queryStr)
	defer xapian.DeleteQuery(query)

	enquire := xapian.NewEnquire(db)
	defer xapian.DeleteEnquire(enquire)

	enquire.Set_query(query)

	// Here I will retreive the results.
	mset := enquire.Get_mset(uint(offset), uint(pageSize))
	defer xapian.DeleteMSet(mset)

	results := make([]*searchpb.SearchResult, 0)

	// Here I will
	for i := 0; i < mset.Size(); i++ {

		docId := mset.Get_docid(uint(i))
		result := new(searchpb.SearchResult)
		result.DocId = Utility.ToString(int(docId))
		doc := mset.Get_document(uint(i))
		defer xapian.DeleteDocument(doc) // release memory
		result.Data = doc.Get_data()     // Set the necessery data to retreive document in it db
		result.Rank = int32(mset.Get_document_percentage(uint(i)))

		it := enquire.Get_matching_terms_begin(mset.Get_hit(uint(i)))

		terms := make([]string, 0)
		for !it.Equals(enquire.Get_matching_terms_end(mset.Get_hit(uint(i)))) {
			term := it.Get_term()
			terms = append(terms, term)
			it.Next()
		}

		infos := make(map[string]interface{}, 0)
		err := json.Unmarshal([]byte(doc.Get_data()), &infos)
		if err != nil {
			return nil, err
		}
		if infos["type"].(string) == "file" {
			result.Snippets, err = self.snippets(terms, infos["path"].(string), infos["mime"].(string), int(snippetLength))
			if err != nil {
				return nil, err
			}
		}

		results = append(results, result)

	}

	return results, nil

}

// Search documents
func (self *server) SearchDocuments(ctx context.Context, rqst *searchpb.SearchDocumentsRequest) (*searchpb.SearchDocumentsResponse, error) {
	var results []*searchpb.SearchResult
	var err error

	results, err = self.searchDocuments(rqst.Paths, rqst.Language, rqst.Fields, rqst.Query, rqst.Offset, rqst.PageSize, rqst.SnippetLength)
	if err != nil {
		return nil, err
	}

	return &searchpb.SearchDocumentsResponse{
		Results: results,
	}, nil

}

/**
 * Index the a dir and it content.
 */
func (self *server) indexDir(dbPath string, dirPath string, language string) error {

	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		return err
	}

	if !dirInfo.IsDir() {
		return errors.New("The file " + dirPath + " is not a directory ")
	}

	// So here I will create the directory entry in the dbPath
	err = Utility.CreateDirIfNotExist(dbPath)

	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	// The database path.
	db := xapian.NewWritableDatabase(dbPath, xapian.DB_CREATE_OR_OPEN)
	defer xapian.DeleteWritableDatabase(db)
	db.Begin_transaction()

	if err != nil {
		db.Cancel_transaction()
		db.Close()
		return err
	}

	// Now I will index files and recursively index dir content.
	for _, file := range files {
		if file.IsDir() {
			err := self.indexDir(dbPath+"/"+file.Name(), dirPath+"/"+file.Name(), language)
			if err != nil {
				return err
			}
		} else {
			// Here I will index the file contain in the directory.
			path := dirPath + "/" + file.Name()
			err := self.indexFile(db, path, language)
			if err != nil {
				fmt.Println(file.Name(), err)
			}
		}
	}

	mime := "folder"

	// set document meta data.
	modified := "D" + dirInfo.ModTime().Format("YYYYMMDD")
	doc := xapian.NewDocument()
	defer xapian.DeleteDocument(doc)
	doc.Add_term(modified)
	doc.Add_term("P" + dirPath)
	doc.Add_term("T" + mime)

	id := "Q" + Utility.GenerateUUID(dirPath)
	doc.Add_boolean_term(id)

	infos := make(map[string]interface{}, 0)

	infos["path"] = dirPath
	infos["type"] = "file"

	infos["mime"] = mime

	jsonStr, _ := Utility.ToJson(infos)
	doc.Set_data(jsonStr)

	// Create the directory information.
	db.Replace_document(id, doc)

	db.Commit_transaction()
	db.Close()

	return err
}

/**
 * Index the content of a dir and it content.
 */
func (self *server) IndexDir(ctx context.Context, rqst *searchpb.IndexDirRequest) (*searchpb.IndexDirResponse, error) {

	err := self.indexDir(rqst.DbPath, rqst.DirPath, rqst.Language)
	if err != nil {
		return nil, err
	}

	return &searchpb.IndexDirResponse{}, nil
}

func (self *server) indexFile(db xapian.WritableDatabase, path string, language string) error {
	path = strings.ReplaceAll(strings.ReplaceAll(path, "\\", string(os.PathSeparator)), "/", string(os.PathSeparator))

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	mime, err := Utility.GetFileContentType(f)
	if err != nil {
		return err
	}

	// create the document.
	doc := xapian.NewDocument()
	defer xapian.DeleteDocument(doc)

	// create the term generator for the file.
	termgenerator := xapian.NewTermGenerator()
	stemmer := xapian.NewStem(language)
	defer xapian.DeleteStem(stemmer)
	termgenerator.Set_stemmer(xapian.NewStem(language))
	defer xapian.DeleteTermGenerator(termgenerator)
	termgenerator.Set_document(doc)

	// Now I will index file metat
	fileStat, err := os.Stat(path)
	if err != nil {
		return err
	}

	// set document meta data.
	modified := "D" + fileStat.ModTime().Format("YYYYMMDD")
	doc.Add_term(modified)
	doc.Add_term("P" + path)
	doc.Add_term("T" + mime)

	if mime == "application/pdf" {
		err = self.indexPdfFile(db, path, doc, termgenerator)
		if err != nil {
			return err
		}
	} else if strings.HasPrefix(mime, "text") {
		text, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		termgenerator.Index_text(strings.ToLower(string(text)))
		termgenerator.Increase_termpos()

	} else {
		return errors.New("Unsuported file type! " + mime)
	}

	id := "Q" + Utility.GenerateUUID(path)
	doc.Add_boolean_term(id)

	infos := make(map[string]interface{}, 0)
	infos["path"] = path
	infos["type"] = "file"
	infos["mime"] = mime

	jsonStr, _ := Utility.ToJson(infos)
	doc.Set_data(jsonStr)

	db.Replace_document(id, doc)

	return nil
}

// pdftotext bin must be install on the server to be able to generate text
// file from pdf file.
// On linux type...
// sudo apt-get install poppler-utils
func (self *server) pdfToText(path string) (string, error) {
	// First of all I will test if pdftotext is install.
	cmd := exec.Command("pdftotext", path)
	_, err := cmd.Output()
	if err != nil {
		return "", err
	}

	_path := path[0:strings.LastIndex(path, ".")] + ".txt"
	defer os.Remove(_path)

	// Here I will index the text file
	text, err := ioutil.ReadFile(_path)
	if err != nil {
		return "", err
	}

	return string(text), err

}

// That function is use to generate a snippet from a text file.
func (self *server) snippets(terms []string, path string, mime string, length int) ([]string, error) {

	length += 15

	// Here I will read the file and try to generate a snippet for it.
	var text string
	snippets := make([]string, 0)
	var err error
	if mime == "application/pdf" {
		text, err = self.pdfToText(path)
		if err != nil {
			return nil, err
		}
	} else if strings.HasPrefix(mime, "text/") {
		_text, err := ioutil.ReadFile(path)
		text = string(_text)
		if err != nil {
			return nil, err
		}
	}

	for i := 0; i < len(terms); i++ {
		index := 0
		_text := strings.ToLower(text)
		for index != -1 {
			_text = _text[index:]
			index = strings.Index(_text, terms[i])
			if index == -1 {
				break
			}
			var snippet string
			if index < length/2 {
				if index+length < len(text) {
					snippet = _text[index : index+length]
				} else {
					snippet = _text[index:]
				}

			} else {
				if index+(length/2) < len(text) {
					snippet = _text[index-(length/2) : index+length]
				} else {
					snippet = _text[index-(length/2):]
				}

			}
			re := regexp.MustCompile("[[:^ascii:]]")
			snippet = re.ReplaceAllLiteralString(snippet, "")
			re = regexp.MustCompile(`\r?\n`)
			snippet = re.ReplaceAllString(snippet, " ")

			// Now I will remove the first and last word as needed.
			words := strings.Split(snippet, " ")
			if len(words) > 0 {
				if words[0] != terms[i] {
					words = words[1:]
				}
				if len(words) > 2 {
					if words[len(words)-1] != terms[i] {
						words = words[0 : len(words)-2]
					}
				}
			}

			snippet = strings.Join(words, " ")

			// Now I will set html balise.
			snippet = "<div class='snippet'>" + strings.ReplaceAll(snippet, terms[i], "<div class='founded-term'>"+terms[i]+"</div>") + " ... </div>"

			snippets = append(snippets, snippet)

			index += len(terms[i])

		}
	}

	return snippets, nil
}

// Indexation of pdf file.
func (self *server) indexPdfFile(db xapian.WritableDatabase, path string, doc xapian.Document, termgenerator xapian.TermGenerator) error {
	text, err := self.pdfToText(path)
	if err != nil {
		return err
	}
	termgenerator.Index_text(strings.ToLower(string(text)))
	termgenerator.Increase_termpos()
	return nil
}

// Indexation of a text (docx, pdf,xlsx...) file.
func (self *server) IndexFile(ctx context.Context, rqst *searchpb.IndexFileRequest) (*searchpb.IndexFileResponse, error) {

	// The file must be accessible on the server side.
	if !Utility.Exists(rqst.FilePath) {
		return nil, errors.New("File " + rqst.FilePath + " was not found!")
	}

	fileInfo, err := os.Stat(rqst.FilePath)
	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		return nil, errors.New("The file " + rqst.FilePath + " is a directory ")
	}

	// The database path.
	db := xapian.NewWritableDatabase(rqst.DbPath, xapian.DB_CREATE_OR_OPEN)
	defer xapian.DeleteWritableDatabase(db)
	db.Begin_transaction()

	err = self.indexFile(db, rqst.FilePath, rqst.Language)
	if err != nil {
		db.Cancel_transaction()
		db.Close()
		return nil, err
	}

	db.Commit_transaction()
	db.Close()

	return &searchpb.IndexFileResponse{}, nil
}

// That function is use to index JSON object/array of object
func (self *server) IndexJsonObject(ctx context.Context, rqst *searchpb.IndexJsonObjectRequest) (*searchpb.IndexJsonObjectResponse, error) {

	db := xapian.NewWritableDatabase(rqst.Path, xapian.DB_CREATE_OR_OPEN)
	defer xapian.DeleteWritableDatabase(db)

	// Begin the transaction.
	db.Begin_transaction(true)

	var obj interface{}
	var err error
	err = json.Unmarshal([]byte(rqst.JsonStr), &obj)
	if err != nil {
		return nil, err
	}

	// Now I will append the object into the database.
	switch v := obj.(type) {
	case map[string]interface{}:
		err = self.indexJsonObject(db, v, rqst.Language, rqst.Id, rqst.Indexs, rqst.Data)

	case []interface{}:
		for i := 0; i < len(v); i++ {
			err := self.indexJsonObject(db, v[i].(map[string]interface{}), rqst.Language, rqst.Id, rqst.Indexs, rqst.Data)
			if err != nil {
				break
			}
		}
	}

	if err != nil {
		db.Cancel_transaction()
		db.Close()
		return nil, err
	}

	// Write object int he database.
	db.Commit_transaction()
	db.Close()

	return &searchpb.IndexJsonObjectResponse{}, nil

}

// That service is use to give access to SQL.
// port number must be pass as argument.
func main() {

	// Set the log information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// The actual server implementation.
	s_impl := new(server)
	s_impl.Name = string(searchpb.File_services_proto_search_proto.Services().Get(0).FullName())
	s_impl.Proto = searchpb.File_services_proto_search_proto.Path()
	s_impl.Port = defaultPort
	s_impl.Proxy = defaultProxy
	s_impl.Protocol = "grpc"
	s_impl.Domain = domain
	s_impl.Version = "0.0.1"
	s_impl.PublisherId = domain
	s_impl.Permissions = make([]interface{}, 0)

	// TODO set it from the program arguments...
	s_impl.AllowAllOrigins = allow_all_origins
	s_impl.AllowedOrigins = allowed_origins
	s_impl.Keywords = make([]string, 0)
	s_impl.Repositories = make([]string, 0)
	s_impl.Discoveries = make([]string, 0)

	// set the logger.

	// Set the root path if is pass as argument.
	if len(os.Args) > 2 {
		s_impl.Root = os.Args[2]
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

			s_impl.PublishService(*publishCommand_domain, *publishCommand_user, *publishCommand_password)
		}
	} else {
		// Register the echo services
		searchpb.RegisterSearchServiceServer(s_impl.grpcServer, s_impl)
		reflection.Register(s_impl.grpcServer)

		// Start the service.
		s_impl.StartService()
	}
}
