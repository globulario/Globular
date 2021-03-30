package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/davecourtois/Utility"
	"github.com/globulario/services/golang/admin/admin_client"
	"github.com/globulario/services/golang/resource/resource_client"
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

	// Start should not block. Do the actual work async.
	go g.run()
	return nil
}

func (g *Globule) run() error {

	// start globular and wait on exit chan...
	go func() {
		g.Serve()
		log.Println("globular serve at domain ", g.Domain)
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
	log.Println("try to stop external services")

	g.stopServices()

	log.Println("locals services are stopped")

	// Stop internal services
	log.Println("try to stop internal services")
	g.stopInternalServices()
	log.Println("internal services  are stopped")

	logger.Infof("Globular has been stopped!")
	pids, err := Utility.GetProcessIdsByName("grpcwebproxy")
	if err == nil {
		for i := 0; i < len(pids); i++ {
			Utility.TerminateProcess(pids[i], 0)
		}
	}
	g.stopMongod()

	// Close kv stores.
	g.logs.Close()
	g.permissions.Close()

	close(g.exit)
	return err
}

func main() {

	g := NewGlobule()
	g.exit = make(chan struct{})
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
	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {

		// Start with sepecific parameter.
		startCommand := flag.NewFlagSet("start", flag.ExitOnError)
		startCommand_domain := startCommand.String("domain", "", "The domain of the service.")

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
		deployCommand_organization := deployCommand.String("o", "", "The name of the organisation that responsible of the application. (Required)")
		deployCommand_path := deployCommand.String("path", "", "You must specify the path that contain the source (bundle.js, index.html...) of the application to deploy. (Required)")
		deployCommand_user := deployCommand.String("u", "", "The user name. (Required)")
		deployCommand_pwd := deployCommand.String("p", "", "The user password. (Required)")
		deployCommand_address := deployCommand.String("a", "", "The domain of the server where to install the appliction (Required)")

		// Publish command.
		publishCommand := flag.NewFlagSet("publish", flag.ExitOnError)
		publishCommand_path := publishCommand.String("path", "", "You must specify the path that contain the config.json, .proto and all dependcies require by the service to run. (Required)")
		publishCommand_user := publishCommand.String("u", "", "The user name. (Required)")
		publishCommand_pwd := publishCommand.String("p", "", "The user password. (Required)")
		publishCommand_address := publishCommand.String("a", "", "The domain of the server where to install the appliction (Required)")

		// *** Those informations are optional they are in the configuration of the service.
		publishCommand_organization := publishCommand.String("o", "", "The Organization that publish the service. (Optional)")
		publishCommand_plaform := publishCommand.String("platform", "", "(Optional)")

		// Install certificates command.
		installCertificatesCommand := flag.NewFlagSet("certificates", flag.ExitOnError)
		installCertificatesCommand_path := installCertificatesCommand.String("path", "", "You must specify where to install certificate (Required)")
		installCertificatesCommand_port := installCertificatesCommand.String("port", "", "You must specify the port where the configuration can be found (Required)")
		installCertificatesCommand_domain := installCertificatesCommand.String("domain", "", "You must specify the domain (Required)")

		switch os.Args[1] {
		case "start":
			startCommand.Parse(os.Args[2:])
		case "dist":
			distCommand.Parse(os.Args[2:])
		case "deploy":
			deployCommand.Parse(os.Args[2:])
		case "publish":
			publishCommand.Parse(os.Args[2:])
		case "install":
			installCommand.Parse(os.Args[2:])
		case "uninstall":
			unstallCommand.Parse(os.Args[2:])
		case "certificates":
			installCertificatesCommand.Parse(os.Args[2:])
		default:
			flag.PrintDefaults()
			os.Exit(1)
		}

		if installCertificatesCommand.Parsed() {
			// Required Flags
			if *installCertificatesCommand_path == "" {
				installCertificatesCommand.PrintDefaults()
				os.Exit(1)
			}

			if *installCertificatesCommand_domain == "" {
				installCertificatesCommand.PrintDefaults()
				os.Exit(1)
			}

			if *installCertificatesCommand_port == "" {
				installCertificatesCommand.PrintDefaults()
				os.Exit(1)
			}

			installCertificates(g, *installCertificatesCommand_domain, Utility.ToInt(*installCertificatesCommand_port), *installCertificatesCommand_path)
		}

		if startCommand.Parsed() {
			// Required Flags

			if len(*startCommand_domain) > 0 {
				g.Domain = *startCommand_domain
			}
			g.run()
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

			deploy(g, *deployCommand_name, *deployCommand_organization, *deployCommand_path, *deployCommand_address, *deployCommand_user, *deployCommand_pwd)
		}

		if publishCommand.Parsed() {

			if *publishCommand_path == "" {
				publishCommand.PrintDefaults()
				fmt.Println("No -path was given!")
				os.Exit(1)
			}

			if *publishCommand_user == "" {
				publishCommand.PrintDefaults()
				fmt.Println("No -u (user) was given!")
				os.Exit(1)
			}

			if *publishCommand_pwd == "" {
				publishCommand.PrintDefaults()
				fmt.Println("No -p (password) was given!")
				os.Exit(1)
			}

			if *publishCommand_address == "" {
				publishCommand.PrintDefaults()
				fmt.Println("No -a (address or domain) was given!")
				os.Exit(1)
			}

			// Here I will read the configuration file...

			if !Utility.Exists(*publishCommand_path) {
				fmt.Println("No configuration file was found at " + *publishCommand_path)
				os.Exit(1)
			}

			// Detect the platform if none was given...
			if *publishCommand_plaform == "" {
				*publishCommand_plaform = runtime.GOOS + "_" + runtime.GOARCH
			}

			// Pulish the services.
			publish(g, *publishCommand_user, *publishCommand_pwd, *publishCommand_address, *publishCommand_organization, *publishCommand_path, *publishCommand_plaform)
		}

	} else {

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
func installCertificates(g *Globule, domain string, port int, path string) error {
	log.Println("Get certificates from ", domain, "...")
	admin_client_, err := admin_client.NewAdminService_Client(domain, "admin.AdminService")
	if err != nil {
		log.Println("fail to get certificates...", err)
		return err
	}

	key, cert, ca, err := admin_client_.InstallCertificates(domain, port, path)
	if err != nil {
		log.Println("fail to get certificates...", err)
	}
	log.Println("Your certificate are installed: ")
	log.Println("cacert: ", ca)
	log.Println("cert: ", cert)
	log.Println("certkey: ", key)
	return nil
}

/**
 * That function can be use to deploy an application on the server...
 */
func deploy(g *Globule, name string, organization string, path string, address string, user string, pwd string) error {

	log.Println("deploy application...", name, " to address ", address, " user ", user)

	// Authenticate the user in order to get the token
	resource_client, err := resource_client.NewResourceService_Client(address, "resource.ResourceService")

	if err != nil {
		log.Println("fail to access resource service at "+address+" with error ", err)
		return err
	}

	token, err := resource_client.Authenticate(user, pwd)
	if err != nil {
		log.Println("fail to authenticate user ", err)
		return err
	}

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	admin_client_, err := admin_client.NewAdminService_Client(address, "admin.AdminService") // create the resource server.
	if err != nil {
		return err
	}

	_, err = admin_client_.DeployApplication(user, name, organization, path, token, address)
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
func publish(g *Globule, user, pwd, domain, organization, path, platform string) error {

	// Authenticate the user in order to get the token
	resource_client_, err := resource_client.NewResourceService_Client(domain, "resource.ResourceService")
	if err != nil {
		log.Panicln(err)
		return err
	}

	token, err := resource_client_.Authenticate(user, pwd)
	if err != nil {
		log.Println(err)
		return err
	}

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	admin_client_, err := admin_client.NewAdminService_Client(domain, "admin.AdminService")
	if err != nil {
		log.Println(err)
		return err
	}

	// first of all I will create and upload the package on the discovery...
	path_, _, err := admin_client_.UploadServicePackage(user, organization, token, domain, path, platform)
	if err != nil {
		log.Println(err)
		return err
	}

	err = admin_client_.PublishService(user, organization, token, domain, path_, path, platform)
	if err != nil {
		return err
	}

	log.Println("Service was pulbish successfully!")
	return nil
}

/**
 * That function is use to install globular from the development environnement.
 * The server must have run at least once before that command is call. Each service must
 * have been run at least one to appear in the installation.
 */
func install(g *Globule, path string) {
	// That function is use to install globular at a given repository.
	fmt.Println("install globular in directory: ", path)

	// I will set the docker file depending of the arch.
	var dockerfile string
	if runtime.GOARCH == "amd64" {
		data, err := ioutil.ReadFile("Dockerfile_amd64")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		dockerfile = string(data)
	} else if runtime.GOARCH == "arm64" {
		data, err := ioutil.ReadFile("Dockerfile_arm64")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		dockerfile = string(data)
	}

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
	// admin, resource, ca and services.
	serviceDir := os.Getenv("GLOBULAR_SERVICES_ROOT")
	Utility.CreateDirIfNotExist(path + "/proto")
	err = Utility.CopyFile(serviceDir+"/proto/admin.proto", path+"/proto/admin.proto")
	if err != nil {
		log.Println("fail to copy with error ", err)
	}
	err = Utility.CopyFile(serviceDir+"/proto/ca.proto", path+"/proto/ca.proto")
	if err != nil {
		log.Println("fail to copy with error ", err)
	}
	err = Utility.CopyFile(serviceDir+"/proto/resource.proto", path+"/proto/resource.proto")
	if err != nil {
		log.Println("fail to copy with error ", err)
	}
	err = Utility.CopyFile(serviceDir+"/proto/packages.proto", path+"/proto/packages.proto")
	if err != nil {
		log.Println("fail to copy with error ", err)
	}
	err = Utility.CopyFile(serviceDir+"/proto/lb.proto", path+"/proto/lb.proto")
	if err != nil {
		log.Println("fail to copy with error ", err)
	}

	for _, f := range files {
		if !f.IsDir() {
			err = os.Chmod(path+"/"+"bin"+"/"+f.Name(), 0755)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	// install services...
	for _, s := range g.getServices() {
		id := getStringVal(s, "Id")
		_, hasName := s.Load("Name")
		if hasName {
			name := getStringVal(s, "Name")

			// I will read the configuration file to have nessecary service information
			// to be able to create the path.
			configPath := getStringVal(s, "Path")
			if len(configPath) > 0 {
				configPath = configPath[:strings.LastIndex(configPath, "/")] + "/config.json"
				if Utility.Exists(configPath) {
					log.Println("install service ", name)
					bytes, err := ioutil.ReadFile(configPath)
					config := make(map[string]interface{}, 0)
					json.Unmarshal(bytes, &config)

					if err == nil {
						_, hasProto := s.Load("Proto")
						_, hasPath := s.Load("Path")
						// set the name.
						if config["PublisherId"] != nil && config["Version"] != nil && hasPath && hasProto {

							execPath := getStringVal(s, "Path")
							protoPath := getStringVal(s, "Proto")

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
						} else if !hasProto {
							fmt.Println(" no proto file was found!")
						} else if !hasPath {
							fmt.Println("no executable was found!")
						}
					} else {
						fmt.Println(err)
					}
				} else {
					fmt.Println("service", name, ":", id, "configuration is incomplete!")
				}
			} else {

				// Internal services here.
				protoPath := getStringVal(s, "Proto")
				// Copy the proto file.
				if Utility.Exists(os.Getenv("GLOBULAR_SERVICES_ROOT") + "/" + protoPath) {
					Utility.Copy(os.Getenv("GLOBULAR_SERVICES_ROOT")+"/"+protoPath, path+"/"+protoPath)
				}
			}
		}
	}

	dockerfile += `WORKDIR /globular
ENTRYPOINT ["/globular/Globular"]`

	// save docker.
	err = ioutil.WriteFile(path+"/"+"Dockerfile", []byte(dockerfile), 0644)
	if err != nil {
		log.Println(err)
	}
}
