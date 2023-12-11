package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
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
	"github.com/globulario/services/golang/resource/resourcepb"
	"github.com/globulario/services/golang/security"
	service_manager_client "github.com/globulario/services/golang/services_manager/services_manager_client"
	"github.com/kardianos/service"
	"github.com/polds/imgbase64"
	//"github.com/pkg/profile"
)

func (g *Globule) Start(s service.Service) error {

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
	g.exit_ = true
	close(g.exit)
	return nil
}

func main() {
	//defer profile.Start(profile.ProfilePath(".")).Stop()
	// be sure no lock is set.
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

	// Create the service logger...
	errs := make(chan error, 5)

	if err == nil {
		g.logger, err = s.Logger(errs)
		if err != nil {
			fmt.Println("fail to create service logger with error: ", err)
		}
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
		connect_peer_command_token := connect_peer_command.String("token", "", "The token valid on the destination peer (Required)")

		// Generate a token from a given globule. The token can be generate as sa or any other valid user.
		generate_token_command := flag.NewFlagSet("generate_token", flag.ExitOnError)
		generate_token_command_address := generate_token_command.String("dest", "", "The address of the peer to connect to, can contain it configuration port (80) by defaut.")
		generate_token_command_user := generate_token_command.String("u", "", "The user name. (Required)")
		generate_token_command_pwd := generate_token_command.String("p", "", "The user password. (Required)")

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
		case "generate_token":
			generate_token_command.Parse(os.Args[2:])
		default:
			flag.PrintDefaults()
			os.Exit(1)
		}

		if generate_token_command.Parsed() {
			address := *generate_token_command_address
			if *generate_token_command_address == "" {
				address = g.getAddress()
				if strings.Contains(address, ":") {
					address = strings.Split(address, ":")[0]
				}
			}

			if *generate_token_command_user == "" {
				generate_token_command.PrintDefaults()
				fmt.Println("no user was given!")
				os.Exit(1)
			}

			if *generate_token_command_pwd == "" {
				generate_token_command.PrintDefaults()
				fmt.Println("no password was given!")
				os.Exit(1)
			}

			err = generate_token(g, address, *generate_token_command_user, *generate_token_command_pwd)
			if err != nil {
				log.Println(err)
			}
		}

		if connect_peer_command.Parsed() {
			if *connect_peer_command_address == "" {
				connect_peer_command.PrintDefaults()
				fmt.Println("No peer address given!")
				os.Exit(1)
			}
			if *connect_peer_command_token == "" {
				connect_peer_command.PrintDefaults()
				fmt.Println("no user was given!")
				os.Exit(1)
			}

			err = connect_peer(g, *connect_peer_command_address, *connect_peer_command_token)
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
				// Here I will keep the start time...
				// set path...
				setSystemPath()

				os.Exit(0) // exit the program.
			} else {
				log.Println(err)
				os.Exit(0) // exit the program.
			}
		}

		if unstallCommand.Parsed() && s != nil {
			fmt.Println("try to remove Globular service...")
			// Required Flags
			err := s.Uninstall()
			if err == nil {
				log.Println("Globular service is now removed!")
			} else {
				log.Println("fail to remove service with error", err)
				os.Exit(0) // exit the program.
			}

			// Be sure all process are stop...
			Utility.KillProcessByName("mongod")
			Utility.KillProcessByName("prometheus")
			Utility.KillProcessByName("torrent")
			Utility.KillProcessByName("envoy")
			Utility.KillProcessByName("etcd")

			// reset environmement...
			resetSystemPath()
			resetRules()

			os.Exit(0) // exit the program.
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

			if *deployCommand_organization != "" {
				if !strings.Contains(*deployCommand_organization, "@") {
					fmt.Print("You must sepcie the organisation domain ex. organization@domain.com where the domain is the domain where the globule where the organization is define.")
					deployCommand.PrintDefaults()
					os.Exit(1)
				}
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
		defer g.cleanup()

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
				fmt.Println("Valid actions: ", service.ControlAction)
			}
			return
		}

		if s != nil {
			err = s.Run()
			if err != nil && g.logger != nil {
				g.logger.Error(err)
			} else if err != nil {
				fmt.Println("fail to run Globular with error: ", err)
			}
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

	token, err := authentication_client.Authenticate(user, pwd)
	if err != nil {
		log.Println("fail to authenticate user ", err)
		return err
	}

	fmt.Println("authentication succeed.")

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	if !strings.HasPrefix(path, "/") {
		path = strings.ReplaceAll(dir, "\\", "/") + "/" + path

	}

	// From the path I will get try to find the package.json file and get information from it...
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	absolutePath = strings.ReplaceAll(absolutePath, "\\", "/")

	if Utility.Exists(absolutePath + "/package.json") {
		absolutePath += "/package.json"
	} else if Utility.Exists(absolutePath[0:strings.LastIndex(absolutePath, "/")] + "/package.json") {
		absolutePath = absolutePath[0:strings.LastIndex(absolutePath, "/")] + "/package.json"
	} else {
		err = errors.New("no package.config file was found")
		return err
	}

	fmt.Println("try to read package.json file at ", absolutePath)

	packageConfig := make(map[string]interface{})
	data, err := os.ReadFile(absolutePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &packageConfig)
	if err != nil {
		return err
	}

	description := packageConfig["description"].(string)
	version := packageConfig["version"].(string)

	alias := name
	if packageConfig["alias"] != nil {
		alias = packageConfig["alias"].(string)
	}

	// Set keywords.
	keywords := make([]string, 0)
	if packageConfig["keywords"] != nil {
		for i := 0; i < len(packageConfig["keywords"].([]interface{})); i++ {
			keywords = append(keywords, packageConfig["keywords"].([]interface{})[i].(string))
		}
	}

	// Now The application is deploy I will set application actions from the
	// package.json file.
	actions := make([]string, 0)
	if packageConfig["actions"] != nil {
		for i := 0; i < len(packageConfig["actions"].([]interface{})); i++ {
			actions = append(actions, packageConfig["actions"].([]interface{})[i].(string))
		}
	}

	// Create roles.
	roles := make([]*resourcepb.Role, 0)
	if packageConfig["roles"] != nil {
		// Here I will create the roles require by the applications.
		roles_ := packageConfig["roles"].([]interface{})
		for i := 0; i < len(roles_); i++ {
			role_ := roles_[i].(map[string]interface{})
			role := new(resourcepb.Role)
			role.Id = role_["id"].(string)
			role.Name = role_["name"].(string)
			role.Domain = role_["domain"].(string)
			role.Description = role_["description"].(string)
			role.Actions = make([]string, 0)
			for j := 0; j < len(role_["actions"].([]interface{})); j++ {
				role.Actions = append(role.Actions, role_["actions"].([]interface{})[j].(string))
			}
			roles = append(roles, role)
		}
	}

	// Create groups.
	groups := make([]*resourcepb.Group, 0)
	if packageConfig["groups"] != nil {
		groups_ := packageConfig["groups"].([]interface{})
		for i := 0; i < len(groups_); i++ {
			group_ := groups_[i].(map[string]interface{})
			group := new(resourcepb.Group)
			group.Id = group_["id"].(string)
			group.Name = group_["name"].(string)
			group.Description = group_["description"].(string)
			group.Domain = group_["domain"].(string)
			groups = append(groups, group)
		}
	}

	var icon string

	// Now the icon...
	if packageConfig["icon"] != nil {
		// The image icon.
		// iconPath := absolutePath[0:strings.LastIndex(absolutePath, "/")] + "/package.json"
		iconPath := strings.ReplaceAll(absolutePath, "\\", "/")
		lastIndex := strings.LastIndex(iconPath, "/")
		iconPath = iconPath[0:lastIndex] + "/" + packageConfig["icon"].(string)

		if Utility.Exists(iconPath) {
			// Convert to png before creating the data url.
			if strings.HasSuffix(strings.ToLower(iconPath), ".svg") {
				pngPath := os.TempDir() + "/output.png"
				defer os.Remove(pngPath)
				err := Utility.SvgToPng(iconPath, pngPath, 128, 128)
				if err == nil {
					iconPath = pngPath
				}
			}
			// So here I will create the b64 string
			icon, _ = imgbase64.FromLocal(iconPath)
		}
	}

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	fmt.Println("try to publish the application...")
	discovery_client_, err := discovery_client.NewDiscoveryService_Client(address, "discovery.PackageDiscovery")
	if err != nil {
		fmt.Println("fail to connecto to discovery service at address", address, "with error", err)
		return err
	}

	err = discovery_client_.PublishApplication(token, user, organization, "/"+name, name, address, version, description, icon, alias, address, address, actions, keywords, roles, groups)
	if err != nil {
		fmt.Println("fail to publish the application with error:", err)
		return err
	}

	repository_client_, err := repository_client.NewRepositoryService_Client(address, "repository.PackageRepository")
	if err != nil {
		fmt.Println("fail to connecto to repository service at address", address, "with error", err)
		return err
	}

	log.Println("Upload application package at path ", path)
	_, err = repository_client_.UploadApplicationPackage(user, organization, path, token, address, name, version)
	if err != nil {
		fmt.Println("fail to upload the application package with error:", err)
		return err
	}

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	log.Println("Connect with application manager at address ", address)
	applications_manager_client_, err := applications_manager_client.NewApplicationsManager_Client(address, "applications_manager.ApplicationManagerService") // create the resource server.
	if err != nil {
		fmt.Println("fail to connecto to application manager service at address", address, "with error", err)
		return err
	}

	// set the publisher id.
	publisherId := user
	if len(organization) > 0 {
		publisherId = organization
	}

	// if no domain was given i will user the local domain.
	if !strings.Contains(publisherId, "@") {
		domain, err := config_.GetDomain()
		if err != nil {
			fmt.Println("fail to get domain with error", err)
			return err
		}
		publisherId += "@" + domain
	}

	fmt.Println("try to install the newly deployed application...")
	err = applications_manager_client_.InstallApplication(token, authentication_client.GetDomain(), user, address, publisherId, name, false)
	if err != nil {
		log.Println("fail to install application with error ", err)
		return err
	}

	log.Println("Application was deployed and installed sucessfully!")
	return nil

}

/**
 * Push globular update on a given server.
 * ex.
 * sudo ./Globular update -path=/home/dave/go/src/github.com/globulario/Globular/Globular -a=globular.cloud -u=sa -p=adminadmin
 */
func update_globular(g *Globule, path, address, user, pwd string, platform string) error {

	// Authenticate the user in order to get the token
	authentication_client, err := authentication_client.NewAuthenticationService_Client(address, "authentication.AuthenticationService")
	if err != nil {
		fmt.Println(err)
		return err
	}

	token, err := authentication_client.Authenticate(user, pwd)
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

	_, err = admin_client_.Update(path, platform, token, address)
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
	fmt.Println("pull globular update from ", src, " to ", dest)

	admin_source, err := admin_client.NewAdminService_Client(src, "admin.AdminService")
	if err != nil {
		return err
	}

	// From the source I will download the new executable and save it in the
	// temp directory...
	path := os.TempDir() + "/Globular_" + Utility.ToString(time.Now().Unix())
	Utility.CreateDirIfNotExist(path)

	fmt.Println("download exec to ", path)
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
func publish(g *Globule, user, pwd, address, organization, path, platform string) error {

	// Authenticate the user in order to get the token
	authentication_client, err := authentication_client.NewAuthenticationService_Client(address, "authentication.AuthenticationService")
	if err != nil {
		fmt.Println(err)
		return err
	}

	token, err := authentication_client.Authenticate(user, pwd)
	if err != nil {
		log.Println(err)
		return err
	}

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	discovery_client_, err := discovery_client.NewDiscoveryService_Client(address, "discovery.PackageDiscovery")
	if err != nil {
		log.Println(err)
		return err
	}

	err = discovery_client_.PublishService(user, organization, token, discovery_client_.GetDomain(), path, platform)
	if err != nil {
		return err
	}

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	repository_client_, err := repository_client.NewRepositoryService_Client(address, "repository.PackageRepository")
	if err != nil {
		log.Println(err)
		return err
	}

	// first of all I will create and upload the package on the discovery...
	err = repository_client_.UploadServicePackage(user, organization, token, address, path, platform)
	if err != nil {
		log.Println(err)
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
		fmt.Println(err)
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
	err = services_manager_client_.InstallService(token, domain, discovery, publisherId, serviceId)
	if err != nil {
		log.Println("fail to install service", serviceId, "with error ", err.Error())
		return err
	}

	log.Println("service was installed")

	return nil
}

func uninstall_service(g *Globule, serviceId, publisherId, version, address, user, pwd string) error {

	// Authenticate the user in order to get the token
	authentication_client, err := authentication_client.NewAuthenticationService_Client(address, "authentication.AuthenticationService")
	if err != nil {
		fmt.Println(err)
		return err
	}

	token, err := authentication_client.Authenticate(user, pwd)
	if err != nil {
		log.Println(err)
		return err
	}

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	services_manager_client_, err := service_manager_client.NewServicesManagerService_Client(address, "services_manager.ServicesManagerService")
	if err != nil {
		log.Println(err)
		return err
	}

	// first of all I will create and upload the package on the discovery...
	err = services_manager_client_.UninstallService(token, address, publisherId, serviceId, version)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

/**
 * Install Globular web application.
 */
func install_application(g *Globule, applicationId, discovery, publisherId, domain, user, pwd string, set_as_default bool) error {

	// Authenticate the user in order to get the token
	authentication_client, err := authentication_client.NewAuthenticationService_Client(domain, "authentication.AuthenticationService")
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
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

	// formalize path
	path = strings.ReplaceAll(path, "\\", "/")

	// That function is use to install globular at a given repository.
	fmt.Println("create distribution for ", runtime.GOOS, runtime.GOARCH, "at", path)

	// first of all I will get the applications...
	// TODO see if those values can be use as parameters...

	// Now I will copy the application icon to the resource.
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	// The debian package...
	// There is a link on how to launch service in macOS
	// https://medium.com/swlh/how-to-use-launchd-to-run-services-in-macos-b972ed1e352
	//
	// But the way to use globular is pretty simple,
	// from the package, go to
	//
	// cd /Globular/Contents/MacOS
	//
	// from the terminal and with admin right you can run globular.
	// sudo ./Globular
	//
	// To install globular permentaly
	// sudo ./Globular install
	//
	// To uninstall Globular,
	// sudo ./Globular uninstall
	//
	// To stop/start the globular deamon,
	//
	// -Start the service
	// sudo launchctl load /Library/LaunchDaemons/Globular.plist
	//
	// -stop the service
	// sudo launchctl unload /Library/LaunchDaemons/Globular.plist
	//
	if runtime.GOOS == "darwin" {

		darwin_package_path := path + "/globular_" + g.Version + "-" + revision + "_" + runtime.GOARCH

		// remove existiong files...
		os.RemoveAll(darwin_package_path)

		// 1. Create the working directory
		Utility.CreateDirIfNotExist(darwin_package_path)

		// 2. Create application directory.
		app_path := darwin_package_path + "/Globular.app"
		app_content := app_path + "/Contents"
		app_bin := app_content + "/MacOS"
		app_resource := app_content + "/Resources"
		config_path := app_bin + "/etc/globular/config"
		applications_path := app_bin + "/var/globular/applications"

		// Copy applications for offline installation...
		Utility.CopyDir(dir+"/applications/.", applications_path)

		// create directories...
		Utility.CreateDirIfNotExist(app_content)
		Utility.CreateDirIfNotExist(app_bin)
		Utility.CreateDirIfNotExist(config_path)
		Utility.CreateDirIfNotExist(app_resource)

		err := Utility.CopyFile(dir+"/assets/icon.icns", app_resource+"/icon.icns")
		if err != nil {
			fmt.Println("fail to copy icon from ", dir+"/assets/icon.icns"+"whit error", err)
		}

		// Create the distribution.
		__dist(g, app_bin, config_path)

		// Now I will create the plist file.
		plistFile := `
		<?xml version="1.0" encoding="UTF-8"?>
		<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
		<plist version="1.0">
		<dict>
			<key>CFBundleExecutable</key>
			<string>globular</string>
			<key>CFBundleIconFile</key>
			<string>icon.icns</string>
			<key>CFBundleIdentifier</key>
			<string>io.globular</string>
			<key>NSHighResolutionCapable</key>
			<true/>
			<key>LSUIElement</key>
			<true/>
		</dict>
		</plist>
		`
		os.WriteFile(app_content+"/Info.plist", []byte(plistFile), 0644)

	} else if runtime.GOOS == "linux" {
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

		// globular data
		applications_path := debian_package_path + "/var/globular/applications"

		// globular configurations
		config_path := debian_package_path + "/etc/globular/config"

		// Create the bin directories.
		Utility.CreateDirIfNotExist(distro_path)
		Utility.CreateDirIfNotExist(data_path)
		Utility.CreateDirIfNotExist(config_path)
		Utility.CreateDirIfNotExist(applications_path)

		// Now the libraries...
		libpath := debian_package_path + "/usr/local/lib"
		Utility.CreateDirIfNotExist(libpath)

		if runtime.GOARCH == "amd64" {

			// zlib
			if !Utility.Exists("/usr/lib/x86_64-linux-gnu/libz.a") {
				fmt.Println("libz.a not found please install it on your computer: sudo apt-get install zlib1g-dev")
				return
			}

			Utility.CopyFile("/usr/lib/x86_64-linux-gnu/libz.a", libpath+"/libz.a")
			Utility.CopyFile("/usr/lib/x86_64-linux-gnu/libz.so.1.2.11", libpath+"/libz.so.1.2.11")
		} else if runtime.GOARCH == "arm64" {
			if !Utility.Exists("/usr/lib/aarch64-linux-gnu/libssl.so.1.1") {
				fmt.Println("libssl.so.1.1 not found on your computer, please install it: ")
				fmt.Println("   wget http://launchpadlibrarian.net/475575244/libssl1.1_1.1.1f-1ubuntu2_arm64.deb")
				fmt.Println("	sudo dpkg -i libssl1.1_1.1.1f-1ubuntu2_arm64.deb ")
				return
			}

			// Copy lib crypto...
			Utility.CopyFile("/usr/lib/aarch64-linux-gnu/libssl.so.1.1", libpath+"/libssl.so.1.1")
			Utility.CopyFile("/usr/lib/aarch64-linux-gnu/libcrypto.so.1.1", libpath+"/libcrypto.so.1.1")

			// zlib
			if !Utility.Exists("/usr/lib/aarch64-linux-gnu/libz.a") {
				fmt.Println("libz.a not found please install it on your computer: sudo apt-get install zlib1g-dev")
				return
			}

			Utility.CopyFile("/usr/lib/aarch64-linux-gnu/libz.a", libpath+"/libz.a")
			Utility.CopyFile("/usr/lib/aarch64-linux-gnu/libz.so.1.2.13", libpath+"/libz.so.1.2.13")
		}

		// ODBC libraries...
		Utility.CopyFile("/usr/local/lib/libodbc.la", libpath+"/libodbc.la")
		Utility.CopyFile("/usr/local/lib/libodbc.so.2.0.0", libpath+"/libodbc.so.2.0.0")

		Utility.CopyFile("/usr/local/lib/libodbccr.la", libpath+"/libodbccr.la")
		Utility.CopyFile("/usr/local/lib/libodbccr.so.2.0.0", libpath+"/libodbccr.so.2.0.0")

		Utility.CopyFile("/usr/local/lib/libodbcinst.la", libpath+"/libodbcinst.la")
		Utility.CopyFile("/usr/local/lib/libodbcinst.so.2.0.0", libpath+"/libodbcinst.so.2.0.0")

		// Now I will create get the configuration files from service and create a copy to /etc/globular/config/services
		// so modification will survice upgrades.

		// Create the distribution.
		configurations := __dist(g, distro_path, config_path)

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

		// - The project homepage
		packageConfig += "Homepage: https://globular.io\n"

		// - The list of dependencies...
		packageConfig += "Depends: python3 (>= 3.8.~), python-is-python3 (>=3.8.~), ffmpeg (>=4.4.~), curl(>=7.8.~), dpkg(>=1.21.~)\n"

		err := os.WriteFile(debian_package_path+"/DEBIAN/control", []byte(packageConfig), 0644)
		if err != nil {
			fmt.Println(err)
		}

		// Here I will set the script to run before the installation...
		// https://www.devdungeon.com/content/debian-package-tutorial-dpkgdeb#toc-17
		// TODO create tow version one for arm7 and one for amd64

		var preinst string

		if runtime.GOARCH == "amd64" {
			preinst = `
		echo "Welcome to Globular!-)"

		# Create the directory where the service will be install.
		mkdir /etc/globular/config/services

		# -- Install prometheus
		wget https://github.com/prometheus/prometheus/releases/download/v2.41.0/prometheus-2.41.0.linux-amd64.tar.gz
		tar -xf prometheus-2.41.0.linux-amd64.tar.gz
		cp prometheus-2.41.0.linux-amd64/prometheus /usr/local/bin/
		cp prometheus-2.41.0.linux-amd64/promtool /usr/local/bin/
		cp -r prometheus-2.41.0.linux-amd64/consoles /etc/prometheus/
		cp -r prometheus-2.41.0.linux-amd64/console_libraries /etc/prometheus/
		rm -rf prometheus-2.41.0.linux-amd64*

		# -- Install alert manager
		wget https://github.com/prometheus/alertmanager/releases/download/v0.25.0/alertmanager-0.25.0.linux-amd64.tar.gz
		tar -xf alertmanager-0.25.0.linux-amd64.tar.gz
		cp alertmanager-0.25.0.linux-amd64/alertmanager /usr/local/bin
		rm -rf alertmanager-0.25.0.linux-amd64*

		# -- Install node exporter
		wget https://github.com/prometheus/node_exporter/releases/download/v1.5.0/node_exporter-1.5.0.linux-amd64.tar.gz
		tar -xf node_exporter-1.5.0.linux-amd64.tar.gz
		cp node_exporter-1.5.0.linux-amd64/node_exporter /usr/local/bin
		rm -rf node_exporter-1.5.0.linux-amd64*

		# -- Install yt-dlp
		curl -L  https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp --output /usr/local/bin/yt-dlp
		chmod a+rx /usr/local/bin/yt-dlp

		# -- Install envoy proxy
		wget https://github.com/envoyproxy/envoy/releases/download/v1.28.0/envoy-1.28.0-linux-x86_64
		mv envoy-1.28.0-linux-x86_64 /usr/local/bin/envoy
		chmod a+rx /usr/local/bin/envoy

		# -- Install etcd
		wget https://github.com/etcd-io/etcd/releases/download/v3.5.10/etcd-v3.5.10-linux-amd64.tar.gz
		tar -xf etcd-v3.5.10-linux-amd64.tar.gz
		cp etcd-v3.5.10-linux-amd64/etcd /usr/local/bin
		cp etcd-v3.5.10-linux-amd64/etcdctl /usr/local/bin
		rm -rf etcd-v3.5.10-linux-amd64*
		

		if [ -f "/usr/local/bin/Globular" ]; then
			rm /usr/local/bin/Globular
			rm /usr/local/bin/torrent
			rm /usr/local/lib/libz.so
			rm /usr/local/lib/libz.so.1
			rm /usr/local/lib/libodbc.so.2
			rm /usr/local/lib/libodbc.so
			rm /usr/local/lib/libodbccr.so.2
			rm /usr/local/lib/libodbccr.so
			rm /usr/local/lib/libodbcinst.so.2
			rm /usr/local/lib/libodbcinst.so
		fi
		`
		} else if runtime.GOARCH == "arm64" {
			preinst = `
		echo "Welcome to Globular!-)"

		# Create the directory where the service will be install.
		mkdir /etc/globular/config/services

		# -- Install prometheus
		wget https://github.com/prometheus/prometheus/releases/download/v2.44.0/prometheus-2.44.0.linux-arm64.tar.gz
		tar -xf prometheus-2.44.0.linux-arm64.tar.gz
		cp prometheus-2.44.0.linux-arm64/prometheus /usr/local/bin/
		cp prometheus-2.44.0.linux-arm64/promtool /usr/local/bin/
		cp -r prometheus-2.44.0.linux-arm64/consoles /etc/prometheus/
		cp -r prometheus-2.44.0.linux-arm64/console_libraries /etc/prometheus/
		rm -rf prometheus-2.44.0.linux-arm64*

		# -- Install alert manager
		wget https://github.com/prometheus/alertmanager/releases/download/v0.25.0/alertmanager-0.25.0.linux-arm64.tar.gz
		tar -xf alertmanager-0.25.0.linux-arm64.tar.gz
		cp alertmanager-0.25.0.linux-arm64/alertmanager /usr/local/bin
		rm -rf alertmanager-0.25.0.linux-arm64*

		# -- Install node exporter
		https://github.com/prometheus/node_exporter/releases/download/v1.6.0/node_exporter-1.6.0.linux-arm64.tar.gz
		tar -xf node_exporter-1.6.0.linux-arm64.tar.gz
		cp node_exporter-1.6.0.linux-arm64/node_exporter /usr/local/bin
		rm -rf node_exporter-1.6.0.linux-arm64*

		# -- Install yt-dlp
		curl -L  https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp --output /usr/local/bin/yt-dlp
		chmod a+rx /usr/local/bin/yt-dlp

		if [ -f "/usr/local/bin/Globular" ]; then
			rm /usr/local/bin/Globular
			rm /usr/local/bin/torrent
			rm /usr/local/lib/libz.so
			rm /usr/local/lib/libz.so.1
			rm /usr/local/lib/libodbc.so.2
			rm /usr/local/lib/libodbc.so
			rm /usr/local/lib/libodbccr.so.2
			rm /usr/local/lib/libodbccr.so
			rm /usr/local/lib/libodbcinst.so.2
			rm /usr/local/lib/libodbcinst.so
		fi
		`

		}

		err = os.WriteFile(debian_package_path+"/DEBIAN/preinst", []byte(preinst), 0755)
		if err != nil {
			fmt.Println(err)
		}

		// the list of list of configurations
		// DEBIAN/conffiles
		conffiles := ""
		for i := 0; i < len(configurations); i++ {
			conffiles += configurations[i] + "\n"
		}

		err = os.WriteFile(debian_package_path+"/DEBIAN/conffiles", []byte(conffiles), 0755)
		if err != nil {
			fmt.Println(err)
		}

		postinst := `
		# Create a link into bin to it original location.
		# the systemd file is /etc/systemd/system/Globular.service
		# the environement variable file is /etc/sysconfig/Globular
	
		 echo "create link to Globular exec"
		 ln -s /usr/local/share/globular/Globular /usr/local/bin/Globular
		 chmod +x /usr/local/bin/Globular

		 echo "create link to torrent exec"
		 ln -s /usr/local/share/globular/bin/torrent /usr/local/bin/torrent
		 chmod +x /usr/local/bin/torrent

		 #create symlink
		 ln -s /usr/local/lib/libz.so.1.2.11 /usr/local/lib/libz.so
		 ln -s /usr/local/lib/libz.so.1.2.11 /usr/local/lib/libz.so.1
		 ln -s /usr/local/lib/libodbc.so.2.0.0 /usr/local/lib/libodbc.so.2
		 ln -s /usr/local/lib/libodbc.so.2.0.0 /usr/local/lib/libodbc.so
		 ln -s /usr/local/lib/libodbccr.so.2 /usr/local/lib/libodbccr.so.2
		 ln -s /usr/local/lib/libodbccr.so.2 /usr/local/lib/libodbccr.so
		 ln -s /usr/local/lib/libodbcinst.so.2.0.0 /usr/local/lib/libodbcinst.so.2
		 ln -s /usr/local/lib/libodbcinst.so.2 /usr/local/lib/libodbcinst.so
		 ln -s /usr/local/lib/libssl.so.1.1 /usr/local/lib/libssl.so.1.1
		 ln -s /usr/local/lib/libcrypto.so.1.1 /usr/local/libcrypto.so.1.1
		 ldconfig

		 echo "install globular as service..."
		 /usr/local/bin/Globular install
		 
		 # here I will modify the /etc/systemd/system/Globular.service file and set 
		 # Restart=always
		 # RestartSec=3
		 echo "set service configuration /etc/systemd/system/Globular.service"
		 sed -i 's/^\(Restart=\).*/\1always/' /etc/systemd/system/Globular.service
		 sed -i 's/^\(RestartSec=\).*/\120/' /etc/systemd/system/Globular.service

		 cd; cd -

		 systemctl daemon-reload
		 systemctl enable Globular

		 echo "start Globular service"
		 echo "run 'service Globular stop' to stop Globular service"

		 service Globular start
		 
		`
		err = os.WriteFile(debian_package_path+"/DEBIAN/postinst", []byte(postinst), 0755)
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
		err = os.WriteFile(debian_package_path+"/DEBIAN/prerm", []byte(prerm), 0755)
		if err != nil {
			fmt.Println(err)
		}

		postrm := `
		# Thing to do after removing
		if [ -f "/usr/local/bin/Globular" ]; then
			echo "remove /usr/local/bin/Globular and links"
			find /usr/local/bin/Globular -xtype l -delete
			rm /etc/systemd/system/Globular.service
			rm /usr/local/bin/Globular
			rm /usr/local/bin/torrent
			rm /usr/local/lib/libz.so
			rm /usr/local/lib/libz.so.1
			rm /usr/local/lib/libodbc.so.2
			rm /usr/local/lib/libodbc.so
			rm /usr/local/lib/libodbccr.so.2
			rm /usr/local/lib/libodbccr.so
			rm /usr/local/lib/libodbcinst.so.2
			rm /usr/local/lib/libodbcinst.so
		fi
		
		echo "Hope to see you again soon!"
		`
		err = os.WriteFile(debian_package_path+"/DEBIAN/postrm", []byte(postrm), 0755)
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

	} else if runtime.GOOS == "windows" {
		fmt.Println("Create the distro at path ", path)
		root := path + "/Globular"
		Utility.CreateIfNotExists(root, 0755)

		app := root + "/app"
		Utility.CreateIfNotExists(app, 0755)

		// Copy globular ditro file.
		__dist(g, app, app+"/config")

		// I will make use of NSIS: Nullsoft Scriptable Install System to create an installer for window.
		// so here I will create the setup.nsi file.

		// copy assets
		Utility.CreateDirIfNotExist(root + "/assets")
		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

		err := Utility.CopyDir(dir+"/assets/.", root+"/assets")
		if err != nil {
			fmt.Println("--> fail to copy assets ", err)
		}

		// redist
		Utility.CreateDirIfNotExist(root + "/redist")

		if Utility.Exists(dir + "/redist") {
			err = Utility.CopyDir(dir+"/redist/.", root+"/redist")
			if err != nil {
				fmt.Println("--> fail to copy redist", root+"/redist", err)
			}
		} else {
			fmt.Println("no directory found with path", dir+"/redist")
		}

		// copy the license
		err = Utility.CopyFile(dir+"/license.txt", root+"/license.txt")
		if err != nil {
			fmt.Println("--> fail to copy license ", err)
			err = nil
		}

		// Now I will create the setup.nsi file.
		setupNsi := `
;---------------------------------
; Includes

  !include "MUI2.nsh"
  !include "logiclib.nsh"
  !include "x64.nsh"

;---------------------------------
; Custom defines
  !define NAME "Globular"
  !define APPFILE "Globular.exe"
  !define VERSION "` + revision + `"
  !define SLUG "${NAME} v${VERSION}"

  ;---------------------------------
  ; General
	Name "${NAME}"
	OutFile "${NAME} Setup.exe"
	InstallDir "$PROGRAMFILES64\${NAME}"
	InstallDirRegKey HKCU  "Software\${NAME}" ""
	RequestExecutionLevel admin
  
  ;--------------------------------
  ; UI
	!define MUI_ICON "assets\globular_logo.ico"
	!define MUI_HEADERIMAGE
	!define MUI_WELCOMEFINISHPAGE_BITMAP "assets\welcome.bmp"
	!define MUI_HEADERIMAGE_BITMAP "assets\head.bmp"
	!define MUI_ABORTWARNING
	!define MUI_WELCOMEPAGE_TITLE "${SLUG} Setup"
  
  ;--------------------------------
  ; Pages
  
  ; Installer pages
	!insertmacro MUI_PAGE_WELCOME
	!insertmacro MUI_PAGE_LICENSE "license.txt"
	!insertmacro MUI_PAGE_INSTFILES
	!insertmacro MUI_PAGE_FINISH
  
  ; Uninstaller pages
	!insertmacro MUI_UNPAGE_CONFIRM
	!insertmacro MUI_UNPAGE_INSTFILES
  
  ; Set UI language
	!insertmacro MUI_LANGUAGE "English"
  
  ;--------------------------------
	Function .onInstSuccess
	  SetOutPath "$INSTDIR"
	  nsExec::ExecToStack '"${APPFILE}" install'
	  DetailPrint "Start Globular service please wait..."
	  nsExec::ExecToStack 'net start ${Name}'
	FunctionEnd
  
  ;--------------------------------
  ; Section - Visual Studio Runtime
   Section "Visual Studio Runtime"
   	SetOutPath "$INSTDIR\Redist"
	File "redist\VC_redist_2019_x64\VC_redist_x64.exe"
	DetailPrint "Running Visual Studio Redistribuatable VC2019"
	ExecWait "$INSTDIR\Redist\VC_redist_x64.exe"
	DetailPrint "Visual Studio Redistribuatable VC2019 is now installed"
  SectionEnd

  ;--------------------------------
  ; Section - Install App
  
	Section "-hidden app"
	  SectionIn RO
	  SetOutPath "$INSTDIR"
	  File /r "app\*.*" 
	  WriteRegStr HKCU  "Software\${NAME}" "" $INSTDIR
	  WriteUninstaller "$INSTDIR\Uninstall.exe"
	SectionEnd
  
  ;--------------------------------
  ; Remove empty parent directories
  
	Function un.RMDirUP
	  !define RMDirUP '!insertmacro RMDirUPCall'
  
	  !macro RMDirUPCall _PATH
		  push '${_PATH}'
		  Call un.RMDirUP
	  !macroend
  
	  ; $0 - current folder
	  ClearErrors
  
	  Exch $0
	  ;DetailPrint "ASDF - $0\.."
	  RMDir "$0\.."
  
	  IfErrors Skip
	  ${RMDirUP} "$0\.."
	  Skip:
  
	  Pop $0
  
	FunctionEnd
  
  
  ;--------------------------------
  ; Section - Uninstaller
  
	Section "Uninstall"
	  SetOutPath "$INSTDIR"
	  
	  nsExec::ExecToStack 'net stop ${Name}'
	  nsExec::ExecToStack '"${APPFILE}" uninstall'
	  
	  ;Delete files
	  Delete "$INSTDIR\Uninstall.exe"
	  Delete "$INSTDIR\Globular.exe"
	  Delete "$INSTDIR\Dockerfile"
  
	  ;Delete Folder's *keep config webroot and data folder.
	  RMDir /r "$INSTDIR\bin"
	  RMDir /r "$INSTDIR\dependencies"
	  RMDir /r "$INSTDIR\Redist"
	  RMDir /r "$INSTDIR\services"
	  RMDir /r "$INSTDIR\applications"
  
	  DeleteRegKey /ifempty HKCU  "Software\${NAME}"
  
	SectionEnd
  
`
		err = Utility.WriteStringToFile(root+"/setup.nsi", setupNsi)
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

func __dist(g *Globule, path, config_path string) []string {

	// Return the configurations list
	configs := make([]string, 0)

	// I will set the docker file depending of the arch.
	var dockerfile string
	if runtime.GOARCH == "amd64" {
		data, err := os.ReadFile("Dockerfile_amd64")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		dockerfile = string(data)
	} else if runtime.GOARCH == "arm64" {
		data, err := os.ReadFile("Dockerfile_arm64")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		dockerfile = string(data)
	}

	Utility.CreateDirIfNotExist(path)

	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir = strings.ReplaceAll(dir, "\\", "/")

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

	err = os.Chmod(path+"/"+destExec, 0755)
	if err != nil {
		fmt.Println(err)
	}


	// Copy the bin file from globular
	if runtime.GOOS == "windows" {
		Utility.CreateDirIfNotExist(path + "/dependencies")
		if Utility.Exists(dir + "/dependencies") {

			err = Utility.CopyDir(dir+"/dependencies/.", path+"/dependencies")
			if err != nil {
				fmt.Println("--> fail to copy dependencies ", err)
			}

			execs := Utility.GetFilePathsByExtension(path+"/dependencies", ".exe")
			for i := 0; i < len(execs); i++ {
				err = os.Chmod(execs[i], 0755)
				if err != nil {
					fmt.Println(err)
				}
			}
		} else {
			fmt.Println("no dir with dependencie was found at path", dir+"/dependencies")
		}
	} else if runtime.GOOS == "darwin" {
		dest := path + "/bin"

		// ffmpeg
		if Utility.Exists("/usr/local/bin/ffmpeg") {
			Utility.CopyFile("/usr/local/bin/ffmpeg", dest)
		}

		if Utility.Exists("/usr/local/bin/ffprobe") {
			Utility.CopyFile("/usr/local/bin/ffprobe", dest)
		}

		// yt-dlp
		if Utility.Exists("/usr/local/bin/yt-dlp") {
			Utility.CopyFile("/usr/local/bin/yt-dlp", dest)
		}

		// unix-odbc
		if Utility.Exists("/usr/local/bin/odbc_config") {
			Utility.CopyFile("/usr/local/bin/odbc_config", dest)
		}

		if Utility.Exists("/usr/local/bin/odbcinst") {
			Utility.CopyFile("/usr/local/bin/odbcinst", dest)
		}

		// Now the libraries...
		dest = path + "/lib"
		Utility.CreateDirIfNotExist(dest)

		if Utility.Exists("/usr/local/lib/libodbc.2.dylib") {
			Utility.CopyFile("/usr/local/lib/libodbc.2.dylib", dest)
		}

		if Utility.Exists("/usr/local/lib/libodbc.dylib") {
			Utility.CopyFile("/usr/local/lib/libodbc.dylib", dest)
		}

		if Utility.Exists("/usr/local/lib/libodbc.la") {
			Utility.CopyFile("/usr/local/lib/libodbc.la", dest)
		}

		if Utility.Exists("/usr/local/lib/libodbccr.2.dylib") {
			Utility.CopyFile("/usr/local/lib/libodbccr.2.dylib", dest)
		}

		if Utility.Exists("/usr/local/lib/libodbccr.dylib") {
			Utility.CopyFile("/usr/local/lib/libodbccr.dylib", dest)
		}

		if Utility.Exists("/usr/local/lib/libodbccr.la") {
			Utility.CopyFile("/usr/local/lib/libodbccr.la", dest)
		}

		if Utility.Exists("/usr/local/lib/libodbcinst.2.dylib") {
			Utility.CopyFile("/usr/local/lib/libodbcinst.2.dylib", dest)
		}

		if Utility.Exists("/usr/local/lib/libodbcinst.dylib") {
			Utility.CopyFile("/usr/local/lib/libodbcinst.dylib", dest)
		}

		if Utility.Exists("/usr/local/lib/libodbcinst.la") {
			Utility.CopyFile("/usr/local/lib/libodbcinst.la", dest)
		}
	}

	// install services...
	services, err := config_.GetServicesConfigurations()
	if err != nil {
		log.Println("fail to retreive services with error ", err)
	}

	var programFilePath string
	// fmt.Println("fail to find service configurations at at path ", serviceConfigDir, "with error ", err)
	if runtime.GOOS == "windows" {
		if runtime.GOARCH == "386" {
			programFilePath, _ = Utility.GetEnvironmentVariable("PROGRAMFILES(X86)")
			programFilePath += "/Globular"
		} else {
			programFilePath, _ = Utility.GetEnvironmentVariable("PROGRAMFILES")
			programFilePath += "/Globular"
		}
	} else {
		programFilePath = "/usr/local/share/globular"
	}

	programFilePath = strings.ReplaceAll(programFilePath, "\\", "/")

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
				configPath := filepath.Dir(execPath) + "/config.json"
				if Utility.Exists(configPath) {
					log.Println("install service ", name)
					bytes, err := os.ReadFile(configPath)
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

								execName := filepath.Base(execPath)
								destPath := path + "/" + serviceDir + "/" + id + "/" + execName

								if Utility.Exists(execPath) {
									Utility.CreateDirIfNotExist(path + "/" + serviceDir + "/" + id)

									err := Utility.Copy(execPath, destPath)
									if err != nil {
										fmt.Println(execPath, destPath, err)
									}

									// Set execute permission
									err = os.Chmod(destPath, 0755)
									if err != nil {
										fmt.Println(err)
									}

									config["Path"] = programFilePath + "/" + serviceDir + "/" + id + "/" + execName
									config["Proto"] = programFilePath + "/" + serviceDir + "/" + config["Name"].(string) + ".proto"

									if config["CacheType"] != nil {
										config["CacheType"] = "BADGER"
									}

									if config["CacheAddress"] != nil {
										config["CacheAddress"] = "localhost"
									}

									if config["Backend_address"] != nil {
										config["Backend_address"] = "localhost"
									}

									if config["Backend_type"] != nil {
										config["Backend_type"] = "SQL"
									}

									// set the security values to nothing...
									config["CertAuthorityTrust"] = ""
									config["Process"] = -1
									config["ProxyProcess"] = -1
									config["CertFile"] = ""
									config["KeyFile"] = ""
									config["TLS"] = false

									if config["Root"] != nil {
										if name == "file.FileService" {
											config["Root"] = config_.GetDataDir() + "/files"

											// I will also copy the mime type directory
											config["Public"] = make([]string, 0)
											Utility.CopyDir(filepath.Dir(execPath)+"/mimetypes", path+"/"+serviceDir+"/"+id)

										} else if name == "conversation.ConversationService" {
											config["Root"] = config_.GetDataDir()
										}
									}

									// Empty the list of connections if connections exist for the services.
									if config["Connections"] != nil {
										config["Connections"] = make(map[string]interface{})
									}

									config["ConfigPath"] = config_.GetConfigDir() + "/" + serviceDir + "/" + id + "/config.json"
									configs = append(configs, "/etc/globular/config/"+serviceDir+"/"+id+"/config.json")
									str, _ := Utility.ToJson(&config)

									if len(config_path) > 0 {
										// So here I will set the service configuration at /etc/globular/config/service... to be sure configuration
										// will survive package upgrades...
										Utility.CreateDirIfNotExist(config_path + "/" + serviceDir + "/" + id)
										os.WriteFile(config_path+"/"+serviceDir+"/"+id+"/config.json", []byte(str), 0644)

									}

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
				if Utility.Exists(os.Getenv("ServicesRoot") + "/" + protoPath) {
					Utility.Copy(os.Getenv("ServicesRoot")+"/"+protoPath, path+"/"+protoPath)
				}
			}
		}

	}

	// save docker.
	err = os.WriteFile(path+"/Dockerfile", []byte(dockerfile), 0644)
	if err != nil {
		log.Println(err)
	}

	//
	return configs
}

/**
 * Generate a token that will be valid for 15 minutes or the session timeout delay.
 */
func generate_token(g *Globule, address, user, pwd string) error {

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

	// simply print the token in the console.
	fmt.Println(token)

	return nil
}

/**
 * Connect one peer's with another. When connected peer's are able to generate token valid for both side.
 * The usr and pwd are the admin password in the destionation (ns1.mycelius.com)
 * ex. ./Globular connect_peer -dest=ns1.mycelius.com -u=sa -p=adminadmin
 */
func connect_peer(g *Globule, address, token string) error {

	// Create the remote ressource service
	remote_resource_client_, err := resource_client.NewResourceService_Client(address, "resource.ResourceService")
	if err != nil {
		return err
	}

	// Get the local peer key
	key, err := security.GetPeerKey(globule.Mac)
	if err != nil {
		log.Println(err)
		return err
	}

	// Register the peer on the remote resourse client...
	hostname, _ := os.Hostname()

	peer, key_, err := remote_resource_client_.RegisterPeer(string(key), &resourcepb.Peer{Hostname: hostname, Mac: g.Mac, Domain: g.Domain, LocalIpAddress: config_.GetLocalIP(), ExternalIpAddress: Utility.MyIP()})
	if err != nil {
		return err
	}

	address_, _ := config_.GetAddress()
	// I will also register the peer to the local server, the local server must running and it domain register,
	// he can be set in /etc/hosts if it's not a public domain.
	local_resource_client_, err := resource_client.NewResourceService_Client(address_, "resource.ResourceService")
	if err != nil {
		return err
	}

	_, _, err = local_resource_client_.RegisterPeer(key_, peer)

	return err
}
