package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/StalkR/httpcache"
	"github.com/StalkR/imdb"
	"github.com/davecourtois/Utility"
	config_ "github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/rbac/rbacpb"
	"github.com/globulario/services/golang/resource/resourcepb"
	"github.com/globulario/services/golang/security"
	colly "github.com/gocolly/colly/v2"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	protoparser "github.com/yoheimuta/go-protoparser/v4"
)

const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"

const cacheTTL = 24 * time.Hour

// client is used by tests to perform cached requests.
// If cache directory exists it is used as a persistent cache.
// Otherwise a volatile memory cache is used.
var client *http.Client

func init() {
	if _, err := os.Stat("cache"); err == nil {
		client, err = httpcache.NewPersistentClient("cache", cacheTTL)
		if err != nil {
			panic(err)
		}
	} else {
		client = httpcache.NewVolatileClient(cacheTTL, 1024)
	}
	client.Transport = &customTransport{client.Transport}
}

// customTransport implements http.RoundTripper interface to add some headers.
type customTransport struct {
	http.RoundTripper
}

func (e *customTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Accept-Language", "en") // avoid IP-based language detection
	r.Header.Set("User-Agent", userAgent)
	return e.RoundTripper.RoundTrip(r)
}

// Find the peer with a given name and redirect the
// the request to it.
func redirectTo(host string) (bool, *resourcepb.Peer) {

	// read the actual configuration.
	__address__, err := config_.GetAddress()

	if err == nil {
		// no redirection if the address is the same...
		if strings.HasPrefix(__address__, host) {
			return false, nil
		}
	}

	var p *resourcepb.Peer

	globule.peers.Range(func(key, value interface{}) bool {
		p_ := value.(*resourcepb.Peer)
		address := p_.Hostname + "." + p_.Domain
		if p_.Protocol == "https" {
			address += ":" + Utility.ToString(p_.PortHttps)
		} else {
			address += ":" + Utility.ToString(p_.PortHttp)
		}

		if strings.HasPrefix(address, host) {
			p = p_
			return false // stop the iteration.
		}
		return true
	})

	return p != nil, p
}

// Redirect the query to a peer one the network
func handleRequestAndRedirect(to *resourcepb.Peer, res http.ResponseWriter, req *http.Request) {

	address := to.Domain
	scheme := "http"
	if to.Protocol == "https" {
		address += ":" + Utility.ToString(to.PortHttps)
	} else {
		address += ":" + Utility.ToString(to.PortHttp)
	}

	// Here I will remove the .localhost part of the address (if it exist)
	address = strings.ReplaceAll(address, ".localhost", "")
	ur, _ := url.Parse(scheme + "://" + address)

	proxy := httputil.NewSingleHostReverseProxy(ur)

	// Update the headers to allow for SSL redirection
	req.URL.Host = ur.Host
	req.URL.Scheme = ur.Scheme
	req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	proxy.ErrorHandler = ErrHandle

	proxy.ServeHTTP(res, req)
}

// Display error message.
func ErrHandle(res http.ResponseWriter, req *http.Request, err error) {
	fmt.Println(err)
}

/**
 * Create a checksum from a given path.
 */
func getChecksumHanldler(w http.ResponseWriter, r *http.Request) {
	// Receive http request...
	redirect, to := redirectTo(r.Host)

	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	// Handle the prefligth oprions...
	setupResponse(&w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	//add prefix and clean
	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(http.StatusCreated)

	execPath := Utility.GetExecName(os.Args[0])
	if Utility.Exists("/usr/local/share/globular/Globular") {
		execPath = "/usr/local/share/globular/Globular"
	}
	fmt.Fprint(w, Utility.CreateFileChecksum(execPath))
}

/**
 * Save the configuration.
 */
func saveConfigHanldler(w http.ResponseWriter, r *http.Request) {

	// Receive http request...
	redirect, to := redirectTo(r.Host)
	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	setupResponse(&w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// I will try to get the token from the header.
	token := r.Header.Get("token")

	// I will validate the token.
	if len(token) == 0 {
		// the token can be given by the url directly...
		token = r.URL.Query().Get("token")
	}

	// If not token was given i will return an error (403).
	if len(token) == 0 {
		http.Error(w, "no token was given!", http.StatusUnauthorized)
		return
	}

	// I will validate the token.
	_, err := security.ValidateToken(token)
	if err != nil {
		http.Error(w, "fail to validate token with error "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Now I will get the configuration from the request.
	decoder := json.NewDecoder(r.Body)
	var config map[string]interface{}
	err = decoder.Decode(&config)
	if err != nil {
		http.Error(w, "fail to decode configuration with error "+err.Error(), http.StatusBadRequest)
		return
	}

	// I will set the globular configuration.
	err = globule.setConfig(config)
	if err != nil {
		http.Error(w, "fail to set configuration with error "+err.Error(), http.StatusBadRequest)
		return
	}

}

func getPublicKeyHanldler(w http.ResponseWriter, r *http.Request) {
	// Here I will get the public key from the configuration.
	redirect, to := redirectTo(r.Host)
	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	// Handle the prefligth oprions...
	setupResponse(&w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	path := config_.GetConfigDir() + "/keys/" + strings.ReplaceAll(globule.Mac, ":", "_") + "_public"
	if !Utility.Exists(path) {
		http.Error(w, "no public key found!", http.StatusBadRequest)
		return
	}

	// read the public key file and return it as text string.
	data, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, "fail to read public key with error "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(data)

	if err != nil {
		http.Error(w, "fail to write public key with error "+err.Error(), http.StatusBadRequest)
		return
	}

}
func getCertificateHanldler(w http.ResponseWriter, r *http.Request) {

	// ... [existing code] ...
	address, err := config_.GetAddress()
	if err != nil {
		http.Error(w, "fail to get address with error "+err.Error(), http.StatusBadRequest)
		return
	}

	domain := strings.Split(address, ":")[0]
	certFilename := config_.GetLocalCertificate()
	path := config_.GetConfigDir() + "/tls/" + domain + "/" + config_.GetLocalCertificate()

	if !Utility.Exists(path) {
		http.Error(w, "no issuer certificate found at path "+path, http.StatusBadRequest)
		return
	}

	data, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, "fail to read public key with error "+err.Error(), http.StatusBadRequest)
		return
	}

	// Set the headers to suggest a download file name and indicate the file type.
	w.Header().Set("Content-Disposition", "attachment; filename=\""+certFilename+"\"")
	w.Header().Set("Content-Type", "application/x-x509-ca-cert")

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "fail to write public key with error "+err.Error(), http.StatusBadRequest)
		return
	}
}

func getIssuerCertificateHandler(w http.ResponseWriter, r *http.Request) {

	// ... [existing code] ...
	address, err := config_.GetAddress()
	if err != nil {
		http.Error(w, "fail to get address with error "+err.Error(), http.StatusBadRequest)
		return
	}

	domain := strings.Split(address, ":")[0]
	certFilename := config_.GetLocalCertificateAuthorityBundle()
	path := config_.GetConfigDir() + "/tls/" + domain + "/" + config_.GetLocalCertificateAuthorityBundle()

	if !Utility.Exists(path) {
		http.Error(w, "no issuer certificate found at path "+path, http.StatusBadRequest)
		return
	}

	data, err := os.ReadFile(path)
	if err != nil {
		http.Error(w, "fail to read public key with error "+err.Error(), http.StatusBadRequest)
		return
	}

	// Set the headers to suggest a download file name and indicate the file type.
	w.Header().Set("Content-Disposition", "attachment; filename=\""+certFilename+"\"")
	w.Header().Set("Content-Type", "application/x-x509-ca-cert")

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		http.Error(w, "fail to write public key with error "+err.Error(), http.StatusBadRequest)
		return
	}
}

/**
 * Return services permissions configuration to be able to manage resources access from rpc request.
 */
func getServicePermissionsHanldler(w http.ResponseWriter, r *http.Request) {
	// Receive http request...
	redirect, to := redirectTo(r.Host)
	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	// Handle the prefligth oprions...
	setupResponse(&w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	//add prefix and clean
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// so here I will retreive the service configuration from the service id given in the query
	serviceId := r.URL.Query().Get("id") // the csr in base64

	//add prefix and clean
	serviceConfig, err := config_.GetServiceConfigurationById(serviceId)
	if err != nil {
		http.Error(w, "fail to get service configuration with error "+err.Error(), http.StatusBadRequest)
		return
	}

	// from the configuration i will read the configuration file...
	data, err := os.ReadFile(serviceConfig["ConfigPath"].(string))

	// reload the configuration with the permissions...
	err = json.Unmarshal(data, &serviceConfig)
	if err != nil {
		http.Error(w, "fail to get service configuration with error "+err.Error(), http.StatusBadRequest)
		return
	}

	// set empty array if not defined...
	if serviceConfig["Permissions"] == nil {
		serviceConfig["Permissions"] = []interface{}{}
	}

	gotJSON, err := json.MarshalIndent(serviceConfig["Permissions"], "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal, err %v\n", err)
	}

	w.Write(gotJSON)
}

/**
 * This function is use to return a json object containing the description of the service.
 */
func getServiceDescriptorHanldler(w http.ResponseWriter, r *http.Request) {

	// Receive http request...
	redirect, to := redirectTo(r.Host)
	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	// Handle the prefligth oprions...
	setupResponse(&w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	//add prefix and clean
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	// so here I will retreive the service configuration from the service id given in the query
	serviceId := r.URL.Query().Get("id") // the csr in base64

	//add prefix and clean
	serviceConfig, err := config_.GetServiceConfigurationById(serviceId)
	if err != nil {
		http.Error(w, "fail to get service configuration with error "+err.Error(), http.StatusBadRequest)
		return
	}

	// from the service configuration I will read it proto file...
	protoFile := serviceConfig["Proto"].(string)

	// I will read the proto file and return it as a json object.
	reader, err := os.Open(protoFile)

	if err != nil {
		http.Error(w, "fail to open proto file with error "+err.Error(), http.StatusBadRequest)
		return
	}

	defer reader.Close()

	// parse the proto file.

	got, err := protoparser.Parse(
		reader,
		protoparser.WithDebug(false),
		protoparser.WithPermissive(false),
		protoparser.WithFilename(filepath.Base(protoFile)),
	)

	var v interface{}
	v = got

	gotJSON, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal, err %v\n", err)
	}

	w.Write(gotJSON)

}

/**
 * Return the service configuration
 */
func getConfigHanldler(w http.ResponseWriter, r *http.Request) {

	// Receive http request...
	redirect, to := redirectTo(r.Host)
	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	setupResponse(&w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// i will redirect to the given host if the host is not the same...
	address, _ := config_.GetAddress()

	// I will redirect the request if host is defined in the query...
	if !strings.HasPrefix(address, r.URL.Query().Get("host")) && len(r.URL.Query().Get("host")) > 0 {

		redirect, to := redirectTo(r.URL.Query().Get("host"))

		if redirect && to != nil {

			// I will get the remote configuration and return it...
			var remoteConfig map[string]interface{}
			var err error
			address := to.LocalIpAddress
			if to.ExternalIpAddress != Utility.MyIP() {
				address = to.ExternalIpAddress
			}

			remoteConfig, err = config_.GetRemoteConfig(address, 0)

			if err != nil {
				http.Error(w, "Fail to get remote configuration with error "+err.Error(), http.StatusBadRequest)
				return
			}

			//add prefix and clean
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)

			json.NewEncoder(w).Encode(remoteConfig)

			return
		} else {
			// I will get the remote configuration and return it.
			remoteConfig, err := config_.GetRemoteConfig(r.URL.Query().Get("host"), Utility.ToInt(r.URL.Query().Get("port")))
			if err != nil {
				// Try again with port 80...
				remoteConfig, err = config_.GetRemoteConfig(r.URL.Query().Get("host"), 80)
				if err != nil {
					http.Error(w, "Fail to get remote configuration with error "+err.Error(), http.StatusBadRequest)
					return
				}
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)

			jsonStr, err := json.MarshalIndent(remoteConfig, "", "  ")
			if err != nil {
				http.Error(w, "fail to encode json with error "+err.Error(), http.StatusBadRequest)

			}

			w.Write(jsonStr)

			return
		}
	}

	setupResponse(&w, r)

	// if the host is not the same...
	serviceId := r.URL.Query().Get("id") // the csr in base64

	//add prefix and clean
	config := globule.getConfig()

	// give list of path...
	config["Root"] = config_.GetRootDir()
	config["DataPath"] = config_.GetDataDir()
	config["ConfigPath"] = config_.GetConfigDir()
	config["WebRoot"] = config_.GetWebRootDir()
	config["Public"] = config_.GetPublicDirs()

	// ask for a service configuration...
	if len(serviceId) > 0 {
		services := config["Services"].(map[string]interface{})
		exist := false
		for _, service := range services {
			if service.(map[string]interface{})["Id"].(string) == serviceId || service.(map[string]interface{})["Name"].(string) == serviceId {
				config = service.(map[string]interface{})
				exist = true
				break
			}
		}
		if !exist {
			http.Error(w, "no service found with name or id "+serviceId, http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	jsonStr, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		http.Error(w, "fail to encode json with error "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(jsonStr)
}

func dealwithErr(err error) {
	if err != nil {
		fmt.Println(err)
		//os.Exit(-1)
	}
}

func getHardwareData(w http.ResponseWriter, r *http.Request) {

	// Receive http request...
	redirect, to := redirectTo(r.Host)
	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

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
	if len(cpuStat) > 0 {
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

	jsonStr, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		http.Error(w, "fail to encode json with error "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(jsonStr)
}

/**
 * Return the ca certificate public key.
 */
func getCaCertificateHanldler(w http.ResponseWriter, r *http.Request) {

	redirect, to := redirectTo(r.Host)
	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	// Handle the prefligth oprions...
	setupResponse(&w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	//add prefix and clean
	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(http.StatusCreated)

	crt, err := os.ReadFile(globule.creds + "/ca.crt")
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
	redirect, to := redirectTo(r.Host)
	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	//add prefix and clean
	w.Header().Set("Content-Type", "application/text")
	setupResponse(&w, r)

	crt, err := os.ReadFile(globule.creds + "/san.conf")
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
		if globule.AllowedOrigins[i] == "*" {
			allowedOrigins = "*"
			break
		}
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

	header := (*w).Header()

	// set the cors header...
	header.Set("Access-Control-Allow-Origin", allowedOrigins)
	header.Set("Access-Control-Allow-Methods", allowedMethods)
	header.Set("Access-Control-Allow-Headers", allowedHeaders)

	if req.Method == http.MethodOptions {
		header.Set("Access-Control-Max-Age", "3600")
		header.Set("Access-Control-Allow-Private-Network", "true")
	}
}

/**
 * Sign ca certificate request and return a certificate.
 */
func signCaCertificateHandler(w http.ResponseWriter, r *http.Request) {

	redirect, to := redirectTo(r.Host)
	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	// Handle the prefligth oprions...
	setupResponse(&w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

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

// Return true if the file is found in the public path...
func isPublic(path string) bool {
	public := config_.GetPublicDirs()
	path = strings.ReplaceAll(path, "\\", "/")

	for i := 0; i < len(public); i++ {
		if strings.HasPrefix(strings.ToLower(path), strings.ReplaceAll(strings.ToLower(public[i]), "\\", "/")) {
			return true
		}
	}

	return false
}

// ImageList is the structure for our response
type ImageList struct {
	Images []string `json:"images"`
}

/**
 * Return a list of images from a given path. The path is given in the query.
 * The path is relative to the web root directory.
 */
func GetImagesHandler(w http.ResponseWriter, r *http.Request) {

	redirect, to := redirectTo(r.Host)

	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	// Handle the prefligth oprions...
	setupResponse(&w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	dir := globule.webRoot

	// If a directory with the same name as the host in the request exist
	// it will be taken as root. Permission will be manage by the resource
	// manager and not simply the name of the directory. If you want to protect
	// a given you need to set permission on it.
	if Utility.Exists(dir + "/" + r.Host) {
		dir += "/" + r.Host
	}

	// so I will get the path from the query...
	path := r.URL.Query().Get("path")

	// If the path is not defined I will return an error.
	if len(path) == 0 {
		http.Error(w, "Failed to get images no path was given", http.StatusInternalServerError)
		return
	}

	// Be sure that the path start with a /.
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if !Utility.Exists(dir + path) {

		http.Error(w, "Failed to get images path not found "+dir+path, http.StatusInternalServerError)
		return

	}

	// Get a list of images
	imageFiles, err := getListOfImages(dir + path)
	if err != nil {
		http.Error(w, "Failed to get images", http.StatusInternalServerError)
		return
	}

	// Create a response structure
	response := ImageList{Images: imageFiles}

	// I will replace all images path by the relative path.
	for i := 0; i < len(response.Images); i++ {
		response.Images[i] = strings.ReplaceAll(response.Images[i], dir, "")
	}

	// Marshal the response structure to JSON
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	// Write the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}

func getListOfImages(dirPath string) ([]string, error) {
	var fileList []string
	err := filepath.Walk(dirPath, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			fileList = append(fileList, path)
		}
		return nil
	})

	return fileList, err
}

/**
 * Evaluate the file size at given url
 */
func GetFileSizeAtUrl(w http.ResponseWriter, r *http.Request) {

	// here in case of file uploaded from other website like pornhub...
	url := r.URL.Query().Get("url")

	fmt.Println("try to get file size for url ", url)
	// we are interested in getting the file or object name
	// so take the last item from the slice
	subStringsSlice := strings.Split(url, "/")
	fileName := subStringsSlice[len(subStringsSlice)-1]

	resp, err := http.Head(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Is our request ok?

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status)
		return
	}

	// the Header "Content-Length" will let us know
	// the total file size to download
	size, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	downloadSize := int64(size)

	fmt.Println("Will be downloading ", fileName, " of ", downloadSize, " bytes.")
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(&map[string]int64{"size": downloadSize})
	if err == nil {
		w.Write(data)
	} else {
		http.Error(w, "Fail to get file size at "+url+" with error "+err.Error(), http.StatusExpectationFailed)
	}
}

/**
 * This code is use to upload a file into the tmp directory of the server
 * via http request.
 */
func FileUploadHandler(w http.ResponseWriter, r *http.Request) {

	redirect, to := redirectTo(r.Host)

	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	// Handle the prefligth oprions...
	setupResponse(&w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	err := r.ParseMultipartForm(32 << 20) // grab the multipart form
	if err != nil {
		fmt.Println("transfert error: ", err)
		http.Error(w, "failed to parse multipart message "+err.Error(), http.StatusBadRequest)
		return
	}

	formdata := r.MultipartForm // ok, no problem so far, read the Form data

	//get the *fileheaders
	files := formdata.File["multiplefiles"] // grab the filenames

	// Get the path where to upload the file.
	path := r.FormValue("path")
	path = strings.ReplaceAll(path, "\\", "/")

	fmt.Println("try to upload file at path: ", path)

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
		if strings.HasPrefix(path, "/applications") {
			hasAccess, hasAccessDenied, err = globule.validateAction("/file.FileService/FileUploadHandler", application, rbacpb.SubjectType_APPLICATION, infos)
			if err != nil {
				http.Error(w, "fail to validate access with error "+err.Error(), http.StatusUnauthorized)
				return
			}
		}
	}

	// get the user id from the token...
	domain := r.URL.Query().Get("domain")
	if len(token) != 0 && !hasAccess {
		var claims *security.Claims
		claims, err := security.ValidateToken(token)
		if err == nil {
			user = claims.Id + "@" + claims.UserDomain
			domain = claims.Domain

			fmt.Println("values found from token are user:", user, "domain", claims.UserDomain)
		} else {
			fmt.Println("fail to validate token with error ", err.Error())
			http.Error(w, "fail to validate token with error "+err.Error(), http.StatusUnauthorized)
			return
		}
	} else {
		fmt.Println("no token was given!")
	}

	if len(user) != 0 {
		if !hasAccess || hasAccessDenied {
			hasAccess, hasAccessDenied, err = globule.validateAction("/file.FileService/FileUploadHandler", user, rbacpb.SubjectType_ACCOUNT, infos)
			if err != nil {
				http.Error(w, "fail to validate access with error "+err.Error(), http.StatusUnauthorized)
				return
			}
		}
	}

	// validate ressource access...
	if !hasAccess || hasAccessDenied {
		http.Error(w, "unable to create the file for writing. Check your access privilege", http.StatusUnauthorized)
		return
	}

	for _, f := range files { // loop through the files one by one
		file, err := f.Open()
		if err != nil {
			return
		}
		defer file.Close()

		// Create the file depending if the path is users, applications or something else...
		path_ := path + "/" + f.Filename
		size, _ := file.Seek(0, 2)
		if len(user) > 0 {
			hasSpace, err := ValidateSubjectSpace(user, rbacpb.SubjectType_ACCOUNT, uint64(size))
			if !hasSpace || err != nil {
				http.Error(w, user+" has no space available to copy file "+path_+" allocated space and try again.", http.StatusUnauthorized)
				return
			}
		}

		file.Seek(0, 0)
		// Now if the os is windows I will remove the leading /
		if len(path_) > 3 {
			if runtime.GOOS == "windows" && path_[0] == '/' && path_[2] == ':' {
				path_ = path_[1:]
			}
		}

		if strings.HasPrefix(path, "/users") || strings.HasPrefix(path, "/applications") {
			path_ = strings.ReplaceAll(globule.data+"/files"+path_, "\\", "/")
		} else if !isPublic(path_) {
			path_ = strings.ReplaceAll(globule.webRoot+path_, "\\", "/")
		}

		out, err := os.Create(path_)
		if err != nil {
			return
		}

		defer out.Close()

		if err != nil {
			http.Error(w, "Unable to create the file for writing. Check your write access privilege", http.StatusUnauthorized)
			return
		}

		_, err = io.Copy(out, file) // file not files[i] !
		if err != nil {
			return
		}

		// Here I will set the ressource owner.
		if len(user) > 0 {
			globule.addResourceOwner(path+"/"+f.Filename, "file", user+"@"+domain, rbacpb.SubjectType_ACCOUNT)
		} else if len(application) > 0 {
			globule.addResourceOwner(path+"/"+f.Filename, "file", application+"@"+domain, rbacpb.SubjectType_APPLICATION)
		}

		// Now from the file extension i will retreive it mime type.
		if strings.LastIndex(path_, ".") != -1 {
			fileExtension := path_[strings.LastIndex(path_, "."):]
			fileType := mime.TypeByExtension(fileExtension)
			path_ = strings.ReplaceAll(path_, "\\", "/")
			if len(fileType) > 0 {
				if strings.HasPrefix(fileType, "video/") {
					// Here I will call convert video...
					globule.publish("generate_video_preview_event", []byte(path_))
				} else if fileType == "application/pdf" || strings.HasPrefix(fileType, "text") {
					// Here I will call convert video...
					globule.publish("index_file_event", []byte(path_))
				}
			}
		}
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

/**
 * findHashedFile search for a hashed file in the same directory as the original file.
 * ex if styles.css not exist i will try to find styles.123456.css in the same directory.
 */
func findHashedFile(originalPath string) (string, error) {

	// Get the directory of the original file
	dir := filepath.Dir(originalPath)

	// Get the base name of the original file without extension
	baseName := strings.TrimSuffix(filepath.Base(originalPath), filepath.Ext(originalPath))

	// Read the files in the directory
	files, err := os.ReadDir(dir)
	if err != nil {
		return "", err
	}

	// Search for a matching hashed file
	for _, file := range files {
		if strings.HasPrefix(file.Name(), baseName) && strings.HasSuffix(file.Name(), filepath.Ext(originalPath)) {
			return filepath.Join(dir, file.Name()), nil
		}
	}

	return "", fmt.Errorf("hashed file not found for %s", originalPath)
}

// Custom file server implementation.
func ServeFileHandler(w http.ResponseWriter, r *http.Request) {

	redirect, to := redirectTo(r.Host)

	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	//add prefix and clean
	rqst_path := path.Clean(r.URL.Path)

	if rqst_path == "/null" {
		http.Error(w, "No file path was given in the file url path!", http.StatusBadRequest)
	}

	// I will test if the requested path is in the reverse proxy list.
	// if it is the case I will redirect the request to the reverse proxy.
	for _, proxy := range globule.ReverseProxies {
		proxyPath_ := strings.TrimSpace(strings.Split(proxy.(string), "|")[1])
		proxyURL_ := strings.TrimSpace(strings.Split(proxy.(string), "|")[0])

		if strings.HasPrefix(rqst_path, proxyPath_) {
			// Create a reverse proxy
			proxyURL, _ := url.Parse(proxyURL_)

			// Connect to the proxy host
			hostUrl, _ := url.Parse(proxyURL.Scheme + "://" + proxyURL.Host)

			reverseProxy := httputil.NewSingleHostReverseProxy(hostUrl)

			// Update the request URL and headers
			r.URL, _ = url.Parse(proxyURL_)

			// Update headers to reflect the forwarded host
			r.Host = proxyURL.Host

			r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))

			// Print request details
			//printRequestInfo(r)

			// Serve the request via the reverse proxy
			reverseProxy.ServeHTTP(w, r)
			return
		}
	}

	setupResponse(&w, r)
	dir := globule.webRoot

	// If a directory with the same name as the host in the request exist
	// it will be taken as root. Permission will be manage by the resource
	// manager and not simply the name of the directory. If you want to protect
	// a given you need to set permission on it.
	if Utility.Exists(dir + "/" + r.Host) {
		dir += "/" + r.Host
	}

	// Now I will test if a token is given in the header and manage it file access.
	application := r.Header.Get("application")
	token := r.Header.Get("token")

	if len(token) == 0 {
		// the token can be given by the url directly...
		token = r.URL.Query().Get("token")

		fmt.Println("try to get token from url ", token)
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
			if Utility.Exists(dir + "/" + rqst_path) {
				rqst_path = "/" + globule.IndexApplication + rqst_path
			}
		}
	}

	hasAccess := true
	var name string
	if strings.HasPrefix(rqst_path, "/users/") || strings.HasPrefix(rqst_path, "/applications/") || strings.HasPrefix(rqst_path, "/templates/") || strings.HasPrefix(rqst_path, "/projects/") {
		dir = globule.data + "/files"
		if !strings.Contains(rqst_path, "/.hidden/") {
			hasAccess = false
		}
	}

	// Now if the os is windows I will remove the leading /
	if len(rqst_path) > 3 {
		if runtime.GOOS == "windows" && rqst_path[0] == '/' && rqst_path[2] == ':' {
			rqst_path = rqst_path[1:]
		}
	}
	// path to file
	if !isPublic(rqst_path) {
		name = path.Join(dir, rqst_path)
	} else {
		name = rqst_path
		hasAccess = false // force validation (denied access...)
	}


	fmt.Println("serve file ", name)

	// stream, the validation is made on the directory containning the playlist...
	if strings.Contains(rqst_path, "/.hidden/") ||
		strings.HasSuffix(rqst_path, ".ts") ||
		strings.HasSuffix(rqst_path, "240p.m3u8") ||
		strings.HasSuffix(rqst_path, "360p.m3u8") ||
		strings.HasSuffix(rqst_path, "480p.m3u8") ||
		strings.HasSuffix(rqst_path, "720p.m3u8") ||
		strings.HasSuffix(rqst_path, "1080p.m3u8") ||
		strings.HasSuffix(rqst_path, "2160p.m3u8") {
		hasAccess = true
	}

	// this is the ca certificate use to sign client certificate.
	if rqst_path == "/ca.crt" {
		name = globule.creds + rqst_path
	}

	hasAccessDenied := false
	var err error
	var userId string

	fmt.Println("try to serve file ", name, " with access ", hasAccess, " and access denied ", hasAccessDenied, " and token ", token)
	if len(token) != 0 && !hasAccess {
		fmt.Println("try to validate token")
		var claims *security.Claims
		claims, err = security.ValidateToken(token)
		userId = claims.Id + "@" + claims.UserDomain
		if err == nil {
			hasAccess, hasAccessDenied, err = globule.validateAccess(userId, rbacpb.SubjectType_ACCOUNT, "read", rqst_path)

			fmt.Println("values found from token are user:", userId, "domain", claims.UserDomain, " hasAccess ", hasAccess, " hasAccessDenied ", hasAccessDenied)
		} else {
			fmt.Println("fail to validate token with error: ", err)
		}
	}

	// Here I will validate applications...
	 if isPublic(rqst_path) && !hasAccessDenied && !hasAccess {
		hasAccess = true
	} else if !hasAccess && !hasAccessDenied && len(application) != 0 {
		hasAccess, hasAccessDenied, err = globule.validateAccess(application, rbacpb.SubjectType_APPLICATION, "read", rqst_path)
	} 

	// validate ressource access...
	if !hasAccess || hasAccessDenied || err != nil {
		msg := "unable to read the file " + rqst_path + " Check your access privilege"
		http.Error(w, msg, http.StatusUnauthorized)
		return
	}

	var code string
	// If the file is a javascript file...
	hasChange := false

	if !Utility.Exists(name) {
		name = "/" + rqst_path // try network path...
	}

	//fmt.Println("serve file ", name)

	f, err := os.Open(name)
	if err != nil {
		if os.IsNotExist(err) {
			name, err = findHashedFile(name)
			if err == nil {
				f, err = os.Open(name)
				if err != nil {
					http.Error(w, "File "+rqst_path+" not found!", http.StatusNoContent)
					return
				}
			} else {
				http.Error(w, "File "+rqst_path+" not found!", http.StatusNoContent)
				return
			}
		}
	}

	f.Close()

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
		http.ServeContent(w, r, name, time.Now(), strings.NewReader(code))
	}
}

/**
 * Return a list of IMDB titles from a keyword...
 */
func getImdbTitlesHanldler(w http.ResponseWriter, r *http.Request) {
	// Receive http request...
	redirect, to := redirectTo(r.Host)
	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	// if the host is not the same...
	query := r.URL.Query().Get("query") // the csr in base64

	titles, err := imdb.SearchTitle(client, query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(titles) == 0 {
		fmt.Fprintf(os.Stderr, "Not found.")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	setupResponse(&w, r)

	w.WriteHeader(http.StatusCreated)

	jsonStr, err := json.MarshalIndent(titles, "", "  ")
	if err != nil {
		http.Error(w, "fail to encode json with error "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(jsonStr)
}

// ////////////////////////// imdb missing sesson and episode number info... /////////////////////////
// get the thumbnail fil with help of youtube dl...
func downloadThumbnail(video_id, video_url, video_path string) (string, error) {

	fmt.Println("download thumbnail for ", video_path)

	if len(video_id) == 0 {
		return "", errors.New("no video id was given")
	}
	if len(video_url) == 0 {
		return "", errors.New("no video url was given")
	}
	if len(video_path) == 0 {
		return "", errors.New("no video path was given")
	}

	lastIndex := -1
	if strings.Contains(video_path, ".mp4") {
		lastIndex = strings.LastIndex(video_path, ".")
	}

	// The hidden folder path...
	path_ := video_path[0:strings.LastIndex(video_path, "/")]

	name_ := video_path[strings.LastIndex(video_path, "/")+1:]
	if lastIndex != -1 {
		name_ = video_path[strings.LastIndex(video_path, "/")+1 : lastIndex]
	}

	thumbnail_path := path_ + "/.hidden/" + name_ + "/__thumbnail__"

	if Utility.Exists(thumbnail_path + "/" + "data_url.txt") {

		thumbnail, err := os.ReadFile(thumbnail_path + "/" + "data_url.txt")
		if err != nil {
			return "", err
		}

		return string(thumbnail), nil
	}

	Utility.CreateDirIfNotExist(thumbnail_path)
	cmd := exec.Command("yt-dlp", video_url, "-o", video_id, "--write-thumbnail", "--skip-download")
	cmd.Dir = thumbnail_path

	err := cmd.Run()
	if err != nil {
		return "", err
	}

	files, err := Utility.ReadDir(thumbnail_path)
	if err != nil {
		return "", err
	}

	thumbnail, err := Utility.CreateThumbnail(filepath.Join(thumbnail_path, files[0].Name()), 300, 180)
	if err != nil {
		return "", err
	}

	err = os.WriteFile(thumbnail_path+"/"+"data_url.txt", []byte(thumbnail), 0664)
	if err != nil {
		return "", err
	}

	// cointain the data url...
	return thumbnail, nil
}

/**
 * Return the cover image...
 */
func GetCoverDataUrl(w http.ResponseWriter, r *http.Request) {
	// here in case of file uploaded from other website like pornhub...
	video_id := r.URL.Query().Get("id")
	video_url := r.URL.Query().Get("url")
	video_path := r.URL.Query().Get("path")

	dataUrl, err := downloadThumbnail(video_id, video_url, video_path)
	if err != nil {
		http.Error(w, "fail to create data url with error'"+err.Error()+"'", http.StatusExpectationFailed)
		return
	}

	w.Write([]byte(dataUrl))
}

func getSeasonAndEpisodeNumber(titleId string, nbCall int) (int, int, string, error) {

	resp, err := client.Get(`https://www.imdb.com/title/` + titleId)
	if err != nil {
		return -1, -1, "", err
	}
	defer resp.Body.Close()

	season := 0
	episode := 0
	serie := ""

	// The first regex to locate the information...
	re_SE := regexp.MustCompile(`>S\d{1,2}<!-- -->\.<!-- -->E\d{1,2}<`)
	page, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, -1, "", err
	}

	matchs_SE := re_SE.FindStringSubmatch(string(page))

	if len(matchs_SE) > 0 {
		re_S := regexp.MustCompile(`S\d{1,2}`)
		matchs_S := re_S.FindStringSubmatch(matchs_SE[0])
		if len(matchs_S) > 0 {
			season = Utility.ToInt(matchs_S[0][1:])
		}

		re_E := regexp.MustCompile(`E\d{1,2}`)
		matchs_E := re_E.FindStringSubmatch(matchs_SE[0])
		if len(matchs_E) > 0 {
			episode = Utility.ToInt(matchs_E[0][1:])
		}
	}

	// Now the serie info..
	re_Serie := regexp.MustCompile(`data-testid="hero-title-block__series-link" href="/title/tt\d{7}/\?ref_=tt_ov_inf"`)
	matchs_Serie := re_Serie.FindStringSubmatch(string(page))

	if len(matchs_Serie) > 0 {
		re_S := regexp.MustCompile(`tt\d{7}`)
		matchs_S := re_S.FindStringSubmatch(matchs_Serie[0])
		if len(matchs_S) > 0 {
			serie = matchs_S[0]
		}
	}

	fmt.Println("Seson ", season, "Episode", episode, "Serie", serie)

	return season, episode, serie, nil
}

/**
 * Return a json object with the movie information from imdb.
 */
func getImdbTitleHanldler(w http.ResponseWriter, r *http.Request) {
	// Receive http request...
	redirect, to := redirectTo(r.Host)
	if redirect {
		handleRequestAndRedirect(to, w, r)
		return
	}

	id := r.URL.Query().Get("id") // the csr in base64

	fmt.Println("get imdb info for ", id)

	title, err := imdb.NewTitle(client, id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	setupResponse(&w, r)

	w.WriteHeader(http.StatusCreated)

	title_, _ := Utility.ToMap(title)

	if title.Type == "TVEpisode" {
		s, e, t, err := getSeasonAndEpisodeNumber(id, 10)
		fmt.Println("get tv episode info ", id)
		if err == nil {
			title_["Season"] = s
			title_["Episode"] = e
			title_["Serie"] = t
		} else {
			fmt.Println("fail to retreive episode info ", err)
		}
	}

	// Now I will try to complete the casting informations...
	if title_["Actors"] != nil {
		for i := 0; i < len(title_["Actors"].([]interface{})); i++ {
			p := title_["Actors"].([]interface{})[i].(map[string]interface{})
			p_, err := setPersonInformation(p)
			if err == nil {
				title_["Actors"].([]interface{})[i] = p_
			}
		}
	}

	if title_["Writers"] != nil {
		for i := 0; i < len(title_["Writers"].([]interface{})); i++ {
			p := title_["Writers"].([]interface{})[i].(map[string]interface{})
			p_, err := setPersonInformation(p)
			if err == nil {
				title_["Writers"].([]interface{})[i] = p_
			}
		}
	}

	if title_["Directors"] != nil {
		for i := 0; i < len(title_["Directors"].([]interface{})); i++ {
			p := title_["Directors"].([]interface{})[i].(map[string]interface{})
			p_, err := setPersonInformation(p)
			if err == nil {
				title_["Directors"].([]interface{})[i] = p_
			}
		}
	}

	jsonStr, err := json.MarshalIndent(title_, "", "  ")
	if err != nil {
		http.Error(w, "fail to encode json with error "+err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(jsonStr)
}

func setPersonInformation(person map[string]interface{}) (map[string]interface{}, error) {
	movieCollector := colly.NewCollector(
		colly.AllowedDomains("www.imdb.com", "imdb.com"),
	)

	// So here I will define collector's...
	biographySelector := `a[name="mini_bio"]`
	movieCollector.OnHTML(biographySelector, func(e *colly.HTMLElement) {

		// keep the text only...
		person["Biography"] = e.DOM.Next().Next().Text()
	})

	// The profile image.
	profilePictureSelector := `#main > div.article.listo > div.subpage_title_block.name-subpage-header-block > a > img`
	movieCollector.OnHTML(profilePictureSelector, func(e *colly.HTMLElement) {
		person["Picture"] = strings.TrimSpace(e.Attr("src"))
	})

	// The birtdate
	birthdateSelector := `#overviewTable > tbody > tr:nth-child(1) > td:nth-child(2) > time`
	movieCollector.OnHTML(birthdateSelector, func(e *colly.HTMLElement) {
		person["BirthDate"] = e.Attr("datetime")
	})

	// The birtplace.
	birthplaceSelector := `#overviewTable > tbody > tr:nth-child(1) > td:nth-child(2) > a`
	movieCollector.OnHTML(birthplaceSelector, func(e *colly.HTMLElement) {
		person["BirthPlace"] = e.Text
	})

	url := person["URL"].(string) + "/bio?ref_=nm_ov_bio_sm"

	err := movieCollector.Visit(url)
	if err != nil {
		return nil, err
	}

	return person, nil
}
