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

	//globular "github.com/davecourtois/Globular/services/golang/globular_client"
	"github.com/davecourtois/Globular/services/golang/admin/admin_client"
	"github.com/davecourtois/Globular/services/golang/ressource/ressource_client"
	"github.com/davecourtois/Utility"
	"github.com/kardianos/service"
)

// This is use to display information to external service manager.
var logger service.Logger

func (g *Globule) Start(s service.Service) error {
	if service.Interactive() {
		logger.Info("Running in terminal.")
	} else {
		logger.Info("Running under service manager.")
	}
	g.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	go g.run()
	return nil
}

func (g *Globule) run() error {
	logger.Infof("Starting Globular as %v.", service.Platform())

	// start globular and wait on exit chan...
	go func() {
		g.Serve()
	}()

	for {
		select {
		case <-g.exit:
			return nil
		}
	}
}

func (g *Globule) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	logger.Info("Globular is stopping!")

	// Stop external services.
	g.stopServices()

	// Stop internal services
	g.stopInternalServices()

	logger.Infof("Globular has been stopped!")

	close(g.exit)
	return nil
}

func main() {

	g := NewGlobule()
	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()

	options := make(service.KeyValue)
	options["Restart"] = "on-success"
	options["SuccessExitStatus"] = "1 2 8 SIGKILL"

	svcConfig := &service.Config{
		Name:         "Globular",
		DisplayName:  "Globular",
		Description:  "gRPC service managers",
		Dependencies: []string{},
		Option:       options,
	}

	s, err := service.New(g, svcConfig)

	if len(os.Args) > 1 {

		// Subcommands

		// Intall globular as service/demon
		installCommand := flag.NewFlagSet("install", flag.ExitOnError)
		installCommand_name := installCommand.String("name", "", "The display name of globular service.")

		// Uninstall globular as service.
		unstallCommand := flag.NewFlagSet("uninstall", flag.ExitOnError)

		// Package development environnement into a given
		distCommand := flag.NewFlagSet("dist", flag.ExitOnError)
		distCommand_path := distCommand.String("path", "", "You must specify the dist path. (Required)")

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
		publishCommand_plaform := publishCommand.String("platform", "", "One of linux32, linux64, win32, win64 (Required)")

		switch os.Args[1] {
		case "package":
			distCommand.Parse(os.Args[2:])
		case "deploy":
			deployCommand.Parse(os.Args[2:])
		case "publish":
			publishCommand.Parse(os.Args[2:])
		case "install":
			installCommand.Parse(os.Args[2:])
		case "uninstall":
			unstallCommand.Parse(os.Args[2:])
		default:
			flag.PrintDefaults()
			os.Exit(1)
		}

		// Check if the command was parsed
		if installCommand.Parsed() {
			if *installCommand_name != "" {
				svcConfig.DisplayName = *installCommand_name
				s, _ = service.New(g, svcConfig)
			}
			// Required Flags
			err := s.Install()
			if err == nil {
				log.Println("Globular service is now installed!")
			} else {
				log.Println(err)
			}
		}

		if unstallCommand.Parsed() {
			// Required Flags
			err := s.Uninstall()
			if err == nil {
				log.Println("Globular service is now removed!")
			} else {
				log.Println(err)
			}
		}

		if distCommand.Parsed() {
			// Required Flags
			if *distCommand_path == "" {
				distCommand.PrintDefaults()
				os.Exit(1)
			}

			install(g, *distCommand_path)
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

			var platform int32
			if *publishCommand_plaform == "linux32" {
				platform = 0
			} else if *publishCommand_plaform == "linux64" {
				platform = 1
			} else if *publishCommand_plaform == "win32" {
				platform = 2
			} else if *publishCommand_plaform == "win64" {
				platform = 3
			} else {
				publishCommand.PrintDefaults()
				os.Exit(1)
			}

			keywords := strings.Split(*publishCommand_keywords, ",")
			for i := 0; i < len(keywords); i++ {
				keywords[i] = strings.TrimSpace(keywords[i])
			}

			// Pulish the services.
			publish(g, *publishCommand_path, *publishCommand_name, *publishCommand_publisher_id, *publishCommand_discovery, *publishCommand_repository, *publishCommand_description, *publishCommand_version, platform, keywords, *publishCommand_address, *publishCommand_user, *publishCommand_pwd)
		}

	} else {

		if err != nil {
			log.Fatal(err)
		}
		errs := make(chan error, 5)
		logger, err = s.Logger(errs)
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			for {
				err := <-errs
				if err != nil {
					log.Print(err)
				}
			}
		}()

		if len(*svcFlag) != 0 {
			err := service.Control(s, *svcFlag)
			if err != nil {
				log.Printf("Valid actions: %q\n", service.ControlAction)
				log.Fatal(err)
			}
			return
		}
		err = s.Run()
		if err != nil {
			logger.Error(err)
		}
	}
}

/**
 * Service interface use to run as Windows Service or Linux deamon...
 */

/**
 * That function can be use to deploy an application on the server...
 */
func deploy(g *Globule, name string, path string, address string, user string, pwd string) error {

	log.Println("deploy application...", name, " to address ", address)

	// Authenticate the user in order to get the token
	ressource_client, err := ressource_client.NewRessourceService_Client(address, "ressource.RessourceService")
	if err != nil {
		log.Println("fail to access ressource service at "+address+" with error ", err)
		return err
	}

	token, err := ressource_client.Authenticate(user, pwd)
	if err != nil {
		log.Println("fail to authenticate user ", err)
		return err
	}

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	admin_client_, err := admin_client.NewAdminService_Client(address, "admin.AdminService") // create the ressource server.
	if err != nil {
		return err
	}

	_, err = admin_client_.DeployApplication(user, name, path, token, address)
	if err != nil {
		log.Println("Fail to deploy applicaiton with error:", err)
		return err
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	log.Println("Application", name, "was deployed successfully!")
	return nil
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
func publish(g *Globule, path string, serviceId string, publisherId string, discoveryId string, repositoryId string, description string, version string, platform int32, keywords []string, address string, user string, pwd string) error {
	log.Println("publish service...", serviceId, "at address", address)

	// Authenticate the user in order to get the token
	ressource_client_, err := ressource_client.NewRessourceService_Client(address, "ressource.RessourceService")
	if err != nil {
		log.Panicln(err)
		return err
	}

	token, err := ressource_client_.Authenticate(user, pwd)
	if err != nil {
		log.Println(err)
		return err
	}

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	admin_client_, err := admin_client.NewAdminService_Client(address, "admin.AdminService")
	if err != nil {
		log.Println(err)
		return err
	}

	// first of all I will create and upload the package on the discovery...
	path_, _, err := admin_client_.UploadServicePackage(path, publisherId, serviceId, version, token, address)
	if err != nil {
		log.Println(err)
		return err
	}

	err = admin_client_.PublishService(user, path_, serviceId, publisherId, discoveryId, repositoryId, description, version, platform, keywords, token, address)
	if err != nil {
		return err
	}

	log.Println("Service", serviceId, "was pulbish successfully!")
	return nil
}

/**
 * That function is use to install globular from the development environnement.
 * The server must have run at least once before that command is call. Each service must
 * have been run at least one to appear in the installation.
 */
/**

ADD /usr/local/lib/libplctag.so /usr/local/lib
ADD /usr/local/lib/libgrpc++.so /usr/local/lib
ADD /usr/local/lib/libgrpc++.so.1 /usr/local/lib
ADD /usr/local/lib/libgrpc++.so.1.20.0 /usr/local/lib
ADD /usr/local/lib/libprotobuf.so /usr/local/lib
ADD /usr/local/lib/libprotobuf.so.20 /usr/local/lib
ADD /usr/local/lib/libprotobuf.so.20.0.1 /usr/local/lib
ADD /usr/local/lib/libgrpc.so /usr/local/lib
ADD /usr/local/lib/libgrpc.so.7 /usr/local/lib
ADD /usr/local/lib/libgrpc.so.7.0.0 /usr/local/lib
ADD /usr/local/lib/libgpr.so /usr/local/lib
ADD /usr/local/lib/libgpr.so.7 /usr/local/lib
ADD /usr/local/lib/libgpr.so.7.0.0 /usr/local/lib
*/
func install(g *Globule, path string) {
	// That function is use to install globular at a given repository.
	fmt.Println("install globular in directory: ", path)

	// I will generate the docker files.
	dockerfile := `#-- Docker install. --
FROM ubuntu
RUN apt-get update && apt-get install -y gnupg2 \
    wget \
  && rm -rf /var/lib/apt/lists/*
RUN wget https://s3-eu-west-1.amazonaws.com/deb.robustperception.io/41EFC99D.gpg && apt-key add 41EFC99D.gpg
RUN apt-get update && apt-get install -y \
  build-essential \
  curl \
  mongodb

# -- Install prometheus
RUN wget https://github.com/prometheus/prometheus/releases/download/v2.17.0/prometheus-2.17.0.linux-amd64.tar.gz
RUN tar -xf prometheus-2.17.0.linux-amd64.tar.gz
RUN cp prometheus-2.17.0.linux-amd64/prometheus /usr/local/bin/
RUN cp prometheus-2.17.0.linux-amd64/promtool /usr/local/bin/
RUN cp -r prometheus-2.17.0.linux-amd64/consoles /etc/prometheus/
RUN cp -r prometheus-2.17.0.linux-amd64/console_libraries /etc/prometheus/
RUN rm -rf prometheus-2.17.0.linux-amd64*

# -- Install alert manager
RUN wget https://github.com/prometheus/alertmanager/releases/download/v0.20.0/alertmanager-0.20.0.linux-amd64.tar.gz
RUN tar -xf alertmanager-0.20.0.linux-amd64.tar.gz
RUN cp alertmanager-0.20.0.linux-amd64/alertmanager /usr/local/bin
RUN rm -rf alertmanager-0.20.0.linux-amd64*

# -- Install node exporter
RUN wget https://github.com/prometheus/node_exporter/releases/download/v0.18.1/node_exporter-0.18.1.linux-amd64.tar.gz
RUN tar -xf node_exporter-0.18.1.linux-amd64.tar.gz
RUN cp node_exporter-0.18.1.linux-amd64/node_exporter /usr/local/bin
RUN rm -rf node_exporter-0.18.1.linux-amd64*

# -- Install unix odbc drivers.
RUN curl http://www.unixodbc.org/unixODBC-2.3.7.tar.gz --output unixODBC-2.3.7.tar.gz
RUN tar -xvf unixODBC-2.3.7.tar.gz
RUN rm unixODBC-2.3.7.tar.gz

# -- Load all newly install libs.
RUN ldconfig

WORKDIR unixODBC-2.3.7
RUN ./configure && make all install clean && ldconfig && mkdir /globular && cd /globular
ADD Globular /globular
COPY bin /globular/bin
COPY proto /globular/proto
COPY services /globular/services
`

	Utility.CreateDirIfNotExist(path)

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	// Here I will copy the proxy.
	globularExec := os.Args[0]
	if "/" == "\\" && !strings.HasSuffix(globularExec, ".exe") {
		globularExec += ".exe" // in case of windows
	}

	err := Utility.Copy(dir+"/"+globularExec, path+"/"+globularExec)
	if err != nil {
		fmt.Println(err)
	}

	err = os.Chmod(path+"/"+globularExec, 0755)
	if err != nil {
		fmt.Println(err)
	}

	// Copy the bin file from globular
	Utility.CreateDirIfNotExist(path + "/" + "bin")
	err = Utility.CopyDir(dir+"/"+"bin", path+"/"+"bin")
	if err != nil {
		log.Panicln("--> fail to copy bin ", err)
	}

	// Change the files permission to add execute write.
	files, err := ioutil.ReadDir(path + "/" + "bin")
	if err != nil {
		log.Fatal(err)
	}

	// Now I will copy the prototype files of the internal gRPC service
	// admin, ressource, ca and services.
	Utility.CreateDirIfNotExist(path + "/" + "proto")
	Utility.CopyFile(dir+"/"+"admin"+"/"+"admin.proto", path+"/"+"proto"+"/"+"admin.proto")
	Utility.CopyFile(dir+"/"+"ca"+"/"+"ca.proto", path+"/"+"proto"+"/"+"ca.proto")
	Utility.CopyFile(dir+"/"+"ressource"+"/"+"ressource.proto", path+"/"+"proto"+"/"+"ressource.proto")
	Utility.CopyFile(dir+"/"+"services"+"/"+"services.proto", path+"/"+"proto"+"/"+"services.proto")

	for _, f := range files {
		if !f.IsDir() {
			err = os.Chmod(path+"/"+"bin"+"/"+f.Name(), 0755)
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
			if s["configPath"] != nil {
				configPath := s["configPath"].(string)

				if Utility.Exists(configPath) {
					log.Println("install service ", name)
					bytes, err := ioutil.ReadFile(configPath)
					config := make(map[string]interface{}, 0)
					json.Unmarshal(bytes, &config)

					if err == nil {

						// set the name.
						if config["PublisherId"] != nil && config["Version"] != nil && s["Proto"] != nil && s["Path"] != nil {

							execPath := s["Path"].(string)
							protoPath := s["Proto"].(string)

							if Utility.Exists(execPath) && Utility.Exists(protoPath) {
								var serviceDir = "services" + "/"
								if len(config["PublisherId"].(string)) == 0 {
									serviceDir += config["Domain"].(string) + "/" + name + "/" + config["Version"].(string)
								} else {
									serviceDir += config["PublisherId"].(string) + "/" + name + "/" + config["Version"].(string)
								}

								lastIndex := strings.LastIndex(execPath, "/")
								if lastIndex == -1 {
									lastIndex = strings.LastIndex(execPath, "/")
								}

								execName := execPath[lastIndex+1:]
								destPath := path + "/" + serviceDir + "/" + id + "/" + execName

								if Utility.Exists(execPath) {
									Utility.CreateDirIfNotExist(path + "/" + serviceDir + "/" + id)

									err := Utility.Copy(execPath, destPath)
									if err != nil {
										log.Panicln(execPath, destPath, err)
									}

									// Set execute permission
									err = os.Chmod(destPath, 0755)
									if err != nil {
										fmt.Println(err)
									}

									config["Path"] = destPath
									config["Proto"] = path + "/" + serviceDir + "/" + name + ".proto"

									// set the security values to nothing...
									config["CertAuthorityTrust"] = ""
									config["CertFile"] = ""
									config["KeyFile"] = ""
									config["TLS"] = false
									delete(config, "configPath")

									str, _ := Utility.ToJson(&config)
									ioutil.WriteFile(path+"/"+serviceDir+"/"+id+"/"+"config.json", []byte(str), 0644)

									// Copy the proto file.
									if Utility.Exists(protoPath) {
										Utility.Copy(protoPath, config["Proto"].(string))
									}
								} else {
									fmt.Println("executable not exist ", execPath)
								}
							} else if !Utility.Exists(execPath) {
								log.Println("no executable found at path " + execPath)
							} else if !Utility.Exists(protoPath) {
								log.Println("no proto file found at path " + protoPath)
							}
						} else if config["PublisherId"] == nil {
							fmt.Println("no publisher was define!")
						} else if config["Version"] == nil {
							fmt.Println("no version was define!")
						} else if s["Proto"] == nil {
							fmt.Println(" no proto file was found!")
						} else if s["Path"] != nil {
							fmt.Println("no executable was found!")
						}
					} else {
						fmt.Println(err)
					}
				} else {
					fmt.Println("service", name, ":", id, "configuration is incomplete!")
				}
			} else {
				fmt.Println("service", name, ":", id, "has no configuration!")
			}
		} else {
			fmt.Println("service Name for", id, " is missing")
		}
	}

	dockerfile += "CMD /globular/Globular\n"
	// save docker.
	err = ioutil.WriteFile(path+"/"+"Dockerfile", []byte(dockerfile), 0644)
	if err != nil {
		log.Println(err)
	}
}
