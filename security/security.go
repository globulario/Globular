package security

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/davecourtois/Utility"
)

// That function will be access via http so event server or client will be able
// to get particular service configuration.
func GetClientConfig(address string, name string, port int, path string) (map[string]interface{}, error) {

	var serverConfig map[string]interface{}
	var config map[string]interface{}
	var err error

	if len(address) == 0 {
		err := errors.New("no address was given for service name " + name)
		return nil, err
	}

	// In case of local service I will get the service value directly from
	// the configuration file.
	serverConfig, err = getLocalConfig()

	isLocal := true
	if err == nil {
		if serverConfig["Domain"] != address {
			isLocal = false
		}
	} else {
		isLocal = false
	}

	if !isLocal {
		// First I will retreive the server configuration.
		serverConfig, err = getRemoteConfig(address, port)
		if err != nil {
			return nil, err
		}
	}

	// get service by id or by name... (take the first service with a given name in case of name.
	for _, s := range serverConfig["Services"].(map[string]interface{}) {
		if s.(map[string]interface{})["Name"].(string) == name || s.(map[string]interface{})["Id"].(string) == name {
			config = s.(map[string]interface{})
			break
		}
	}

	// No service with name or id was found...
	if config == nil {
		return nil, errors.New("No service found whit name " + name + " exist on the server.")
	}

	// Set the config tls...
	config["TLS"] = serverConfig["Protocol"].(string) == "https"
	config["Domain"] = address

	// get / init credential values.
	if config["TLS"] == false {
		// set the credential function here
		config["KeyFile"] = ""
		config["CertFile"] = ""
		config["CertAuthorityTrust"] = ""
	} else {
		// Here I will retreive the credential or create it if not exist.
		var country string
		if serverConfig["Country"] != nil {
			country = serverConfig["Country"].(string)
		}

		var state string
		if serverConfig["State"] != nil {
			state = serverConfig["State"].(string)
		}

		var city string
		if serverConfig["City"] != nil {
			city = serverConfig["City"].(string)
		}

		var organization string
		if serverConfig["Organization"] != nil {
			state = serverConfig["Organization"].(string)
		}

		var alternateDomains []interface{}
		if serverConfig["AlternateDomains"] != nil {
			alternateDomains = serverConfig["AlternateDomains"].([]interface{})
		}

		if !isLocal {

			keyPath, certPath, caPath, err := getCredentialConfig(path, serverConfig["Domain"].(string), country, state, city, organization, alternateDomains, port)
			if err != nil {
				return nil, err
			}
			// set the credential function here
			config["KeyFile"] = keyPath
			config["CertFile"] = certPath
			config["CertAuthorityTrust"] = caPath
		}
	}
	return config, nil
}

func InstallCertificates(domain string, port int, path string) (string, string, string, error) {
	return getCredentialConfig(path, domain, "", "", "", "", []interface{}{}, port)
}

/**
 * Return the server local configuration if one exist.
 */
func getLocalConfig() (map[string]interface{}, error) {

	if !Utility.Exists(os.TempDir() + "/GLOBULAR_ROOT") {
		return nil, errors.New("No local Globular instance found!")
	}

	config := make(map[string]interface{}, 0)
	root, err := ioutil.ReadFile(os.TempDir() + "/GLOBULAR_ROOT")
	if err != nil {
		return nil, err
	}

	index := strings.LastIndex(string(root), ":")
	if index == -1 {
		return nil, errors.New("File contain does not contain ':' separator! ")
	}

	root_ := string(root)[0:index]
	data, err := ioutil.ReadFile(root_ + "/config/config.json")

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

/**
 * Get the remote client configuration.
 */
func getRemoteConfig(address string, port int) (map[string]interface{}, error) {

	// Here I will get the configuration information from http...
	var resp *http.Response
	var err error
	var configAddress = "http://" + address + ":" + Utility.ToString(port) + "/config"
	resp, err = http.Get(configAddress)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var config map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

/**
 * Get the ca certificate
 */
func getCaCertificate(address string, port int) (string, error) {

	if len(address) == 0 {
		return "", errors.New("No address was given!")
	}

	// Here I will get the configuration information from http...
	var resp *http.Response
	var err error
	var caAddress = "http://" + address + ":" + Utility.ToString(port) + "/get_ca_certificate"
	log.Println("get ca certificate from ", caAddress)
	resp, err = http.Get(caAddress)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusCreated {

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(bodyBytes))
		return string(bodyBytes), nil
	}

	return "", errors.New("fail to retreive ca certificate with error " + Utility.ToString(resp.StatusCode))
}

func signCaCertificate(address string, csr string, port int) (string, error) {

	if len(address) == 0 {
		return "", errors.New("No address was given!")
	}

	csr_str := base64.StdEncoding.EncodeToString([]byte(csr))
	// Here I will get the configuration information from http...
	var resp *http.Response
	var err error

	var signCertificateAddress = "http://" + address + ":" + Utility.ToString(port) + "/sign_ca_certificate"
	log.Println("----------> 228: ", signCertificateAddress)

	resp, err = http.Get(signCertificateAddress + "?csr=" + csr_str)
	if err != nil {
		log.Println("---> 232")
		return "", err
	}

	log.Println("----------> 236: ", signCertificateAddress)
	defer resp.Body.Close()
	log.Println("---> 232")
	if resp.StatusCode == http.StatusCreated {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("certificate are now signed!")
		return string(bodyBytes), nil
	}

	return "", errors.New("fail to sign ca certificate with error " + Utility.ToString(resp.StatusCode))
}

/**
 * Return the credential configuration.
 */
func getCredentialConfig(basePath string, address string, country string, state string, city string, organization string, alternateDomains []interface{}, port int) (keyPath string, certPath string, caPath string, err error) {

	// TODO Clarify the use of the password here.
	pwd := "1111"

	// use the temp dir to store the certificate in that case.
	path := basePath + "/" + "config" + "/" + "tls"

	// must have write access of file.
	_, err = ioutil.ReadFile(path + "/" + address + "/" + "client.pem")
	if err != nil {
		path = basePath + "/" + "config" + "/" + "tls"
		err = nil
	}

	// Create a new directory to put the credential.
	creds := path + "/" + address

	// Return the existing paths...
	if Utility.Exists(creds) &&
		Utility.Exists(creds+"/"+"client.pem") &&
		Utility.Exists(creds+"/"+"client.crt") &&
		Utility.Exists(creds+"/"+"ca.crt") {
		info, _ := os.Stat(creds)

		// test if the certificate are older than 5 mount.
		if info.ModTime().Add(24*30*5*time.Hour).Unix() < time.Now().Unix() {
			os.RemoveAll(creds)
		} else {

			keyPath = creds + "/" + "client.pem"
			certPath = creds + "/" + "client.crt"
			caPath = creds + "/" + "ca.crt"
			return
		}
	}

	Utility.CreateDirIfNotExist(creds)

	// I will connect to the certificate authority of the server where the application must
	// be deployed. Certificate autority run wihtout tls.

	// Get the ca.crt certificate.
	ca_crt, err := getCaCertificate(address, port)
	if err != nil {
		return "", "", "", err
	}

	// Write the ca.crt file on the disk
	err = ioutil.WriteFile(creds+"/"+"ca.crt", []byte(ca_crt), 0664)
	if err != nil {
		return "", "", "", err
	}

	// Now I will generate the certificate for the client...
	// Step 1: Generate client private key.
	err = GenerateClientPrivateKey(creds, pwd)
	if err != nil {
		return "", "", "", err
	}

	alternateDomains_ := make([]string, 0)
	alternateDomains_ = append(alternateDomains_, address)
	for i := 0; i < len(alternateDomains); i++ {
		alternateDomains_ = append(alternateDomains_, alternateDomains[i].(string))
	}

	// generate the SAN file
	err = GenerateSanConfig(creds, country, state, city, organization, alternateDomains_)

	// Step 2: Generate the client signing request.
	err = GenerateClientCertificateSigningRequest(creds, pwd, address)
	if err != nil {
		return "", "", "", err
	}

	// Step 3: Generate client signed certificate.
	client_csr, err := ioutil.ReadFile(creds + "/" + "client.csr")
	if err != nil {
		return "", "", "", err
	}

	// Sign the certificate from the server ca...
	client_crt, err := signCaCertificate(address, string(client_csr), Utility.ToInt(port))
	if err != nil {
		return "", "", "", err
	}

	// Write bact the client certificate in file on the disk
	err = ioutil.WriteFile(creds+"/"+"client.crt", []byte(client_crt), 0664)
	if err != nil {
		return "", "", "", err
	}

	// Now ask the ca to sign the certificate.

	// Step 4: Convert to pem format.
	err = KeyToPem("client", creds, pwd)
	if err != nil {
		return "", "", "", err
	}

	// set the credential paths.
	keyPath = creds + "/" + "client.pem"
	certPath = creds + "/" + "client.crt"
	caPath = creds + "/" + "ca.crt"

	return
}

//////////////////////////////// Certificate Authority /////////////////////////

// Generate the Certificate Authority private key file (this shouldn't be shared in real life)
func GenerateAuthorityPrivateKey(path string, pwd string) error {
	if Utility.Exists(path + "/" + "ca.key") {
		return nil
	}

	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "genrsa")
	args = append(args, "-passout")
	args = append(args, "pass:"+pwd)
	args = append(args, "-des3")
	args = append(args, "-out")
	args = append(args, path+"/"+"ca.key")
	args = append(args, "4096")

	err := exec.Command(cmd, args...).Run()
	if err != nil || !Utility.Exists(path+"/"+"ca.key") {
		return errors.New("Fail to generate the Authority private key")
	}
	return nil
}

// Certificate Authority trust certificate (this should be shared whit users)
func GenerateAuthorityTrustCertificate(path string, pwd string, expiration_delay int, domain string) error {
	if Utility.Exists(path + "/" + "ca.crt") {
		return nil
	}

	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "req")
	args = append(args, "-passin")
	args = append(args, "pass:"+pwd)
	args = append(args, "-new")
	args = append(args, "-x509")
	args = append(args, "-days")
	args = append(args, strconv.Itoa(expiration_delay))
	args = append(args, "-key")
	args = append(args, path+"/"+"ca.key")
	args = append(args, "-out")
	args = append(args, path+"/"+"ca.crt")
	args = append(args, "-subj")
	args = append(args, "/CN=Root CA")

	err := exec.Command(cmd, args...).Run()
	if err != nil || !Utility.Exists(path+"/"+"ca.crt") {
		log.Println(err)
		return errors.New("Fail to generate the trust certificate")
	}

	return nil
}

/////////////////////// Server Keys //////////////////////////////////////////

// Server private key, password protected (this shoudn't be shared)
func GenerateSeverPrivateKey(path string, pwd string) error {
	if Utility.Exists(path + "/" + "server.key") {
		return nil
	}
	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "genrsa")
	args = append(args, "-passout")
	args = append(args, "pass:"+pwd)
	args = append(args, "-des3")
	args = append(args, "-out")
	args = append(args, path+"/"+"server.key")
	args = append(args, "4096")

	err := exec.Command(cmd, args...).Run()
	if err != nil || !Utility.Exists(path+"/"+"server.key") {
		log.Println(err)
		return errors.New("Fail to generate server private key")
	}
	return nil
}

// Generate client private key and certificate.
func GenerateClientPrivateKey(path string, pwd string) error {
	if Utility.Exists(path + "/" + "client.key") {
		return nil
	}

	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "genrsa")
	args = append(args, "-passout")
	args = append(args, "pass:"+pwd)
	args = append(args, "-des3")
	args = append(args, "-out")
	args = append(args, path+"/"+"client.pass.key")
	args = append(args, "4096")

	err := exec.Command(cmd, args...).Run()
	if err != nil || !Utility.Exists(path+"/"+"client.pass.key") {
		log.Println(err)
		return errors.New("Fail to generate client private key " + err.Error())
	}

	args = make([]string, 0)
	args = append(args, "rsa")
	args = append(args, "-passin")
	args = append(args, "pass:"+pwd)
	args = append(args, "-in")
	args = append(args, path+"/"+"client.pass.key")
	args = append(args, "-out")
	args = append(args, path+"/"+"client.key")

	err = exec.Command(cmd, args...).Run()
	if err != nil || !Utility.Exists(path+"/"+"client.key") {
		return errors.New("Fail to generate client private key " + err.Error())
	}

	// Remove the file.
	err = os.Remove(path + "/" + "client.pass.key")
	if err != nil {
		return errors.New("fail to remove intermediate key client.pass.key")
	}
	return nil
}

func GenerateClientCertificateSigningRequest(path string, pwd string, domain string) error {
	if Utility.Exists(path + "/" + "client.csr") {
		return nil
	}

	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "req")
	args = append(args, "-new")
	args = append(args, "-key")
	args = append(args, path+"/"+"client.key")
	args = append(args, "-out")
	args = append(args, path+"/"+"client.csr")
	args = append(args, "-subj")
	args = append(args, "/CN="+domain)
	args = append(args, "-config")
	args = append(args, path+"/"+"san.conf")

	err := exec.Command(cmd, args...).Run()

	if err != nil || !Utility.Exists(path+"/"+"client.csr") {
		log.Println(args)
		return errors.New("Fail to generate client certificate signing request.")
	}

	return nil
}

func GenerateSignedClientCertificate(path string, pwd string, expiration_delay int) error {

	if Utility.Exists(path + "/" + "client.crt") {
		return nil
	}

	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "x509")
	args = append(args, "-req")
	args = append(args, "-passin")
	args = append(args, "pass:"+pwd)
	args = append(args, "-days")
	args = append(args, strconv.Itoa(expiration_delay))
	args = append(args, "-in")
	args = append(args, path+"/"+"client.csr")
	args = append(args, "-CA")
	args = append(args, path+"/"+"ca.crt")
	args = append(args, "-CAkey")
	args = append(args, path+"/"+"ca.key")
	args = append(args, "-set_serial")
	args = append(args, "01")
	args = append(args, "-out")
	args = append(args, path+"/"+"client.crt")
	args = append(args, "-extfile")
	args = append(args, path+"/"+"san.conf")
	args = append(args, "-extensions")
	args = append(args, "v3_req")

	err := exec.Command(cmd, args...).Run()
	if err != nil || !Utility.Exists(path+"/"+"client.crt") {
		log.Println("fail to get the signed server certificate")
	}

	return nil
}

func GenerateSanConfig(path string, country string, state string, city string, organization string, domains []string) error {

	config := fmt.Sprintf(`
[req]
distinguished_name = req_distinguished_name
req_extensions = v3_req
prompt = no

[req_distinguished_name]
C = %s
ST =  %s
L =  %s
O	=  %s
CN =  %s

[v3_req]
# Extensions to add to a certificate request
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
subjectAltName = @alt_names

[alt_names]
`, country, state, city, organization, domains[0])

	// set alternate domain
	for i := 0; i < len(domains); i++ {
		config += fmt.Sprintf("DNS.%d = %s \n", i, domains[i])
	}

	if Utility.Exists(path + "/san.conf") {
		return nil
	}

	f, err := os.Create(path + "/san.conf")
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(config)

	return err
}

// Server certificate signing request (this should be shared with the CA owner)
func GenerateServerCertificateSigningRequest(path string, pwd string, domain string) error {

	if Utility.Exists(path + "/" + "server.crs") {
		return nil
	}

	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "req")
	args = append(args, "-passin")
	args = append(args, "pass:"+pwd)
	args = append(args, "-new")
	args = append(args, "-key")
	args = append(args, path+"/"+"server.key")
	args = append(args, "-out")
	args = append(args, path+"/"+"server.csr")
	args = append(args, "-subj")
	args = append(args, "/CN="+domain)
	args = append(args, "-config")
	args = append(args, path+"/"+"san.conf")

	err := exec.Command(cmd, args...).Run()
	if err != nil || !Utility.Exists(path+"/"+"server.csr") {
		return errors.New("Fail to generate server certificate signing request.")
	}

	return nil
}

// Server certificate signed by the CA (this would be sent back to the client by the CA owner)
func GenerateSignedServerCertificate(path string, pwd string, expiration_delay int) error {

	if Utility.Exists(path + "/" + "server.crt") {
		return nil
	}

	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "x509")
	args = append(args, "-req")
	args = append(args, "-passin")
	args = append(args, "pass:"+pwd)
	args = append(args, "-days")
	args = append(args, strconv.Itoa(expiration_delay))
	args = append(args, "-in")
	args = append(args, path+"/"+"server.csr")
	args = append(args, "-CA")
	args = append(args, path+"/"+"ca.crt")
	args = append(args, "-CAkey")
	args = append(args, path+"/"+"ca.key")
	args = append(args, "-set_serial")
	args = append(args, "01")
	args = append(args, "-out")
	args = append(args, path+"/"+"server.crt")
	args = append(args, "-extfile")
	args = append(args, path+"/"+"san.conf")
	args = append(args, "-extensions")
	args = append(args, "v3_req")

	err := exec.Command(cmd, args...).Run()
	if err != nil || !Utility.Exists(path+"/"+"server.crt") {
		log.Println(args)
		log.Println("fail to get the signed server certificate")

	}

	return nil
}

// Conversion of server.key into a format gRpc likes (this shouldn't be shared)
func KeyToPem(name string, path string, pwd string) error {
	if Utility.Exists(path + "/" + name + ".pem") {
		return nil
	}

	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "pkcs8")
	args = append(args, "-topk8")
	args = append(args, "-nocrypt")
	args = append(args, "-passin")
	args = append(args, "pass:"+pwd)
	args = append(args, "-in")
	args = append(args, path+"/"+name+".key")
	args = append(args, "-out")
	args = append(args, path+"/"+name+".pem")

	err := exec.Command(cmd, args...).Run()
	if err != nil || !Utility.Exists(path+"/"+name+".key") {
		return errors.New("Fail to generate " + name + ".pem key from " + name + ".key")
	}

	return nil
}

/**
 * That function is use to generate services certificates.
 * Private ca.key, server.key, server.pem, server.crt
 * Share ca.crt (needed by the client), server.csr (needed by the CA)
 */
func GenerateServicesCertificates(pwd string, expiration_delay int, domain string, path string, country string, state string, city string, organization string, alternateDomains []interface{}) error {
	if Utility.Exists(path + "/" + "client.crt") {
		return nil // certificate are already created.
	}

	alternateDomains_ := make([]string, len(alternateDomains))
	for i := 0; i < len(alternateDomains); i++ {
		alternateDomains_[i] = alternateDomains[i].(string)
	}
	// Alternate domain must contain the domain (CN=domain)...
	if !Utility.Contains(alternateDomains_, domain) {
		alternateDomains_ = append(alternateDomains_, domain)
	}

	// Generate the SAN configuration.
	err := GenerateSanConfig(path, country, state, city, organization, alternateDomains_)
	if err != nil {
		return err
	}

	log.Println("Step 1: Generate Certificate Authority + Trust Certificate (ca.crt)")
	err = GenerateAuthorityPrivateKey(path, pwd)
	if err != nil {

		return err
	}

	err = GenerateAuthorityTrustCertificate(path, pwd, expiration_delay, domain)
	if err != nil {

		return err
	}

	log.Println("Setp 2: Generate the server Private Key (server.key)")
	err = GenerateSeverPrivateKey(path, pwd)
	if err != nil {
		return err
	}

	log.Println("Setp 3: Get a certificate signing request from the CA (server.csr)")
	err = GenerateServerCertificateSigningRequest(path, pwd, domain)
	if err != nil {
		return err
	}

	log.Println("Step 4: Sign the certificate with the CA we create(it's called self signing) - server.crt")
	err = GenerateSignedServerCertificate(path, pwd, expiration_delay)
	if err != nil {
		return err
	}

	log.Println("Step 5: Convert the server Certificate to .pem format (server.pem) - usable by gRpc")
	err = KeyToPem("server", path, pwd)
	if err != nil {
		return err
	}

	log.Println("Step 6: Generate client private key.")
	err = GenerateClientPrivateKey(path, pwd)
	if err != nil {
		return err
	}

	log.Println("Step 7: Generate the client signing request.")
	err = GenerateClientCertificateSigningRequest(path, pwd, domain)
	if err != nil {
		return err
	}

	log.Println("Step 8: Generate client signed certificate.")
	err = GenerateSignedClientCertificate(path, pwd, expiration_delay)
	if err != nil {
		return err
	}

	log.Println("Step 9: Convert to pem format.")
	err = KeyToPem("client", path, pwd)
	if err != nil {
		return err
	}

	return nil
}
