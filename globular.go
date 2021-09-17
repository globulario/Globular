package main

import (
	"bytes"
	"context"
	"crypto"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/globulario/services/golang/admin/admin_client"
	"github.com/globulario/services/golang/authentication/authentication_client"
	"github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/dns/dns_client"
	"github.com/globulario/services/golang/event/event_client"
	"github.com/globulario/services/golang/event/eventpb"
	"github.com/globulario/services/golang/globular_service"
	"github.com/globulario/services/golang/interceptors"
	"github.com/globulario/services/golang/log/log_client"
	"github.com/globulario/services/golang/log/logpb"
	"github.com/globulario/services/golang/persistence/persistence_client"
	"github.com/globulario/services/golang/process"
	"github.com/globulario/services/golang/rbac/rbac_client"
	"github.com/globulario/services/golang/rbac/rbacpb"
	"github.com/globulario/services/golang/resource/resource_client"
	"github.com/globulario/services/golang/security"
	"github.com/gookit/color"

	// Interceptor for authentication, event, log...

	// Client services.
	"github.com/davecourtois/Utility"
	"github.com/go-acme/lego/certcrypto"
	"github.com/go-acme/lego/challenge/http01"
	"github.com/go-acme/lego/lego"
	"github.com/go-acme/lego/registration"
	"github.com/go-acme/lego/v4/challenge/dns01"
)

// Global variable.
var (
	globule *Globule
)

/**
 * The web server.
 */
type Globule struct {
	// The share part of the service.
	Name string // The service name
	Mac  string // The Mac addresse

	// Globualr specifics ports.

	// can be https or http.
	Protocol   string
	PortHttp   int    // The port of the http file server.
	PortHttps  int    // The secure port
	PortsRange string // The range of grpc ports.

	// Cors policy
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string

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
	Version  string
	Build    int64
	Platform string

	// Admin informations.
	AdminEmail   string
	RootPassword string

	SessionTimeout int // The time before session expire.

	// Service discoveries.
	Discoveries []string // Contain the list of discovery service use to keep globular up to date.

	// Update delay in second...
	WatchUpdateDelay int

	// DNS stuff.
	DNS              []interface{} // Domain name server use to located the server.
	DnsUpdateIpInfos []interface{} // The internet provader SetA info to keep ip up to date.

	// Directories.
	path    string // The path of the exec...
	webRoot string // The root of the http file server.
	data    string // the data directory
	creds   string // tls certificates
	config  string // configuration directory

	users        string // the users files directory
	applications string // The applications
	templates    string // the html/css templates
	projects     string // the web projects

	// ACME protocol registration
	registration *registration.Resource

	// exit channel.
	exit  chan bool
	exit_ bool

	// The http server
	http_server  *http.Server
	https_server *http.Server

	// The http client(s)
	https_clients map[string]*http.Client

	// Keep track of the strart time...
	startTime time.Time
}

/**
 * Globule constructor.
 */
func NewGlobule() *Globule {

	// Here I will keep the start time...


	// Here I will initialyse configuration.
	g := new(Globule)
	g.startTime = time.Now()
	g.exit_ = false
	g.exit = make(chan bool)
	g.Mac = Utility.MyMacAddr()
	g.Version = "1.0.0" // Automate version...
	g.Build = 0
	g.Platform = runtime.GOOS + ":" + runtime.GOARCH
	g.IndexApplication = "" // I will use the installer as defaut.

	g.PortHttp = 80              // The default http port
	g.PortHttps = 443            // The default https port number
	g.PortsRange = "10000-10100" // The default port range.

	if g.AllowedOrigins == nil {
		g.AllowedOrigins = []string{"*"}
	}

	if g.AllowedMethods == nil {
		g.AllowedMethods = []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"}
	}

	if g.AllowedHeaders == nil {
		g.AllowedHeaders = []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "domain", "application", "token"}
	}

	// Set the default checksum...
	g.Protocol = "http"
	g.Domain = "localhost"

	// set default values.
	g.CertExpirationDelay = 365
	g.CertPassword = "1111"
	g.AdminEmail = "root@globular.app"
	g.RootPassword = "adminadmin"

	// keep up to date by default.
	g.WatchUpdateDelay = 30 // seconds...
	g.SessionTimeout = 15 * 60 * 1000

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

	// Stop previous running process.
	g.stopProxies()

	g.initDirectories()

	return g
}

func (globule *Globule) registerAdminAccount() error {
	resource_client_, err := GetResourceClient(globule.Domain)
	if err != nil {
		return err
	}

	// Create the admin account.
	err = resource_client_.RegisterAccount("sa", globule.AdminEmail, globule.RootPassword, globule.RootPassword)
	if err != nil {

		return err
	}

	// Set admin role to that account.
	err = resource_client_.AddAccountRole("sa", "admin")
	if err != nil {

		return err
	}

	return nil
}

/**
 * Find http client associated with a given domain. Http server must be register
 * as a peer's to be able to process http request on the same domain.
 */
func (globule *Globule) getHttpClient(domain string) (*http.Client, error) {

	domain = strings.Split(domain, ":")[0]
	port := 80
	if len(strings.Split(domain, ":")) > 1 {
		port = Utility.ToInt(strings.Split(domain, ":")[1])
	}

	resource_client_, err := GetResourceClient(domain)
	if err != nil {
		return nil, err
	}

	peers, err := resource_client_.GetPeers(`{"domain":"` + domain + `"}`)
	if err != nil {
		return nil, err
	}

	if len(peers) == 0 {
		return nil, errors.New("no peers was found for domain " + domain)
	}

	p := peers[0]

	if globule.https_clients == nil {
		globule.https_clients = make(map[string]*http.Client, 0)
	}

	// if a http client was already register then I will made use of it.
	if globule.https_clients[p.Domain] != nil {
		return globule.https_clients[p.Domain], nil
	}

	clientCertFile := os.TempDir() + "/config/tls/" + domain + "/client.crt"
	clientKeyFile := os.TempDir() + "/config/tls/" + domain + "/client.pem"
	caCertFile := os.TempDir() + "/config/tls/" + domain + "/ca.crt"

	if !Utility.Exists(os.TempDir() + "/config/tls/" + domain) {
		admin_client_, err := admin_client.NewAdminService_Client(domain, "admin.AdminService")
		if err != nil {
			return nil, err
		}

		clientKeyFile, clientCertFile, caCertFile, err = admin_client_.GetCertificates(domain, port, os.TempDir())
		if err != nil {
			return nil, err
		}
	}

	if !Utility.Exists(clientKeyFile) {
		return nil, errors.New("no client key file found at path " + clientKeyFile)
	}

	if !Utility.Exists(clientKeyFile) {
		return nil, errors.New("no client certificate file found at path " + clientCertFile)
	}

	if !Utility.Exists(clientKeyFile) {
		return nil, errors.New("no client CA certificate file found at path " + caCertFile)
	}

	// So here I will create the client.
	// I will made use of the certifcate to connect.
	t := &http.Transport{
		TLSClientConfig: globular_service.GetTLSConfig(clientKeyFile, clientCertFile, caCertFile),
	}

	// open the client connection
	client := http.Client{Transport: t, Timeout: 15 * time.Second}

	// Keep the client for further use
	globule.https_clients[p.Domain] = &client

	return &client, nil
}

/**
 * Return globular configuration.
 */
func (globule *Globule) getConfig() map[string]interface{} {
	// TODO filter unwanted attributes...
	config_, _ := Utility.ToMap(globule)
	services, _ := config.GetServicesConfigurations()

	// Get the array of service and set it back in the configurations.
	config_["Services"] = make(map[string]interface{})

	// Here I will set in a map and put in the Services key
	for i := 0; i < len(services); i++ {
		s := make(map[string]interface{})
		s["AllowAllOrigins"] = services[i]["AllowAllOrigins"]
		s["AllowedOrigins"] = services[i]["AllowedOrigins"]
		s["Description"] = services[i]["Description"]
		s["Discoveries"] = services[i]["Discoveries"]
		s["Domain"] = services[i]["Domain"]
		s["Id"] = services[i]["Id"]
		s["Keywords"] = services[i]["Keywords"]
		s["Name"] = services[i]["Name"]
		s["Mac"] = services[i]["Mac"]
		s["Port"] = services[i]["Port"]
		s["Proxy"] = services[i]["Proxy"]
		s["PublisherId"] = services[i]["PublisherId"]
		s["State"] = services[i]["State"]
		s["TLS"] = services[i]["TLS"]
		s["Dependencies"] = services[i]["Dependencies"]
		s["Version"] = services[i]["Version"]
		config_["Services"].(map[string]interface{})[s["Id"].(string)] = s
	}

	return config_
}

/**
 * That function will wath if the configuration file has change.
 */
func (globule *Globule) watchConfig() {
	go func() {
		checksum := Utility.CreateFileChecksum(globule.config + "/config.json")

		for {
			checksum_ := Utility.CreateFileChecksum(globule.config + "/config.json")

			if checksum_ != checksum {
				file, _ := ioutil.ReadFile(globule.config + "/config.json")
				config := make(map[string]interface{})

				err := json.Unmarshal(file, &config)

				if err != nil {
					globule.saveConfig() // write back the configuration...
				} else {

					// Here I will make some validation...
					if config["Protocol"].(string) == "https" && config["Domain"].(string) == "localhost" {
						log.Println("The domain localhost cannot be use with https, domain must contain dot's")
					} else {

						hasProtocolChange := globule.Protocol != config["Protocol"].(string)
						hasDomainChange := globule.Domain != config["Domain"].(string)
						certificateChange := globule.CertificateAuthorityBundle != config["CertificateAuthorityBundle"].(string)
						json.Unmarshal(file, &globule)

						// stop the http server
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

						if err = globule.http_server.Shutdown(ctx); err != nil {
							log.Println("fail to stop the http server with error ", err)
						}

						if globule.https_server != nil {
							if err := globule.https_server.Shutdown(ctx); err != nil {
								log.Println("fail to stop the https server with error ", err)
							}
						}

						// restart
						globule.serve()

						if hasProtocolChange || hasDomainChange || certificateChange {
							// stop services...
							fmt.Println("Stop gRpc Services")
							globule.stopServices()
							// restart it...
							globule.startServices()
							// start proxies
							globule.startProxies()

							// restart watching
							//process.ManageServicesProcess(globule.exit)
						}

						// clear context
						cancel()

						checksum = checksum_
					}
				}
			}

			time.Sleep(time.Duration(10) * time.Second)
		}
	}()
}

/**
 * Save the configuration
 */
func (globule *Globule) saveConfig() error {

	jsonStr, err := Utility.ToJson(globule)
	if err != nil {
		return err
	}

	configPath := globule.config + "/config.json"

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
 * That function must be use to generate public
 */

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

/** That interface withe the let's encrypt DNS chanlenge. **/
type DNSProviderGlobularDNS struct {
	apiAuthToken string
}

func NewDNSProviderGlobularDNS(apiAuthToken string) (*DNSProviderGlobularDNS, error) {
	return &DNSProviderGlobularDNS{apiAuthToken: apiAuthToken}, nil
}

func (d *DNSProviderGlobularDNS) Present(domain, token, keyAuth string) error {
	key, value := dns01.GetRecord(domain, keyAuth)

	if len(globule.DNS) > 0 {
		dns_client_, err := dns_client.NewDnsService_Client(globule.DNS[0].(string), "dns.DnsService")
		if err != nil {
			return err
		}

		jwtKey, err := security.GetPeerKey(dns_client_.GetMac())
		if err != nil {
			return err
		}
		token, err := interceptors.GenerateToken(jwtKey, time.Duration(globule.SessionTimeout), Utility.MyMacAddr(), "", "", globule.AdminEmail)

		if err != nil {

			return err
		}

		err = dns_client_.SetText(token, key, []string{value}, 30)
		if err != nil {
			log.Println("fail to set let's encrypt dns chalenge key with error ", err)
			return err
		}
	}

	return nil
}

func (d *DNSProviderGlobularDNS) CleanUp(domain, token, keyAuth string) error {
	// clean up any state you created in Present, like removing the TXT record

	key, _ := dns01.GetRecord(domain, keyAuth)

	if len(globule.DNS) > 0 {
		dns_client_, err := dns_client.NewDnsService_Client(globule.DNS[0].(string), "dns.DnsService")
		if err != nil {
			return err
		}

		jwtKey, err := security.GetPeerKey(dns_client_.GetMac())
		if err != nil {
			return err
		}
		token, err := interceptors.GenerateToken(jwtKey, time.Duration(globule.SessionTimeout), Utility.MyMacAddr(), "", "", globule.AdminEmail)

		if err != nil {

			return err
		}

		err = dns_client_.RemoveText(token, key)
		if err != nil {
			log.Println("fail to remove challenge key with error ", err)
			return err
		}
	}
	return nil
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
	// Dns registration will be use in case dns service are available.
	if len(globule.DNS) > 0 {

		// Get the local token.
		token, err := globule.getLocalToken(globule.DNS[0].(string))
		if err != nil {
			return err
		}

		globularDNS, err := NewDNSProviderGlobularDNS(token)
		if err != nil {

			return err
		}

		client.Challenge.SetDNS01Provider(globularDNS)

	} else {
		err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", strconv.Itoa(globule.PortHttp)))
		if err != nil {

			log.Fatal(err)
		}
	}

	if err != nil {

		return err
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
	rqstForCsr, err := x509.ParseCertificateRequest(csrBlock.Bytes)
	if err != nil {

		return err
	}

	resource, err := client.Certificate.ObtainForCSR(*rqstForCsr, true)
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
func (globule *Globule) initDirectories() error {

	// DNS info.
	globule.DNS = make([]interface{}, 0)
	globule.DnsUpdateIpInfos = make([]interface{}, 0)

	// Set the list of discorvery service avalaible...
	globule.Discoveries = make([]string, 0)

	//////////////////////////////////////////////////////////////////////////////////////
	// There is the default directory initialisation...
	//////////////////////////////////////////////////////////////////////////////////////

	// Create the directory if is not exist.
	globule.data = config.GetDataDir()

	err := Utility.CreateDirIfNotExist(globule.data)
	if err != nil {
		return err
	}

	globule.webRoot = config.GetWebRootDir()

	err = Utility.CreateDirIfNotExist(globule.webRoot)
	if err != nil {
		return err
	}

	globule.templates = globule.data + "/files/templates"

	err = Utility.CreateDirIfNotExist(globule.templates)
	if err != nil {
		return err
	}

	globule.projects = globule.data + "/files/projects"
	err = Utility.CreateDirIfNotExist(globule.projects)
	if err != nil {
		return err
	}

	globule.config = config.GetConfigDir()
	err = Utility.CreateDirIfNotExist(globule.config)

	if err != nil {
		log.Println("fail to create configuration directory  with error", err)
		return err
	}

	// Create the tokens directory
	err = Utility.CreateDirIfNotExist(globule.config + "/tokens")
	if err != nil {
		log.Println("fail to create tokens directory  with error", err)
		return err
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
			log.Println("fail to init configuation with error ", err)
			return err
		}

	} else {
		jsonStr, err := Utility.ToJson(&globule)
		if err == nil {
			err := os.WriteFile(globule.config+"/config.json", []byte(jsonStr), 0644)
			if err != nil {
				return err
			}
		}
	}

	// save config...
	globule.saveConfig()

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
		convertVideo() // call once and at each minutes....
	}()

	return nil
}


func (globule *Globule) stopProxies(){
	execName := "grpcwebproxy"
	if runtime.GOOS == "windows" {
		execName += ".exe" // in case of windows
	}

	// Kill all proxies..
	Utility.KillProcessByName(execName)

}

/**
 * Start proxies
 */
func (globule *Globule) startProxies() {
	fmt.Println("Start gRpc proxies")

	services, err := config.GetServicesConfigurations()
	if err == nil {
		for i := 0; i < len(services); i++ {
			// Here I will start the proxy
			process.StartServiceProxyProcess(services[i]["Id"].(string), globule.CertificateAuthorityBundle, globule.Certificate, globule.PortsRange, Utility.ToInt(services[i]["Process"]))
		}
	}
}

/**
 * Here I will start the services manager who will start all microservices
 * installed on that computer.
 */
func (globule *Globule) startServices() error {
	fmt.Println("Start gRpc Services")
	// Retreive all configurations
	services, err := config.GetServicesConfigurations()
	if err != nil {
		return err
	}

	// I will try to get the services manager configuration from the
	// services configurations list.
	for i := 0; i < len(services); i++ {
		// Set the
		if globule.Protocol == "https" {

			// set tls file...
			services[i]["TLS"] = true
			services[i]["KeyFile"] = globule.creds + "/client.pem"
			services[i]["CertFile"] = globule.creds + "/client.crt"
			services[i]["CertAuthorityTrust"] = globule.creds + "/ca.crt"

			if services[i]["CertificateAuthorityBundle"] != nil {
				services[i]["CertificateAuthorityBundle"] = globule.CertificateAuthorityBundle
			}

			if services[i]["Certificate"] != nil {
				services[i]["Certificate"] = globule.Certificate
			}

		} else {
			services[i]["TLS"] = false
		}

		// Save back the values...
		services[i]["Domain"] = globule.getDomain()
		services[i]["Mac"] = globule.Mac

		config.SaveServiceConfiguration(services[i]) // save service values.

		// Create the service process.
		_, err = process.StartServiceProcess(services[i]["Id"].(string), globule.PortsRange)
		if err != nil {
			log.Println("fail to start service ", services[i]["Name"], err)
		}

		// Here I will listen for logger event...
		go func() {
			// subscribe to log events
			globule.subscribe("new_log_evt", logListener)

			// subscribe to serive change event.
			globule.subscribe("update_globular_service_configuration_evt", updateServiceConfigurationListener)
			
		}()
	}

	// Create the admin account.
	globule.registerAdminAccount()

	// Create application connection
	globule.createApplicationConnection()

	return nil
}

/**
 * Here I will create application backend connection.
 */
func (globule *Globule) createApplicationConnection() error {
	resource_client_, err := GetResourceClient(globule.Domain)
	if err != nil {
		return err
	}

	persistence_client_, err := GetPersistenceClient(globule.Domain)
	if err != nil {
		return err
	}

	if err == nil {
		applications, err := resource_client_.GetApplications("{}")
		if err == nil {
			for i := 0; i < len(applications); i++ {
				app := applications[i]
				//err := persistence_client_.CreateConnection("local_resource", "local_resource",  "localhost", 27017, 0, "sa", "adminadmin", 500, "", true)
				err := persistence_client_.CreateConnection(app.Id+"_db", app.Id, globule.Domain, 27017, 0, app.Id, app.Password, 500, "", true)
				if err != nil {
					fmt.Println("fail to create application connection  : ", app.Id, err)
				}
			}
		}
	}

	return err
}

/**
 * Start refresh local token...
 */
func (globule *Globule) startRefreshLocalTokens() {
	globule.refreshLocalTokens()
	ticker := time.NewTicker(time.Duration(globule.SessionTimeout) * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				// Connect to service update events...
				// I will iterate over the list token and close expired session...
				globule.refreshLocalTokens()
			case <-globule.exit:

				return // exit from the loop when the service exit.
			}
		}
	}()
}

/**
 * Stop all services.
 */
func (globule *Globule) stopServices() error {
	services, err := config.GetServicesConfigurations()
	if err != nil {
		return err
	}

	// exit channel.
	globule.exit <- true

	for i := 0; i < len(services); i++ {
		process.KillServiceProcess(services[i])
	}

	return nil
}

// Start http/https server...
func (globule *Globule) serve() error {
	// start listen
	err := globule.Listen()
	if err != nil {
		return err
	}

	url := globule.Protocol + "://" + globule.getDomain()

	if globule.Protocol == "https" {
		if globule.PortHttps != 443 {
			url += ":" + Utility.ToString(globule.PortHttps)
		}
	} else if globule.Protocol == "http" {
		if globule.PortHttp != 80 {
			url += ":" + Utility.ToString(globule.PortHttp)
		}
	}

	elapsed := time.Since(globule.startTime)

	log.Println("globular version " + globule.Version + " build " + Utility.ToString(globule.Build) + " listen at address " + url)
	log.Printf("startup took %s", elapsed)

	return nil

}

/**
 * Start serving the content.
 */
func (globule *Globule) Serve() error {

	// Initialyse directories.
	globule.initDirectories()

	// Start microservice manager.
	globule.startServices()

	// start proxies
	globule.startProxies()

	// Here I will remove the local token and recreate it...
	globule.startRefreshLocalTokens()

	globule.watchConfig()

	// Set the log information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Start managing process.
	//process.ManageServicesProcess(globule.exit)

	return globule.serve()
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
		domain = globule.Name + "." + domain
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

				key, err := security.GetPeerKey(dns_client_.GetMac())
				if err != nil {
					return err
				}

				// Here the token must be generated for the dns server...
				// That peer must be register on the dns to be able to generate a valid token.
				token, err := interceptors.GenerateToken(key, time.Duration(globule.SessionTimeout), Utility.MyMacAddr(), "", "", globule.AdminEmail)

				if err != nil {

					return err
				}
				// The domain is the parent domain and getDomain the sub-domain
				_, err = dns_client_.SetA(token, globule.Domain, globule.getDomain(), Utility.MyIP(), 60)

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

/**
 *  Watch if globular need to be update.
 */
func (globule *Globule) watchForUpdate() {
	go func() {
		for !globule.exit_ {

			// stop watching if exit was call...
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
				if Utility.Exists(config.GetRootDir() + "/Globular") {
					execPath = config.GetRootDir() + "/Globular"
				}
				if err == nil {
					if checksum != Utility.CreateFileChecksum(execPath) {

						err := update_globular_from(globule, discovery, globule.getDomain(), "sa", globule.RootPassword, runtime.GOOS+":"+runtime.GOARCH)
						if err != nil {
							log.Println("fail to update globular from " + discovery + " with error " + err.Error())
						}

					}
				}
			}

			// The time here can be set to higher value.
			time.Sleep(time.Duration(globule.WatchUpdateDelay) * time.Second)
		}
	}()
}

// received when service configuration change.
func updateServiceConfigurationListener(evt *eventpb.Event) {
	s := make(map[string]interface{})
	err := json.Unmarshal(evt.Data, &s)
	if err == nil {
		config.SetServiceConfiguration(s)
	}
}

// Try to display application message in a nice way
func logListener(evt *eventpb.Event) {
	info := make(map[string]interface{})
	err := json.Unmarshal(evt.Data, &info)
	if err == nil {
		// So here Will display message
		var header string
		if info["application"] != nil {
			header = info["application"].(string)
		}
		messageTime := time.Unix(int64(Utility.ToInt(info["date"])), 0)
		method := "NA"
		if info["method"] != nil {
			method = info["method"].(string)
		}

		if info["functionName"] != nil {
			method += ":" + info["functionName"].(string)
		}

		header += " " + messageTime.Format("2006-01-02 15:04:05") + " " + method

		if info["level"].(string) == "ERROR_MESSAGE" {
			color.Error.Println(header)
		} else if info["level"].(string) == "DEBUG_MESSAGE" || info["level"].(string) == "INFO_MESSAGE" {
			color.Info.Println(header)
		} else {
			color.Warn.Println(header)
		}

		if info["message"] != nil {
			// Now I will process the message itself...
			msg := info["message"].(string)
			// if the message is grpc error I will parse it content and display it content...
			if strings.HasPrefix(msg, "rpc") {
				errorDescription := make(map[string]interface{})
				startIndex := strings.Index(msg, "{")
				endIndex := strings.Index(msg, "}")
				if startIndex >= 0 && endIndex > startIndex {
					jsonStr := msg[startIndex : endIndex+1]
					err := json.Unmarshal([]byte(jsonStr), &errorDescription)
					if err == nil {
						if errorDescription["FileLine"] != nil {
							fmt.Println(errorDescription["FileLine"])
						}
						if errorDescription["ErrorMsg"] != nil {
							color.Comment.Println(errorDescription["ErrorMsg"])
						}
					}
				}
			} else {
				if info["line"] != nil {
					fmt.Println(info["line"])
				}
				color.Comment.Println(msg)
			}
		}
	}
}

/**
 * Listen for new connection.
 */
func (globule *Globule) Listen() error {

	var err error

	// if no certificates are specified I will try to get one from let's encrypts.
	// Start https server.
	if len(globule.Certificate) == 0 && globule.Protocol == "https" {

		// Here is the command to be execute in order to ge the certificates.
		// ./lego --email="admin@globular.app" --accept-tos --key-type=rsa4096 --path=../config/http_tls --http --csr=../config/tls/server.csr run
		// I need to remove the gRPC certificate and recreate it.
		err = Utility.RemoveDirContents(globule.creds)
		if err != nil {
			return err
		}

		// recreate the certificates.
		err = security.GenerateServicesCertificates(globule.CertPassword, globule.CertExpirationDelay, globule.getDomain(), globule.creds, globule.Country, globule.State, globule.City, globule.Organization, globule.AlternateDomains)
		if err != nil {
			return err
		}

		// Register that peer with the dns.
		err := globule.registerIpToDns()
		if err != nil {
			return err
		}

		err = globule.obtainCertificateForCsr()
		if err != nil {
			return err
		}
	}

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

var (
	rbac_client_           *rbac_client.Rbac_Client
	event_client_          *event_client.Event_Client
	authentication_client_ *authentication_client.Authentication_Client
	log_client_            *log_client.Log_Client
	resource_client_       *resource_client.Resource_Client
	persistence_client_    *persistence_client.Persistence_Client
)

//////////////////////// Resource Client ////////////////////////////////////////////
func GetResourceClient(domain string) (*resource_client.Resource_Client, error) {
	var err error
	if resource_client_ == nil {
		resource_client_, err = resource_client.NewResourceService_Client(domain, "resource.ResourceService")
		if err != nil {
			return nil, err
		}

	}

	return resource_client_, nil
}

//////////////////////// Resource Client ////////////////////////////////////////////
func GetPersistenceClient(domain string) (*persistence_client.Persistence_Client, error) {
	var err error
	if persistence_client_ == nil {
		persistence_client_, err = persistence_client.NewPersistenceService_Client(domain, "persistence.PersistenceService")
		if err != nil {
			log.Println("fail to get persistence client with error ", err)
			return nil, err
		}

	}

	return persistence_client_, nil
}

//////////////////////// RBAC function //////////////////////////////////////////////

/**
 * Get the rbac client.
 */
func GetRbacClient(domain string) (*rbac_client.Rbac_Client, error) {
	var err error
	if rbac_client_ == nil {
		rbac_client_, err = rbac_client.NewRbacService_Client(domain, "rbac.RbacService")
		if err != nil {
			return nil, err
		}

	}
	return rbac_client_, nil
}

// Use rbac client here...
func (globule *Globule) addResourceOwner(path string, subject string, subjectType rbacpb.SubjectType) error {
	rbac_client_, err := GetRbacClient(globule.getDomain())
	if err != nil {
		return err
	}
	return rbac_client_.AddResourceOwner(path, subject, subjectType)
}

func (globule *Globule) validateAction(method string, subject string, subjectType rbacpb.SubjectType, infos []*rbacpb.ResourceInfos) (bool, error) {
	rbac_client_, err := GetRbacClient(globule.getDomain())
	if err != nil {
		return false, err
	}

	return rbac_client_.ValidateAction(method, subject, subjectType, infos)
}

func (globule *Globule) validateAccess(subject string, subjectType rbacpb.SubjectType, name string, path string) (bool, bool, error) {
	rbac_client_, err := GetRbacClient(globule.getDomain())
	if err != nil {
		return false, false, err
	}

	return rbac_client_.ValidateAccess(subject, subjectType, name, path)
}

///////////////////// event service functions ////////////////////////////////////
func (globule *Globule) getEventClient() (*event_client.Event_Client, error) {
	var err error
	if event_client_ != nil {
		return event_client_, nil
	}
	event_client_, err = event_client.NewEventService_Client(globule.getDomain(), "event.EventService")
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

func (globule *Globule) subscribe(evt string, listener func(evt *eventpb.Event)) error {
	eventClient, err := globule.getEventClient()
	if err != nil {
		return err
	}

	// register a listener...
	return eventClient.Subscribe(evt, globule.Name, listener)
}

///////////////////////  Log Services functions ////////////////////////////////////////////////

/**
 * Get the log client.
 */
func (globule *Globule) GetLogClient() (*log_client.Log_Client, error) {
	var err error
	if log_client_ == nil {
		log_client_, err = log_client.NewLogService_Client(globule.Domain, "log.LogService")
		if err != nil {
			return nil, err
		}

	}
	return log_client_, nil
}

func (globule *Globule) log(fileLine, functionName, message string, level logpb.LogLevel) {
	log_client_, err := globule.GetLogClient()
	if err != nil {
		return
	}
	log_client_.Log(globule.Name, globule.Domain, functionName, level, message, fileLine, functionName)
}

/**
 * Generate local token key, that key is use by internal service.
 */
func (globule *Globule) refreshLocalTokens() error {

	tokensPath := globule.config + "/tokens"

	// Here I will get the list token files in the folder...
	files, err := ioutil.ReadDir(tokensPath)
	if err != nil {
		log.Fatal(err)
	}

	// Remove expired token files.
	for _, f := range files {
		authentication_client_, err = authentication_client.NewAuthenticationService_Client(strings.ReplaceAll(f.Name(), "_token", ""), "authentication.AuthenticationService")
		if err == nil {
			path := tokensPath + "/" + f.Name()
			data, err := ioutil.ReadFile(path)
			if err == nil {
				err := authentication_client_.ValidateToken(string(data))
				if err != nil {
					// remove path from the list.
					os.Remove(path)
				}
			}
		}
	}

	// The local key.
	key, err := security.GetLocalKey()
	if err != nil {
		return err
	}

	// This is the local token...
	tokenString, err := interceptors.GenerateToken(key, time.Duration(globule.SessionTimeout), Utility.MyMacAddr(), "sa", "sa", globule.AdminEmail)
	if err != nil {
		return err
	}

	path := tokensPath + "/" + globule.getDomain() + "_token"
	return ioutil.WriteFile(path, []byte(tokenString), 0644)
}

/**
 * Return the local token string
 */
func (globule *Globule) getLocalToken(domain string) (string, error) {
	tokensPath := globule.config + "/tokens"
	path := tokensPath + "/" + domain + "_token"

	token, err := ioutil.ReadFile(path)

	if err != nil {
		return "", nil
	}

	return string(token), nil
}
