package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/davecourtois/Utility"
	"github.com/globulario/Globular/Interceptors"
	"github.com/globulario/services/golang/rbac/rbacpb"
)

/**
 * Return the service configuration
 */
func getConfigHanldler(w http.ResponseWriter, r *http.Request) {

	//add prefix and clean
	config := globule.getConfig()
	w.Header().Set("Content-Type", "application/json")
	setupResponse(&w, r)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(config)
}

/**
 * Return the ca certificate public key.
 */
func getCaCertificateHanldler(w http.ResponseWriter, r *http.Request) {
	//add prefix and clean
	w.Header().Set("Content-Type", "application/text")
	setupResponse(&w, r)
	w.WriteHeader(http.StatusCreated)

	crt, err := ioutil.ReadFile(globule.creds + "/" + "ca.crt")
	if err != nil {
		http.Error(w, "Client ca cert not found!", http.StatusBadRequest)
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

	crt, err := ioutil.ReadFile(globule.creds + "/" + "san.conf")
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
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, domain, application, token")
}

/**
 * Sign ca certificate request and return a certificate.
 */
func signCaCertificateHandler(w http.ResponseWriter, r *http.Request) {
	//add prefix and clean
	w.Header().Set("Content-Type", "application/text")
	setupResponse(&w, r)

	w.WriteHeader(http.StatusCreated)
	log.Println("sign Ca Certificate Handler was call")
	// sign the certificate.
	csr_str := r.URL.Query().Get("csr") // the csr in base64
	csr, err := base64.StdEncoding.DecodeString(csr_str)

	if err != nil {
		log.Println(err)
		http.Error(w, "Fail to decode csr base64 string", http.StatusBadRequest)
		return
	}

	// Now I will sign the certificate.
	crt, err := globule.signCertificate(string(csr))
	if err != nil {
		log.Println(err)
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

	log.Println("upload file request was called!")
	setupResponse(&w, r)

	log.Println("upload was called... ")
	// I will
	err := r.ParseMultipartForm(200000) // grab the multipart form
	if err != nil {
		log.Println(w, err)
		return
	}

	formdata := r.MultipartForm // ok, no problem so far, read the Form data

	//get the *fileheaders
	files := formdata.File["multiplefiles"] // grab the filenames
	var path string                         // grab the filenames

	// Get the path where to upload the file.
	path = r.FormValue("path")
	log.Println("----------> path: ", path)
	log.Println("---------->", 143)
	// If application is defined.
	token := r.Header.Get("token")
	application := r.Header.Get("application")
	hasAccess := true //TODO set it back to false when the
	hasAccessDenied := false
	user := ""
	infos := []*rbacpb.ResourceInfos{}

	// Here I will validate applications...
	if len(application) != 0 {
		// Test if the requester has the permission to do the upload...
		// Here I will named the methode /file.FileService/FileUploadHandler
		// I will be threaded like a file service methode.
		hasAccess, err = globule.validateAction("/file.FileService/FileUploadHandler", application, rbacpb.SubjectType_APPLICATION, infos)
		if err != nil || !hasAccess {
			// http.Error(w, "Unable to create the file for writing. Check your access privilege", http.StatusUnauthorized)
			log.Println("---------->", 160)
			//return
		}

		// validate ressource access...
		hasAccess, hasAccessDenied, err = globule.validateAccess(application, rbacpb.SubjectType_APPLICATION, "write", path)
		if !hasAccess || hasAccessDenied || err != nil {
			// http.Error(w, "Unable to create the file for writing. Check your access privilege", http.StatusUnauthorized)
			//return
			log.Println("---------->", 169)
		}
	}

	// domain := r.Header.Get("domain")
	if len(token) != 0 && !hasAccess {
		id, username, _, expiresAt, err := Interceptors.ValidateToken(token)
		user = username
		if err != nil || time.Now().Before(time.Unix(expiresAt, 0)) {
			// http.Error(w, "Unable to create the file for writing. Check your access privilege", http.StatusUnauthorized)
			//return
			log.Println("---------->", 180)
		} else {
			hasAccess, err = globule.validateAction("/file.FileService/FileUploadHandler", id, rbacpb.SubjectType_ACCOUNT, infos)
			if err != nil || !hasAccess {
				//http.Error(w, "Unable to create the file for writing. Check your access privilege", http.StatusUnauthorized)
				//return
				log.Println("---------->", 186)
			}

			hasAccess, hasAccessDenied, err = globule.validateAccess(id, rbacpb.SubjectType_ACCOUNT, "write", path)
			if !hasAccess || hasAccessDenied || err != nil {
				//http.Error(w, "Unable to create the file for writing. Check your access privilege", http.StatusUnauthorized)
				//return
				log.Println("---------->", 192)
			}
		}
	}

	// Here the path dosent exist.
	if !Utility.Exists(globule.webRoot + path) {
		// TODO validate ressource access here
		Utility.CreateDirIfNotExist(globule.webRoot + path)
	}

	for i, _ := range files { // loop through the files one by one
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			log.Println(w, err)
			return
		}

		// Here I will set the ressource owner.
		if len(user) > 0 {
			globule.addResourceOwner(path+"/"+files[i].Filename, user, rbacpb.SubjectType_ACCOUNT)
		} else if len(application) > 0 {
			globule.addResourceOwner(path+"/"+files[i].Filename, application, rbacpb.SubjectType_APPLICATION)
		}

		// Create the file depending if the path is users, applications or something else...
		path_ := path + "/" + files[i].Filename
		if strings.HasPrefix(path, "/users") || strings.HasPrefix(path, "/applications") {
			path_ = strings.ReplaceAll(globule.data+"/files"+path_, "\\", "/")
		} else {
			path_ = strings.ReplaceAll(globule.webRoot+path_, "\\", "/")
		}

		out, err := os.Create(path_)

		defer out.Close()
		if err != nil {
			http.Error(w, "Unable to create the file for writing. Check your write access privilege", http.StatusUnauthorized)
			return
		}
		_, err = io.Copy(out, file) // file not files[i] !
		if err != nil {
			log.Println(w, err)
			return
		}

		// Now from the file extension i will retreive it mime type.
		fileExtension := path_[strings.LastIndex(path_, "."):]
		fileType := mime.TypeByExtension(fileExtension)
		path_ = strings.ReplaceAll(path_, "\\", "/")
		if len(fileType) > 0 {
			if strings.HasPrefix(fileType, "text/") {
				indexFile(path_, fileType)
			} else if strings.HasPrefix(fileType, "video/") {
				if strings.HasSuffix(path_, ".mp4") {
					err := createVideoPreview(path_, 20, 128)
					if err != nil {
						log.Println(err)
					}

				} else {
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
			log.Println(err)
			return nil
		}

		if err != nil {
			log.Println(err)
			return nil
		}
		mimeType := ""
		if strings.Index(info.Name(), ".") != -1 {
			fileExtension := info.Name()[strings.LastIndex(info.Name(), "."):]
			mimeType = mime.TypeByExtension(fileExtension)
		} else {
			f_, err := os.Open(path)
			defer f_.Close()
			if err != nil {
				return nil
			}
			mimeType, _ = Utility.GetFileContentType(f_)
		}

		if strings.HasPrefix(mimeType, "video/") && !strings.HasSuffix(info.Name(), ".mp4") {
			*files = append(*files, path)
		} else if strings.HasPrefix(mimeType, "video/") && strings.HasSuffix(info.Name(), ".mp4") {
			createVideoPreview(path, 20, 128)
		}
		return nil
	}
}

func convertVideo() {

	pids, err := Utility.GetProcessIdsByName("ffmpeg")
	if err == nil {
		if len(pids) > 0 {
			return // already running...
		}
	}
	var files []string

	err = filepath.Walk(globule.data+"/files", visit(&files))
	if err != nil {
		log.Println(err)
		return
	}
	for _, file := range files {
		log.Println("-------> convert video: ", file)
		file = strings.ReplaceAll(file, "\\", "/")
		createVideoStream(file)
	}

}

// Set file indexation to be able to search text file on the server.
func indexFile(path string, fileType string) error {
	log.Println("---------> index file ", path, fileType)
	return nil
}

/**
 * Convert all kind of video to mp4 so all browser will be able to read it.
 */
func createVideoStream(path string) error {

	path_ := path[0:strings.LastIndex(path, "/")]
	name_ := path[strings.LastIndex(path, "/"):strings.LastIndex(path, ".")]
	output := path_ + "/" + name_ + ".mp4"

	defer os.RemoveAll(path)

	var cmd *exec.Cmd

	// ffmpeg -i input.mkv -c:v libx264 -c:a aac output.mp4
	cmd = exec.Command("ffmpeg", "-i", path, "-c:v", "libx264", "-c:a", "aac", output)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	// Create a video preview
	log.Println("create preview ", path)
	err = createVideoPreview(output, 20, 128)
	if err != nil {
		log.Println(err)
	}

	return nil
}

// Here I will create video
func createVideoPreview(path string, nb int, height int) error {

	duration := getVideoDuration(path)
	if duration == 0 {
		return errors.New("the video lenght is 0 sec")
	}

	path_ := path[0:strings.LastIndex(path, "/")]
	name_ := path[strings.LastIndex(path, "/")+1 : strings.LastIndex(path, ".")]
	output := path_ + "/.hidden/" + name_ + "/__preview__"
	if Utility.Exists(output) {
		return nil
	}

	Utility.CreateDirIfNotExist(output)
	log.Println("Create preview in ", output)

	// ffmpeg -i bob_ross_img-0-Animated.mp4 -ss 15 -t 16 -f image2 preview_%05d.jpg
	start := .1 * duration
	laps := 120 // 1 minutes

	cmd := exec.Command("ffmpeg", "-i", path, "-ss", Utility.ToString(start), "-t", Utility.ToString(laps), "-vf", "scale="+Utility.ToString(height)+":-1,fps=.250", "preview_%05d.jpg")
	cmd.Dir = output // the output directory...

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return err
	}

	evtHub, err := globule.getEventHub()
	if err == nil {
		path_ := strings.ReplaceAll(path, globule.data+"/files", "")
		path_ = path_[0:strings.LastIndex(path_, "/")]

		evtHub.Publish("reload_dir_event", []byte(path_))
	}

	return nil
}

func getVideoDuration(path string) float64 {

	// ffprobe -v quiet -print_format compact=print_section=0:nokey=1:escape=csv -show_entries format=duration bob_ross_img-0-Animated.mp4
	cmd := exec.Command("ffprobe", `-v`, `quiet`, `-print_format`, `compact=print_section=0:nokey=1:escape=csv`, `-show_entries`, `format=duration`, path)

	cmd.Dir = os.TempDir()

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		return 0.0
	}

	duration, _ := strconv.ParseFloat(strings.TrimSpace(string(out.Bytes())), 64)

	return duration
}

// Custom file server implementation.
func ServeFileHandler(w http.ResponseWriter, r *http.Request) {

	setupResponse(&w, r)

	//if empty, set current directory

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

	if strings.HasPrefix(rqst_path, "/users/") || strings.HasPrefix(rqst_path, "/applications/") {
		dir = globule.data + "/files"
	}

	//path to file
	name := path.Join(dir, rqst_path)
	log.Println("Try to access file...", name)

	// this is the ca certificate use to sign client certificate.
	if rqst_path == "/ca.crt" {
		name = globule.creds + rqst_path
	}

	token := r.Header.Get("token")

	// domain := r.Header.Get("domain")
	hasAccess := true // TODO set it back to false when permission system will be completed.
	hasAccessDenied := false
	var err error
	infos := []*rbacpb.ResourceInfos{}

	if len(application) != 0 {
		// Test if the requester has the permission to do the upload...
		// Here I will named the methode /file.FileService/FileUploadHandler
		// I will be threaded like a file service methode.
		hasAccess, err = globule.validateAction("/file.FileService/ServeFileHandler", application, rbacpb.SubjectType_APPLICATION, infos)
		if err != nil || !hasAccess {
			//http.Error(w, "Unable to create the file for writing. Check your access privilege", http.StatusUnauthorized)
			//return
		}

		// validate ressource access...
		hasAccess, hasAccessDenied, err = globule.validateAccess(application, rbacpb.SubjectType_APPLICATION, "read", rqst_path)
		if !hasAccess || hasAccessDenied || err != nil {
			//http.Error(w, "Unable to create the file for writing. Check your access privilege", http.StatusUnauthorized)
			//return
		}
	}

	// domain := r.Header.Get("domain")
	if len(token) != 0 && !hasAccess {
		id /*username*/, _, _, expiresAt, err := Interceptors.ValidateToken(token)
		if err != nil || time.Now().Before(time.Unix(expiresAt, 0)) {
			//http.Error(w, "Unable to create the file for writing. Check your access privilege", http.StatusUnauthorized)
			//return
		} else {
			hasAccess, err = globule.validateAction("/file.FileService/ServeFileHandler", id, rbacpb.SubjectType_ACCOUNT, infos)
			if err != nil || !hasAccess {
				//http.Error(w, "Unable to create the file for writing. Check your access privilege", http.StatusUnauthorized)
				//return
			}

			hasAccess, hasAccessDenied, err = globule.validateAccess(id, rbacpb.SubjectType_ACCOUNT, "read", rqst_path)
			if !hasAccess || hasAccessDenied || err != nil {
				//http.Error(w, "Unable to create the file for writing. Check your access privilege", http.StatusUnauthorized)
				//return
			}
		}
	}

	// if the file dosent exist... I will try to get it from the index application...
	if !Utility.Exists(name) && len(globule.IndexApplication) > 0 {
		name = path.Join(dir, globule.IndexApplication+"/"+rqst_path)
	}

	//check if file exists
	f, err := os.Open(name)
	if err != nil {
		if os.IsNotExist(err) {

			http.Error(w, "File "+rqst_path+" not found!", http.StatusNoContent)
			return

		}
	}

	defer f.Close()

	// If the file is a javascript file...
	var code string
	hasChange := false
	if strings.HasSuffix(name, ".js") {
		w.Header().Add("Content-Type", "application/javascript")
		if err == nil {
			//hasChange = true
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "import") {
					if strings.Index(line, `'@`) > -1 {
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
		http.ServeFile(w, r, name)
	} else {
		http.ServeContent(w, r, name, time.Now(), strings.NewReader(code))
	}
}
