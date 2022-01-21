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

	"github.com/globulario/services/golang/authentication/authentication_client"
	"github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/config/config_client"
	"github.com/globulario/services/golang/dns/dns_client"
	"github.com/globulario/services/golang/event/event_client"
	"github.com/globulario/services/golang/event/eventpb"
	"github.com/globulario/services/golang/log/log_client"
	"github.com/globulario/services/golang/log/logpb"
	"github.com/globulario/services/golang/persistence/persistence_client"
	"github.com/globulario/services/golang/process"
	"github.com/globulario/services/golang/rbac/rbac_client"
	"github.com/globulario/services/golang/rbac/rbacpb"
	"github.com/globulario/services/golang/resource/resource_client"
	"github.com/globulario/services/golang/resource/resourcepb"
	"github.com/globulario/services/golang/security"
	"github.com/gookit/color"
	"github.com/kardianos/service"
	"github.com/txn2/txeh"

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
	Build    int
	Platform string

	// Admin informations.
	AdminEmail   string
	RootPassword string

	SessionTimeout int // The time before session expire.

	// Service discoveries.
	Discoveries []string // Contain the list of discovery service use to keep globular up to date.

	// Update delay in second...
	WatchUpdateDelay int64

	// DNS stuff.
	DNS              []interface{} // External dns.
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

	// List of peers
	peers []*resourcepb.Peer

	// Keep track of the strart time...
	startTime time.Time

	// This is use to display information to external service manager.
	logger service.Logger
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
	g.IndexApplication = ""      // I will use the installer as defaut.
	g.PortHttp = 8080            // The default http port 80 is almost already use by other http server...
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
	g.Name, _ = config.GetHostName()

	// set default values.
	g.CertExpirationDelay = 365
	g.CertPassword = "1111"
	g.AdminEmail = "root@globular.app"
	g.RootPassword = "adminadmin"

	// keep up to date by default.
	g.WatchUpdateDelay = 30 // seconds...
	g.SessionTimeout = 15   // in minutes

	// Keep in global var to by http handlers.
	globule = g

	// Set the list of http handler.

	// The configuration handler.
	http.HandleFunc("/config", getConfigHanldler)

	// The checksum handler.
	http.HandleFunc("/checksum", getChecksumHanldler)

	// Handle the get ca certificate function
	http.HandleFunc("/get_ca_certificate", getCaCertificateHanldler)

	// Return info about the server
	http.HandleFunc("/stats", getHardwareData)

	// Return the san server configuration.
	http.HandleFunc("/get_san_conf", getSanConfigurationHandler)

	// Handle the signing certificate function.
	http.HandleFunc("/sign_ca_certificate", signCaCertificateHandler)

	// Start listen for http request.
	http.HandleFunc("/", ServeFileHandler)

	// The file upload handler.
	http.HandleFunc("/uploads", FileUploadHandler)

	g.path, _ = filepath.Abs(filepath.Dir(os.Args[0]))

	g.initDirectories()

	return g
}

func (globule *Globule) registerAdminAccount() error {
	resource_client_, err := GetResourceClient(globule.getAddress())
	if err != nil {
		return err
	}

	// Create the admin account.
	err = resource_client_.RegisterAccount(globule.getDomain(), "sa", "sa", globule.AdminEmail, globule.RootPassword, globule.RootPassword)
	if err != nil {
		fmt.Println("fail to create admin user")
		return err
	}

	// Set admin role to that account.
	err = resource_client_.AddAccountRole("sa", "admin")
	if err != nil {
		fmt.Println("fail to create admin role")
		return err
	}

	fmt.Println("Admin User create!")
	return nil
}

/**
 * Return globular configuration.
 */
func (globule *Globule) getConfig() map[string]interface{} {

	// TODO filter unwanted attributes...
	config_, _ := Utility.ToMap(globule)
	config_["Domain"], _ = config.GetDomain()
	config_["Name"], _ = config.GetHostName()
	services, _ := config_client.GetServicesConfigurations()

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
		s["CertAuthorityTrust"] = services[i]["CertAuthorityTrust"]
		s["CertFile"] = services[i]["CertFile"]
		s["KeyFile"] = services[i]["KeyFile"]
		s["ConfigPath"] = services[i]["ConfigPath"]
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
						fmt.Println("The domain localhost cannot be use with https, domain must contain dot's")
					} else {

						hasProtocolChange := globule.Protocol != config["Protocol"].(string)
						hasDomainChange := globule.getDomain() != config["Domain"].(string)
						certificateChange := globule.CertificateAuthorityBundle != config["CertificateAuthorityBundle"].(string)
						json.Unmarshal(file, &globule)

						// stop the http server
						ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
						if globule.http_server != nil {
							if err = globule.http_server.Shutdown(ctx); err != nil {
								fmt.Println("fail to stop the http server with error ", err)
							}

							if globule.https_server != nil {
								if err := globule.https_server.Shutdown(ctx); err != nil {
									fmt.Println("fail to stop the https server with error ", err)
								}
							}
						}

						// restart
						globule.serve()

						if hasProtocolChange || hasDomainChange || certificateChange {
							// stop services...
							fmt.Println("Stop gRpc Services")

							err := globule.stopServices()
							if err != nil {
								log.Panicln(err)
							}

							// restart it...
							err = globule.startServices()
							if err != nil {
								log.Panicln(err)
							}

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
		fmt.Println("Let's encrypt dns challenge...")
		dns_client_, err := dns_client.NewDnsService_Client(globule.DNS[0].(string), "dns.DnsService")
		if err != nil {
			return err
		}

		token, err := security.GenerateToken(globule.SessionTimeout, dns_client_.GetMac(), "sa", "", globule.AdminEmail)

		if err != nil {
			fmt.Println("fail to connect with the dns server")
			return err
		}

		err = dns_client_.SetText(token, key, []string{value}, 30)

		if err != nil {
			fmt.Println("fail to set let's encrypt dns chalenge key with error ", err)
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

		token, err := security.GenerateToken(globule.SessionTimeout, dns_client_.GetMac(), "sa", "", globule.AdminEmail)

		if err != nil {

			return err
		}

		err = dns_client_.RemoveText(token, key)
		if err != nil {
			fmt.Println("fail to remove challenge key with error ", err)
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
	config_ := lego.NewConfig(globule)
	config_.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(config_)
	if err != nil {

		return err
	}
	// Dns registration will be use in case dns service are available.
	// TODO dns challenge give JWS has invalid anti-replay nonce error... at the moment
	// http chanllenge do the job but wildcald domain name are not allowed...
	if len(globule.DNS) > 0 {

		// Get the local token.
		dns_client_, err := dns_client.NewDnsService_Client(globule.DNS[0].(string), "dns.DnsService")
		if err != nil {
			return err
		}
		defer dns_client_.Close()

		token, err := security.GenerateToken(globule.SessionTimeout, dns_client_.GetMac(), "sa", "", globule.AdminEmail)
		if err != nil {
			return err
		}

		globularDNS, err := NewDNSProviderGlobularDNS(token)
		if err != nil {
			fmt.Println("fail to create new Dns provider")
			return err
		}

		client.Challenge.SetDNS01Provider(globularDNS)

	} else {
		err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", strconv.Itoa(globule.PortHttp)))
		if err != nil {

			fmt.Println(err)
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

	// initilayse configurations...
	// it must be call here in order to initialyse a sync map...
	config_client.GetServicesConfigurations()

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
		fmt.Println("fail to create configuration directory  with error", err)
		return err
	}

	// Create the tokens directory
	err = Utility.CreateDirIfNotExist(globule.config + "/tokens")
	if err != nil {
		fmt.Println("fail to create tokens directory  with error", err)
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
			fmt.Println("fail to init configuation with error ", err)
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

	// I will put the domain into the
	if globule.AlternateDomains == nil {
		globule.AlternateDomains = make([]interface{}, 0)
		globule.AlternateDomains = append(globule.AlternateDomains, globule.getDomain())
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

	return nil
}

/**
 * Here I will start the services manager who will start all microservices
 * installed on that computer.
 */
func (globule *Globule) startServices() error {

	// Here I will generate the keys for this server if not already exist.
	security.GeneratePeerKeys(Utility.MyMacAddr())

	// This is the local token...
	tokenString, err := security.GenerateToken(globule.SessionTimeout, Utility.MyMacAddr(), "sa", "sa", globule.AdminEmail)
	if err != nil {
		fmt.Println("fail to generate token with error: ", err)
		return err
	}

	err = ioutil.WriteFile(globule.config+"/tokens/"+globule.getDomain()+"_token", []byte(tokenString), 0644)
	if err != nil {
		return err
	}

	// Retreive all configurations
	services, err := config.GetOrderedServicesConfigurations()
	if err != nil {
		return err
	}

	// Register that peer with the dns.
	err = globule.registerIpToDns()
	if err != nil {
		fmt.Println("Fail to write Ip to hosts. ", err)
		return err
	}

	// I will try to get the services manager configuration from the
	// services configurations list.
	for i := 0; i < len(services); i++ {
		services[i]["State"] = "starting"
		config_client.SaveServiceConfiguration(services[i])

		if err != nil {
			fmt.Println("fail to save service configuration with error ", err)
		} else if (len(globule.Certificate) > 0 && globule.Protocol == "https") || (globule.Protocol == "http") {

			// Create the service process.
			_, err = process.StartServiceProcess(services[i], globule.PortsRange)

			if err != nil {
				fmt.Println("fail to start service ", services[i]["Name"], err)

			} else {
				_, err = process.StartServiceProxyProcess(services[i], globule.CertificateAuthorityBundle, globule.Certificate, globule.PortsRange, Utility.ToInt(services[i]["Process"]))
				if err != nil {
					fmt.Println("fail to start service proxy ", services[i]["Name"], err)
				}
			}
		}
	}

	// So here I will register services permissions.
	for i := 0; i < len(services); i++ {
		s := services[i]
		if s["Permissions"] != nil {
			permissions := s["Permissions"].([]interface{})
			for j := 0; j < len(permissions); j++ {
				globule.setActionResourcesPermissions(permissions[j].(map[string]interface{}))
			}
		}
	}

	// Here I will listen for logger event...
	go func() {
		// subscribe to log events
		globule.subscribe("new_log_evt", logListener(globule))

	}()

	// Convert video file, set permissions...
	go func() {
		// So here I will remove files that with no subject...
		files, err := ioutil.ReadDir(config.GetDataDir() + "/files/users")
		for i := 0; i < len(files); i++ {
			f := files[i]
			if f.Name() != "sa" {
				if !globule.accountExist(f.Name()) {
					// Here I will delete the permissission...
					globule.deleteResourcePermissions("/users/" + f.Name())
					os.RemoveAll(config.GetDataDir() + "/files/users/" + f.Name())
				}
			}
		}

		if err != nil {
			fmt.Println(err)
		}
		processFiles() // Process files...
	}()

	// Start process monitoring with prometheus.
	process.StartProcessMonitoring(globule.PortHttp, globule.exit)

	return nil
}

/**
 * Update peers list.
 */
func updatePeersEvent(evt *eventpb.Event) {

	// re-init the list of peers.
	globule.initPeers()

}

/**
 * Here I will init the list of peers.
 */
func (globule *Globule) initPeers() error {

	resource_client_, err := GetResourceClient(globule.getAddress())
	if err != nil {
		return err
	}

	// Return the registered peers
	peers, err := resource_client_.GetPeers(`{}`)
	if err != nil {
		return err
	}

	// Now I will set peers in the host file.
	for i := 0; i < len(peers); i++ {
		// Here I will try to set the peer ip...
		if Utility.IsLocal(peers[i].Address) {
			globule.setHost(peers[i].LocalIpAddress, peers[i].Domain)
		} else {
			globule.setHost(peers[i].ExternalIpAddress, peers[i].Domain)
		}
	}

	globule.peers = peers

	// Subscribe to new peers event...
	globule.subscribe("update_peers_evt", updatePeersEvent)

	return nil
}

// func (globule *Globule) getHttpClient

/**
 * Here I will create application backend connection.
 */
func (globule *Globule) createApplicationConnection() error {
	resource_client_, err := GetResourceClient(globule.getAddress())
	if err != nil {
		return err
	}

	persistence_client_, err := GetPersistenceClient(globule.getAddress())
	if err != nil {
		return err
	}

	if err == nil {
		applications, err := resource_client_.GetApplications("{}")
		if err == nil {
			for i := 0; i < len(applications); i++ {
				app := applications[i]
				err := persistence_client_.CreateConnection(app.Id, app.Id+"_db", globule.getDomain(), 27017, 0, app.Id, app.Password, 500, "", true)
				if err != nil {
					fmt.Println("fail to create application connection  : ", app.Id, err)
				}
			}
		}
	}

	return err
}

/**
 * Stop all services.
 */
func (globule *Globule) stopServices() error {

	services, err := config_client.GetServicesConfigurations()
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

	// Create the admin account.
	globule.registerAdminAccount()

	// Create application connection
	globule.createApplicationConnection()

	// Initialise the list of peers...
	go func() {
		globule.initPeers()
	}()

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

	fmt.Println("globular version " + globule.Version + " build " + Utility.ToString(globule.Build) + " listen at address " + url)
	fmt.Printf("startup took %s", elapsed)

	return nil

}

/**
 * Start serving the content.
 */
func (globule *Globule) Serve() error {

	// So here if another instance of the server exist I will kill it.
	pids, err := Utility.GetProcessIdsByName("Globular")
	if err == nil {
		for i := 0; i < len(pids); i++ {
			if pids[i] != os.Getpid() {
				Utility.TerminateProcess(pids[i], 0)
			}
		}
	}

	// Initialyse directories.
	globule.initDirectories()

	// start listen to http(s)
	// service must be able to get their configuration via http...
	err = globule.Listen()
	if err != nil {
		return err
	}

	// Start microservice manager.
	globule.startServices()

	globule.watchConfig()

	// Set the fmt information in case of crash...
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	return globule.serve()
}

/**
 * Return the domain of the Globule. The name can be empty. If the name is empty
 * it will mean that the domain is entirely control by the globule so it must be
 * able to do it own validation, other wise the domain validation will be done by
 * the globule asscosiate with that domain.
 */
func (globule *Globule) getDomain() string {
	domain, _ := config.GetDomain()

	// if no hostname or domain are found it will be use as localhost.
	return domain
}

/**
 * Return the globule address
 */
func (globule *Globule) getAddress() string {
	address, _ := config.GetAddress() // return the address with the http port.
	return address
}

func (globule *Globule) setHost(ipv4, domain string) error {
	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		return err
	}

	hosts.AddHost(ipv4, domain)
	return hosts.Save()
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
			// Here I will set dns in the resolv.conf file
			resolv_conf := "# That file was generated by globular at server startup. To reset to it original move the file resolv.conf_ to resolv.conf\n"
			resolv_conf += "nameserver 8.8.8.8\n"

			for i := 0; i < len(globule.DNS); i++ {
				dns_client_, err := dns_client.NewDnsService_Client(globule.DNS[i].(string), "dns.DnsService")
				if err != nil {
					return err
				}
				defer dns_client_.Close()

				ipv4, err := Utility.GetIpv4(globule.DNS[i].(string))
				if err == nil {
					resolv_conf += "nameserver " + ipv4 + "\n"
				}

				// Here the token must be generated for the dns server...
				// That peer must be register on the dns to be able to generate a valid token.
				token, err := security.GenerateToken(globule.SessionTimeout, dns_client_.GetMac(), "sa", "", globule.AdminEmail)
				if err != nil {
					return err
				}

				// if the dns address is a local address i will register the local ip...
				if Utility.IsLocal(globule.DNS[i].(string)) {
					_, err = dns_client_.SetA(token, globule.getDomain(), Utility.MyLocalIP(), 60)
				} else {
					_, err = dns_client_.SetA(token, globule.getDomain(), Utility.MyIP(), 60)
				}

				for j := 0; j < len(globule.AlternateDomains); j++ {
					if Utility.IsLocal(globule.DNS[i].(string)) {
						_, err = dns_client_.SetA(token, globule.AlternateDomains[j].(string), Utility.MyLocalIP(), 60)
					} else {
						_, err = dns_client_.SetA(token, globule.AlternateDomains[j].(string), Utility.MyIP(), 60)
					}
				}

				if err != nil {
					return err
				}

			}

			// save the file to /etc/resolv.conf
			if Utility.Exists("/etc/resolv.conf") {
				Utility.MoveFile("/etc/resolv.conf", "/etc/resolv.conf_")
				Utility.WriteStringToFile("/etc/resolv.conf", resolv_conf)
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

	return globule.setHost(Utility.MyLocalIP(), globule.getDomain())
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
			fmt.Println(err)
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
							fmt.Println("fail to update globular from " + discovery + " with error " + err.Error())
						}

					}
				}
			}

			// The time here can be set to higher value.
			time.Sleep(time.Duration(globule.WatchUpdateDelay) * time.Second)
		}
	}()
}

// Try to display application message in a nice way
func logListener(g *Globule) func(evt *eventpb.Event) {
	return func(evt *eventpb.Event) {
		info := make(map[string]interface{})
		err := json.Unmarshal(evt.Data, &info)
		if err == nil {
			// So here Will display message
			var header string
			if info["application"] != nil {
				header = info["application"].(string)
			}

			occurences := info["occurences"].([]interface{})
			occurence := occurences[len(occurences)-1].(map[string]interface{})

			// Set the occurence date.
			messageTime := time.Unix(int64(Utility.ToInt(occurence["date"])), 0)
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

				// I will also display the message in the system logger.
				if info["level"].(string) == "ERROR_MESSAGE" {
					g.logger.Error(msg)
				}else if info["level"].(string) == "WARNING_MESSAGE" {
					g.logger.Warning(msg)
				}else if info["level"].(string) == "INFO_MESSAGE" {
					g.logger.Info(msg)
				}
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

		// start / restart services
		fmt.Println("Succed to receive certificates you need to restart the server...")
		os.Exit(0)

		//globule.startServices()
	}

	// Must be started before other services.
	go func() {
		// local - non secure connection.
		globule.http_server = &http.Server{
			Addr: "0.0.0.0:" + strconv.Itoa(globule.PortHttp),
		}
		err = globule.http_server.ListenAndServe()
	}()

	// Start the http server.
	if globule.Protocol == "https" {

		globule.https_server = &http.Server{
			Addr: "0.0.0.0:" + strconv.Itoa(globule.PortHttps),
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
			resource_client_ = nil
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
			persistence_client_ = nil
			fmt.Println("fail to get persistence client with error ", err)
			return nil, err
		}

	}

	return persistence_client_, nil
}

/**
 * Return an application with a given id
 */
func (globule *Globule) getAccount(accountId string) (*resourcepb.Account, error) {
	resourceClient, err := GetResourceClient(globule.getAddress())
	if err != nil {
		return nil, err
	}

	return resourceClient.GetAccount(accountId)
}

func (globule *Globule) accountExist(id string) bool {
	a, err := globule.getAccount(id)
	if err != nil || a == nil {
		return false
	}
	return true
}

/**
 * Return a group with a given id
 */
func (globule *Globule) getGroup(groupId string) (*resourcepb.Group, error) {
	resourceClient, err := GetResourceClient(globule.getAddress())
	if err != nil {
		return nil, err
	}

	groups, err := resourceClient.GetGroups(`{"$or":[{"_id":"` + groupId + `"},{"name":"` + groupId + `"} ]}`)
	if err != nil {
		return nil, err
	}

	if len(groups) == 0 {
		return nil, errors.New("no group found wiht name or _id " + groupId)
	}

	return groups[0], nil
}

/**
 * Test if a group exist.
 */
func (globule *Globule) groupExist(id string) bool {
	g, err := globule.getGroup(id)
	if err != nil || g == nil {
		return false
	}
	return true
}

/**
 * Return an application with a given id
 */
func (globule *Globule) getApplication(applicationId string) (*resourcepb.Application, error) {
	resourceClient, err := GetResourceClient(globule.getAddress())
	if err != nil {
		return nil, err
	}

	applications, err := resourceClient.GetApplications(`{"$or":[{"_id":"` + applicationId + `"},{"name":"` + applicationId + `"} ]}`)
	if err != nil {
		return nil, err
	}

	if len(applications) == 0 {
		return nil, errors.New("no application found wiht name or _id " + applicationId)
	}

	return applications[0], nil
}

/**
 * Test if a application exist.
 */
func (globule *Globule) applicationExist(id string) bool {
	g, err := globule.getApplication(id)
	if err != nil || g == nil {
		return false
	}
	return true
}

/**
 * Return a peer with a given id
 */
func (globule *Globule) getPeer(peerId string) (*resourcepb.Peer, error) {
	resourceClient, err := GetResourceClient(globule.getAddress())
	if err != nil {
		return nil, err
	}

	peers, err := resourceClient.GetPeers(`{"$or":[{"domain":"` + peerId + `"},{"mac":"` + peerId + `"} ]}`)
	if err != nil {
		return nil, err
	}

	if len(peers) == 0 {
		return nil, errors.New("no peer found wiht name or _id " + peerId)
	}

	return peers[0], nil
}

/**
 * Test if a peer exist.
 */
func (globule *Globule) peerExist(id string) bool {
	g, err := globule.getPeer(id)
	if err != nil || g == nil {
		return false
	}
	return true
}

/**
 * Return a peer with a given id
 */
func (globule *Globule) getOrganization(organisationId string) (*resourcepb.Organization, error) {
	resourceClient, err := GetResourceClient(globule.getAddress())
	if err != nil {
		return nil, err
	}

	organisations, err := resourceClient.GetOrganizations(`{"$or":[{"_id":"` + organisationId + `"},{"name":"` + organisationId + `"} ]}`)
	if err != nil {
		return nil, err
	}

	if len(organisations) == 0 {
		return nil, errors.New("no organization found wiht name or _id " + organisationId)
	}

	return organisations[0], nil
}

/**
 * Test if a organisation exist.
 */
func (globule *Globule) organisationExist(id string) bool {
	o, err := globule.getOrganization(id)
	if err != nil || o == nil {
		return false
	}
	return true
}

/**
 * Return a role with a given id
 */
func (globule *Globule) getRole(roleId string) (*resourcepb.Role, error) {
	resourceClient, err := GetResourceClient(globule.getAddress())
	if err != nil {
		return nil, err
	}

	roles, err := resourceClient.GetRoles(`{"$or":[{"_id":"` + roleId + `"},{"name":"` + roleId + `"} ]}`)
	if err != nil {
		return nil, err
	}

	if len(roles) == 0 {
		return nil, errors.New("no role found wiht name or _id " + roleId)
	}

	return roles[0], nil
}

/**
 * Test if a role exist.
 */
func (globule *Globule) roleExist(id string) bool {
	r, err := globule.getRole(id)
	if err != nil || r == nil {
		return false
	}
	return true
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
			rbac_client_ = nil
			return nil, err
		}
	}

	return rbac_client_, nil
}

// Use rbac client here...
func (globule *Globule) addResourceOwner(path string, subject string, subjectType rbacpb.SubjectType) error {

	rbac_client_, err := GetRbacClient(globule.getAddress())
	if err != nil {
		return err
	}
	return rbac_client_.AddResourceOwner(path, subject, subjectType)
}

func (globule *Globule) validateAction(method string, subject string, subjectType rbacpb.SubjectType, infos []*rbacpb.ResourceInfos) (bool, error) {
	rbac_client_, err := GetRbacClient(globule.getAddress())
	if err != nil {
		return false, err
	}

	return rbac_client_.ValidateAction(method, subject, subjectType, infos)
}

func (globule *Globule) validateAccess(subject string, subjectType rbacpb.SubjectType, name string, path string) (bool, bool, error) {
	rbac_client_, err := GetRbacClient(globule.getAddress())
	if err != nil {
		return false, false, err
	}
	hasAccess, hasAccessDenied, err := rbac_client_.ValidateAccess(subject, subjectType, name, path)
	return hasAccess, hasAccessDenied, err
}

func (globule *Globule) setActionResourcesPermissions(permissions map[string]interface{}) error {
	rbac_client_, err := GetRbacClient(globule.getAddress())
	if err != nil {
		return err
	}
	return rbac_client_.SetActionResourcesPermissions(permissions)
}

func (globule *Globule) deleteResourcePermissions(path string) error {
	rbac_client_, err := GetRbacClient(globule.getAddress())
	if err != nil {
		return err
	}
	return rbac_client_.DeleteResourcePermissions(path)
}

///////////////////// event service functions ////////////////////////////////////
func (globule *Globule) getEventClient() (*event_client.Event_Client, error) {
	var err error
	if event_client_ != nil {
		return event_client_, nil
	}
	event_client_, err = event_client.NewEventService_Client(globule.getAddress(), "event.EventService")
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

	err = eventClient.Subscribe(evt, globule.Name, listener)
	// register a listener...
	return err
}

///////////////////////  fmt Services functions ////////////////////////////////////////////////

/**
 * Get the fmt client.
 */
func (globule *Globule) GetLogClient() (*log_client.Log_Client, error) {
	var err error
	if log_client_ == nil {
		log_client_, err = log_client.NewLogService_Client(globule.getAddress(), "log.LogService")
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
	log_client_.Log(globule.Name, globule.getAddress(), functionName, level, message, fileLine, functionName)
}

//////////////////////////////////////////////////////////////////////////////
// Process files.
//////////////////////////////////////////////////////////////////////////////
func processFiles() {
	convertVideo()

	// sleep a minute...
	time.Sleep(1 * time.Minute)

	processFiles()
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
		return
	}

	// Also convert video from public file...
	for i := 0; i < len(config.GetPublicDirs()); i++ {
		err = filepath.Walk(config.GetPublicDirs()[i], visit(&files))
		if err != nil {
			return
		}
	}

	for _, file := range files {
		file = strings.ReplaceAll(file, "\\", "/")
		createVideoStream(file)
	}
}

// Set file indexation to be able to search text file on the server.
func indexFile(path string, fileType string) error {
	return nil
}

func getStreamInfos(path string) (map[string]interface{}, error) {

	cmd := exec.Command("ffprobe", "-v", "error", "-show_format", "-show_streams", "-print_format", "json", path)
	data, _ := cmd.CombinedOutput()
	infos := make(map[string]interface{})
	err := json.Unmarshal(data, &infos)
	if err != nil {
		return nil, err
	}

	return infos, nil
}

/**
 * Convert all kind of video to mp4 so all browser will be able to read it.
 */
func createVideoStream(path string) error {

	path_ := path[0:strings.LastIndex(path, "/")]
	name_ := path[strings.LastIndex(path, "/"):strings.LastIndex(path, ".")]
	output := path_ + "/" + name_ + ".mp4"

	defer os.RemoveAll(path)

	// Test if cuda is available.
	getVersion := exec.Command("ffmpeg", "-version")
	version, _ := getVersion.CombinedOutput()

	var cmd *exec.Cmd

	streamInfos, err := getStreamInfos(path)
	if err != nil {
		return err
	}

	// Here I will test if the encoding is valid
	encoding := streamInfos["streams"].([]interface{})[0].(map[string]interface{})["codec_long_name"].(string)

	//  https://docs.nvidia.com/video-technologies/video-codec-sdk/ffmpeg-with-nvidia-gpu/
	if strings.Index(string(version), "--enable-cuda-nvcc") > -1 {
		fmt.Println("use gpu for convert ", path)
		if strings.HasPrefix(encoding, "H.264") || strings.HasPrefix(encoding, "MPEG-4 part 2") {
			cmd = exec.Command("ffmpeg", "-i", path, "-c:v", "h264_nvenc", "-c:a", "aac", output)
		} else if strings.HasPrefix(encoding, "H.265") {
			// in future when all browser will support H.265 I will compile it with this line instead.
			//cmd = exec.Command("ffmpeg", "-i", path, "-c:v", "hevc_nvenc",  "-c:a", "aac", output)
			cmd = exec.Command("ffmpeg", "-i", path, "-c:v", "h264_nvenc", "-c:a", "aac", "-pix_fmt", "yuv420p", output)

		} else {
			err := errors.New("no encoding command foud for " + encoding)
			return err
		}

	} else {
		fmt.Println("use cpu for convert ", path)
		// ffmpeg -i input.mkv -c:v libx264 -c:a aac output.mp4
		if strings.HasPrefix(encoding, "H.264") {
			cmd = exec.Command("ffmpeg", "-i", path, "-c:v", "libx264", "-c:a", "aac", output)
		} else if strings.HasPrefix(encoding, "H.265") {
			// in future when all browser will support H.265 I will compile it with this line instead.
			// cmd = exec.Command("ffmpeg", "-i", path, "-c:v", "libx265", "-c:a", "aac", output)
			cmd = exec.Command("ffmpeg", "-i", path, "-c:v", "libx264", "-c:a", "aac", "-pix_fmt", "yuv420p", output)
		} else {
			err := errors.New("no encoding command foud for " + encoding)
			fmt.Println(err.Error())
			return err
		}
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	// Create a video preview
	return createVideoPreview(output, 20, 128)
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
		return err
	}

	path_ = strings.ReplaceAll(path, globule.data+"/files", "")
	path_ = path_[0:strings.LastIndex(path_, "/")]

	globule.publish("reload_dir_event", []byte(path_))

	return nil
}

func getVideoDuration(path string) float64 {
	// original command...
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

	duration, _ := strconv.ParseFloat(strings.TrimSpace(out.String()), 64)

	return duration
}
