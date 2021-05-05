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
		// ex. Globular install_service -publisher=globulario -discovery=globular.io -service=echo.EchoService -a=globular1.globular.io -u=sa -p=*******
		install_service_command := flag.NewFlagSet("install_service", flag.ExitOnError)
		install_service_command_publisher := install_service_command.String("publisher", "", "The publisher id (Required)")
		install_service_command_discovery := install_service_command.String("discovery", "", "The addresse where the service was publish (Required)")
		install_service_command_service := install_service_command.String("service", "", " the service name ex file.FileService (Required)")
		install_service_command_address := install_service_command.String("a", "", "The domain of the server where to install the service (Required)")
		install_service_command_user := install_service_command.String("u", "", "The user name. (Required)")
		install_service_command_pwd := install_service_command.String("p", "", "The user password. (Required)")

		// Uninstall a service on the server.
		uninstall_service_command := flag.NewFlagSet("uninstall_service", flag.ExitOnError)
		uninstall_service_command_service := uninstall_service_command.String("service", "", " the service name ex file.FileService (Required)")
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

		// pull globular update
		update_globular_from_command := flag.NewFlagSet("update_from", flag.ExitOnError)
		update_globular_command_from_source := update_globular_from_command.String("source", "", " the address of the server from where to update the a given server.")
		update_globular_from_command_dest := update_globular_from_command.String("a", "", "The domain of the server to update (Required)")
		update_globular_from_command_user := update_globular_from_command.String("u", "", "The user name. (Required)")
		update_globular_from_command_pwd := update_globular_from_command.String("p", "", "The user password. (Required)")
		update_globular_from_command_platform := update_globular_from_command.String("platform", "", "The os and arch info ex: linux:arm64 (optional)")

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

		default:
			flag.PrintDefaults()
			os.Exit(1)
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
func deploy(g *Globule, name string, organization string, path string, address string, user string, pwd string, set_as_default bool) error {

	log.Println("deploy application", name, " to address ", address, " user ", user)

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

	_, err = admin_client_.DeployApplication(user, name, organization, path, token, address, set_as_default)
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

	admin_source.DownloadGlobular(src, platform, path)

	path_ := path
	path_ += "/Globular"
	if runtime.GOOS == "windows" {
		path_ += ".exe"
	}

	if !Utility.Exists(path) {
		err := errors.New(path_ + " not found!")
		return err
		log.Println("---> the executable path was not found!")
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

func install_service(g *Globule, serviceId, discovery, publisherId, domain, user, pwd string) error {
	log.Println("try to install service",serviceId, "on", domain)
	// Authenticate the user in order to get the token
	resource_client_, err := resource_client.NewResourceService_Client(domain, "resource.ResourceService")
	if err != nil {
		log.Panicln(err)
		return err
	}
	token, err := resource_client_.Authenticate(user, pwd)
	if err != nil {

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
	err = admin_client_.InstallService(token, domain, user, discovery, publisherId, serviceId)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func uninstall_service(g *Globule, serviceId, publisherId, version, domain, user, pwd string) error {

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
	err = admin_client_.UninstallService(token, domain, user, publisherId, serviceId, version)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func install_application(g *Globule, applicationId, discovery, publisherId, domain, user, pwd string, set_as_default bool) error {

	// Authenticate the user in order to get the token
	resource_client_, err := resource_client.NewResourceService_Client(domain, "resource.ResourceService")
	if err != nil {
		log.Panicln(err)
		return err
	}

	token, err := resource_client_.Authenticate(user, pwd)
	if err != nil {
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
	err = admin_client_.InstallApplication(token, domain, user, discovery, publisherId, applicationId, set_as_default)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func uninstall_application(g *Globule, applicationId, publisherId, version, domain, user, pwd string) error {

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
	err = admin_client_.UninstallApplication(token, domain, user, publisherId, applicationId, version)
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

		// test the
		/*
			installMongoScript :=`
			cd ~/
			wget http://downloads.mongodb.org/linux/mongodb-linux-x86_64-1.2.2.tgz
			tar -xf mongodb-linux-x86_64-1.2.2.tgz
			mkdir -p /var/mongodb
			mv mongodb-linux-x86_64-1.2.2/* /var/mongodb/
			ln -nfs /var/mongodb/bin/mongod /usr/local/sbin

			# Ensures the required folders exist for MongoDB to run properly
			mkdir -p /data/db
			mkdir -p /usr/local/mongodb/logs

			# Downloads the MongoDB init script to the /etc/init.d/ folder
			# Renames it to mongodb and makes it executable
			cd /etc/init.d/
			wget http://gist.github.com/raw/162954/f5d6434099b192f2da979a0356f4ec931189ad07/gistfile1.sh
			mv gistfile1.sh mongodb
			chmod +x mongodb

			# Ensure MongoDB starts up automatically after every system (re)boot
			update-rc.d mongodb start 51 S .

			# Starts up MongoDB right now
			/etc/init.d/mongodb start
			`
		*/

		// Here I will set the script to run before the installation...
		// https://www.devdungeon.com/content/debian-package-tutorial-dpkgdeb#toc-17
		preinst := `
		echo "Welcome to Globular!-)"
		echo "Note on Dependencies"
		echo "Those dependencies can be install after globular installation. It's strongly recommended"
		echo "to have MogonDB installed, all Globular applications made use of it. But if you plan to"
		echo "use Globular as simple webserver you can live whitout MongoDB."
		echo "Prometheus is use to keep Globular monitoring informations."
		echo "Youtube-dl, Transmission-cli and FFmpeg give are required if you plan to use media files (mp3, mp4...) in"
		echo "your applications."
		echo "- 1. MongoDB https://docs.mongodb.com/manual/tutorial/install-mongodb-on-ubuntu/"
		echo "- 2. Prometheus https://prometheus.io/download/"
		echo "- 3. Youtube-dl sudo apt-get install youtube-dl"
		echo "- 4. Transmission-cli https://phil.tech/2009/How-to-Install-Transmission-CLI-to-Ubuntu-Server/"
		echo "- 5. FFmpeg sudo apt install ffmpeg"
		`

		err = ioutil.WriteFile(debian_package_path+"/DEBIAN/preinst", []byte(preinst), 0755)
		if err != nil {
			fmt.Println(err)
		}

		postinst := `
		# Create a link into bin to it original location.
		 echo "install globular as service..."
		 ln -s /usr/local/share/globular/Globular /usr/local/bin/Globular
		 chmod ugo+x /usr/local/bin/Globular
		 /usr/local/bin/Globular install
		 systemctl start Globular
		 systemctl enable Globular
		 echo "To complete your server setup go to http://localhost"
		`
		err = ioutil.WriteFile(debian_package_path+"/DEBIAN/postinst", []byte(postinst), 0755)
		if err != nil {
			fmt.Println(err)
		}

		prerm := `
		# Thing to do before removing
		# Stop, Disable and Uninstall Globular service.
		echo "Stop runing globular service..."
		systemctl stop Globular
		systemctl disable Globular
		echo "Unistall globular service..."
		/usr/local/bin/Globular uninstall
		`
		err = ioutil.WriteFile(debian_package_path+"/DEBIAN/prerm", []byte(prerm), 0755)
		if err != nil {
			fmt.Println(err)
		}

		postrm := `
		# Thing to do after removing
		find /usr/local/bin/Globular -xtype l -delete
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
		fmt.Println("Only linux debian package are available as distro at the moment...")
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

	// Here I will copy the proxy.
	globularExec := os.Args[0]
	if runtime.GOOS == "windows" && !strings.HasSuffix(globularExec, ".exe") {
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
			err = os.Chmod(path+"/bin/"+f.Name(), 0755)
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

									config["Path"] = "/usr/local/share/globular/" + serviceDir + "/" + id + "/" + execName
									config["Proto"] = "/usr/local/share/globular/" + serviceDir + "/" + name + ".proto"

									// set the security values to nothing...
									config["CertAuthorityTrust"] = ""
									config["CertFile"] = ""
									config["KeyFile"] = ""
									config["TLS"] = false

									if config["Root"] != nil {
										if name == "file.FileService" {
											config["Root"] = "/usr/local/share/globular/data/files"
										} else if name == "conversation.ConversationService" {
											config["Root"] = "/usr/local/share/globular/data"
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
	err = ioutil.WriteFile(path+"/Dockerfile", []byte(dockerfile), 0644)
	if err != nil {
		log.Println(err)
	}
}
