package main

import (
	"bytes"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/globulario/services/golang/dns/dns_client"
	"github.com/globulario/services/golang/event/event_client"
	"github.com/globulario/services/golang/rbac/rbac_client"
	"github.com/globulario/services/golang/rbac/rbacpb"

	// Interceptor for authentication, event, log...

	// Client services.
	"github.com/davecourtois/Utility"
	"github.com/go-acme/lego/certcrypto"
	"github.com/go-acme/lego/challenge/http01"
	"github.com/go-acme/lego/lego"
	"github.com/go-acme/lego/registration"
)

// Global variable.
var (
	globule *Globule
	configPath = "/etc/globular/config/config.json"
)



/**
 * The web server.
 */
type Globule struct {
	// The share part of the service.
	Name string // The service name

	// Globualr specifics ports.

	// can be https or http.
	Protocol  string
	PortHttp  int // The port of the http file server.
	PortHttps int // The secure port

	Domain           string        // The principale domain
	AlternateDomains []interface{} // Alternate domain for multiple domains
	IndexApplication string        // If defined It will be use as the entry point where not application path was given in the url.

	// Certificate generation variables.
	CertExpirationDelay int
	CertPassword        string
	Country             string // tow letter.
	State               string // Full state name
	City                string
	Organization        string

	// https certificate info.
	Certificate                string
	CertificateAuthorityBundle string
	CertURL                    string
	CertStableURL              string

	// Keep the version number.
	Version      string
	Build        int64
	Platform     string

	// Admin informations.
	AdminEmail   string
	RootPassword string

	// There's are Directory
	// The user's directory
	UsersDirectory string

	// The application directory
	ApplicationDirectory string

	// Service discoveries.
	Discoveries []string // Contain the list of discovery service use to keep globular up to date.

	// Update delay in second...
	WatchUpdateDelay int

	// DNS stuff.
	DNS              []interface{} // Domain name server use to located the server.
	DnsUpdateIpInfos []interface{} // The internet provader SetA info to keep ip up to date.

	// Directories.
	path         string // The path of the exec...
	webRoot      string // The root of the http file server.
	data         string // the data directory
	creds        string // tls certificates
	config       string // configuration directory
	users        string // the users files directory
	applications string // The applications

	// ACME protocol registration
	registration *registration.Resource

	// exit channel.
	exit  chan struct{}
	exit_ bool

	// The http server
	http_server  *http.Server
	https_server *http.Server
}

/**
 * Globule constructor.
 */
func NewGlobule() *Globule {

	// Here I will initialyse configuration.
	g := new(Globule)

	g.Version = "1.0.0" // Automate version...
	g.Build = 0
	g.Platform = runtime.GOOS + ":" + runtime.GOARCH
	g.IndexApplication = "globular_installer" // I will use the installer as defaut.

	g.PortHttp = 80   // The default http port
	g.PortHttps = 443 // The default https port number
	execPath := Utility.GetExecName(os.Args[0])
	g.Name = strings.Replace(execPath, ".exe", "", -1)

	// Set the default checksum...
	g.Protocol = "http"
	g.Domain = "localhost"

	// set default values.
	// g.SessionTimeout = 15 * 60 * 1000 // miliseconds.
	g.CertExpirationDelay = 365
	g.CertPassword = "1111"
	g.AdminEmail = "root@globular.app"
	g.RootPassword = "adminadmin"

	// keep up to date by default.
	g.WatchUpdateDelay = 30 // seconds...

	// Keep in global var to by http handlers.
	globule = g

	// Set the list of http handler.

	// The configuration handler.
	http.HandleFunc("/config", getConfigHanldler)

	// The checksum handler.
	http.HandleFunc("/checksum", getChecksumHanldler)

	// Handle the get ca certificate function
	http.HandleFunc("/get_ca_certificate", getCaCertificateHanldler)

	// Return the san server configuration.
	http.HandleFunc("/get_san_conf", getSanConfigurationHandler)

	// Handle the signing certificate function.
	http.HandleFunc("/sign_ca_certificate", signCaCertificateHandler)

	// Start listen for http request.
	http.HandleFunc("/", ServeFileHandler)

	// The file upload handler.
	http.HandleFunc("/uploads", FileUploadHandler)

	g.path, _ = filepath.Abs(filepath.Dir(os.Args[0]))

	if Utility.Exists(g.path+"/bin/grpcwebproxy") || Utility.Exists(g.path+"/bin/grpcwebproxy.exe") {
		// TODO test restart with initDirectories
		g.initDirectories()
	}

	return g
}

func (globule *Globule) getConfig() map[string]interface{} {
	config := make(map[string]interface{})

	// TODO implement it.

	return config
}

/**
 * Save the configuration
 */
func (globule *Globule) saveConfig() error {
	jsonStr, err := Utility.ToJson(globule)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(configPath, []byte(jsonStr), 0644)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Return the admin email.
 */
func (globule *Globule) GetEmail() string {
	return globule.AdminEmail
}

/**
 * Use the time of registration... Nil other wise.
 */
func (globule *Globule) GetRegistration() *registration.Resource {
	return globule.registration
}

/**
 * I will reuse the client public key here as key instead of generate another key
 * and manage it...
 */
func (globule *Globule) GetPrivateKey() crypto.PrivateKey {
	keyPem, err := ioutil.ReadFile(globule.creds + "/client.pem")
	if err != nil {
		return nil
	}

	keyBlock, _ := pem.Decode(keyPem)
	privateKey, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil
	}
	return privateKey
}

/**
 * That function work correctly, but the DNS fail time to time to give the
 * IP address that result in a fail request... The DNS must be fix!
 */
func (globule *Globule) obtainCertificateForCsr() error {

	config := lego.NewConfig(globule)
	config.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config)
	if err != nil {
		return err
	}

	err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", strconv.Itoa(globule.PortHttp)))
	if err != nil {
		log.Fatal(err)
	}

	reg, err := client.Registration.Register(registration.RegisterOptions{TermsOfServiceAgreed: true})
	globule.registration = reg
	if err != nil {
		return err
	}

	csrPem, err := ioutil.ReadFile(globule.creds + "/server.csr")
	if err != nil {
		return err
	}

	csrBlock, _ := pem.Decode(csrPem)
	csr, err := x509.ParseCertificateRequest(csrBlock.Bytes)
	if err != nil {
		return err
	}

	resource, err := client.Certificate.ObtainForCSR(*csr, true)
	if err != nil {
		return err
	}

	// Keep certificates url in the config.
	globule.CertURL = resource.CertURL
	globule.CertStableURL = resource.CertStableURL

	// Set the certificates paths...
	globule.Certificate = globule.getDomain() + ".crt"
	globule.CertificateAuthorityBundle = globule.getDomain() + ".issuer.crt"

	// Save the certificate in the cerst folder.
	ioutil.WriteFile(globule.creds+"/"+globule.Certificate, resource.Certificate, 0400)
	ioutil.WriteFile(globule.creds+"/"+globule.CertificateAuthorityBundle, resource.IssuerCertificate, 0400)

	// save the config with the values.
	return globule.saveConfig()
}

func (globule *Globule) signCertificate(client_csr string) (string, error) {

	// first of all I will save the incomming file into a temporary file...
	client_csr_path := os.TempDir() + "/" + Utility.RandomUUID()
	err := ioutil.WriteFile(client_csr_path, []byte(client_csr), 0644)
	if err != nil {
		return "", err

	}

	client_crt_path := os.TempDir() + "/" + Utility.RandomUUID()

	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "x509")
	args = append(args, "-req")
	args = append(args, "-passin")
	args = append(args, "pass:"+globule.CertPassword)
	args = append(args, "-days")
	args = append(args, Utility.ToString(globule.CertExpirationDelay))
	args = append(args, "-in")
	args = append(args, client_csr_path)
	args = append(args, "-CA")
	args = append(args, globule.creds+"/ca.crt") // use certificate
	args = append(args, "-CAkey")
	args = append(args, globule.creds+"/ca.key") // and private key to sign the incommin csr
	args = append(args, "-set_serial")
	args = append(args, "01")
	args = append(args, "-out")
	args = append(args, client_crt_path)
	args = append(args, "-extfile")
	args = append(args, globule.creds+"/san.conf")
	args = append(args, "-extensions")
	args = append(args, "v3_req")
	err = exec.Command(cmd, args...).Run()
	if err != nil {
		return "", err
	}

	// I will read back the crt file.
	client_crt, err := ioutil.ReadFile(client_crt_path)

	// remove the tow temporary files.
	defer os.Remove(client_crt_path)
	defer os.Remove(client_csr_path)

	if err != nil {
		return "", err
	}

	return string(client_crt), nil

}

/**
 * Initialize the server directories config, data, webroot...
 */
func (globule *Globule) initDirectories() {

	// DNS info.
	globule.DNS = make([]interface{}, 0)
	globule.DnsUpdateIpInfos = make([]interface{}, 0)

	// Set the list of discorvery service avalaible...
	globule.Discoveries = make([]string, 0)

	//////////////////////////////////////////////////////////////////////////////////////
	// There is the default directory initialisation...
	//////////////////////////////////////////////////////////////////////////////////////

	// Create the directory if is not exist.
	if globule.path == "/usr/local/share/globular" {
		// Here we have a linux standard installation.
		globule.data = "/var/globular/data"
		globule.webRoot = "/var/globular/webroot"
		globule.config = "/etc/globular/config"
	} else {

		globule.data = globule.path + "/data"
		Utility.CreateDirIfNotExist(globule.data)

		// Configuration directory
		globule.config = globule.path + "/config"
		Utility.CreateDirIfNotExist(globule.config)

		globule.webRoot = globule.path + "/webroot" // The default directory to server.
		// keep the root in global variable for the file handler.
		Utility.CreateDirIfNotExist(globule.webRoot) // Create the directory if it not exist.
	}

	// Create the creds directory if it not already exist.
	globule.creds = globule.config + "/tls"
	Utility.CreateDirIfNotExist(globule.creds)

	// Files directorie that contain user's directories and application's directory
	globule.users = globule.data + "/files/users"
	Utility.CreateDirIfNotExist(globule.users)

	// Contain the application directory.
	globule.applications = globule.data + "/files/applications"
	Utility.CreateDirIfNotExist(globule.applications)

	// Initialyse globular from it configuration file.
	file, err := ioutil.ReadFile(globule.config + "/config.json")

	// Init the service with the default port address
	if err == nil {
		err := json.Unmarshal(file, &globule)
		if err != nil {
			log.Println("fail to read init from ", globule.config+"/config.json", err)
		}

	} else {
		log.Println("fail to read configuration ", globule.config+"/config.json", err)
		jsonStr, err := Utility.ToJson(&globule)
		if err == nil {
			err := os.WriteFile(globule.config+"/config.json", []byte(jsonStr), 0644 )
			if err != nil {
				log.Println("fail to write file ", globule.config+"/config.json", err)
			}
		}
	}

	if !Utility.Exists(globule.webRoot + "/index.html") {

		// in that case I will create a new index.html file.
		ioutil.WriteFile(globule.webRoot+"/index.html", []byte(
			`<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">
		<html lang="en">
		
			<head>
				<meta http-equiv="content-type" content="text/html; charset=utf-8">
				<title>Title Goes Here</title>
			</head>
		
			<body>
				<p>Welcome to Globular `+globule.Version+`</p>
			</body>
		
		</html>`), 0644)
	}

	// Convert video file if there some to be convert.
	go func() {
		convertVideo()
	}()
	log.Println("Globular is running!")
}

/**
 * Start serving the content.
 */
func (globule *Globule) Serve() {

	globule.initDirectories()

	// TODO start admin control plane services...

	// Set the log information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// lisen
	err := globule.Listen()

	url := globule.Protocol + "://" + globule.Domain

	if globule.Protocol == "https" {
		if globule.PortHttps != 443 {
			url += ":" + Utility.ToString(globule.PortHttps)
		}
	} else if globule.Protocol == "http" {
		if globule.PortHttp != 80 {
			url += ":" + Utility.ToString(globule.PortHttp)
		}
	}

	log.Println("Globular is running at address " + url)
	if err != nil {
		log.Println(err)
	}
}

/**
 * Return the domain of the Globule. The name can be empty. If the name is empty
 * it will mean that the domain is entirely control by the globule so it must be
 * able to do it own validation, other wise the domain validation will be done by
 * the globule asscosiate with that domain.
 */
func (globule *Globule) getDomain() string {
	domain := globule.Domain
	if len(globule.Name) > 0 && domain != "localhost" {
		domain = /*globule.Name + "." +*/ domain
	}
	return domain
}

/**
 * Set the ip for a given domain or sub-domain
 */
func (globule *Globule) registerIpToDns() error {

	// Globular DNS is use to create sub-domain.
	// ex: globular1.globular.io here globular.io is the domain and globular1 is
	// the sub-domain. Domain must be manage by dns provider directly, by using
	// the DnsSetA (set ip api call)... see the next part of that function
	// for more information.
	if globule.DNS != nil {
		if len(globule.DNS) > 0 {
			for i := 0; i < len(globule.DNS); i++ {
				dns_client_, err := dns_client.NewDnsService_Client(globule.DNS[i].(string), "dns.DnsService")
				if err != nil {
					return err
				}
				// The domain is the parent domain and getDomain the sub-domain
				_, err = dns_client_.SetA(globule.Domain, globule.getDomain(), Utility.MyIP(), 60)

				if err != nil {
					// return the setA error
					return err
				}

				// TODO also register the ipv6 here...
				dns_client_.Close()
			}
		}
	}

	// Here If the DNS provides has api to update the ip address I will use it.
	// TODO test it for different internet provider's

	for i := 0; i < len(globule.DnsUpdateIpInfos); i++ {
		// the api call "https://api.godaddy.com/v1/domains/globular.io/records/A/@"
		setA := globule.DnsUpdateIpInfos[i].(map[string]interface{})["SetA"].(string)
		key := globule.DnsUpdateIpInfos[i].(map[string]interface{})["Key"].(string)
		secret := globule.DnsUpdateIpInfos[i].(map[string]interface{})["Secret"].(string)

		// set the data to the actual ip address.
		data := `[{"data":"` + Utility.MyIP() + `"}]`

		// initialize http client
		client := &http.Client{}

		// set the HTTP method, url, and request body
		req, err := http.NewRequest(http.MethodPut, setA, bytes.NewBuffer([]byte(data)))
		if err != nil {
			return err
		}

		// set the request header Content-Type for json
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
		req.Header.Set("Authorization", "sso-key "+key+":"+secret)

		// execute the request.
		_, err = client.Do(req)
		if err != nil {
			return (err)
		}

	}
	domains := globule.AlternateDomains

	for i := 0; i < len(domains); i++ {
		if !testDomainIp(domains[i].(string), Utility.MyIP(), 3) {
			return errors.New("The domain " + domains[i].(string) + "is not associated with ip " + Utility.MyIP())
		}
	}

	return nil
}

// Test if a domain is asscociated with a given ip.
func testDomainIp(domain string, ip string, try int) bool {
	if try == 0 {
		return false
	}

	if Utility.DomainHasIp(domain, Utility.MyIP()) {
		return true
	} else {
		time.Sleep(5 * time.Second)
		try--
		return testDomainIp(domain, ip, try)
	}

}

func (globule *Globule) getBasePath() string {
	// Each service contain a file name config.json that describe service.
	// I will keep services info in services map and also it running process.
	basePath, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	// Start from development environnement.
	if Utility.Exists("README.md") {
		// GLOBULAR_SERVICES_ROOT is the path of the globular service executables.
		// if not set the services must be in the same folder as Globurar executable.
		globularServicesRoot := os.Getenv("GLOBULAR_SERVICES_ROOT")
		log.Println("GLOBULAR_SERVICES_ROOT ", globularServicesRoot)
		if len(globularServicesRoot) > 0 {
			basePath = globularServicesRoot
		}
	}
	return basePath
}

/**
 * retreive checksum from the server.
 */
func getChecksum(address string, port int) (string, error) {
	if len(address) == 0 {
		return "", errors.New("no address was given")
	}

	// Here I will get the configuration information from http...
	var resp *http.Response
	var err error
	var url = "http://" + address + ":" + Utility.ToString(port) + "/checksum"
	resp, err = http.Get(url)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusCreated {

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		return string(bodyBytes), nil
	}

	return "", errors.New("fail to retreive checksum with error " + Utility.ToString(resp.StatusCode))
}

func (globule *Globule) watchForUpdate() {
	go func() {
		for {

			// stop watching if exit was call...
			if globule.exit_ {
				return
			}

			if len(globule.Discoveries) > 0 {
				// Here I will retreive the checksum information from it parent.
				discovery := globule.Discoveries[0]
				address := strings.Split(discovery, ":")[0]
				port := 80
				if strings.Contains(discovery, ":") {
					port = Utility.ToInt(strings.Split(discovery, ":")[1])
				}

				// Here I will test if the checksum has change...
				checksum, err := getChecksum(address, port)
				execPath := Utility.GetExecName(os.Args[0])
				if Utility.Exists("/usr/local/share/globular/Globular") {
					execPath = "/usr/local/share/globular/Globular"
				}
				if err == nil {
					if checksum != Utility.CreateFileChecksum(execPath) {

						err := update_globular_from(globule, discovery, globule.Domain, "sa", globule.RootPassword, runtime.GOOS+":"+runtime.GOARCH)
						if err != nil {
							log.Println("fail to update globular from " + discovery + " with error " + err.Error())
						} else {
							log.Println("update globular checksum is ", checksum)
						}

					}
				}
			}

			// The time here can be set to higher value.
			time.Sleep(time.Duration(globule.WatchUpdateDelay) * time.Second)
		}
	}()
}

// check if the process is actually running
// However, on Unix systems, os.FindProcess always succeeds and returns
// a Process for the given pid...regardless of whether the process exists
// or not.
func getProcessRunningStatus(pid int) (*os.Process, error) {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return nil, err
	}

	//double check if process is running and alive
	//by sending a signal 0
	//NOTE : syscall.Signal is not available in Windows

	err = proc.Signal(syscall.Signal(0))
	if err == nil {
		return proc, nil
	}

	if err == syscall.ESRCH {
		return nil, errors.New("process not running")
	}

	// default
	return nil, errors.New("process running but query operation not permitted")
}

/**
 * Listen for new connection.
 */
func (globule *Globule) Listen() error {

	var err error
	// Must be started before other services.
	go func() {
		// local - non secure connection.
		globule.http_server = &http.Server{
			Addr: ":" + strconv.Itoa(globule.PortHttp),
		}
		err = globule.http_server.ListenAndServe()
	}()

	// Start the http server.
	if globule.Protocol == "https" {

		// if no certificates are specified I will try to get one from let's encrypts.
		// Start https server.
		globule.https_server = &http.Server{
			Addr: ":" + strconv.Itoa(globule.PortHttps),
			TLSConfig: &tls.Config{
				ServerName: globule.getDomain(),
			},
		}

		// get the value from the configuration files.
		go func() {
			err = globule.https_server.ListenAndServeTLS(globule.creds+"/"+globule.Certificate, globule.creds+"/server.pem")
		}()
	}

	//////////////////////////////////////////////////////////
	// Keep the server state in sync with configuration and
	// version.
	//////////////////////////////////////////////////////////

	// Keep globular up to date subscription.
	globule.watchForUpdate()

	return err
}

//////////////////////// RBAC function //////////////////////////////////////////////
var (
	rbac_client_ *rbac_client.Rbac_Client
	event_client_ *event_client.Event_Client
)

/**
 * Get the rbac client.
 */
 func GetRbacClient(domain string) (*rbac_client.Rbac_Client, error) {
	var err error
	if rbac_client_ == nil {
		rbac_client_, err = rbac_client.NewRbacService_Client(domain, "rbac.RbacService")
		if err != nil {
			log.Println("fail to get RBAC client with error ", err)
			return nil, err
		}

	}
	return rbac_client_, nil
}

// Use rbac client here...
func (globule *Globule) addResourceOwner(path string, subject string, subjectType rbacpb.SubjectType) error {
	rbac_client_, err := GetRbacClient(globule.Domain)
	if err != nil {
		return err
	}
	return rbac_client_.AddResourceOwner(path, subject, subjectType)
}

func (globule *Globule) validateAction(method string, subject string, subjectType rbacpb.SubjectType, infos []*rbacpb.ResourceInfos) (bool, error) {
	rbac_client_, err := GetRbacClient(globule.Domain)
	if err != nil {
		return false, err
	}

	return rbac_client_.ValidateAction(method, subject, subjectType, infos)
}

func (globule *Globule) validateAccess(subject string, subjectType rbacpb.SubjectType, name string, path string) (bool, bool, error) {
	rbac_client_, err := GetRbacClient(globule.Domain)
	if err != nil {
		return false,false, err
	}

	return rbac_client_.ValidateAccess(subject, subjectType, name, path)
}

///////////////////// event service functions ////////////////////////////////////
func (globule *Globule) getEventClient() (*event_client.Event_Client, error) {
	var err error
	if event_client_ != nil {
		return event_client_, nil
	}
	event_client_, err = event_client.NewEventService_Client(globule.Domain, "event.EventService")
	if err != nil {
		return nil, err
	}

	return event_client_, nil
}

func (globule *Globule) publish(event string, data []byte) error {
	eventClient, err := globule.getEventClient()
	if err != nil {
		return err
	}
	return eventClient.Publish(event, data)
}


/////////////////////// services manager functions ///////////////////////////////
/**
 * Return an array of all services available on the globule
 */
func (globular *Globule) getServices() ([]map[string]interface{}, error) {

	services := make([]map[string]interface{}, 0)

	return services, errors.New("--------------> not implemented")
}
