package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(config)
}

/**
 * Return the ca certificate public key.
 */
func getCaCertificateHanldler(w http.ResponseWriter, r *http.Request) {
	//add prefix and clean
	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(http.StatusCreated)

	crt, err := ioutil.ReadFile(globule.creds + string(os.PathSeparator) + "ca.crt")
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
	w.WriteHeader(http.StatusCreated)

	crt, err := ioutil.ReadFile(globule.creds + string(os.PathSeparator) + "san.conf")
	if err != nil {
		http.Error(w, "Client Subject Alernate Name configuration found!", http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, string(crt))
}

/**
 * Sign ca certificate request and return a certificate.
 */
func signCaCertificateHandler(w http.ResponseWriter, r *http.Request) {
	//add prefix and clean
	w.Header().Set("Content-Type", "application/text")
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

	// If application is defined.
	token := r.Header.Get("token")
	application := r.Header.Get("application")
	hasAccess := false
	hasAccessDenied := false
	user := ""
	infos := []*rbacpb.ResourceInfos{}

	if len(application) != 0 {
		// Test if the requester has the permission to do the upload...
		// Here I will named the methode /file.FileService/FileUploadHandler
		// I will be threaded like a file service methode.
		hasAccess, err = globule.validateAction("/file.FileService/FileUploadHandler", application, rbacpb.SubjectType_APPLICATION, infos)
		if err != nil || !hasAccess {
			http.Error(w, "Unable to create the file for writing. Check your access privilege", http.StatusUnauthorized)
			return
		}

		// validate ressource access...
		hasAccess, hasAccessDenied, err = globule.validateAccess(application, rbacpb.SubjectType_APPLICATION, "write", path)
		if !hasAccess || hasAccessDenied || err != nil {
			http.Error(w, "Unable to create the file for writing. Check your access privilege", http.StatusUnauthorized)
			return
		}
	}

	// domain := r.Header.Get("domain")
	if len(token) != 0 && !hasAccess {
		username, _, expiresAt, err := Interceptors.ValidateToken(token)
		user = username
		if err != nil || time.Now().Before(time.Unix(expiresAt, 0)) {
			http.Error(w, "Unable to create the file for writing. Check your access privilege", http.StatusUnauthorized)
			return
		} else {
			hasAccess, err = globule.validateAction("/file.FileService/FileUploadHandler", username, rbacpb.SubjectType_ACCOUNT, infos)
			if err != nil || !hasAccess {
				http.Error(w, "Unable to create the file for writing. Check your access privilege", http.StatusUnauthorized)
				return
			}

			hasAccess, hasAccessDenied, err = globule.validateAccess(username, rbacpb.SubjectType_ACCOUNT, "write", path)
			if !hasAccess || hasAccessDenied || err != nil {
				http.Error(w, "Unable to create the file for writing. Check your access privilege", http.StatusUnauthorized)
				return
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

		// Create the file.
		out, err := os.Create(globule.webRoot + path + "/" + files[i].Filename)

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
	}
}

// Custom file server implementation.
func ServeFileHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")

	//if empty, set current directory
	dir := string(webRoot)
	if dir == "" {
		dir = "."
	}

	// If a directory with the same name as the host in the request exist
	// it will be taken as root. Permission will be manage by the resource
	// manager and not simply the name of the directory. If you want to protect
	// a given you need to set permission on it.
	if Utility.Exists(dir + "/" + r.Host) {
		dir += "/" + r.Host
	}

	//add prefix and clean
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
		r.URL.Path = upath
	}

	upath = path.Clean(upath)

	//path to file
	name := path.Join(dir, upath)

	// this is the ca certificate use to sign client certificate.
	if upath == "/ca.crt" {
		name = globule.creds + upath
	}

	// Now I will test if a token is given in the header and manage it file access.
	application := r.Header.Get("application")
	token := r.Header.Get("token")

	// domain := r.Header.Get("domain")
	hasAccess := false
	hasAccessDenied := false
	var err error
	infos := []*rbacpb.ResourceInfos{}

	if len(application) != 0 {
		// Test if the requester has the permission to do the upload...
		// Here I will named the methode /file.FileService/FileUploadHandler
		// I will be threaded like a file service methode.
		hasAccess, err = globule.validateAction("/file.FileService/ServeFileHandler", application, rbacpb.SubjectType_APPLICATION, infos)
		if err != nil || !hasAccess {
			http.Error(w, "Unable to create the file for writing. Check your access privilege", http.StatusUnauthorized)
			return
		}

		// validate ressource access...
		hasAccess, hasAccessDenied, err = globule.validateAccess(application, rbacpb.SubjectType_APPLICATION, "read", upath)
		if !hasAccess || hasAccessDenied || err != nil {
			http.Error(w, "Unable to create the file for writing. Check your access privilege", http.StatusUnauthorized)
			return
		}
	}

	// domain := r.Header.Get("domain")
	if len(token) != 0 && !hasAccess {
		username, _, expiresAt, err := Interceptors.ValidateToken(token)
		if err != nil || time.Now().Before(time.Unix(expiresAt, 0)) {
			http.Error(w, "Unable to create the file for writing. Check your access privilege", http.StatusUnauthorized)
			return
		} else {
			hasAccess, err = globule.validateAction("/file.FileService/ServeFileHandler", username, rbacpb.SubjectType_ACCOUNT, infos)
			if err != nil || !hasAccess {
				http.Error(w, "Unable to create the file for writing. Check your access privilege", http.StatusUnauthorized)
				return
			}

			hasAccess, hasAccessDenied, err = globule.validateAccess(username, rbacpb.SubjectType_ACCOUNT, "read", upath)
			if !hasAccess || hasAccessDenied || err != nil {
				http.Error(w, "Unable to create the file for writing. Check your access privilege", http.StatusUnauthorized)
				return
			}
		}
	}
	//check if file exists
	f, err := os.Open(name)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "File "+upath+" not found!", http.StatusBadRequest)
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
						path_, err := resolveImportPath(upath, line)
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
