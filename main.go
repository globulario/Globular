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

	// Authenticate the user in order to get the token
	ressource_client := ressource.NewRessource_Client(address, "ressource")
	_, err := ressource_client.Authenticate(user, pwd)
	if err != nil {
		log.Panicln(err)
		return
	}

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	admin_client := admin.NewAdmin_Client(address, "admin")
	err = admin_client.DeployApplication(name, path)
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
