package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/davecourtois/Globular/admin"
	"github.com/davecourtois/Globular/ressource"
	"github.com/davecourtois/Utility"
)

func main() {
	g := NewGlobule()
	if len(os.Args) > 1 {

		// Subcommands

		// Intall command
		installCommand := flag.NewFlagSet("install", flag.ExitOnError)
		installCommand_path := installCommand.String("path", "", "You must specefied the intallation path. (Required)")

		// Deploy command
		deployCommand := flag.NewFlagSet("deploy", flag.ExitOnError)
		deployCommand_name := deployCommand.String("name", "", "You must specify an application name. (Required)")
		deployCommand_path := deployCommand.String("path", "", "You must specify the path that contain the source (bundle.js, index.html...) of the application to deploy. (Required)")
		deployCommand_user := deployCommand.String("u", "", "The user name. (Required)")
		deployCommand_pwd := deployCommand.String("p", "", "The user password. (Required)")
		deployCommand_address := deployCommand.String("a", "", "The domain of the server where to install the appliction (Required)")

		// Publish command.
		publishCommand := flag.NewFlagSet("publish", flag.ExitOnError)
		publishCommand_name := publishCommand.String("name", "", "You must specify an service name. (Required)")
		publishCommand_publisher_id := publishCommand.String("publisher", "", "The publisher id. (Required)")
		publishCommand_path := publishCommand.String("path", "", "You must specify the path that contain the config.json, .proto and all dependcies require by the service to run. (Required)")
		publishCommand_discovery := publishCommand.String("discovery", "", "You must specified the domain of the discovery service where to publish your service (Required)")
		publishCommand_repository := publishCommand.String("repository", "", "You must specified the domain of the repository service where to publish your service (Required)")
		publishCommand_user := publishCommand.String("u", "", "The user name. (Required)")
		publishCommand_pwd := publishCommand.String("p", "", "The user password. (Required)")
		publishCommand_address := publishCommand.String("a", "", "The domain of the server where to install the appliction (Required)")
		publishCommand_description := publishCommand.String("description", "", "You must specify a service description. (Required)")
		publishCommand_version := publishCommand.String("version", "", "You must specified the version of the service. (Required)")
		publishCommand_keywords := publishCommand.String("keywords", "", "You must give keywords. (Required)")
		publishCommand_plaform := publishCommand.Int("platform", 0, "1 = Linux32; 2 = Linux64; 3 = windows32; 4 = windows64 (Required)")

		switch os.Args[1] {
		case "install":
			installCommand.Parse(os.Args[2:])
		case "deploy":
			deployCommand.Parse(os.Args[2:])
		case "publish":
			publishCommand.Parse(os.Args[2:])
		default:
			flag.PrintDefaults()
			os.Exit(1)
		}

		// Check if the command was parsed
		if installCommand.Parsed() {
			// Required Flags
			if *installCommand_path == "" {
				installCommand.PrintDefaults()
				os.Exit(1)
			}

			install(g, *installCommand_path)
		}

		if deployCommand.Parsed() {
			// Required Flags
			if *deployCommand_path == "" {
				fmt.Print("No application 'dist' path was given")
				deployCommand.PrintDefaults()
				os.Exit(1)
			}
			if *deployCommand_name == "" {
				fmt.Print("No applicaiton 'name' was given")
				deployCommand.PrintDefaults()
				os.Exit(1)
			}

			if *deployCommand_user == "" {
				fmt.Print("You must authenticate yourself")
				deployCommand.PrintDefaults()
				os.Exit(1)
			}

			if *deployCommand_pwd == "" {
				fmt.Print("You must specifie the user password.")
				deployCommand.PrintDefaults()
				os.Exit(1)
			}

			if *deployCommand_address == "" {
				fmt.Print("You must sepcie the server address")
				deployCommand.PrintDefaults()
				os.Exit(1)
			}

			deploy(g, *deployCommand_name, *deployCommand_path, *deployCommand_address, *deployCommand_user, *deployCommand_pwd)
		}

		if publishCommand.Parsed() {
			if *publishCommand_name == "" {
				publishCommand.PrintDefaults()
				os.Exit(1)
			}

			if *publishCommand_publisher_id == "" {
				publishCommand.PrintDefaults()
				os.Exit(1)
			}

			if *publishCommand_path == "" {
				publishCommand.PrintDefaults()
				os.Exit(1)
			}

			if *publishCommand_discovery == "" {
				publishCommand.PrintDefaults()
				os.Exit(1)
			}

			if *publishCommand_repository == "" {
				publishCommand.PrintDefaults()
				os.Exit(1)
			}

			if *publishCommand_user == "" {
				publishCommand.PrintDefaults()
				os.Exit(1)
			}

			if *publishCommand_pwd == "" {
				publishCommand.PrintDefaults()
				os.Exit(1)
			}

			if *publishCommand_address == "" {
				publishCommand.PrintDefaults()
				os.Exit(1)
			}

			if *publishCommand_description == "" {
				publishCommand.PrintDefaults()
				os.Exit(1)
			}

			if *publishCommand_keywords == "" {
				publishCommand.PrintDefaults()
				os.Exit(1)
			}

			if *publishCommand_version == "" {
				publishCommand.PrintDefaults()
				os.Exit(1)
			}

			if *publishCommand_plaform == -1 {
				publishCommand.PrintDefaults()
				os.Exit(1)
			}

			keywords := strings.Split(*publishCommand_keywords, ",")
			for i := 0; i < len(keywords); i++ {
				keywords[i] = strings.TrimSpace(keywords[i])
			}

			// Pulish the services.
			publish(g, *publishCommand_path, *publishCommand_name, *publishCommand_publisher_id, *publishCommand_discovery, *publishCommand_repository, *publishCommand_description, *publishCommand_version, int32(*publishCommand_plaform-1), keywords, *publishCommand_address, *publishCommand_user, *publishCommand_pwd)
		}

	} else {
		g.Serve()
	}
}

/**
 * That function can be use to deploy an application on the server...
 */
func deploy(g *Globule, name string, path string, address string, user string, pwd string) {

	log.Println("deploy application...", name)

	// Authenticate the user in order to get the token
	ressource_client := ressource.NewRessource_Client(address, "ressource")
	token, err := ressource_client.Authenticate(user, pwd)
	if err != nil {
		log.Panicln(err)
		return
	}

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	admin_client := admin.NewAdmin_Client(address, "admin")
	err = admin_client.DeployApplication(name, path, token)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Application", name, "was deployed successfully!")
}

/**
 * That function is use to publish a service on the network.
 * The service must contain
 * - A .proto file that define it gRPC services interface.
 * - A config.json file that define it service configuration.
 * All file in the service directory will be part of the package, so take
 * care to include all dependencies, dll... to be sure your services will run
 * as expected.
 */
func publish(g *Globule, path string, serviceId string, publisherId string, discoveryId string, repositoryId string, description string, version string, platform int32, keywords []string, address string, user string, pwd string) {
	log.Println("publish service...", serviceId, "at address", address)

	// Authenticate the user in order to get the token
	ressource_client := ressource.NewRessource_Client(address, "ressource")
	token, err := ressource_client.Authenticate(user, pwd)
	if err != nil {
		log.Panicln(err)
		return
	}

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	admin_client := admin.NewAdmin_Client(address, "admin")

	// first of all I will create and upload the package on the discovery...
	path_, err := admin_client.UploadServicePackage(path, publisherId, serviceId, version, token)
	if err != nil {
		log.Panicln(err)
		return
	}
	if platform < 0 {
		log.Println("Plaform must be a number from 1 to 4")
		return
	}

	err = admin_client.PublishService(path_, serviceId, publisherId, discoveryId, repositoryId, description, version, platform, keywords, token)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Service", serviceId, "was pulbish successfully!")
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

	// Now I will copy the prototype files of the internal gRPC service
	// admin, ressource, ca and services.
	Utility.CreateDirIfNotExist(path + string(os.PathSeparator) + "proto")
	Utility.CopyFile(dir+string(os.PathSeparator)+"admin"+string(os.PathSeparator)+"admin.proto", path+string(os.PathSeparator)+"proto"+string(os.PathSeparator)+"admin.proto")
	Utility.CopyFile(dir+string(os.PathSeparator)+"ca"+string(os.PathSeparator)+"ca.proto", path+string(os.PathSeparator)+"proto"+string(os.PathSeparator)+"ca.proto")
	Utility.CopyFile(dir+string(os.PathSeparator)+"ressource"+string(os.PathSeparator)+"ressource.proto", path+string(os.PathSeparator)+"proto"+string(os.PathSeparator)+"ressource.proto")
	Utility.CopyFile(dir+string(os.PathSeparator)+"services"+string(os.PathSeparator)+"services.proto", path+string(os.PathSeparator)+"proto"+string(os.PathSeparator)+"services.proto")

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
		if s["Name"] != nil {
			name := s["Name"].(string)
			// I will read the configuration file to have nessecary service information
			// to be able to create the path.
			configPath := dir + s["configPath"].(string)
			if Utility.Exists(configPath) {
				log.Println("install service ", name)
				bytes, err := ioutil.ReadFile(configPath)
				config := make(map[string]interface{}, 0)
				json.Unmarshal(bytes, &config)

				if err == nil {
					// set the name.
					if config["PublisherId"] != nil && config["Version"] != nil && s["protoPath"] != nil && s["servicePath"] != nil {

						execPath := dir + s["servicePath"].(string)
						protoPath := dir + s["protoPath"].(string)

						if Utility.Exists(execPath) && Utility.Exists(protoPath) {

							var serviceDir = path + string(os.PathSeparator) + "globular_services"
							if len(config["PublisherId"].(string)) == 0 {
								serviceDir += string(os.PathSeparator) + config["Domain"].(string) + string(os.PathSeparator) + id + string(os.PathSeparator) + config["Version"].(string)
							} else {
								serviceDir += string(os.PathSeparator) + config["PublisherId"].(string) + string(os.PathSeparator) + id + string(os.PathSeparator) + config["Version"].(string)
							}

							Utility.CreateDirIfNotExist(serviceDir)
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
							if Utility.Exists(protoPath) {
								Utility.Copy(protoPath, serviceDir+string(os.PathSeparator)+protoPath[strings.LastIndex(protoPath, "/"):])
							}
						}
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
