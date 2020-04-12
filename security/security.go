package security

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/davecourtois/Utility"
)

// That function will be access via http so event server or client will be able
// to get particular service configuration.
func GetClientConfig(address string, name string) (map[string]interface{}, error) {

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
		if config["Domain"] != address {
			isLocal = false
		}
	} else {
		isLocal = false
	}

	if !isLocal {
		// First I will retreive the server configuration.
		serverConfig, err = getRemoteConfig(address)
		if err != nil {
			return nil, err
		}
	}

	if serverConfig["Services"].(map[string]interface{})[name] != nil {
		config = serverConfig["Services"].(map[string]interface{})[name].(map[string]interface{})
	} else {
		return nil, errors.New("No service found whit name " + name + " exist on the server.")
	}

	// Set the config tls...
	config["TLS"] = serverConfig["Protocol"].(string) == "https"

	// get / init credential values.
	if config["TLS"] == false {
		// set the credential function here
		config["KeyFile"] = ""
		config["CertFile"] = ""
		config["CertAuthorityTrust"] = ""
	} else {
		// Here I will retreive the credential or create it if not exist.
		keyPath, certPath, caPath, err := getCredentialConfig(address)
		if err != nil {
			return nil, err
		}
		// set the credential function here
		config["KeyFile"] = keyPath
		config["CertFile"] = certPath
		config["CertAuthorityTrust"] = caPath
	}

	return config, nil
}

/**
 * Return the server local configuration if one exist.
 */
func getLocalConfig() (map[string]interface{}, error) {
	config := make(map[string]interface{}, 0)
	root, _ := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + "GLOBULAR_ROOT")
	data, err := ioutil.ReadFile(string(root) + string(os.PathSeparator) + "config" + string(os.PathSeparator) + "config.json")

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
func getRemoteConfig(address string) (map[string]interface{}, error) {

	// Here I will get the configuration information from http...
	var resp *http.Response
	var err error
	resp, err = http.Get("http://" + address + ":10000/config")
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
func getCaCertificate(address string) (string, error) {
	if len(address) == 0 {
		return "", errors.New("No address was given!")
	}

	// Here I will get the configuration information from http...
	var resp *http.Response
	var err error
	resp, err = http.Get("http://" + address + ":10000/get_ca_certificate")

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

	return "", errors.New("fail to retreive ca certificate with error " + Utility.ToString(resp.StatusCode))
}

func signCaCertificate(address string, csr string) (string, error) {

	if len(address) == 0 {
		return "", errors.New("No address was given!")
	}

	csr_str := base64.StdEncoding.EncodeToString([]byte(csr))
	// Here I will get the configuration information from http...
	var resp *http.Response
	var err error
	resp, err = http.Get("http://" + address + ":10000/sign_ca_certificate?csr=" + csr_str)
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

	return "", errors.New("fail to sign ca certificate with error " + Utility.ToString(resp.StatusCode))
}

/**
 * Return the credential configuration.
 */
func getCredentialConfig(address string) (keyPath string, certPath string, caPath string, err error) {

	root, _ := ioutil.ReadFile(os.TempDir() + string(os.PathSeparator) + "GLOBULAR_ROOT")
	path := string(root) + string(os.PathSeparator) + "config" + string(os.PathSeparator) + "grpc_tls"
	pwd := "1111"

	// Here I will get the local configuration...
	var config map[string]interface{}
	config, err = getLocalConfig()
	isLocal := true
	if err == nil {
		if config["Domain"] != address {
			isLocal = false
		}
	} else {
		isLocal = false
	}

	if isLocal {
		pwd = config["CertPassword"].(string)
		keyPath = path + string(os.PathSeparator) + "client.pem"
		certPath = path + string(os.PathSeparator) + "client.crt"
		caPath = path + string(os.PathSeparator) + "ca.crt"
		return
	} else {
		// use the temp dir to store the certificate in that case.
		path = os.TempDir() + string(os.PathSeparator) + "config" + string(os.PathSeparator) + "grpc_tls"
		err = nil
	}

	// must have write access of file.
	_, err = ioutil.ReadFile(path + string(os.PathSeparator) + address + string(os.PathSeparator) + "client.pem")
	if err != nil {
		path = os.TempDir() + string(os.PathSeparator) + "config" + string(os.PathSeparator) + "grpc_tls"
		err = nil
	}

	// Create a new directory to put the credential.
	creds := path + string(os.PathSeparator) + address

	// Return the existing paths...
	if Utility.Exists(creds) &&
		Utility.Exists(creds+string(os.PathSeparator)+"client.pem") &&
		Utility.Exists(creds+string(os.PathSeparator)+"client.crt") &&
		Utility.Exists(creds+string(os.PathSeparator)+"ca.crt") {
		info, _ := os.Stat(creds)

		// test if the certificate are older than 5 mount.
		if info.ModTime().Add(24*30*5*time.Hour).Unix() < time.Now().Unix() {
			os.RemoveAll(creds)
		} else {

			keyPath = creds + string(os.PathSeparator) + "client.pem"
			certPath = creds + string(os.PathSeparator) + "client.crt"
			caPath = creds + string(os.PathSeparator) + "ca.crt"
			return
		}
	}

	Utility.CreateDirIfNotExist(creds)

	// I will connect to the certificate authority of the server where the application must
	// be deployed. Certificate autority run wihtout tls.

	// Get the ca.crt certificate.
	ca_crt, err := getCaCertificate(address)
	if err != nil {
		log.Println(err)
		return "", "", "", err
	}

	// Write the ca.crt file on the disk
	err = ioutil.WriteFile(creds+string(os.PathSeparator)+"ca.crt", []byte(ca_crt), 0664)
	if err != nil {
		log.Println(err)
		return "", "", "", err
	}

	// Now I will generate the certificate for the client...
	// Step 1: Generate client private key.
	err = GenerateClientPrivateKey(creds, pwd)
	if err != nil {
		log.Println(err)
		return "", "", "", err
	}

	// Step 2: Generate the client signing request.
	err = GenerateClientCertificateSigningRequest(creds, pwd, address)
	if err != nil {
		log.Println(err)
		return "", "", "", err
	}

	// Step 3: Generate client signed certificate.
	client_csr, err := ioutil.ReadFile(creds + string(os.PathSeparator) + "client.csr")
	if err != nil {
		log.Println(err)
		return "", "", "", err
	}

	// Sign the certificate from the server ca...
	client_crt, err := signCaCertificate(address, string(client_csr))
	if err != nil {
		log.Println(err)
		return "", "", "", err
	}

	// Write bact the client certificate in file on the disk
	err = ioutil.WriteFile(creds+string(os.PathSeparator)+"client.crt", []byte(client_crt), 0664)
	if err != nil {
		log.Println(err)
		return "", "", "", err
	}

	// Now ask the ca to sign the certificate.

	// Step 4: Convert to pem format.
	err = KeyToPem("client", creds, pwd)
	if err != nil {
		log.Println(err)
		return "", "", "", err
	}

	// set the credential paths.
	keyPath = creds + string(os.PathSeparator) + "client.pem"
	certPath = creds + string(os.PathSeparator) + "client.crt"
	caPath = creds + string(os.PathSeparator) + "ca.crt"

	return
}

//////////////////////////////// Certificate Authority /////////////////////////

// Generate the Certificate Authority private key file (this shouldn't be shared in real life)
func GenerateAuthorityPrivateKey(path string, pwd string) error {
	if Utility.Exists(path + string(os.PathSeparator) + "ca.key") {
		return nil
	}

	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "genrsa")
	args = append(args, "-passout")
	args = append(args, "pass:"+pwd)
	args = append(args, "-des3")
	args = append(args, "-out")
	args = append(args, path+string(os.PathSeparator)+"ca.key")
	args = append(args, "4096")

	err := exec.Command(cmd, args...).Run()
	if err != nil {
		return errors.New("Fail to generate the Authority private key")
	}
	return nil
}

// Certificate Authority trust certificate (this should be shared whit users)
func GenerateAuthorityTrustCertificate(path string, pwd string, expiration_delay int, domain string) error {
	if Utility.Exists(path + string(os.PathSeparator) + "ca.crt") {
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
	args = append(args, path+string(os.PathSeparator)+"ca.key")
	args = append(args, "-out")
	args = append(args, path+string(os.PathSeparator)+"ca.crt")
	args = append(args, "-subj")
	args = append(args, "/CN="+domain)

	err := exec.Command(cmd, args...).Run()
	if err != nil {
		return errors.New("Fail to generate the trust certificate")
	}

	return nil
}

/////////////////////// Server Keys //////////////////////////////////////////

// Server private key, password protected (this shoudn't be shared)
func GenerateSeverPrivateKey(path string, pwd string) error {
	if Utility.Exists(path + string(os.PathSeparator) + "server.key") {
		return nil
	}
	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "genrsa")
	args = append(args, "-passout")
	args = append(args, "pass:"+pwd)
	args = append(args, "-des3")
	args = append(args, "-out")
	args = append(args, path+string(os.PathSeparator)+"server.key")
	args = append(args, "4096")

	err := exec.Command(cmd, args...).Run()
	if err != nil {
		return errors.New("Fail to generate server private key")
	}
	return nil
}

// Generate client private key and certificate.
func GenerateClientPrivateKey(path string, pwd string) error {
	if Utility.Exists(path + string(os.PathSeparator) + "client.key") {
		return nil
	}

	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "genrsa")
	args = append(args, "-passout")
	args = append(args, "pass:"+pwd)
	args = append(args, "-des3")
	args = append(args, "-out")
	args = append(args, path+string(os.PathSeparator)+"client.pass.key")
	args = append(args, "4096")

	err := exec.Command(cmd, args...).Run()
	if err != nil {
		return errors.New("Fail to generate client private key " + err.Error())
	}

	args = make([]string, 0)
	args = append(args, "rsa")
	args = append(args, "-passin")
	args = append(args, "pass:"+pwd)
	args = append(args, "-in")
	args = append(args, path+string(os.PathSeparator)+"client.pass.key")
	args = append(args, "-out")
	args = append(args, path+string(os.PathSeparator)+"client.key")

	err = exec.Command(cmd, args...).Run()
	if err != nil {
		return errors.New("Fail to generate client private key " + err.Error())
	}

	// Remove the file.
	os.Remove(path + string(os.PathSeparator) + "client.pass.key")
	return nil
}

func GenerateClientCertificateSigningRequest(path string, pwd string, domain string) error {
	if Utility.Exists(path + string(os.PathSeparator) + "client.csr") {
		return nil
	}

	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "req")
	args = append(args, "-new")
	args = append(args, "-key")
	args = append(args, path+string(os.PathSeparator)+"client.key")
	args = append(args, "-out")
	args = append(args, path+string(os.PathSeparator)+"client.csr")
	args = append(args, "-subj")
	args = append(args, "/CN="+domain)
	err := exec.Command(cmd, args...).Run()
	if err != nil {
		log.Println(args)
		return errors.New("Fail to generate client certificate signing request.")
	}

	return nil
}

func GenerateSignedClientCertificate(path string, pwd string, expiration_delay int) error {

	if Utility.Exists(path + string(os.PathSeparator) + "client.crt") {
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
	args = append(args, path+string(os.PathSeparator)+"client.csr")
	args = append(args, "-CA")
	args = append(args, path+string(os.PathSeparator)+"ca.crt")
	args = append(args, "-CAkey")
	args = append(args, path+string(os.PathSeparator)+"ca.key")
	args = append(args, "-set_serial")
	args = append(args, "01")
	args = append(args, "-out")
	args = append(args, path+string(os.PathSeparator)+"client.crt")

	err := exec.Command(cmd, args...).Run()
	if err != nil {
		log.Println("fail to get the signed server certificate")
	}

	return nil
}

// Server certificate signing request (this should be shared with the CA owner)
func GenerateServerCertificateSigningRequest(path string, pwd string, domain string) error {

	if Utility.Exists(path + string(os.PathSeparator) + "sever.crs") {
		return nil
	}

	cmd := "openssl"
	args := make([]string, 0)
	args = append(args, "req")
	args = append(args, "-passin")
	args = append(args, "pass:"+pwd)
	args = append(args, "-new")
	args = append(args, "-key")
	args = append(args, path+string(os.PathSeparator)+"server.key")
	args = append(args, "-out")
	args = append(args, path+string(os.PathSeparator)+"server.csr")
	args = append(args, "-subj")
	args = append(args, "/CN="+domain)

	err := exec.Command(cmd, args...).Run()
	if err != nil {
		return errors.New("Fail to generate server certificate signing request.")
	}
	return nil
}

// Server certificate signed by the CA (this would be sent back to the client by the CA owner)
func GenerateSignedServerCertificate(path string, pwd string, expiration_delay int) error {

	if Utility.Exists(path + string(os.PathSeparator) + "sever.crt") {
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
	args = append(args, path+string(os.PathSeparator)+"server.csr")
	args = append(args, "-CA")
	args = append(args, path+string(os.PathSeparator)+"ca.crt")
	args = append(args, "-CAkey")
	args = append(args, path+string(os.PathSeparator)+"ca.key")
	args = append(args, "-set_serial")
	args = append(args, "01")
	args = append(args, "-out")
	args = append(args, path+string(os.PathSeparator)+"server.crt")

	err := exec.Command(cmd, args...).Run()
	if err != nil {
		log.Println("fail to get the signed server certificate")
	}

	return nil
}

// Conversion of server.key into a format gRpc likes (this shouldn't be shared)
func KeyToPem(name string, path string, pwd string) error {
	if Utility.Exists(path + string(os.PathSeparator) + name + ".pem") {
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
	args = append(args, path+string(os.PathSeparator)+name+".key")
	args = append(args, "-out")
	args = append(args, path+string(os.PathSeparator)+name+".pem")

	err := exec.Command(cmd, args...).Run()
	if err != nil {
		return errors.New("Fail to generate server.pem key from server.key")
	}

	return nil
}

/**
 * That function is use to generate services certificates.
 * Private ca.key, server.key, server.pem, server.crt
 * Share ca.crt (needed by the client), server.csr (needed by the CA)
 */
func GenerateServicesCertificates(pwd string, expiration_delay int, domain string, path string) {

	// Step 1: Generate Certificate Authority + Trust Certificate (ca.crt)
	err := GenerateAuthorityPrivateKey(path, pwd)
	if err != nil {
		log.Println(err)
		return
	}

	err = GenerateAuthorityTrustCertificate(path, pwd, expiration_delay, domain)
	if err != nil {
		log.Println(err)
		return
	}

	// Setp 2: Generate the server Private Key (server.key)
	err = GenerateSeverPrivateKey(path, pwd)
	if err != nil {
		log.Println(err)
		return
	}

	// Setp 3: Get a certificate signing request from the CA (server.csr)
	err = GenerateServerCertificateSigningRequest(path, pwd, domain)
	if err != nil {
		log.Println(err)
		return
	}

	// Step 4: Sign the certificate with the CA we create(it's called self signing) - server.crt
	err = GenerateSignedServerCertificate(path, pwd, expiration_delay)
	if err != nil {
		log.Println(err)
		return
	}

	// Step 5: Convert the server Certificate to .pem format (server.pem) - usable by gRpc
	err = KeyToPem("server", path, pwd)
	if err != nil {
		log.Println(err)
		return
	}

	// Step 6: Generate client private key.
	err = GenerateClientPrivateKey(path, pwd)
	if err != nil {
		log.Println(err)
		return
	}

	// Step 7: Generate the client signing request.
	err = GenerateClientCertificateSigningRequest(path, pwd, domain)
	if err != nil {
		log.Println(err)
		return
	}

	// Step 8: Generate client signed certificate.
	err = GenerateSignedClientCertificate(path, pwd, expiration_delay)
	if err != nil {
		log.Println(err)
		return
	}

	// Step 9: Convert to pem format.
	err = KeyToPem("client", path, pwd)
	if err != nil {
		log.Println(err)
		return
	}
}
