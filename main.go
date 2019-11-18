package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/davecourtois/Globular/admin"
	"github.com/davecourtois/Utility"
)

func main() {
	g := NewGlobule()
	if len(os.Args) > 1 {

		// Subcommands
		installCommand := flag.NewFlagSet("install", flag.ExitOnError)
		deployCommand := flag.NewFlagSet("deploy", flag.ExitOnError)

		path := installCommand.String("path", "", "You must specefied the intallation path. (Required)")

		name := deployCommand.String("name", "", "You must specify an application name. (Required)")
		applicationPath := deployCommand.String("path", "", "You must specify the path that contain the source (bundle.js, index.html...) of the application to deploy. (Required)")
		user := deployCommand.String("u", "", "The user name. (Required)")
		pwd := deployCommand.String("p", "", "The user password. (Required)")
		address := deployCommand.String("a", "", "The address (domain:port default admin port is 10001) of the server where to install the appliction (Required)")

		switch os.Args[1] {
		case "install":
			installCommand.Parse(os.Args[2:])
		case "deploy":
			deployCommand.Parse(os.Args[2:])
		default:
			flag.PrintDefaults()
			os.Exit(1)
		}

		// Check if the command was parsed
		if installCommand.Parsed() {
			// Required Flags
			if *path == "" {
				installCommand.PrintDefaults()
				os.Exit(1)
			}

			install(g, *path)
		}

		if deployCommand.Parsed() {
			// Required Flags
			if *applicationPath == "" {
				deployCommand.PrintDefaults()
				os.Exit(1)
			}

			if *name == "" {
				deployCommand.PrintDefaults()
				os.Exit(1)
			}

			if *user == "" {
				deployCommand.PrintDefaults()
				os.Exit(1)
			}

			if *pwd == "" {
				deployCommand.PrintDefaults()
				os.Exit(1)
			}

			if *address == "" {
				deployCommand.PrintDefaults()
				os.Exit(1)
			}

			deploy(g, *name, *applicationPath, *address, *user, *pwd)
		}

	} else {
		g.Listen()
	}
}

/**
 * That function can be use to deploy an application on the server...
 */
func deploy(g *Globule, name string, path string, address string, user string, pwd string) {

	log.Println("deploy application...", name)
	var port int
	var domain string

	if strings.Index(address, ":") > 0 {
		domain = strings.Split(address, ":")[0]
		port = Utility.ToInt(strings.Split(address, ":")[1])
	} else {
		port = g.AdminPort
		domain = address
	}

	// hasTLS bool, keyFile string, certFile string, caFile string, token string

	// First of all I will create a client...
	log.Println("--> domain:", domain, port)
	var keyPath string
	var certPath string
	var caPath string
	var tls bool
	if domain == "localhost" {
		keyPath = g.creds + string(os.PathSeparator) + "client.pem"
		certPath = g.creds + string(os.PathSeparator) + "client.crt"
		caPath = g.creds + string(os.PathSeparator) + "ca.crt"
	} else {
		creds := g.creds + string(os.PathSeparator) + domain
		Utility.CreateDirIfNotExist(creds)
		log.Println("---> here I will generate the certificate for the client... ", creds)
		// So here I will get the ca.crt file...
		// Get the data
		resp, err := http.Get("https://" + address + "/ca.crt")

		if err != nil {
			log.Println(err)
			return
		}

		defer resp.Body.Close()

		// Create the ca file.
		out, err := os.Create(creds + string(os.PathSeparator) + "ca.crt")
		if err != nil {
			log.Println(err)
			return
		}

		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Println(err)
			return
		}

		// Now I will generate the certificate for the client...
		// Step 1: Generate client private key.
		err = g.GenerateClientPrivateKey(creds, pwd)
		if err != nil {
			log.Println(err)
			return
		}

		// Step 2: Generate the client signing request.
		err = g.GenerateClientCertificateSigningRequest(creds, pwd, domain)
		if err != nil {
			log.Println(err)
			return
		}

		// Step 3: Generate client signed certificate.
		/* err = g.GenerateSignedClientCertificate(creds, pwd, g.CertExpirationDelay)
		if err != nil {
			log.Println(err)
			return
		}*/

		// Now ask the ca to sign the certificate.

		// Step 4: Convert to pem format.
		err = g.KeyToPem("client", creds, pwd)
		if err != nil {
			log.Println(err)
			return
		}

		// Set tls to true.
		tls = true

		keyPath = creds + string(os.PathSeparator) + "client.pem"
		certPath = creds + string(os.PathSeparator) + "client.crt"
		caPath = creds + string(os.PathSeparator) + "ca.crt"
		log.Println("-----------------------------------> ", keyPath)
		/*keyPath = g.creds + string(os.PathSeparator) + "client.pem"
		certPath = g.creds + string(os.PathSeparator) + "client.crt"
		caPath = g.creds + string(os.PathSeparator) + "ca.crt"*/
	}

	// Token...
	token := ""

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	client := admin.NewAdmin_Client(domain, g.AdminPort, tls, keyPath, certPath, caPath, token)

	err := client.DeployApplication(name, path)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Application", name, "was deployed successfully!")
}

/**
 * That function is use to install globular from the development environnement.
 * The server must have run at least once before that command is call. Each service must
 * have been run at least one to appear in the installation.
 */
func install(g *Globule, path string) {
	// That function is use to install globular at a given repository.
	fmt.Println("install globular in directory: ", path)

	Utility.CreateDirIfNotExist(path)

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	// Here I will copy the proxy.
	globularExec := os.Args[0]
	if string(os.PathSeparator) == "\\" && !strings.HasSuffix(globularExec, ".exe") {
		globularExec += ".exe" // in case of windows
	}

	err := Utility.Copy(dir+string(os.PathSeparator)+globularExec, path+string(os.PathSeparator)+globularExec)
	if err != nil {
		fmt.Println(err)
	}

	err = os.Chmod(path+string(os.PathSeparator)+globularExec, 0755)
	if err != nil {
		fmt.Println(err)
	}

	// Copy the bin file from globular
	Utility.CreateDirIfNotExist(path + string(os.PathSeparator) + "bin")
	err = Utility.CopyDir(dir+string(os.PathSeparator)+"bin", path+string(os.PathSeparator)+"bin")
	if err != nil {
		log.Panicln("--> fail to copy bin ", err)
	}

	// Change the files permission to add execute write.
	files, err := ioutil.ReadDir(path + string(os.PathSeparator) + "bin")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if !f.IsDir() {
			err = os.Chmod(path+string(os.PathSeparator)+"bin"+string(os.PathSeparator)+f.Name(), 0755)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	// install services...
	for id, service := range g.Services {
		s := service.(map[string]interface{})
		name := s["Name"].(string)
		// I will read the configuration file to have nessecary service information
		// to be able to create the path.
		configPath := dir + s["configPath"].(string)
		if Utility.Exists(configPath) {
			bytes, err := ioutil.ReadFile(configPath)
			config := make(map[string]interface{}, 0)
			json.Unmarshal(bytes, &config)

			if err == nil {
				// set the name.
				if config["PublisherId"] != nil && config["Version"] != nil {
					var serviceDir = path + string(os.PathSeparator) + "globular_services"
					if len(config["PublisherId"].(string)) == 0 {
						serviceDir += string(os.PathSeparator) + config["Domain"].(string) + string(os.PathSeparator) + id + string(os.PathSeparator) + config["Version"].(string)
					} else {
						serviceDir += string(os.PathSeparator) + config["PublisherId"].(string) + string(os.PathSeparator) + id + string(os.PathSeparator) + config["Version"].(string)
					}
					Utility.CreateDirIfNotExist(serviceDir)

					execPath := dir + s["servicePath"].(string)
					destPath := serviceDir + string(os.PathSeparator) + name
					if string(os.PathSeparator) == "\\" {
						execPath += ".exe" // in case of windows
						destPath += ".exe"
					}

					err := Utility.Copy(execPath, destPath)
					if err != nil {

						log.Panicln(execPath, destPath, err)
					}

					// Set execute permission
					err = os.Chmod(destPath, 0755)
					if err != nil {
						fmt.Println(err)
					}

					// Copy the service config file.
					if Utility.Exists(configPath) {
						Utility.Copy(configPath, serviceDir+string(os.PathSeparator)+"config.json")
					}

					// Copy the proto file.
					protoPath := dir + s["protoPath"].(string)
					if Utility.Exists(protoPath) {
						Utility.Copy(protoPath, serviceDir+string(os.PathSeparator)+protoPath[strings.LastIndex(protoPath, "/"):])
					}

				} else {
					log.Println("service", id, "configuration is incomplete!")
				}
			} else {
				log.Println("service", id, "has no configuration!")
			}
		}
	}
}
