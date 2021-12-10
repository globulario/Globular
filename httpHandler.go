package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/davecourtois/Utility"
	config_ "github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/rbac/rbacpb"
	"github.com/globulario/services/golang/security"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

func getChecksumHanldler(w http.ResponseWriter, r *http.Request) {

	//add prefix and clean
	w.Header().Set("Content-Type", "application/text")
	setupResponse(&w, r)
	w.WriteHeader(http.StatusCreated)

	execPath := Utility.GetExecName(os.Args[0])
	if Utility.Exists("/usr/local/share/globular/Globular") {
		execPath = "/usr/local/share/globular/Globular"
	}
	fmt.Fprint(w, Utility.CreateFileChecksum(execPath))
}

/**
 * Return the service configuration
 */
func getConfigHanldler(w http.ResponseWriter, r *http.Request) {
	// if the host is not the same...
	/*
		if globule.Domain != r.Host {
			//log.Println("------------> request redirected " + globule.Protocol+"://"+r.Host + r.URL.String())
			//http.Redirect(w, r, globule.Protocol+"://"+r.Host + r.URL.String(), http.StatusMovedPermanently)
			// return
			client, err := globule.getHttpClient(r.Host)
			if err == nil {
				rsp, err := client.Get(r.URL.String())
				if err == nil {
				 log.Println("--------------> response found from ", r.Host, " with status ", rsp.StatusCode)
				}else{
					log.Println("60 ----> error ", err)
				}
			}else{
				log.Println("63 ----> error ", err)
			}
		}
	*/

	//add prefix and clean
	config := globule.getConfig()

	// give list of path...
	config["Root"] = config_.GetRootDir()
	config["DataPath"] = config_.GetDataDir()
	config["ConfigPath"] = config_.GetConfigDir()
	config["WebRoot"] = config_.GetWebRootDir()

	w.Header().Set("Content-Type", "application/json")
	setupResponse(&w, r)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(config)
}

func dealwithErr(err error) {
	if err != nil {
		fmt.Println(err)
		//os.Exit(-1)
	}
}

func getHardwareData(w http.ResponseWriter, r *http.Request) {
	runtimeOS := runtime.GOOS

	// memory
	vmStat, err := mem.VirtualMemory()
	dealwithErr(err)

	stats := make(map[string]interface{})

	// disk - start from "/" mount point for Linux
	// might have to change for Windows!!
	// don't have a Window to test this out, if detect OS == windows
	// then use "\" instead of "/"

	diskStat, err := disk.Usage("/")
	dealwithErr(err)

	// cpu - get CPU number of cores and speed
	cpuStat, err := cpu.Info()
	dealwithErr(err)
	percentage, err := cpu.Percent(0, true)
	dealwithErr(err)

	// host or machine kernel, uptime, platform Info
	hostStat, err := host.Info()
	dealwithErr(err)

	// get interfaces MAC/hardware address
	interfStat, err := net.Interfaces()
	dealwithErr(err)

	stats["os"] = runtimeOS
	stats["memory"] = make(map[string]interface{}, 0)
	stats["memory"].(map[string]interface{})["total"] = strconv.FormatUint(vmStat.Total, 10)
	stats["memory"].(map[string]interface{})["free"] = strconv.FormatUint(vmStat.Free, 10)
	stats["memory"].(map[string]interface{})["used"] = strconv.FormatUint(vmStat.Used, 10)
	stats["memory"].(map[string]interface{})["used_percent"] = strconv.FormatFloat(vmStat.UsedPercent, 'f', 2, 64)

	// get disk serial number.... strange... not available from disk package at compile time
	// undefined: disk.GetDiskSerialNumber
	//serial := disk.GetDiskSerialNumber("/dev/sda")
	stats["disk"] = make(map[string]interface{}, 0)
	stats["disk"].(map[string]interface{})["total"] = strconv.FormatUint(diskStat.Total, 10)
	stats["disk"].(map[string]interface{})["free"] = strconv.FormatUint(diskStat.Used, 10)
	stats["disk"].(map[string]interface{})["used_bytes"] = strconv.FormatUint(diskStat.Used, 10)

	// since my machine has one CPU, I'll use the 0 index
	// if your machine has more than 1 CPU, use the correct index
	// to get the proper data

	// cpu infos.
	stats["cpu"] = make(map[string]interface{}, 0)
	stats["cpu"].(map[string]interface{})["index_number"] = strconv.FormatInt(int64(cpuStat[0].CPU), 10)
	stats["cpu"].(map[string]interface{})["vendor_id"] = cpuStat[0].VendorID
	stats["cpu"].(map[string]interface{})["family"] = cpuStat[0].Family
	stats["cpu"].(map[string]interface{})["number_of_cores"] = strconv.FormatInt(int64(cpuStat[0].Cores), 10)
	stats["cpu"].(map[string]interface{})["model_name"] = cpuStat[0].ModelName
	stats["cpu"].(map[string]interface{})["speed"] = strconv.FormatFloat(cpuStat[0].Mhz, 'f', 2, 64)
	stats["cpu"].(map[string]interface{})["utilizations"] = make([]map[string]interface{}, 0)
	for idx, cpupercent := range percentage {
		stats["cpu"].(map[string]interface{})["utilizations"] = append(stats["cpu"].(map[string]interface{})["utilizations"].([]map[string]interface{}), map[string]interface{}{"idx": strconv.Itoa(idx), "utilization": strconv.FormatFloat(cpupercent, 'f', 2, 64)})
	}

	stats["hostname"] = hostStat.Hostname
	stats["uptime"] = strconv.FormatUint(hostStat.Uptime, 10)
	stats["number_of_running_processes"] = strconv.FormatUint(hostStat.Procs, 10)

	// another way to get the operating system name
	// both darwin for Mac OSX, For Linux, can be ubuntu as platform
	// and linux for OS
	stats["os"] = hostStat.OS
	stats["platform"] = hostStat.Platform

	stats["network_interfaces"] = make([]map[string]interface{}, 0)

	// the unique hardware id for this machine
	for _, interf := range interfStat {
		network_interface := make(map[string]interface{}, 0)
		network_interface["mac"] = interf.HardwareAddr

		network_interface["flags"] = interf.Flags
		network_interface["addresses"] = make([]string, 0)
		for _, addr := range interf.Addrs {
			network_interface["addresses"] = append(network_interface["addresses"].([]string), addr.String())
		}

		stats["network_interfaces"] = append(stats["network_interfaces"].([]map[string]interface{}), network_interface)

	}

	// generate a json output.
	w.Header().Set("Content-Type", "application/json")
	setupResponse(&w, r)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(stats)

}

/**
 * Return the ca certificate public key.
 */
func getCaCertificateHanldler(w http.ResponseWriter, r *http.Request) {
	//add prefix and clean
	w.Header().Set("Content-Type", "application/text")
	setupResponse(&w, r)
	w.WriteHeader(http.StatusCreated)

	crt, err := ioutil.ReadFile(globule.creds + "/ca.crt")
	if err != nil {
		http.Error(w, "Client ca cert not found", http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, string(crt))
}

/**
 * Return the server SAN configuration file.
 */
func getSanConfigurationHandler(w http.ResponseWriter, r *http.Request) {
	//add prefix and clean
	w.Header().Set("Content-Type", "application/text")
	setupResponse(&w, r)

	crt, err := ioutil.ReadFile(globule.creds + "/san.conf")
	if err != nil {
		http.Error(w, "Client Subject Alernate Name configuration found!", http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, string(crt))
}

/**
 * Setup allow Cors policies.
 */
func setupResponse(w *http.ResponseWriter, req *http.Request) {
	var allowedOrigins string
	for i := 0; i < len(globule.AllowedOrigins); i++ {
		allowedOrigins += globule.AllowedOrigins[i]
		if i < len(globule.AllowedOrigins)-1 {
			allowedOrigins += ","
		}
	}

	var allowedMethods string
	for i := 0; i < len(globule.AllowedMethods); i++ {
		allowedMethods += globule.AllowedMethods[i]
		if i < len(globule.AllowedMethods)-1 {
			allowedMethods += ","
		}
	}

	var allowedHeaders string
	for i := 0; i < len(globule.AllowedHeaders); i++ {
		allowedHeaders += globule.AllowedHeaders[i]
		if i < len(globule.AllowedHeaders)-1 {
			allowedHeaders += ","
		}
	}

	(*w).Header().Set("Access-Control-Allow-Origin", allowedOrigins)
	(*w).Header().Set("Access-Control-Allow-Methods", allowedMethods)
	(*w).Header().Set("Access-Control-Allow-Headers", allowedHeaders)
}

/**
 * Sign ca certificate request and return a certificate.
 */
func signCaCertificateHandler(w http.ResponseWriter, r *http.Request) {

	//add prefix and clean
	w.Header().Set("Content-Type", "application/text")
	setupResponse(&w, r)

	w.WriteHeader(http.StatusCreated)
	// sign the certificate.
	csr_str := r.URL.Query().Get("csr") // the csr in base64
	csr, err := base64.StdEncoding.DecodeString(csr_str)

	if err != nil {
		http.Error(w, "Fail to decode csr base64 string", http.StatusBadRequest)
		return
	}

	// Now I will sign the certificate.
	crt, err := globule.signCertificate(string(csr))
	if err != nil {
		http.Error(w, "fail to sign certificate!", http.StatusBadRequest)
		return
	}

	// Return the result as text string.
	fmt.Fprint(w, crt)
}

/**
 * This code is use to upload a file into the tmp directory of the server
 * via http request.
 */
func FileUploadHandler(w http.ResponseWriter, r *http.Request) {

	setupResponse(&w, r)

	// I will
	err := r.ParseMultipartForm(200000) // grab the multipart form
	if err != nil {
		return
	}

	formdata := r.MultipartForm // ok, no problem so far, read the Form data

	//get the *fileheaders
	files := formdata.File["multiplefiles"] // grab the filenames

	// Get the path where to upload the file.
	path := r.FormValue("path")
	path = strings.ReplaceAll(path, "\\", "/")

	
	// If application is defined.
	token := r.Header.Get("token")
	application := r.Header.Get("application")

	// If the header dosent contain the required values i I will try to get it from the 
	// http query instead...
	if len(token) == 0 {
		// the token can be given by the url directly...
		token = r.URL.Query().Get("token")
	}

	if len(application) == 0 {
		// the token can be given by the url directly...
		application = r.URL.Query().Get("application")
	}

	user := ""
	hasAccess := false

	// TODO fix it and uncomment it...
	hasAccessDenied := false
	infos := []*rbacpb.ResourceInfos{}

	// Here I will validate applications...
	if len(application) != 0 {
		// Test if the requester has the permission to do the upload...
		// Here I will named the methode /file.FileService/FileUploadHandler
		// I will be threaded like a file service methode.
		hasAccess, err = globule.validateAction("/file.FileService/FileUploadHandler", application, rbacpb.SubjectType_APPLICATION, infos)
		if hasAccess && err == nil {
			hasAccess, hasAccessDenied, err = globule.validateAccess(application, rbacpb.SubjectType_APPLICATION, "write", path)
		}
	}

	// get the user id from the token...
	if len(token) != 0 {
		user, _, _, _, _, err = security.ValidateToken(token)
		if err != nil {
			user = ""
		}
	}

	if len(user) != 0 {

		if !hasAccess {
			hasAccess, err = globule.validateAction("/file.FileService/FileUploadHandler", user, rbacpb.SubjectType_ACCOUNT, infos)
		}
		fmt.Println("Validate accees for user with id: ", user)
		if hasAccess && err == nil {
			hasAccess, hasAccessDenied, err = globule.validateAccess(user, rbacpb.SubjectType_ACCOUNT, "write", path)
			if err != nil {
				log.Println("Fail to validate action with error ", err)
			}
		} else {
			log.Println("Fail to validate action with error ", err)
		}

	}

	// validate ressource access...
	if !hasAccess || hasAccessDenied || err != nil {
		http.Error(w, "unable to create the file for writing. Check your access privilege", http.StatusUnauthorized)
		return
	}

	for _, f := range files { // loop through the files one by one
		file, err := f.Open()
		if err != nil {
			return
		}

		defer file.Close()

		// Here I will set the ressource owner.
		if len(user) > 0 {
			globule.addResourceOwner(path+"/"+f.Filename, user, rbacpb.SubjectType_ACCOUNT)
		} else if len(application) > 0 {
			globule.addResourceOwner(path+"/"+f.Filename, application, rbacpb.SubjectType_APPLICATION)
		}
		
		// Create the file depending if the path is users, applications or something else...
		path_ := path + "/" + f.Filename
		if strings.HasPrefix(path, "/users") || strings.HasPrefix(path, "/applications") {
			path_ = strings.ReplaceAll(globule.data+"/files"+path_, "\\", "/")
		} else {
			path_ = strings.ReplaceAll(globule.webRoot+path_, "\\", "/")
		}

		out, err := os.Create(path_)
		if err != nil {
			return
		}

		defer out.Close()

		if err != nil {
			log.Println("fail to create dir ", path_, err)
			file.Close()
			http.Error(w, "Unable to create the file for writing. Check your write access privilege", http.StatusUnauthorized)
			return
		}
		_, err = io.Copy(out, file) // file not files[i] !
		if err != nil {
			log.Println("fail to copy file  ", path_, "with error", err)
			file.Close()
			return
		}
		file.Close()

		// Now from the file extension i will retreive it mime type.
		if strings.LastIndex(path_, ".") != -1 {
			fileExtension := path_[strings.LastIndex(path_, "."):]
			fileType := mime.TypeByExtension(fileExtension)
			path_ = strings.ReplaceAll(path_, "\\", "/")
			if len(fileType) > 0 {
				if strings.HasPrefix(fileType, "text/") {
					indexFile(path_, fileType)
				} else if strings.HasPrefix(fileType, "video/") {
					// Here I will call convert video...
					go func() {
						convertVideo()
					}()
				}
			}
		}
	}
}

func visit(files *[]string) filepath.WalkFunc {

	return func(path string, info os.FileInfo, err error) error {

		path = strings.ReplaceAll(path, "\\", "/")
		if err != nil {
			return nil
		}

		if strings.HasPrefix(path, config_.GetDataDir()+"/files/users/") && !strings.Contains(path, ".hidden") {

			// Here I will set the owner write to file inside it directory...
			//userId :=  [0:strings.Index(path[len(config.GetDataDir() + "/files/users/"):], "/")]
			path_ := path[len(config_.GetDataDir()+"/files/users/"):]
			index := strings.Index(path_, "/")
			userId := ""
			if index > 0 {
				userId = path_[0:index]
			} else {
				userId = path_
			}

			if len(userId) > 0 {
				globule.addResourceOwner("/users/"+path_, userId, rbacpb.SubjectType_ACCOUNT)
			}
		} else if strings.HasPrefix(path, config_.GetDataDir()+"/files/applications/") && !strings.Contains(path, ".hidden") {
			path_ := path[len(config_.GetDataDir()+"/files/applications/"):]
			index := strings.Index(path_, "/")
			id := ""
			if index > 0 {
				id = path_[0:index]
			} else {
				id = path_
			}

			if len(id) > 0 {
				globule.addResourceOwner("/applications/"+path_, id, rbacpb.SubjectType_APPLICATION)
			}
		}

		if err != nil {
			return nil
		}
		mimeType := ""
		if strings.Contains(info.Name(), ".") {
			fileExtension := info.Name()[strings.LastIndex(info.Name(), "."):]
			mimeType = mime.TypeByExtension(fileExtension)
		} else {
			f_, err := os.Open(path)
			if err != nil {
				f_.Close()
				return nil
			}
			mimeType, _ = Utility.GetFileContentType(f_)
			f_.Close()
		}

		if strings.HasPrefix(mimeType, "video/") && !strings.HasSuffix(info.Name(), ".mp4") {
			*files = append(*files, path)
		} else if strings.HasPrefix(mimeType, "video/") && strings.HasSuffix(info.Name(), ".mp4") {
			createVideoPreview(path, 20, 128)
		}

		return nil
	}
}

// That function resolve import path.
func resolveImportPath(path string, importPath string) (string, error) {

	// firt of all i will keep only the path part of the import...
	startIndex := strings.Index(importPath, `'@`) + 1
	endIndex := strings.LastIndex(importPath, `'`)
	importPath_ := importPath[startIndex:endIndex]

	filepath.Walk(globule.webRoot+path[0:strings.Index(path, "/")],
		func(path string, info os.FileInfo, err error) error {
			path = strings.ReplaceAll(path, "\\", "/")
			if err != nil {
				return err
			}

			if strings.HasSuffix(path, importPath_) {
				importPath_ = path
				return io.EOF
			}

			return nil
		})

	importPath_ = strings.Replace(importPath_, strings.Replace(globule.webRoot, "\\", "/", -1), "", -1)

	// Now i will make the path relative.
	importPath__ := strings.Split(importPath_, "/")
	path__ := strings.Split(path, "/")

	var index int
	for ; importPath__[index] == path__[index]; index++ {
	}

	importPath_ = ""

	// move up part..
	for i := index; i < len(path__)-1; i++ {
		importPath_ += "../"
	}

	// go down to the file.
	for i := index; i < len(importPath__); i++ {
		importPath_ += importPath__[i]
		if i < len(importPath__)-1 {
			importPath_ += "/"
		}
	}

	// remove the
	importPath_ = strings.Replace(importPath_, globule.webRoot, "", 1)

	// remove the root path part and the leading / caracter.
	return importPath_, nil
}

// Custom file server implementation.
func ServeFileHandler(w http.ResponseWriter, r *http.Request) {

	setupResponse(&w, r)
	dir := globule.webRoot

	// If a directory with the same name as the host in the request exist
	// it will be taken as root. Permission will be manage by the resource
	// manager and not simply the name of the directory. If you want to protect
	// a given you need to set permission on it.
	if Utility.Exists(dir + "/" + r.Host) {
		dir += "/" + r.Host
	}

	//add prefix and clean
	rqst_path := path.Clean(r.URL.Path)

	if rqst_path == "/null" {
		http.Error(w, "No file path was given in the file url path!", http.StatusBadRequest)
	}

	// Now I will test if a token is given in the header and manage it file access.
	application := r.Header.Get("application")
	token := r.Header.Get("token")

	if len(token) == 0 {
		// the token can be given by the url directly...
		token = r.URL.Query().Get("token")
	}

	if len(application) == 0 {
		// the token can be given by the url directly...
		application = r.URL.Query().Get("application")
	}

	// If the path is '/' it mean's no application name was given and we are
	// at the root.
	if rqst_path == "/" {
		// if a default application is define in the globule i will use it.
		if len(globule.IndexApplication) > 0 {
			rqst_path += globule.IndexApplication
			application = globule.IndexApplication
		}

	} else if strings.Count(rqst_path, "/") == 1 {
		if strings.HasSuffix(rqst_path, ".js") ||
			strings.HasSuffix(rqst_path, ".json") ||
			strings.HasSuffix(rqst_path, ".css") ||
			strings.HasSuffix(rqst_path, ".htm") ||
			strings.HasSuffix(rqst_path, ".html") {
			rqst_path = "/" + globule.IndexApplication + rqst_path
		}
	}

	hasAccess := true
	if strings.HasPrefix(rqst_path, "/users/") || strings.HasPrefix(rqst_path, "/applications/") || strings.HasPrefix(rqst_path, "/templates/") || strings.HasPrefix(rqst_path, "/projects/") {
		dir = globule.data + "/files"
		if !strings.Contains(rqst_path, "/.hidden/") {
			hasAccess = false
		}

	}

	//path to file
	name := path.Join(dir, rqst_path)

	// this is the ca certificate use to sign client certificate.
	if rqst_path == "/ca.crt" {
		name = globule.creds + rqst_path
	}

	hasAccessDenied := false
	var err error
	var userId string

	//hasAccess = false
	// Here I will validate applications...
	if len(application) != 0 && !hasAccess {
		hasAccess, hasAccessDenied, err = globule.validateAccess(application, rbacpb.SubjectType_APPLICATION, "read", rqst_path)
	}

	if len(token) != 0 && !hasAccess {
		userId, _, _, _, _, err = security.ValidateToken(token)
		if err == nil {
			log.Println("validate access for ", userId, rqst_path)
			hasAccess, hasAccessDenied, err = globule.validateAccess(userId, rbacpb.SubjectType_ACCOUNT, "read", rqst_path)
		}
	}

	// validate ressource access...
	if !hasAccess || hasAccessDenied || err != nil {
		http.Error(w, "unable to read the file "+rqst_path+" Check your access privilege", http.StatusUnauthorized)
		return
	}

	// if the file dosent exist... I will try to get it from the index application...
	if !Utility.Exists(name) && len(globule.IndexApplication) > 0 {
		name = path.Join(dir, globule.IndexApplication+"/"+rqst_path)
	}

	var code string
	// If the file is a javascript file...
	hasChange := false

	//check if file exists
	f, err := os.Open(name)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File "+rqst_path+" not found!", http.StatusNoContent)
			return
		}
	}

	defer f.Close()

	if strings.HasSuffix(name, ".js") {
		w.Header().Add("Content-Type", "application/javascript")
		if err == nil {
			//hasChange = true
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "import") {
					if strings.Contains(line, `'@`) {
						path_, err := resolveImportPath(rqst_path, line)
						if err == nil {
							line = line[0:strings.Index(line, `'@`)] + `'` + path_ + `'`
							hasChange = true
						}
					}
				}
				code += line + "\n"
			}
		}

	} else if strings.HasSuffix(name, ".css") {
		w.Header().Add("Content-Type", "text/css")

	} else if strings.HasSuffix(name, ".html") || strings.HasSuffix(name, ".htm") {
		w.Header().Add("Content-Type", "text/html")
	}

	// if the file has change...
	if !hasChange {
		//log.Println("server file ", name)
		http.ServeFile(w, r, name)
	} else {
		//log.Println("server content ", name)
		http.ServeContent(w, r, name, time.Now(), strings.NewReader(code))
	}
}
