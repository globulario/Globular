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
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/StalkR/imdb"
	"github.com/davecourtois/Utility"
	"github.com/dhowden/tag"
	config_ "github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/rbac/rbacpb"
	"github.com/globulario/services/golang/resource/resourcepb"
	"github.com/globulario/services/golang/security"
	"github.com/globulario/services/golang/title/titlepb"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

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
		address := p_.Domain
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
func handleRequestAndRedirect(address string, res http.ResponseWriter, req *http.Request) {

	// So here I will require a little more info about the peers...
	address_ := ""
	port := 0
	scheme := "http"

	if strings.Contains(address, ":") {
		address_ = strings.Split(address, ":")[0]
		port = Utility.ToInt(strings.Split(address, ":")[1])
	}

	config__, err := config_.GetRemoteConfig(address_, port, "")

	if err == nil {
		// if
		if config__["Protocol"].(string) == "https" && len(config__["Certificate"].(string)) != 0 {
			scheme = "https"
			address_ += ":" + Utility.ToString(config__["PortHttps"])
		} else if config__["Protocol"].(string) == "http" {
			address_ += ":" + Utility.ToString(config__["PortHttp"])
		}
	} else {
		address_ = address
	}

	ur, _ := url.Parse(scheme + "://" + address_)
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
	// Handle the prefligth oprions...
	setupResponse(&w, r)

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Receive http request...
	redirect, to := redirectTo(r.Host)

	if redirect {
		address := to.Domain
		if to.Protocol == "https" {
			address += ":" + Utility.ToString(to.PortHttps)
		} else {
			address += ":" + Utility.ToString(to.PortHttp)
		}
		handleRequestAndRedirect(address, w, r)
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
 * Return the service configuration
 */
func getConfigHanldler(w http.ResponseWriter, r *http.Request) {
	// Handle the prefligth oprions...
	setupResponse(&w, r)

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Receive http request...
	/*redirect, to := redirectTo(r.Host)


	if redirect {
		address := to.Domain
		if to.Protocol == "https" {
			address += ":" + Utility.ToString(to.PortHttps)
		} else {
			address += ":" + Utility.ToString(to.PortHttp)
		}
		handleRequestAndRedirect(address, w, r)
		return
	}*/

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
	json.NewEncoder(w).Encode(config)
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
		address := to.Domain
		if to.Protocol == "https" {
			address += ":" + Utility.ToString(to.PortHttps)
		} else {
			address += ":" + Utility.ToString(to.PortHttp)
		}
		handleRequestAndRedirect(address, w, r)
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
	// Handle the prefligth oprions...
	setupResponse(&w, r)

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	redirect, to := redirectTo(r.Host)
	if redirect {
		address := to.Domain
		if to.Protocol == "https" {
			address += ":" + Utility.ToString(to.PortHttps)
		} else {
			address += ":" + Utility.ToString(to.PortHttp)
		}
		handleRequestAndRedirect(address, w, r)
		return
	}

	//add prefix and clean
	w.Header().Set("Content-Type", "application/text")
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
	redirect, to := redirectTo(r.Host)
	if redirect {
		address := to.Domain
		if to.Protocol == "https" {
			address += ":" + Utility.ToString(to.PortHttps)
		} else {
			address += ":" + Utility.ToString(to.PortHttp)
		}
		handleRequestAndRedirect(address, w, r)
		return
	}

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

	// globule.peers.

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

	// Other policies...
	(*w).Header().Set("Cross-Origin-Resource-Policy", "cross-origin")
}

/**
 * Sign ca certificate request and return a certificate.
 */
func signCaCertificateHandler(w http.ResponseWriter, r *http.Request) {
	// Handle the prefligth oprions...
	setupResponse(&w, r)

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	redirect, to := redirectTo(r.Host)
	if redirect {
		address := to.Domain
		if to.Protocol == "https" {
			address += ":" + Utility.ToString(to.PortHttps)
		} else {
			address += ":" + Utility.ToString(to.PortHttp)
		}

		handleRequestAndRedirect(address, w, r)
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
		if strings.HasPrefix(path, public[i]) {
			return true
		}
	}

	return false
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
 * Index video handler...
 */
func IndexVideoHandler(w http.ResponseWriter, r *http.Request) {
	// Handle the prefligth oprions...
	setupResponse(&w, r)

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	redirect, to := redirectTo(r.Host)

	if redirect {
		address := to.Domain
		if to.Protocol == "https" {
			address += ":" + Utility.ToString(to.PortHttps)
		} else {
			address += ":" + Utility.ToString(to.PortHttp)
		}

		handleRequestAndRedirect(address, w, r)
		return
	}

	var err error

	// Get the path where to upload the file.
	video_path := r.URL.Query().Get("video-path")
	video_path = strings.ReplaceAll(video_path, "\\", "/")

	// If application is defined.
	token := r.Header.Get("token")

	// If the header dosent contain the required values i I will try to get it from the
	// http query instead...
	if len(token) == 0 {
		// the token can be given by the url directly...
		token = r.URL.Query().Get("token")
	}

	// here in case of file uploaded from other website like pornhub...
	video_url := r.URL.Query().Get("video-url")
	index_path := r.URL.Query().Get("index-path")

	if len(video_url) > 0 {

		// Now if the os is windows I will remove the leading /
		if len(video_path) > 3 {
			if runtime.GOOS == "windows" && video_path[0] == '/' && video_path[2] == ':' {
				video_path = video_path[1:]
			}
		}

		if strings.HasPrefix(video_path, "/users") || strings.HasPrefix(video_path, "/applications") {
			video_path = strings.ReplaceAll(globule.data+"/files"+video_path, "\\", "/")
		} else if !isPublic(video_path) {
			video_path = strings.ReplaceAll(globule.webRoot+video_path, "\\", "/")
		}

		// Scrapper...
		if strings.Contains(video_url, "pornhub") {
			err = indexPornhubVideo(token, video_url, index_path, r.URL.Query().Get("video-path"))
		} else if strings.Contains(video_url, "xnxx") {
			err = indexXnxxVideo(token, video_url, index_path, r.URL.Query().Get("video-path"), video_path)
		} else if strings.Contains(video_url, "xvideo") {
			err = indexXvideosVideo(token, video_url, index_path, r.URL.Query().Get("video-path"), video_path)
		} else if strings.Contains(video_url, "xhamster") {
			err = indexXhamsterVideo(token, video_url, index_path, r.URL.Query().Get("video-path"), video_path)
		} else if strings.Contains(video_url, "youtube") {
			err = indexYoutubeVideo(token, video_url, index_path, r.URL.Query().Get("video-path"), video_path)
		}

		if err != nil {
			http.Error(w, "Fail to upload the file "+video_url+" with error "+err.Error(), http.StatusExpectationFailed)
		}

	}
}

func toBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func readMetadata(path string) (map[string]interface{}, error) {

	f_, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	m, err := tag.ReadFrom(f_)
	var metadata map[string]interface{}

	if err == nil {
		metadata = make(map[string]interface{})
		metadata["Album"] = m.Album()
		metadata["AlbumArtist"] = m.AlbumArtist()
		metadata["Artist"] = m.Artist()
		metadata["Comment"] = m.Comment()
		metadata["Composer"] = m.Composer()
		metadata["FileType"] = m.FileType()
		metadata["Format"] = m.Format()
		metadata["Genre"] = m.Genre()
		metadata["Lyrics"] = m.Lyrics()
		metadata["Picture"] = m.Picture()
		metadata["Raw"] = m.Raw()
		metadata["Title"] = m.Title()
		metadata["Year"] = m.Year()
		
		metadata["DisckNumber"], _ = m.Disc()
		_, metadata["DiscTotal"] = m.Disc()

		metadata["TrackNumber"], _ = m.Track()
		_, metadata["TrackTotal"] = m.Track()

		if m.Picture() != nil {

			// Determine the content type of the image file
			mimeType := m.Picture().MIMEType

			// Prepend the appropriate URI scheme header depending
			fileName := Utility.RandomUUID()

			
			// on the MIME type
			switch mimeType {
			case "image/jpg":
				fileName += ".jpg"
			case "image/jpeg":
				fileName += ".jpg"
			case "image/png":
				fileName += ".png"
			}

			imagePath := os.TempDir() + "/" + fileName
			defer os.Remove(imagePath)

			os.WriteFile(imagePath, m.Picture().Data, 0664)

			if Utility.Exists(imagePath) {
				metadata["ImageUrl"], _ = Utility.CreateThumbnail(imagePath, 300, 300)
			}
		}
	} else {
		return nil, err
	}
	return metadata, nil
}

/**
 * Index audio handler.
 */
 func IndexAudioHandler(w http.ResponseWriter, r *http.Request) {
	// Handle the prefligth oprions...
	setupResponse(&w, r)

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	redirect, to := redirectTo(r.Host)

	if redirect {
		address := to.Domain
		if to.Protocol == "https" {
			address += ":" + Utility.ToString(to.PortHttps)
		} else {
			address += ":" + Utility.ToString(to.PortHttp)
		}

		handleRequestAndRedirect(address, w, r)
		return
	}


	// Get the path where to upload the file.
	audio_path := r.URL.Query().Get("audio-path")
	audio_path = strings.ReplaceAll(audio_path, "\\", "/")

	// If application is defined.
	token := r.Header.Get("token")

	// If the header dosent contain the required values i I will try to get it from the
	// http query instead...
	if len(token) == 0 {
		// the token can be given by the url directly...
		token = r.URL.Query().Get("token")
	}

	// here in case of file uploaded from other website like pornhub...
	audio_url := r.URL.Query().Get("audio-url")
	index_path := r.URL.Query().Get("index-path")

	if len(audio_url) > 0 {

		// Now if the os is windows I will remove the leading /
		if len(audio_path) > 3 {
			if runtime.GOOS == "windows" && audio_path[0] == '/' && audio_path[2] == ':' {
				audio_path = audio_path[1:]
			}
		}

		if strings.Contains(audio_path, "/users/") || strings.Contains(audio_path, "/applications/") {
			audio_path = strings.ReplaceAll(globule.data+"/files"+audio_path, "\\", "/")
		} else if !isPublic(audio_path) {
			audio_path = strings.ReplaceAll(globule.webRoot+audio_path, "\\", "/")
		}

		if !Utility.Exists(audio_path){
			fmt.Println("no file found with path ", audio_path)
			return
		}


		// So here I will take information from the metada...
		metadata, err := readMetadata(audio_path)
		if err != nil {
			fmt.Println("fail to index ", audio_path, " with error ", err)
			return;
		}

		// so here I go the metadata...
		title_client_, err := getTitleClient()
		if err != nil {
			fmt.Println("no title client found with error ", err)
			return
		}

		track := new(titlepb.Audio)
		track.ID = Utility.GenerateUUID(metadata["Album"].(string) + ":" + metadata["Title"].(string) + ":" + metadata["AlbumArtist"].(string))
		track.Album = metadata["Album"].(string)
		track.AlbumArtist = metadata["AlbumArtist"].(string)
		track.Artist = metadata["Artist"].(string)
		track.Comment = metadata["Comment"].(string)
		track.Composer = metadata["Composer"].(string)

		if metadata["Genres"] != nil {
			fmt.Println("---------> file genre: ", metadata["Genres"] )
			track.Genres = metadata["Genres"].([]string)
		}

		track.Lyrics = metadata["Lyrics"].(string)
		track.Title = metadata["Title"].(string)
		track.Year = int32(Utility.ToInt(metadata["Year"]))
		track.DiscNumber = int32(Utility.ToInt(metadata["DiscNumber"]))
		track.DiscTotal = int32(Utility.ToInt(metadata["DiscTotal"]))
		track.TrackNumber = int32(Utility.ToInt(metadata["TrackNumber"]))
		track.TrackTotal = int32(Utility.ToInt(metadata["TrackTotal"]))

		imageUrl := ""
		if metadata["ImageUrl"] != nil {
			imageUrl = metadata["ImageUrl"].(string)
		}

		track.Poster = &titlepb.Poster{ID: track.ID, URL: "", TitleId: track.ID, ContentUrl: imageUrl}


		err = title_client_.CreateAudio(token, index_path, track)
		if err != nil {
			fmt.Println("fail to create audio with error ", err)
			return
		}else{
			fmt.Println(metadata)
		}

		// Now I will associate the file.
		err = title_client_.AssociateFileWithTitle(index_path, track.ID, audio_path)
		if err != nil {
			fmt.Println("fail to associate audio", track.Title, " with file ", audio_path, " error: ", err)
			return
		}
	}
}


/**
 * This code is use to upload a file into the tmp directory of the server
 * via http request.
 */
func FileUploadHandler(w http.ResponseWriter, r *http.Request) {

	setupResponse(&w, r)

	// Handle the prefligth oprions...
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	redirect, to := redirectTo(r.Host)

	if redirect {
		address := to.Domain
		if to.Protocol == "https" {
			address += ":" + Utility.ToString(to.PortHttps)
		} else {
			address += ":" + Utility.ToString(to.PortHttp)
		}

		handleRequestAndRedirect(address, w, r)
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
			var err error
			hasAccess, err = globule.validateAction("/file.FileService/FileUploadHandler", application, rbacpb.SubjectType_APPLICATION, infos)
			if hasAccess && err == nil {
				hasAccess, hasAccessDenied, err = globule.validateAccess(application, rbacpb.SubjectType_APPLICATION, "write", path)
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
		}else{
			fmt.Println("fail to validate token with error ", err.Error())
			http.Error(w, "fail to validate token with error " + err.Error(), http.StatusUnauthorized)
			return
		}
	} else {
		fmt.Println("no token was given!")
	}

	if len(user) != 0 {
		var err error
		if !hasAccess {
			hasAccess, err = globule.validateAction("/file.FileService/FileUploadHandler", user, rbacpb.SubjectType_ACCOUNT, infos)
		}
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
				http.Error(w, user + " has no space available to copy file "+path_+" allocated space and try again.", http.StatusUnauthorized)
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
			fmt.Println("714")
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

// Custom file server implementation.
func ServeFileHandler(w http.ResponseWriter, r *http.Request) {

	redirect, to := redirectTo(r.Host)

	if redirect {
		address := to.Domain
		if to.Protocol == "https" {
			address += ":" + Utility.ToString(to.PortHttps)
		} else {
			address += ":" + Utility.ToString(to.PortHttp)
		}
		handleRequestAndRedirect(address, w, r)
		return
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
		hasAccess = true
	}

	// this is the ca certificate use to sign client certificate.
	if rqst_path == "/ca.crt" {
		name = globule.creds + rqst_path
	}

	hasAccessDenied := false
	var err error
	var userId string

	// Here I will validate applications...
	if len(application) != 0 && !hasAccess {
		hasAccess, hasAccessDenied, err = globule.validateAccess(application, rbacpb.SubjectType_APPLICATION, "read", rqst_path)
	}

	if len(token) != 0 && !hasAccess {
		var claims *security.Claims
		claims, err = security.ValidateToken(token)
		userId = claims.Id + "@" + claims.UserDomain
		if err == nil {
			hasAccess, hasAccessDenied, err = globule.validateAccess(userId, rbacpb.SubjectType_ACCOUNT, "read", rqst_path)
		}
	}

	// validate ressource access...
	if !hasAccess || hasAccessDenied || err != nil {
		log.Println(err)
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

/**
 * Return a list of IMDB titles from a keyword...
 */
func getImdbTitlesHanldler(w http.ResponseWriter, r *http.Request) {
	// Receive http request...
	redirect, to := redirectTo(r.Host)
	if redirect {
		address := to.Domain
		if to.Protocol == "https" {
			address += ":" + Utility.ToString(to.PortHttps)
		} else {
			address += ":" + Utility.ToString(to.PortHttp)
		}
		handleRequestAndRedirect(address, w, r)
		return
	}

	// if the host is not the same...
	query := r.URL.Query().Get("query") // the csr in base64

	client := http.DefaultClient
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
	json.NewEncoder(w).Encode(titles)
}

/**
 * Return a json object with the movie information from imdb.
 */
func getImdbTitleHanldler(w http.ResponseWriter, r *http.Request) {
	// Receive http request...
	redirect, to := redirectTo(r.Host)
	if redirect {
		address := to.Domain
		if to.Protocol == "https" {
			address += ":" + Utility.ToString(to.PortHttps)
		} else {
			address += ":" + Utility.ToString(to.PortHttp)
		}

		handleRequestAndRedirect(address, w, r)
		return
	}

	id := r.URL.Query().Get("id") // the csr in base64

	fmt.Println("get imdb info for ", id)

	client := http.DefaultClient
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
		}else{
			fmt.Println("fail to retreive episode info ", err)
		}
	}

	json.NewEncoder(w).Encode(title_)
}
