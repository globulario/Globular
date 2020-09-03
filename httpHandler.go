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

	"github.com/davecourtois/Globular/Interceptors"
	"github.com/davecourtois/Utility"
)

/**
 * Return the service configuration
 */
func getConfigHanldler(w http.ResponseWriter, r *http.Request) {
	//add prefix and clean
	config := globule.getConfig()
	w.Header().Set("Content-Type", "application/json")
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
	domain := r.Header.Get("domain")
	hasPermission := false
	user := ""

	log.Println("-------> validate path '", path, "' for ", application)
	if len(application) != 0 {

		err := Interceptors.ValidateApplicationRessourceAccess(domain, application, "/file.FileService/FileUploadHandler", path, 2)
		if err == nil {
			hasPermission = true
		}
	}

	if len(token) != 0 && !hasPermission {
		// Test if the requester has the permission to do the upload...
		// Here I will named the methode /file.FileService/FileUploadHandler
		// I will be threaded like a file service methode.
		err := Interceptors.ValidateUserRessourceAccess(domain, token, "/file.FileService/FileUploadHandler", path, 2)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		user, _, _, _ = Interceptors.ValidateToken(token)
		hasPermission = true
	}

	if !hasPermission {
		http.Error(w, "Unable to create the file for writing. Check your write access privilege", http.StatusUnauthorized)
		return
	}

	if !Utility.Exists(globule.webRoot + path) {
		hasPermission := false

		if len(token) != 0 && !hasPermission {
			err := Interceptors.ValidateUserRessourceAccess(domain, token, "/file.FileService/CreateDir", path, 2)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			user, _, _, _ = Interceptors.ValidateToken(token)
			hasPermission = true
		}

		if !hasPermission {
			http.Error(w, "Unable to create the file for writing. Check your write access privilege", http.StatusUnauthorized)
			return
		}

		if len(user) > 0 {
			globule.setRessourceOwner(user, path)
		}
		Utility.CreateDirIfNotExist(globule.webRoot + path)
	}

	for i, _ := range files { // loop through the files one by one
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			log.Println(w, err)
			return
		}
		// set the file owner if the length of the user if greather than 0
		if len(user) > 0 {

			log.Println("--------> path ", path)
			log.Println("--------> file name ", files[i].Filename)

			globule.setRessourceOwner(user, path+"/"+files[i].Filename)
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

	//if empty, set current directory
	dir := string(root)
	if dir == "" {
		dir = "."
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
	domain := r.Header.Get("domain")
	hasPermission := false

	if len(application) != 0 {
		err := Interceptors.ValidateApplicationRessourceAccess(domain, application, "/file.FileService/ServeFileHandler", name, 4)
		if err != nil && len(token) == 0 {
			log.Println("Fail to download the file with error ", err.Error())
			return
		}
		hasPermission = err == nil
	}

	if len(token) != 0 && !hasPermission {
		// Test if the requester has the permission to do the upload...
		// Here I will named the methode /file.FileService/FileUploadHandler
		// I will be threaded like a file service methode.
		err := Interceptors.ValidateUserRessourceAccess(domain, token, "/file.FileService/ServeFileHandler", name, 4)
		if err != nil {
			log.Println("Fail to dowload the file with error ", err.Error())
			return
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
