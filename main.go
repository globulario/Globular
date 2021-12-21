package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/davecourtois/Utility"
	"github.com/globulario/services/golang/admin/admin_client"
	"github.com/globulario/services/golang/applications_manager/applications_manager_client"
	"github.com/globulario/services/golang/authentication/authentication_client"
	config_ "github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/discovery/discovery_client"
	"github.com/globulario/services/golang/repository/repository_client"
	"github.com/globulario/services/golang/resource/resource_client"
	"github.com/globulario/services/golang/security"
	service_manager_client "github.com/globulario/services/golang/services_manager/services_manager_client"
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

	// stop mongo demon if it running
	Utility.KillProcessByName("mongod")

	// Start should not block. Do the actual work async.
	go g.run()
	return nil
}

func (g *Globule) run() error {

	// start globular and wait on exit chan...
	go func() {
		g.Serve()
	}()

	// wait for exit.
	<-g.exit

	return nil
}

func (g *Globule) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	logger.Info("Globular is stopping!")
	g.exit_ = true

	// Close all proxy
	g.stopProxies()

	// Close all services.
	g.stopServices()

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
	if err != nil {
		log.Fatal(err)
	}

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
		distCommand_path := distCommand.String("path", "", "You must specify the distribution path. (Required)")
		distCommand_revision := distCommand.String("revision", "", "You must specify the package revision. (Required)")

		// Deploy command
		deployCommand := flag.NewFlagSet("deploy", flag.ExitOnError)
		deployCommand_name := deployCommand.String("name", "", "You must specify an application name. (Required)")
		deployCommand_organization := deployCommand.String("o", "", "The name of the organisation that responsible of the application. (Required)")
		deployCommand_path := deployCommand.String("path", "", "You must specify the path that contain the source (bundle.js, index.html...) of the application to deploy. (Required)")
		deployCommand_user := deployCommand.String("u", "", "The user name. (Required)")
		deployCommand_pwd := deployCommand.String("p", "", "The user password. (Required)")
		deployCommand_address := deployCommand.String("a", "", "The domain of the server where to install the appliction (Required)")
		deployCommand_index := deployCommand.String("set_as_default", "", "The value is true the application will be set as default (Optional false by default)")

		// Publish Service.
		// Service can be written in various language, however all service must contain a config.json file and a .proto file.
		// The config.json file contain field to be inform, the service version, the discovery and the repostiory.
		// ex. ./Globular publish -a=globular.io -path=/tmp/echo_v_1  -o=globulario  -u=userid -p=*******
		publishCommand := flag.NewFlagSet("publish", flag.ExitOnError)
		publishCommand_path := publishCommand.String("path", "", "You must specify the path that contain the config.json, .proto and all dependcies require by the service to run. (Required)")
		publishCommand_user := publishCommand.String("u", "", "The user name. (Required)")
		publishCommand_pwd := publishCommand.String("p", "", "The user password. (Required)")
		publishCommand_address := publishCommand.String("a", "", "The domain of the server where to publish the service (Required)")
		publishCommand_organization := publishCommand.String("o", "", "The Organization that publish the service. (Optional)")
		publishCommand_plaform := publishCommand.String("platform", "", "(Optional it take your current platform as default.)")

		// Install certificates on a server from a local service command.
		// TODO test it... (the path is not working properly, be sure permission are correctly assigned for certificates...)
		installCertificatesCommand := flag.NewFlagSet("certificates", flag.ExitOnError)
		installCertificatesCommand_path := installCertificatesCommand.String("path", "", "You must specify where to install certificate (Required)")
		installCertificatesCommand_port := installCertificatesCommand.String("port", "", "You must specify the port where the configuration can be found (Required)")
		installCertificatesCommand_domain := installCertificatesCommand.String("domain", "", "You must specify the domain (Required)")

		// Install a service on the server.
		// That function must be run as sa.
		// ex. Globular install_service -publisher=globulario -discovery=globular.io -service=echo.EchoService -a=globular1.globular.io -u=sa -p=*******
		install_service_command := flag.NewFlagSet("install_service", flag.ExitOnError)
		install_service_command_publisher := install_service_command.String("publisher", "", "The publisher id (Required)")
		install_service_command_discovery := install_service_command.String("discovery", "", "The addresse where the service was publish (Required)")
		install_service_command_service := install_service_command.String("service", "", " the service id (uuid) (Required)")
		install_service_command_address := install_service_command.String("a", "", "The domain of the server where to install the service (Required)")
		install_service_command_user := install_service_command.String("u", "", "The user name. (Required)")
		install_service_command_pwd := install_service_command.String("p", "", "The user password. (Required)")

		// Uninstall a service on the server.
		// That function must be run as sa.
		uninstall_service_command := flag.NewFlagSet("uninstall_service", flag.ExitOnError)
		uninstall_service_command_service := uninstall_service_command.String("service", "", " the service uuid (Required)")
		uninstall_service_command_publisher := uninstall_service_command.String("publisher", "", "The publisher id (Required)")
		uninstall_service_command_version := uninstall_service_command.String("version", "", " The service vesion(Required)")
		uninstall_service_command_address := uninstall_service_command.String("a", "", "The domain of the server where to install the service (Required)")
		uninstall_service_command_user := uninstall_service_command.String("u", "", "The user name. (Required)")
		uninstall_service_command_pwd := uninstall_service_command.String("p", "", "The user password. (Required)")

		// Install a application on the server.
		install_application_command := flag.NewFlagSet("install_application", flag.ExitOnError)
		install_application_command_publisher := install_application_command.String("publisher", "", "The publisher id (Required)")
		install_application_command_discovery := install_application_command.String("discovery", "", "The addresse where the application was publish (Required)")
		install_application_command_name := install_application_command.String("application", "", " the application name (Required)")
		install_application_command_address := install_application_command.String("a", "", "The domain of the server where to install the application (Required)")
		install_application_command_user := install_application_command.String("u", "", "The user name. (Required)")
		install_application_command_pwd := install_application_command.String("p", "", "The user password. (Required)")
		install_application_command_index := install_application_command.String("set_as_default", "", "The value is true the application will be set as default (Optional false by default)")

		// Uninstall a service on the server.
		uninstall_application_command := flag.NewFlagSet("uninstall_application", flag.ExitOnError)
		uninstall_application_command_name := uninstall_application_command.String("application", "", " the application name (Required)")
		uninstall_application_command_publisher := uninstall_application_command.String("publisher", "", "The publisher id (Required)")
		uninstall_application_command_version := uninstall_application_command.String("version", "", " The application vesion(Required)")
		uninstall_application_command_address := uninstall_application_command.String("a", "", "The domain where the application is runing (Required)")
		uninstall_application_command_user := uninstall_application_command.String("u", "", "The user name. (Required)")
		uninstall_application_command_pwd := uninstall_application_command.String("p", "", "The user password. (Required)")

		// push globular update.
		// Update a given globular server with a new executable file. That command must be call as sa.
		// ex. ./Globular update -path=/home/dave/go/src/github.com/globulario/Globular/Globular -a=globular.io -u=sa -p=adminadmin
		update_globular_command := flag.NewFlagSet("update", flag.ExitOnError)
		update_globular_command_exec_path := update_globular_command.String("path", "", " the path to the new executable to update from")
		update_globular_command_address := update_globular_command.String("a", "", "The domain of the server where to push the update(Required)")
		update_globular_command_user := update_globular_command.String("u", "", "The user name. (Required)")
		update_globular_command_pwd := update_globular_command.String("p", "", "The user password. (Required)")
		update_globular_command_platform := update_globular_command.String("platform", "", "The os and arch info ex: linux:arm64 (optional)")

		// pull globular update.
		update_globular_from_command := flag.NewFlagSet("update_from", flag.ExitOnError)
		update_globular_command_from_source := update_globular_from_command.String("source", "", " the address of the server from where to update the a given server.")
		update_globular_from_command_dest := update_globular_from_command.String("a", "", "The domain of the server to update (Required)")
		update_globular_from_command_user := update_globular_from_command.String("u", "", "The user name. (Required)")
		update_globular_from_command_pwd := update_globular_from_command.String("p", "", "The user password. (Required)")
		update_globular_from_command_platform := update_globular_from_command.String("platform", "", "The os and arch info ex: linux:arm64 (optional)")

		// Connect peer one to another. The peer Domain must be set before the calling that function.
		connect_peer_command := flag.NewFlagSet("connect_peer", flag.ExitOnError)
		connect_peer_command_address := connect_peer_command.String("dest", "", "The address of the peer to connect to, can contain it configuration port (80) by defaut.")
		connect_peer_command_user := connect_peer_command.String("u", "", "The user name. (Required)")
		connect_peer_command_pwd := connect_peer_command.String("p", "", "The user password. (Required)")

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
		case "update":
			update_globular_command.Parse(os.Args[2:])
		case "update_from":
			update_globular_from_command.Parse(os.Args[2:])
		case "install_service":
			install_service_command.Parse(os.Args[2:])
		case "uninstall_service":
			uninstall_service_command.Parse(os.Args[2:])
		case "install_application":
			install_application_command.Parse(os.Args[2:])
		case "uninstall_application":
			uninstall_application_command.Parse(os.Args[2:])
		case "certificates":
			installCertificatesCommand.Parse(os.Args[2:])
		case "connect_peer":
			connect_peer_command.Parse(os.Args[2:])
		default:
			flag.PrintDefaults()
			os.Exit(1)
		}

		if connect_peer_command.Parsed() {
			if *connect_peer_command_address == "" {
				connect_peer_command.PrintDefaults()
				fmt.Println("No peer address given!")
				os.Exit(1)
			}
			if *connect_peer_command_user == "" {
				connect_peer_command.PrintDefaults()
				fmt.Println("no user was given!")
				os.Exit(1)
			}
			if *connect_peer_command_pwd == "" {
				connect_peer_command.PrintDefaults()
				fmt.Println("no password was given!")
				os.Exit(1)
			}
			err = connect_peer(g, *connect_peer_command_address, *connect_peer_command_user, *connect_peer_command_pwd)
			if err != nil {
				log.Println(err)
			}
		}

		if install_service_command.Parsed() {
			if *install_service_command_service == "" {
				install_service_command.PrintDefaults()
				fmt.Println("no service name was given!")
				os.Exit(1)
			}
			if *install_service_command_discovery == "" {
				install_service_command.PrintDefaults()
				fmt.Println("no discovery adress was given!")
				os.Exit(1)
			}
			if *install_service_command_publisher == "" {
				install_service_command.PrintDefaults()
				fmt.Println("no publiser was given!")
				os.Exit(1)
			}
			if *install_service_command_address == "" {
				install_service_command.PrintDefaults()
				fmt.Println("no domain was given!")
				os.Exit(1)
			}
			if *install_service_command_user == "" {
				install_service_command.PrintDefaults()
				fmt.Println("no user was given!")
				os.Exit(1)
			}
			if *install_service_command_pwd == "" {
				install_service_command.PrintDefaults()
				fmt.Println("no password was given!")
				os.Exit(1)
			}
			install_service(g, *install_service_command_service, *install_service_command_discovery, *install_service_command_publisher, *install_service_command_address, *install_service_command_user, *install_service_command_pwd)
		}

		if uninstall_service_command.Parsed() {
			if *uninstall_service_command_service == "" {
				install_service_command.PrintDefaults()
				os.Exit(1)
			}
			if *uninstall_service_command_publisher == "" {
				install_service_command.PrintDefaults()
				os.Exit(1)
			}

			if *uninstall_service_command_address == "" {
				install_service_command.PrintDefaults()
				os.Exit(1)
			}
			if *uninstall_service_command_user == "" {
				install_service_command.PrintDefaults()
				os.Exit(1)
			}
			if *uninstall_service_command_pwd == "" {
				install_service_command.PrintDefaults()
				os.Exit(1)
			}
			if *uninstall_service_command_version == "" {
				install_service_command.PrintDefaults()
				os.Exit(1)
			}
			uninstall_service(g, *uninstall_service_command_service, *uninstall_service_command_publisher, *uninstall_service_command_version, *uninstall_service_command_address, *uninstall_service_command_user, *uninstall_service_command_pwd)
		}

		if update_globular_command.Parsed() {
			if *update_globular_command_exec_path == "" {
				update_globular_command.PrintDefaults()
				fmt.Println("no executable path was given!")
				os.Exit(1)
			}

			if *update_globular_command_address == "" {
				update_globular_command.PrintDefaults()
				fmt.Println("no domain was given!")
				os.Exit(1)
			}
			if *update_globular_command_user == "" {
				update_globular_command.PrintDefaults()
				fmt.Println("no user was given!")
				os.Exit(1)
			}
			if *update_globular_command_pwd == "" {
				update_globular_command.PrintDefaults()
				fmt.Println("no password was given!")
				os.Exit(1)
			}
			if *update_globular_command_platform == "" {
				*update_globular_command_platform = runtime.GOOS + ":" + runtime.GOARCH
			}

			update_globular(g, *update_globular_command_exec_path, *update_globular_command_address, *update_globular_command_user, *update_globular_command_pwd, *update_globular_command_platform)
		}

		if update_globular_from_command.Parsed() {
			if *update_globular_command_from_source == "" {
				update_globular_from_command.PrintDefaults()
				fmt.Println("no address was given to update Globular from")
				os.Exit(1)
			}

			if *update_globular_from_command_dest == "" {
				update_globular_command.PrintDefaults()
				fmt.Println("no domain was given!")
				os.Exit(1)
			}
			if *update_globular_from_command_user == "" {
				update_globular_from_command.PrintDefaults()
				fmt.Println("no user (for domain) was given!")
				os.Exit(1)
			}

			if *update_globular_from_command_pwd == "" {
				update_globular_from_command.PrintDefaults()
				fmt.Println("no password (for domain) was given!")
				os.Exit(1)
			}

			if *update_globular_from_command_platform == "" {
				*update_globular_from_command_platform = runtime.GOOS + ":" + runtime.GOARCH
			}

			update_globular_from(g, *update_globular_command_from_source, *update_globular_from_command_dest, *update_globular_from_command_user, *update_globular_from_command_pwd, *update_globular_from_command_platform)
		}

		if install_application_command.Parsed() {
			if *install_application_command_name == "" {
				install_application_command.PrintDefaults()
				fmt.Println("no application name was given!")
				os.Exit(1)
			}
			if *install_application_command_discovery == "" {
				install_application_command.PrintDefaults()
				fmt.Println("no discovery adress was given!")
				os.Exit(1)
			}
			if *install_application_command_publisher == "" {
				install_application_command.PrintDefaults()
				fmt.Println("no publiser was given!")
				os.Exit(1)
			}
			if *install_application_command_address == "" {
				install_application_command.PrintDefaults()
				fmt.Println("no domain was given!")
				os.Exit(1)
			}
			if *install_application_command_user == "" {
				install_application_command.PrintDefaults()
				fmt.Println("no user was given!")
				os.Exit(1)
			}
			if *install_application_command_pwd == "" {
				install_application_command.PrintDefaults()
				fmt.Println("no password was given!")
				os.Exit(1)
			}

			var set_as_default bool
			if *install_application_command_index != "" {
				set_as_default = *install_application_command_index == "true"
			}

			install_application(g, *install_application_command_name, *install_application_command_discovery, *install_application_command_publisher, *install_application_command_address, *install_application_command_user, *install_application_command_pwd, set_as_default)
		}

		if uninstall_application_command.Parsed() {
			if *uninstall_application_command_name == "" {
				uninstall_application_command.PrintDefaults()
				fmt.Println("no application name was given!")
				os.Exit(1)
			}
			if *uninstall_application_command_publisher == "" {
				uninstall_application_command.PrintDefaults()
				fmt.Println("no publisher was given!")
				os.Exit(1)
			}

			if *uninstall_application_command_address == "" {
				uninstall_application_command.PrintDefaults()
				fmt.Println("no domain was given!")
				os.Exit(1)
			}
			if *uninstall_application_command_user == "" {
				uninstall_application_command.PrintDefaults()
				fmt.Println("no user was given!")
				os.Exit(1)
			}
			if *uninstall_application_command_pwd == "" {
				uninstall_application_command.PrintDefaults()
				fmt.Println("no password was given!")
				os.Exit(1)
			}
			if *uninstall_application_command_version == "" {
				uninstall_application_command.PrintDefaults()
				fmt.Println("no version was given!")
				os.Exit(1)
			}
			uninstall_application(g, *uninstall_application_command_name, *uninstall_application_command_publisher, *uninstall_application_command_version, *uninstall_application_command_address, *uninstall_application_command_user, *uninstall_application_command_pwd)
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
				fmt.Println("No path was given!")
				distCommand.PrintDefaults()
				os.Exit(1)
			}
			if *distCommand_revision == "" {
				fmt.Println("No revision number was given!")
				distCommand.PrintDefaults()
				os.Exit(1)
			}
			dist(g, *distCommand_path, *distCommand_revision)
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

			var set_as_default bool
			if *deployCommand_index != "" {
				set_as_default = *deployCommand_index == "true"
			}

			deploy(g, *deployCommand_name, *deployCommand_organization, *deployCommand_path, *deployCommand_address, *deployCommand_user, *deployCommand_pwd, set_as_default)
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

	key, cert, ca, err := admin_client_.GetCertificates(domain, port, path)
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
func deploy(g *Globule, name string, organization string, path string, address string, user string, pwd string, set_as_default bool) error {

	log.Println("try to deploy application", name, " to address ", address, " with user ", user)

	// Authenticate the user in order to get the token
	authentication_client, err := authentication_client.NewAuthenticationService_Client(address, "authentication.AuthenticationService")
	if err != nil {
		log.Println("fail to access resource service at "+address+" with error ", err)
		return err
	}

	log.Println("authenticate user ", user, " at adress ", address)
	token, err := authentication_client.Authenticate(user, pwd)
	if err != nil {
		log.Println("fail to authenticate user ", err)
		return err
	}

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	log.Println("Connect with application manager")
	applications_manager_client_, err := applications_manager_client.NewApplicationsManager_Client(address, "applications_manager.ApplicationManagerService") // create the resource server.
	if err != nil {
		log.Println("fail to connect to application manager service ", address)
		return err
	}

	log.Println("Connect with application manager")
	_, err = applications_manager_client_.DeployApplication(user, name, organization, path, token, address, set_as_default)
	if err != nil {
		log.Println("Fail to deploy applicaiton with error:", err)
		return err
	}

	log.Println("Application", name, "was deployed successfully!")
	return nil
}

/**
 * Push globular update on a given server.
 * ex.
 * sudo ./Globular update -path=/home/dave/go/src/github.com/globulario/Globular/Globular -a=globular.cloud -u=sa -p=adminadmin
 */
func update_globular(g *Globule, path, domain, user, pwd string, platform string) error {

	// Authenticate the user in order to get the token
	authentication_client, err := authentication_client.NewAuthenticationService_Client(domain, "authentication.AuthenticationService")
	if err != nil {
		log.Panicln(err)
		return err
	}

	token, err := authentication_client.Authenticate(user, pwd)
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

	_, err = admin_client_.Update(path, platform, token, domain)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Pull globular update from a given server to a given destination.
 * ex.
 * ./Globular update_from -a=globular.cloud -p=adminadmin -source=globular.cloud -u=sa
 */
func update_globular_from(g *Globule, src, dest, user, pwd string, platform string) error {
	log.Println("pull globular update from ", src, " to ", dest)

	admin_source, err := admin_client.NewAdminService_Client(src, "admin.AdminService")
	if err != nil {
		return err
	}

	// From the source I will download the new executable and save it in the
	// temp directory...
	path := os.TempDir() + "/Globular_" + Utility.ToString(time.Now().Unix())
	Utility.CreateDirIfNotExist(path)
	err = admin_source.DownloadGlobular(src, platform, path)
	if err != nil {
		log.Println("fail to download new globular executable file with error: ", err)
		return err
	}

	path_ := path
	path_ += "/Globular"
	if runtime.GOOS == "windows" {
		path_ += ".exe"
	}

	if !Utility.Exists(path_) {
		err := errors.New(path_ + " not found! ")
		return err
	}

	defer os.RemoveAll(path)

	err = update_globular(g, path_, dest, user, pwd, platform)
	if err != nil {
		log.Println(err)
		return err
	}

	// Send the command.
	return nil

}

/**
 * That function is use to publish a service on the network.
 * The service must contain
 * - A .proto file that define it gRPC services interface.
 * - A config.json file that define it service configuration.
 * - (optional) A preinst shell script file that contain the script to be run before the installation
 * - (optional) A postinst shell script file tha contain the script to be run after the installation
 * - (optional) A prerm shell script file that contain the script to be run before the service removal
 * - (optional) A postrm shell script file tha contain the script to be run after the service removal
 * All file in the service directory will be part of the package, so take
 * care to include all dependencies, dll... to be sure your services will run
 * as expected.
 */
func publish(g *Globule, user, pwd, domain, organization, path, platform string) error {

	// Authenticate the user in order to get the token
	authentication_client, err := authentication_client.NewAuthenticationService_Client(domain, "authentication.AuthenticationService")
	if err != nil {
		log.Panicln(err)
		return err
	}

	token, err := authentication_client.Authenticate(user, pwd)
	if err != nil {
		log.Println(err)
		return err
	}

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	repository_client_, err := repository_client.NewRepositoryService_Client(domain, "discovery.PackageDiscovery")
	if err != nil {
		log.Println(err)
		return err
	}

	// first of all I will create and upload the package on the discovery...
	err = repository_client_.UploadServicePackage(user, organization, token, domain, path, platform)
	if err != nil {
		log.Println(err)
		return err
	}

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	discovery_client_, err := discovery_client.NewDiscoveryService_Client(domain, "discovery.PackageDiscovery")
	if err != nil {
		log.Println(err)
		return err
	}

	err = discovery_client_.PublishService(user, organization, token, domain, path, platform)
	if err != nil {
		return err
	}

	log.Println("Service was pulbish successfully have a nice day folk's!")
	return nil
}

/**
 * That function is use to intall a service at given adresse. The service Id is the unique identifier of
 * the service to be install.
 */
func install_service(g *Globule, serviceId, discovery, publisherId, domain, user, pwd string) error {
	log.Println("try to install service", serviceId, "from", publisherId, "on", domain)
	// Authenticate the user in order to get the token
	authentication_client, err := authentication_client.NewAuthenticationService_Client(domain, "authentication.AuthenticationService")
	if err != nil {
		log.Panicln(err)
		return err
	}
	token, err := authentication_client.Authenticate(user, pwd)
	if err != nil {
		log.Println("fail to authenticate with error ", err.Error())
		return err
	}

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	services_manager_client_, err := service_manager_client.NewServicesManagerService_Client(domain, "services_manager.ServicesManagerService")
	if err != nil {
		log.Println("fail to connect to services manager at ", domain, " with error ", err.Error())
		return err
	}

	// first of all I will create and upload the package on the discovery...
	err = services_manager_client_.InstallService(token, domain, user, discovery, publisherId, serviceId)
	if err != nil {
		log.Println("fail to install service", serviceId, "with error ", err.Error())
		return err
	}

	log.Println("service was installed")
	return nil
}

func uninstall_service(g *Globule, serviceId, publisherId, version, domain, user, pwd string) error {

	// Authenticate the user in order to get the token
	authentication_client, err := authentication_client.NewAuthenticationService_Client(domain, "authentication.AuthenticationService")
	if err != nil {
		log.Panicln(err)
		return err
	}

	token, err := authentication_client.Authenticate(user, pwd)
	if err != nil {
		log.Println(err)
		return err
	}

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	services_manager_client_, err := service_manager_client.NewServicesManagerService_Client(domain, "services_manager.ServicesManagerService")
	if err != nil {
		log.Println(err)
		return err
	}

	// first of all I will create and upload the package on the discovery...
	err = services_manager_client_.UninstallService(token, domain, user, publisherId, serviceId, version)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func install_application(g *Globule, applicationId, discovery, publisherId, domain, user, pwd string, set_as_default bool) error {

	// Authenticate the user in order to get the token
	authentication_client, err := authentication_client.NewAuthenticationService_Client(domain, "authentication.AuthenticationService")
	if err != nil {
		log.Panicln(err)
		return err
	}

	token, err := authentication_client.Authenticate(user, pwd)
	if err != nil {
		return err
	}

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	applications_manager_client_, err := applications_manager_client.NewApplicationsManager_Client(domain, "applications_manager.ApplicationManagerService")
	if err != nil {
		log.Println(err)
		return err
	}

	// first of all I will create and upload the package on the discovery...
	err = applications_manager_client_.InstallApplication(token, domain, user, discovery, publisherId, applicationId, set_as_default)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func uninstall_application(g *Globule, applicationId, publisherId, version, domain, user, pwd string) error {

	// Authenticate the user in order to get the token
	authentication_client, err := authentication_client.NewAuthenticationService_Client(domain, "authentication.AuthenticationService")
	if err != nil {
		log.Panicln(err)
		return err
	}

	token, err := authentication_client.Authenticate(user, pwd)
	if err != nil {
		log.Println(err)
		return err
	}

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	applications_manager_client_, err := applications_manager_client.NewApplicationsManager_Client(domain, "applications_manager.ApplicationManagerService")
	if err != nil {
		log.Println(err)
		return err
	}

	// first of all I will create and upload the package on the discovery...
	err = applications_manager_client_.UninstallApplication(token, domain, user, publisherId, applicationId, version)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

/**
 * That function is use to install globular from the development environnement.
 * The server must have run at least once before that command is call. Each service must
 * have been run at least one to appear in the installation.
 *
 * The basic sequence are describe at
 * https://www.internalpointers.com/post/build-binary-deb-package-practical-guide
 *
 * Dependencie script are describe at
 *
 * Utilisation:
 * Globular dist -path=/tmp -revision=1
 */
func dist(g *Globule, path string, revision string) {
	// That function is use to install globular at a given repository.
	fmt.Println("create distribution in ", path)
	fmt.Println(Utility.IsLocal("globular.cloud"))
	// The debian package...
	if runtime.GOOS == "linux" {
		debian_package_path := path + "/globular_" + g.Version + "-" + revision + "_" + runtime.GOARCH

		// remove existiong files...
		os.RemoveAll(debian_package_path)

		// 1. Create the working directory
		Utility.CreateDirIfNotExist(debian_package_path)

		// 2. Create the internal structure

		// globular exec and other services exec
		distro_path := debian_package_path + "/usr/local/share/globular"

		// globular data
		data_path := debian_package_path + "/var/globular/data"

		// globular webroot
		webroot_path := debian_package_path + "/var/globular/webroot"

		// globular configurations
		config_path := debian_package_path + "/etc/globular/config"

		// set the web installer in the webroot
		if Utility.Exists(g.webRoot + "/globular_installer") {
			Utility.CreateDirIfNotExist(webroot_path)
			Utility.CopyDir(g.webRoot+"/globular_installer", webroot_path)
		} else {
			log.Println("** There's no gobular installer project found!")
		}

		// Create the bin directories.
		Utility.CreateDirIfNotExist(distro_path)
		Utility.CreateDirIfNotExist(data_path)
		Utility.CreateDirIfNotExist(config_path)

		// Create the distribution.
		__dist(g, distro_path)

		// 3. Create the control file
		Utility.CreateDirIfNotExist(debian_package_path + "/DEBIAN")

		packageConfig := ""
		// Now I will create the debian.

		// - Package – the name of your program;
		packageConfig += "Package:globular\n"

		// - Version – the version of your program;
		packageConfig += "Version:" + g.Version + "\n"

		// - Architecture – the target architecture
		packageConfig += "Architecture:" + runtime.GOARCH + "\n"

		// - Maintainer – the name and the email address of the person in charge of the package maintenance;
		packageConfig += "Maintainer: Project developed and maitained by Globular.io for more infos <info@globular.io>\n"

		// - Description - a brief description of the program.
		packageConfig += "Description: Globular is a complete web application developement suite. Globular is based on microservices architecture and implemented with help of gRPC.\n"

		err := ioutil.WriteFile(debian_package_path+"/DEBIAN/control", []byte(packageConfig), 0644)
		if err != nil {
			fmt.Println(err)
		}

		// Here I will set the script to run before the installation...
		// https://www.devdungeon.com/content/debian-package-tutorial-dpkgdeb#toc-17
		// TODO create tow version one for arm7 and one for amd64
		preinst := `
		echo "Welcome to Globular!-)"

		echo "insall dependencies..."

		apt-get update && apt-get install -y gnupg2 \
		wget \		
		build-essential \
		curl \
		nano \
		python3 \
		python-is-python3 \
		openssh-server \
		&& rm -rf /var/lib/apt/lists/*

		# install mongo db..
		curl -O https://fastdl.mongodb.org/linux/mongodb-linux-x86_64-ubuntu2004-5.0.5.tgz
		tar -zxvf mongodb-linux-x86_64-ubuntu2004-5.0.5.tgz
		cp -R -n mongodb-linux-x86_64-ubuntu2004-5.0.5/bin/* /usr/local/bin
		rm mongodb-linux-x86_64-ubuntu2004-5.0.5.tgz
		rm -R mongodb-linux-x86_64-ubuntu2004-5.0.5

		apt-get install python3
		update-alternatives --install -y  /usr/bin/python python /usr/bin/python3 1000 

		apt-get install -y ffmpeg

		# -- Install prometheus
		wget https://github.com/prometheus/prometheus/releases/download/v2.32.0/prometheus-2.32.0.linux-amd64.tar.gz
		tar -xf prometheus-2.32.0.linux-amd64.tar.gz
		cp prometheus-2.32.0.linux-amd64/prometheus /usr/local/bin/
		cp prometheus-2.32.0.linux-amd64/promtool /usr/local/bin/
		cp -r prometheus-2.32.0.linux-amd64/consoles /etc/prometheus/
		cp -r prometheus-2.32.0.linux-amd64/console_libraries /etc/prometheus/
		rm -rf prometheus-2.32.0.linux-amd64*

		# -- Install alert manager
		wget https://github.com/prometheus/alertmanager/releases/download/v0.23.0/alertmanager-0.23.0.linux-amd64.tar.gz
		tar -xf alertmanager-0.23.0.linux-amd64.tar.gz
		cp alertmanager-0.23.0.linux-amd64/alertmanager /usr/local/bin
		rm -rf alertmanager-0.23.0.linux-amd64*

		# -- Install node exporter
		wget https://github.com/prometheus/node_exporter/releases/download/v1.3.1/node_exporter-1.3.1.linux-amd64.tar.gz
		tar -xf node_exporter-1.3.1.linux-amd64.tar.gz
		cp node_exporter-1.3.1.linux-amd64/node_exporter /usr/local/bin
		rm -rf node_exporter-1.3.1.linux-amd64*

		# -- Install unix odbc drivers.
		curl http://www.unixodbc.org/unixODBC-2.3.9.tar.gz --output unixODBC-2.3.9.tar.gz
		tar -xvf unixODBC-2.3.9.tar.gz
		rm unixODBC-2.3.9.tar.gz
		cd unixODBC-2.3.9
		./configure && make all install clean && ldconfig
		cd ..

		# -- Install zlib
		curl  https://zlib.net/zlib-1.2.11.tar.gz --output zlib-1.2.11.tar.gz
		tar -xvf zlib-1.2.11.tar.gz
		rm zlib-1.2.11.tar.gz
		cd zlib-1.2.11
		./configure && make all install clean && ldconfig
		cd ..

		# -- Install xapian the search engine.

		curl  https://oligarchy.co.uk/xapian/1.4.18/xapian-core-1.4.18.tar.xz --output xapian-core-1.4.18.tar.xz
		tar -xvf xapian-core-1.4.18.tar.xz
		rm xapian-core-1.4.18.tar.xz
		cd xapian-core-1.4.18
		./configure && make all install clean && ldconfig
		cd ..

		# -- Install youtube-dl
		curl -L https://yt-dl.org/downloads/latest/youtube-dl --output /usr/local/bin/youtube-dl
		chmod a+rx /usr/local/bin/youtube-dl

		if [ -f "/usr/local/bin/Globular" ]; then
			rm /usr/local/bin/Globular
		fi
		`

		err = ioutil.WriteFile(debian_package_path+"/DEBIAN/preinst", []byte(preinst), 0755)
		if err != nil {
			fmt.Println(err)
		}

		postinst := `
		# Create a link into bin to it original location.
		# the systemd file is /etc/systemd/system/Globular.service
		# the environement variable file is /etc/sysconfig/Globular
		 echo "install globular as service..."
		 ln -s /usr/local/share/globular/Globular /usr/local/bin/Globular
		 chmod ugo+x /usr/local/bin/Globular
		 /usr/local/bin/Globular install
		 # here I will modify the /etc/systemd/system/Globular.service file and set 
		 # Restart=always
		 # RestartSec=3
		 echo "set service configuration /etc/systemd/system/Globular.service"
		 sed -i 's/^\(Restart=\).*/\1always/' /etc/systemd/system/Globular.service
		 sed -i 's/^\(RestartSec=\).*/\120/' /etc/systemd/system/Globular.service
		 systemctl daemon-reload
		 systemctl enable Globular
		 echo "To complete your server setup go to http://localhost"
		 echo "To start globular service 'sudo systemctl start Globular'"
		`
		err = ioutil.WriteFile(debian_package_path+"/DEBIAN/postinst", []byte(postinst), 0755)
		if err != nil {
			fmt.Println(err)
		}

		prerm := `
		# Thing to do before removing

		if [ -f "/etc/systemd/system/Globular.service" ]; then
			# Stop, Disable and Uninstall Globular service.
			echo "Stop runing globular service..."
			systemctl stop Globular
			systemctl disable Globular
			systemctl daemon-reload
			rm /etc/systemd/system/Globular.service
		fi

		if [ -f "/usr/local/bin/Globular" ]; then
			echo "Unistall globular service..."
			/usr/local/bin/Globular uninstall
		fi
		`
		err = ioutil.WriteFile(debian_package_path+"/DEBIAN/prerm", []byte(prerm), 0755)
		if err != nil {
			fmt.Println(err)
		}

		postrm := `
		# Thing to do after removing
		if [ -f "/usr/local/bin/Globular" ]; then
			find /usr/local/bin/Globular -xtype l -delete
			rm /etc/systemd/system/Globular.service
		fi
		echo "Hope to see you again soon!"
		`
		err = ioutil.WriteFile(debian_package_path+"/DEBIAN/postrm", []byte(postrm), 0755)
		if err != nil {
			fmt.Println(err)
		}

		// 5. Build the deb package
		cmd := exec.Command("dpkg-deb", "--build", "--root-owner-group", debian_package_path)
		cmdOutput := &bytes.Buffer{}
		cmd.Stdout = cmdOutput

		err = cmd.Run()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Print(cmdOutput.String())

	} else {
		fmt.Println("Create the distro at path ", path)
		// Create the distribution.
		__dist(g, path)
	}

}

// TODO summit to debian https://www.debian.org/doc/manuals/developers-reference/pkgs.html#newpackage
// The globular distribution directory.
//
// There's some note about globular running as service...
// ** service info https://phoenixnap.com/kb/start-stop-restart-linux-services
// sudo systemctl start Globular
// restart the service.
// sudo systemctl restart Globular
// stop the service
// sudo systemctl stop Globular
// enable at boot.
// sudo systemctl enable Globular
// disable at boot time
// sudo systemctl disable Globular
// cd /etc/systemd/system/
// Permissions
// https://ma.ttias.be/auto-restart-crashed-service-systemd/
// https://www.digitalocean.com/community/questions/proper-permissions-for-web-server-s-directory

func __dist(g *Globule, path string) {

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

	// Copy the exec to...

	// Here I will copy the proxy.
	globularExec := os.Args[0]
	destExec := "Globular"
	if runtime.GOOS == "windows" && !strings.HasSuffix(globularExec, ".exe") {
		globularExec += ".exe" // in case of windows
	}

	if runtime.GOOS == "windows" {
		destExec += ".exe"
	}

	err := Utility.Copy(globularExec, path+"/"+destExec)
	if err != nil {
		fmt.Println(err)
	}

	err = os.Chmod(path+"/"+globularExec, 0755)
	if err != nil {
		fmt.Println(err)
	}

	// Copy the bin file from globular
	Utility.CreateDirIfNotExist(path + "/bin")

	err = Utility.CopyDir(dir+"/bin/.", path+"/bin")
	if err != nil {
		log.Panicln("--> fail to copy bin ", err)
	}

	// Change the files permission to add execute write.
	files, err := ioutil.ReadDir(path + "/bin")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if !f.IsDir() {
			err = os.Chmod(path+"/bin/"+f.Name(), 0755)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	// install services...
	services, err := config_.GetServicesConfigurations()
	if err != nil {
		log.Println("fail to retreive services with error ", err)
	}

	for i := 0; i < len(services); i++ {

		// set the service configuration...
		s := services[i]
		id := s["Id"].(string)
		name := s["Name"].(string)

		// I will read the configuration file to have nessecary service information
		// to be able to create the path.
		hasPath := s["Path"] != nil
		if hasPath {
			execPath := s["Path"].(string)
			if len(execPath) > 0 {
				configPath := execPath[:strings.LastIndex(execPath, "/")] + "/config.json"
				if Utility.Exists(configPath) {
					log.Println("install service ", name)
					bytes, err := ioutil.ReadFile(configPath)
					config := make(map[string]interface{}, 0)
					json.Unmarshal(bytes, &config)

					if err == nil {
						hasProto := s["Proto"] != nil

						// set the name.
						if config["PublisherId"] != nil && config["Version"] != nil && hasProto {
							protoPath := s["Proto"].(string)

							if Utility.Exists(execPath) && Utility.Exists(protoPath) {
								var serviceDir = "services/"
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

									config["Path"] = config_.GetRootDir() + "/" + serviceDir + "/" + id + "/" + execName
									config["Proto"] = config_.GetRootDir() + "/" + serviceDir + "/" + name + ".proto"

									// set the security values to nothing...
									config["CertAuthorityTrust"] = ""
									config["CertFile"] = ""
									config["KeyFile"] = ""
									config["TLS"] = false

									if config["Root"] != nil {
										if name == "file.FileService" {
											config["Root"] = config_.GetDataDir() + "/files"

											// I will also copy the mime type directory
											config["Public"] = make([]string, 0)
											Utility.CopyDir(execPath[0:lastIndex]+"/mimetypes", path+"/"+serviceDir+"/"+id)

										} else if name == "conversation.ConversationService" {
											config["Root"] = config_.GetDataDir()
										}
									}

									str, _ := Utility.ToJson(&config)
									ioutil.WriteFile(path+"/"+serviceDir+"/"+id+"/config.json", []byte(str), 0644)

									// Copy the proto file.
									if Utility.Exists(protoPath) {
										Utility.Copy(protoPath, path+"/"+serviceDir+"/"+name+".proto")
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
				protoPath := s["Proto"].(string)

				// Copy the proto file.
				if Utility.Exists(os.Getenv("GLOBULAR_SERVICES_ROOT") + "/" + protoPath) {
					Utility.Copy(os.Getenv("GLOBULAR_SERVICES_ROOT")+"/"+protoPath, path+"/"+protoPath)
				}
			}
		}
	}

	// save docker.
	err = ioutil.WriteFile(path+"/Dockerfile", []byte(dockerfile), 0644)
	if err != nil {
		log.Println(err)
	}
}

/**
 * Connect one peer's with another. When connected peer's are able to generate token valid for both side.
 * The usr and pwd are the admin password in the destionation (ns1.mycelius.com)
 * ex. ./Globular connect_peer -dest=ns1.mycelius.com -u=sa -p=adminadmin
 */
func connect_peer(g *Globule, address, user, pwd string) error {

	// get the local token.
	local_token, err := security.GetLocalToken(g.getDomain())
	if err != nil {
		return nil
	}

	// Authenticate the user in order to get the token
	authentication_client, err := authentication_client.NewAuthenticationService_Client(address, "authentication.AuthenticationService")
	if err != nil {
		return err
	}

	// Get the remote token
	token, err := authentication_client.Authenticate(user, pwd)
	if err != nil {
		return err
	}

	// Create the remote ressource service
	remote_resource_client_, err := resource_client.NewResourceService_Client(address, "resource.ResourceService")
	if err != nil {
		return err
	}

	// Get the local peer key
	key, err := security.GetPeerKey(Utility.MyMacAddr())
	if err != nil {
		log.Println(err)
		return err
	}

	// Register the peer on the remote resourse client...
	peer, key_, err := remote_resource_client_.RegisterPeer(token, g.Mac, g.getDomain(), Utility.MyIP(), Utility.MyLocalIP(), string(key))
	if err != nil {
		return err
	}

	// I will also register the peer to the local server, the local server must running and it domain register,
	// he can be set in /etc/hosts if it's not a public domain.
	local_resource_client_, err := resource_client.NewResourceService_Client(g.getDomain(), "resource.ResourceService")
	if err != nil {
		return err
	}

	_, _, err = local_resource_client_.RegisterPeer(local_token, peer.Mac, peer.Domain, peer.ExternalIpAddress, peer.LocalIpAddress, key_)

	return err
}
