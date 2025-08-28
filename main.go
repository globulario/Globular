// main package
package main

import (
	"Globular/internal/observability"
	"bytes"
	"context"
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

	"github.com/globulario/services/golang/admin/admin_client"
	"github.com/globulario/services/golang/applications_manager/applications_manager_client"
	"github.com/globulario/services/golang/authentication/authentication_client"
	configpkg "github.com/globulario/services/golang/config"
	"github.com/globulario/services/golang/discovery/discovery_client"
	"github.com/globulario/services/golang/repository/repository_client"
	"github.com/globulario/services/golang/resource/resource_client"
	"github.com/globulario/services/golang/resource/resourcepb"
	"github.com/globulario/services/golang/security"
	serviceManagerClient "github.com/globulario/services/golang/services_manager/services_manager_client"
	Utility "github.com/globulario/utility"
	"github.com/kardianos/service"
	"github.com/polds/imgbase64"
	//"github.com/pkg/profile"
)

// Start initiates the Globule service asynchronously without blocking the caller.
// It initializes a platform logger via the service manager (Event Log on Windows,
// syslog/journal on Linux/macOS) and logs lifecycle events and run-loop errors.
func (g *Globule) Start(s service.Service) error {

	// Initialize the platform logger once if possible.
	if g.logger == nil && s != nil {
		if l, err := s.Logger(nil); err == nil {
			g.logger = l
		} else {
			// Fall back to standard logger if the system logger can't be created.
			log.Println("failed to create service logger:", err)
		}
	}

	// Log a lifecycle message.
	if g.logger != nil {
		_ = g.logger.Info("Starting Globular service...")
	} else {
		log.Println("Starting Globular service...")
	}

	// Run the main loop in the background and report any errors via the logger.
	go func() {
		if err := g.run(); err != nil {
			if g.logger != nil {
				_ = g.logger.Error(err)
			} else {
				log.Println("Globular run error:", err)
			}
		}
	}()

	return nil
}

// Stop gracefully shuts down the Globule service by setting the exit flag and
// closing the exit channel. It logs the stop lifecycle events using the
// platform logger when available.
func (g *Globule) Stop(s service.Service) error {
	// Ensure we have a logger to report shutdown.
	if g.logger == nil && s != nil {
		if l, err := s.Logger(nil); err == nil {
			g.logger = l
		} else {
			log.Println("failed to create service logger:", err)
		}
	}

	if g.logger != nil {
		_ = g.logger.Info("Stopping Globular service...")
	} else {
		log.Println("Stopping Globular service...")
	}

	// Perform fast, graceful shutdown.
	g.isExit = true
	close(g.exit)

	if g.logger != nil {
		_ = g.logger.Info("Globular service stopped.")
	} else {
		log.Println("Globular service stopped.")
	}
	return nil
}

func (g *Globule) run() error {
	errs := make(chan error, 1)

	go func() {
		if err := g.Serve(); err != nil {
			errs <- err
		}
	}()

	select {
	case <-g.exit:
		return nil
	case err := <-errs:
		return err
	}
}

func main() {

	fmt.Println("Globular server is starting...")

	ctx := context.Background()
	shutdown := observability.Init(ctx)

	defer func() { _ = shutdown(context.Background()) }()

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

		// Start with specific parameter.
		startCommand := flag.NewFlagSet("start", flag.ExitOnError)
		startCommandDomain := startCommand.String("domain", "", "The domain of the service.")

		// Install globular as service/daemon
		installCommand := flag.NewFlagSet("install", flag.ExitOnError)
		installCommandName := installCommand.String("name", "", "The display name of globular service.")

		// Uninstall globular as service.
		uninstallCommand := flag.NewFlagSet("uninstall", flag.ExitOnError)

		// Package development environment into a given path
		distCommand := flag.NewFlagSet("dist", flag.ExitOnError)
		distCommandPath := distCommand.String("path", "", "You must specify the distribution path. (Required)")
		distCommandRevision := distCommand.String("revision", "", "You must specify the package revision. (Required)")

		// Deploy command
		deployCommand := flag.NewFlagSet("deploy", flag.ExitOnError)
		deployCommandName := deployCommand.String("name", "", "You must specify an application name. (Required)")
		deployCommandOrganization := deployCommand.String("o", "", "The name of the organisation that responsible of the application. (Required)")
		deployCommandPath := deployCommand.String("path", "", "You must specify the path that contain the source (bundle.js, index.html...) of the application to deploy. (Required)")
		deployCommandUser := deployCommand.String("u", "", "The user name. (Required)")
		deployCommandPwd := deployCommand.String("p", "", "The user password. (Required)")
		deployCommandAddress := deployCommand.String("a", "", "The domain of the server where to install the appliction (Required)")
		deployCommandIndex := deployCommand.String("set_as_default", "", "The value is true the application will be set as default (Optional false by default)")

		// Publish Service.
		publishCommand := flag.NewFlagSet("publish", flag.ExitOnError)
		publishCommandPath := publishCommand.String("path", "", "You must specify the path that contain the config.json, .proto and all dependcies require by the service to run. (Required)")
		publishCommandUser := publishCommand.String("u", "", "The user name. (Required)")
		publishCommandPwd := publishCommand.String("p", "", "The user password. (Required)")
		publishCommandAddress := publishCommand.String("a", "", "The domain of the server where to publish the service (Required)")
		publishCommandOrganization := publishCommand.String("o", "", "The Organization that publish the service. (Optional)")
		publishCommandPlatform := publishCommand.String("platform", "", "(Optional it take your current platform as default.)")

		// Install certificates on a server from a local service command.
		installCertificatesCommand := flag.NewFlagSet("certificates", flag.ExitOnError)
		installCertificatesCommandPath := installCertificatesCommand.String("path", "", "You must specify where to install certificate (Required)")
		installCertificatesCommandPort := installCertificatesCommand.String("port", "", "You must specify the port where the configuration can be found (Required)")
		installCertificatesCommandDomain := installCertificatesCommand.String("domain", "", "You must specify the domain (Required)")

		// Install a service on the server.
		installServiceCommand := flag.NewFlagSet("installService", flag.ExitOnError)
		installServiceCommandPublisher := installServiceCommand.String("publisher", "", "The publisher id (Required)")
		installServiceCommandDiscovery := installServiceCommand.String("discovery", "", "The addresse where the service was publish (Required)")
		installServiceCommandService := installServiceCommand.String("service", "", " the service id (uuid) (Required)")
		installServiceCommandAddress := installServiceCommand.String("a", "", "The domain of the server where to install the service (Required)")
		installServiceCommandUser := installServiceCommand.String("u", "", "The user name. (Required)")
		installServiceCommandPwd := installServiceCommand.String("p", "", "The user password. (Required)")

		// Uninstall a service on the server.
		uninstallServiceCommand := flag.NewFlagSet("uninstallService", flag.ExitOnError)
		uninstallServiceCommandService := uninstallServiceCommand.String("service", "", " the service uuid (Required)")
		uninstallServiceCommandPublisher := uninstallServiceCommand.String("publisher", "", "The publisher id (Required)")
		uninstallServiceCommandVersion := uninstallServiceCommand.String("version", "", " The service vesion(Required)")
		uninstallServiceCommandAddress := uninstallServiceCommand.String("a", "", "The domain of the server where to install the service (Required)")
		uninstallServiceCommandUser := uninstallServiceCommand.String("u", "", "The user name. (Required)")
		uninstallServiceCommandPwd := uninstallServiceCommand.String("p", "", "The user password. (Required)")

		// Install an application on the server.
		installApplicationCommand := flag.NewFlagSet("installApplication", flag.ExitOnError)
		installApplicationCommandPublisher := installApplicationCommand.String("publisher", "", "The publisher id (Required)")
		installApplicationCommandDiscovery := installApplicationCommand.String("discovery", "", "The addresse where the application was publish (Required)")
		installApplicationCommandName := installApplicationCommand.String("application", "", " the application name (Required)")
		installApplicationCommandAddress := installApplicationCommand.String("a", "", "The domain of the server where to install the application (Required)")
		installApplicationCommandUser := installApplicationCommand.String("u", "", "The user name. (Required)")
		installApplicationCommandPwd := installApplicationCommand.String("p", "", "The user password. (Required)")
		installApplicationCommandIndex := installApplicationCommand.String("set_as_default", "", "The value is true the application will be set as default (Optional false by default)")

		// Uninstall an application on the server.
		uninstallApplicationCommand := flag.NewFlagSet("uninstallApplication", flag.ExitOnError)
		uninstallApplicationCommandName := uninstallApplicationCommand.String("application", "", " the application name (Required)")
		uninstallApplicationCommandPublisher := uninstallApplicationCommand.String("publisher", "", "The publisher id (Required)")
		uninstallApplicationCommandVersion := uninstallApplicationCommand.String("version", "", " The application vesion(Required)")
		uninstallApplicationCommandAddress := uninstallApplicationCommand.String("a", "", "The domain where the application is running (Required)")
		uninstallApplicationCommandUser := uninstallApplicationCommand.String("u", "", "The user name. (Required)")
		uninstallApplicationCommandPwd := uninstallApplicationCommand.String("p", "", "The user password. (Required)")

		// push globular update.
		updateGlobularCommand := flag.NewFlagSet("update", flag.ExitOnError)
		updateGlobularCommandExecPath := updateGlobularCommand.String("path", "", " the path to the new executable to update from")
		updateGlobularCommandAddress := updateGlobularCommand.String("a", "", "The domain of the server where to push the update(Required)")
		updateGlobularCommandUser := updateGlobularCommand.String("u", "", "The user name. (Required)")
		updateGlobularCommandPwd := updateGlobularCommand.String("p", "", "The user password. (Required)")
		updateGlobularCommandPlatform := updateGlobularCommand.String("platform", "", "The os and arch info ex: linux:arm64 (optional)")

		// pull globular update.
		updateGlobularFromCommand := flag.NewFlagSet("update_from", flag.ExitOnError)
		updateGlobularCommandFromSource := updateGlobularFromCommand.String("source", "", " the address of the server from where to update the a given server.")
		updateGlobularFromCommandDest := updateGlobularFromCommand.String("a", "", "The domain of the server to update (Required)")
		updateGlobularFromCommandUser := updateGlobularFromCommand.String("u", "", "The user name. (Required)")
		updateGlobularFromCommandPwd := updateGlobularFromCommand.String("p", "", "The user password. (Required)")
		updateGlobularFromCommandPlatform := updateGlobularFromCommand.String("platform", "", "The os and arch info ex: linux:arm64 (optional)")

		// Connect peer one to another. The peer Domain must be set before calling that function.
		connectPeerCommand := flag.NewFlagSet("connectPeer", flag.ExitOnError)
		connectPeerCommandAddress := connectPeerCommand.String("dest", "", "The address of the peer to connect to, can contain it configuration port (80) by defaut.")
		connectPeerCommandToken := connectPeerCommand.String("token", "", "The token valid on the destination peer (Required)")

		// Generate a token from a given globule.
		generateTokenCommand := flag.NewFlagSet("generate_token", flag.ExitOnError)
		generateTokenCommandAddress := generateTokenCommand.String("dest", "", "The address of the peer to connect to, can contain it configuration port (80) by defaut.")
		generateTokenCommandUser := generateTokenCommand.String("u", "", "The user name. (Required)")
		generateTokenCommandPwd := generateTokenCommand.String("p", "", "The user password. (Required)")

		switch os.Args[1] {
		case "start":
			err = startCommand.Parse(os.Args[2:])
		case "dist":
			err = distCommand.Parse(os.Args[2:])
		case "deploy":
			err = deployCommand.Parse(os.Args[2:])
		case "publish":
			err = publishCommand.Parse(os.Args[2:])
		case "install":
			err = installCommand.Parse(os.Args[2:])
		case "uninstall":
			err = uninstallCommand.Parse(os.Args[2:])
		case "update":
			err = updateGlobularCommand.Parse(os.Args[2:])
		case "update_from":
			err = updateGlobularFromCommand.Parse(os.Args[2:])
		case "installService":
			err = installServiceCommand.Parse(os.Args[2:])
		case "uninstallService":
			err = uninstallServiceCommand.Parse(os.Args[2:])
		case "installApplication":
			err = installApplicationCommand.Parse(os.Args[2:])
		case "uninstallApplication":
			err = uninstallApplicationCommand.Parse(os.Args[2:])
		case "certificates":
			err = installCertificatesCommand.Parse(os.Args[2:])
		case "connectPeer":
			err = connectPeerCommand.Parse(os.Args[2:])
		case "generate_token":
			err = generateTokenCommand.Parse(os.Args[2:])
		default:
			flag.PrintDefaults()
			os.Exit(1)
		}

		// Check for errors
		if err != nil {
			fmt.Println("Error parsing command:", err)
			os.Exit(1)
		}

		if generateTokenCommand.Parsed() {
			address := *generateTokenCommandAddress
			if *generateTokenCommandAddress == "" {
				address = g.getAddress()
				if strings.Contains(address, ":") {
					address = strings.Split(address, ":")[0]
				}
			}

			if *generateTokenCommandUser == "" {
				generateTokenCommand.PrintDefaults()
				fmt.Println("no user was given!")
				os.Exit(1)
			}

			if *generateTokenCommandPwd == "" {
				generateTokenCommand.PrintDefaults()
				fmt.Println("no password was given!")
				os.Exit(1)
			}

			err = generateToken(address, *generateTokenCommandUser, *generateTokenCommandPwd)
			if err != nil {
				log.Println("fail to generate token:", err)
			}
		}

		if connectPeerCommand.Parsed() {
			if *connectPeerCommandAddress == "" {
				connectPeerCommand.PrintDefaults()
				fmt.Println("No peer address given!")
				os.Exit(1)
			}
			if *connectPeerCommandToken == "" {
				connectPeerCommand.PrintDefaults()
				fmt.Println("no user was given!")
				os.Exit(1)
			}

			err = connectPeer(g, *connectPeerCommandAddress)
			if err != nil {
				log.Println("fail to connect peer:", err)
			}
		}

		if installServiceCommand.Parsed() {
			if *installServiceCommandService == "" {
				installServiceCommand.PrintDefaults()
				fmt.Println("no service name was given!")
				os.Exit(1)
			}
			if *installServiceCommandDiscovery == "" {
				installServiceCommand.PrintDefaults()
				fmt.Println("no discovery address was given!")
				os.Exit(1)
			}
			if *installServiceCommandPublisher == "" {
				installServiceCommand.PrintDefaults()
				fmt.Println("no publisher was given!")
				os.Exit(1)
			}
			if *installServiceCommandAddress == "" {
				installServiceCommand.PrintDefaults()
				fmt.Println("no domain was given!")
				os.Exit(1)
			}
			if *installServiceCommandUser == "" {
				installServiceCommand.PrintDefaults()
				fmt.Println("no user was given!")
				os.Exit(1)
			}
			if *installServiceCommandPwd == "" {
				installServiceCommand.PrintDefaults()
				fmt.Println("no password was given!")
				os.Exit(1)
			}
			err = installService(*installServiceCommandService, *installServiceCommandDiscovery, *installServiceCommandPublisher, *installServiceCommandAddress, *installServiceCommandUser, *installServiceCommandPwd)
			if err != nil {
				log.Println("fail to install service:", err)
			}
		}

		if uninstallServiceCommand.Parsed() {
			if *uninstallServiceCommandService == "" {
				uninstallServiceCommand.PrintDefaults()
				os.Exit(1)
			}
			if *uninstallServiceCommandPublisher == "" {
				uninstallServiceCommand.PrintDefaults()
				os.Exit(1)
			}

			if *uninstallServiceCommandAddress == "" {
				uninstallServiceCommand.PrintDefaults()
				os.Exit(1)
			}
			if *uninstallServiceCommandUser == "" {
				uninstallServiceCommand.PrintDefaults()
				os.Exit(1)
			}
			if *uninstallServiceCommandPwd == "" {
				uninstallServiceCommand.PrintDefaults()
				os.Exit(1)
			}
			if *uninstallServiceCommandVersion == "" {
				uninstallServiceCommand.PrintDefaults()
				os.Exit(1)
			}
			err = uninstallService(*uninstallServiceCommandService, *uninstallServiceCommandPublisher, *uninstallServiceCommandVersion, *uninstallServiceCommandAddress, *uninstallServiceCommandUser, *uninstallServiceCommandPwd)
			if err != nil {
				log.Println("fail to uninstall service:", err)
			}
		}

		if updateGlobularCommand.Parsed() {
			if *updateGlobularCommandExecPath == "" {
				updateGlobularCommand.PrintDefaults()
				fmt.Println("no executable path was given!")
				os.Exit(1)
			}

			if *updateGlobularCommandAddress == "" {
				updateGlobularCommand.PrintDefaults()
				fmt.Println("no domain was given!")
				os.Exit(1)
			}
			if *updateGlobularCommandUser == "" {
				updateGlobularCommand.PrintDefaults()
				fmt.Println("no user was given!")
				os.Exit(1)
			}
			if *updateGlobularCommandPwd == "" {
				updateGlobularCommand.PrintDefaults()
				fmt.Println("no password was given!")
				os.Exit(1)
			}
			if *updateGlobularCommandPlatform == "" {
				*updateGlobularCommandPlatform = runtime.GOOS + ":" + runtime.GOARCH
			}

			err = updateGlobular(*updateGlobularCommandExecPath, *updateGlobularCommandAddress, *updateGlobularCommandUser, *updateGlobularCommandPwd, *updateGlobularCommandPlatform)
			if err != nil {
				log.Println("fail to update Globular:", err)
			}
		}

		if updateGlobularFromCommand.Parsed() {
			if *updateGlobularCommandFromSource == "" {
				updateGlobularFromCommand.PrintDefaults()
				fmt.Println("no address was given to update Globular from")
				os.Exit(1)
			}

			if *updateGlobularFromCommandDest == "" {
				updateGlobularFromCommand.PrintDefaults()
				fmt.Println("no domain was given!")
				os.Exit(1)
			}
			if *updateGlobularFromCommandUser == "" {
				updateGlobularFromCommand.PrintDefaults()
				fmt.Println("no user (for domain) was given!")
				os.Exit(1)
			}

			if *updateGlobularFromCommandPwd == "" {
				updateGlobularFromCommand.PrintDefaults()
				fmt.Println("no password (for domain) was given!")
				os.Exit(1)
			}

			if *updateGlobularFromCommandPlatform == "" {
				*updateGlobularFromCommandPlatform = runtime.GOOS + ":" + runtime.GOARCH
			}

			err = updateGlobularFrom(*updateGlobularCommandFromSource, *updateGlobularFromCommandDest, *updateGlobularFromCommandUser, *updateGlobularFromCommandPwd, *updateGlobularFromCommandPlatform)
			if err != nil {
				log.Println("fail to update Globular from:", err)
			}
		}

		if installApplicationCommand.Parsed() {
			if *installApplicationCommandName == "" {
				installApplicationCommand.PrintDefaults()
				fmt.Println("no application name was given!")
				os.Exit(1)
			}
			if *installApplicationCommandDiscovery == "" {
				installApplicationCommand.PrintDefaults()
				fmt.Println("no discovery address was given!")
				os.Exit(1)
			}
			if *installApplicationCommandPublisher == "" {
				installApplicationCommand.PrintDefaults()
				fmt.Println("no publisher was given!")
				os.Exit(1)
			}
			if *installApplicationCommandAddress == "" {
				installApplicationCommand.PrintDefaults()
				fmt.Println("no domain was given!")
				os.Exit(1)
			}
			if *installApplicationCommandUser == "" {
				installApplicationCommand.PrintDefaults()
				fmt.Println("no user was given!")
				os.Exit(1)
			}
			if *installApplicationCommandPwd == "" {
				installApplicationCommand.PrintDefaults()
				fmt.Println("no password was given!")
				os.Exit(1)
			}

			var setAsDefault bool
			if *installApplicationCommandIndex != "" {
				setAsDefault = *installApplicationCommandIndex == "true"
			}

			err = installApplication(*installApplicationCommandName, *installApplicationCommandDiscovery, *installApplicationCommandPublisher, *installApplicationCommandAddress, *installApplicationCommandUser, *installApplicationCommandPwd, setAsDefault)
			if err != nil {
				log.Println("fail to install application:", err)
			}
		}

		if uninstallApplicationCommand.Parsed() {
			if *uninstallApplicationCommandName == "" {
				uninstallApplicationCommand.PrintDefaults()
				fmt.Println("no application name was given!")
				os.Exit(1)
			}
			if *uninstallApplicationCommandPublisher == "" {
				uninstallApplicationCommand.PrintDefaults()
				fmt.Println("no publisher was given!")
				os.Exit(1)
			}

			if *uninstallApplicationCommandAddress == "" {
				uninstallApplicationCommand.PrintDefaults()
				fmt.Println("no domain was given!")
				os.Exit(1)
			}
			if *uninstallApplicationCommandUser == "" {
				uninstallApplicationCommand.PrintDefaults()
				fmt.Println("no user was given!")
				os.Exit(1)
			}
			if *uninstallApplicationCommandPwd == "" {
				uninstallApplicationCommand.PrintDefaults()
				fmt.Println("no password was given!")
				os.Exit(1)
			}
			if *uninstallApplicationCommandVersion == "" {
				uninstallApplicationCommand.PrintDefaults()
				fmt.Println("no version was given!")
				os.Exit(1)
			}
			err = uninstallApplication(*uninstallApplicationCommandName, *uninstallApplicationCommandPublisher, *uninstallApplicationCommandVersion, *uninstallApplicationCommandAddress, *uninstallApplicationCommandUser, *uninstallApplicationCommandPwd)
			if err != nil {
				log.Println("fail to uninstall application:", err)
			}
		}

		if installCertificatesCommand.Parsed() {
			// Required Flags
			if *installCertificatesCommandPath == "" {
				installCertificatesCommand.PrintDefaults()
				os.Exit(1)
			}

			if *installCertificatesCommandDomain == "" {
				installCertificatesCommand.PrintDefaults()
				os.Exit(1)
			}

			if *installCertificatesCommandPort == "" {
				installCertificatesCommand.PrintDefaults()
				os.Exit(1)
			}

			err = installCertificates(*installCertificatesCommandDomain, Utility.ToInt(*installCertificatesCommandPort), *installCertificatesCommandPath)
			if err != nil {
				log.Println("fail to install certificates:", err)
			}
		}

		if startCommand.Parsed() {
			if len(*startCommandDomain) > 0 {
				g.Domain = *startCommandDomain
			}
			if err := g.run(); err != nil {
				log.Println("Globular run error:", err)
				os.Exit(1)
			}
		}

		// Check if the command was parsed
		if installCommand.Parsed() {
			if *installCommandName != "" {
				svcConfig.DisplayName = *installCommandName
				s, _ = service.New(g, svcConfig)
			}
			// Required Flags
			err := s.Install()
			if err != nil {
				log.Println("fail to install service:", err)
				os.Exit(1) // exit the program.
			}

			log.Println("Globular service is now installed!")
			// Here I will keep the start time...
			// set path...
			err = setSystemPath()
			if err != nil {
				log.Println("fail to set system path with error", err)
				os.Exit(1) // exit the program.

			}

			os.Exit(0) // exit the program.

		}

		if uninstallCommand.Parsed() && s != nil {
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
			err = Utility.KillProcessByName("mongod")
			if err != nil {
				log.Println("fail to kill mongod process with error", err)
			}
			err = Utility.KillProcessByName("prometheus")
			if err != nil {
				log.Println("fail to kill prometheus process with error", err)
			}
			err = Utility.KillProcessByName("torrent")
			if err != nil {
				log.Println("fail to kill torrent process with error", err)
			}
			err = Utility.KillProcessByName("envoy")
			if err != nil {
				log.Println("fail to kill envoy process with error", err)
			}
			err = Utility.KillProcessByName("etcd")
			if err != nil {
				log.Println("fail to kill etcd process with error", err)
			}

			// reset environmement...
			err = resetSystemPath()
			if err != nil {
				log.Println("fail to reset system path with error", err)
			}

			err = resetRules()
			if err != nil {
				log.Println("fail to reset rules with error", err)
			}

			os.Exit(0) // exit the program.
		}

		if distCommand.Parsed() {
			// Required Flags
			if *distCommandPath == "" {
				fmt.Println("No path was given!")
				distCommand.PrintDefaults()
				os.Exit(1)
			}
			if *distCommandRevision == "" {
				fmt.Println("No revision number was given!")
				distCommand.PrintDefaults()
				os.Exit(1)
			}
			dist(g, *distCommandPath, *distCommandRevision)
		}

		if deployCommand.Parsed() {
			// Required Flags
			if *deployCommandPath == "" {
				fmt.Print("No application 'dist' path was given")
				deployCommand.PrintDefaults()
				os.Exit(1)
			}
			if *deployCommandName == "" {
				fmt.Print("No application 'name' was given")
				deployCommand.PrintDefaults()
				os.Exit(1)
			}

			if *deployCommandUser == "" {
				fmt.Print("You must authenticate yourself")
				deployCommand.PrintDefaults()
				os.Exit(1)
			}

			if *deployCommandPwd == "" {
				fmt.Print("You must specifie the user password.")
				deployCommand.PrintDefaults()
				os.Exit(1)
			}

			if *deployCommandAddress == "" {
				fmt.Print("You must sepcie the server address")
				deployCommand.PrintDefaults()
				os.Exit(1)
			}

			if *deployCommandOrganization != "" {
				if !strings.Contains(*deployCommandOrganization, "@") {
					fmt.Print("You must sepcie the organisation domain ex. organization@domain.com where the domain is the domain where the globule where the organization is define.")
					deployCommand.PrintDefaults()
					os.Exit(1)
				}
			}

			var setAsDefault bool
			if *deployCommandIndex != "" {
				setAsDefault = *deployCommandIndex == "true"
			}

			err = deploy(*deployCommandName, *deployCommandOrganization, *deployCommandPath, *deployCommandAddress, *deployCommandUser, *deployCommandPwd, setAsDefault)
			if err != nil {
				log.Println("fail to deploy application:", err)
			}
		}

		if publishCommand.Parsed() {

			if *publishCommandPath == "" {
				publishCommand.PrintDefaults()
				fmt.Println("No -path was given!")
				os.Exit(1)
			}

			if *publishCommandUser == "" {
				publishCommand.PrintDefaults()
				fmt.Println("No -u (user) was given!")
				os.Exit(1)
			}

			if *publishCommandPwd == "" {
				publishCommand.PrintDefaults()
				fmt.Println("No -p (password) was given!")
				os.Exit(1)
			}

			if *publishCommandAddress == "" {
				publishCommand.PrintDefaults()
				fmt.Println("No -a (address or domain) was given!")
				os.Exit(1)
			}

			// Here I will read the configuration file...

			if !Utility.Exists(*publishCommandPath) {
				fmt.Println("No configuration file was found at " + *publishCommandPath)
				os.Exit(1)
			}

			// Detect the platform if none was given...
			if *publishCommandPlatform == "" {
				*publishCommandPlatform = runtime.GOOS + "_" + runtime.GOARCH
			}

			fmt.Println("publish for platform ", *publishCommandPlatform)
			err = publish(*publishCommandUser, *publishCommandPwd, *publishCommandAddress, *publishCommandOrganization, *publishCommandPath, *publishCommandPlatform)
			if err != nil {
				log.Println("fail to publish application:", err)
			}
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
				logError := g.logger.Error(err)
				if logError != nil {
					fmt.Println("fail to log error: ", logError)
				}
			} else if err != nil {
				fmt.Println("fail to run Globular with error: ", err)
			}
		}
	}

}

/**
 * Service interface use to run as Windows Service or Linux deamon...
 */
func installCertificates(domain string, port int, path string) error {

	log.Println("Get certificates from ", domain, "...")

	adminClient, err := admin_client.NewAdminService_Client(domain, "admin.AdminService")
	if err != nil {
		log.Println("fail to get certificates...", err)
		return err
	}

	key, cert, ca, err := adminClient.GetCertificates(domain, port, path)
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
func deploy(name string, organization string, path string, address string, user string, pwd string, setAsDefault bool) error {

	log.Println("try to deploy application", name, " to address ", address, " with user ", user)

	// Authenticate the user in order to get the token
	authenticationClient, err := authentication_client.NewAuthenticationService_Client(address, "authentication.AuthenticationService")
	if err != nil {
		log.Println("fail to access resource service at "+address+" with error ", err)
		return err
	}

	token, err := authenticationClient.Authenticate(user, pwd)
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
	// #nosec G304 -- File path is constructed from validated input and constant strings.
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
		rolesData := packageConfig["roles"].([]interface{})
		for i := 0; i < len(rolesData); i++ {
			roleMap := rolesData[i].(map[string]interface{})
			role := new(resourcepb.Role)
			role.Id = roleMap["id"].(string)
			role.Name = roleMap["name"].(string)
			role.Domain = roleMap["domain"].(string)
			role.Description = roleMap["description"].(string)
			role.Actions = make([]string, 0)
			for j := 0; j < len(roleMap["actions"].([]interface{})); j++ {
				role.Actions = append(role.Actions, roleMap["actions"].([]interface{})[j].(string))
			}
			roles = append(roles, role)
		}
	}

	// Create groups.
	groups := make([]*resourcepb.Group, 0)
	if packageConfig["groups"] != nil {
		groupsData := packageConfig["groups"].([]interface{})
		for i := 0; i < len(groupsData); i++ {
			groupMap := groupsData[i].(map[string]interface{})
			group := new(resourcepb.Group)
			group.Id = groupMap["id"].(string)
			group.Name = groupMap["name"].(string)
			group.Description = groupMap["description"].(string)
			group.Domain = groupMap["domain"].(string)
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

				defer func() {
					err := os.Remove(pngPath)
					if err != nil {
						log.Println("fail to remove temporary png file with error", err)
					}
				}()

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
	discoveryClient, err := discovery_client.NewDiscoveryService_Client(address, "discovery.PackageDiscovery")
	if err != nil {
		fmt.Println("fail to connecto to discovery service at address", address, "with error", err)
		return err
	}

	err = discoveryClient.PublishApplication(token, user, organization, "/"+name, name, address, version, description, icon, alias, address, address, actions, keywords, roles, groups)
	if err != nil {
		fmt.Println("fail to publish the application with error:", err)
		return err
	}

	repositoryClient, err := repository_client.NewRepositoryService_Client(address, "repository.PackageRepository")
	if err != nil {
		fmt.Println("fail to connecto to repository service at address", address, "with error", err)
		return err
	}

	log.Println("Upload application package at path ", path)
	_, err = repositoryClient.UploadApplicationPackage(user, organization, path, token, address, name, version)
	if err != nil {
		fmt.Println("fail to upload the application package with error:", err)
		return err
	}

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	log.Println("Connect with application manager at address ", address)
	applicationsManagerClient, err := applications_manager_client.NewApplicationsManager_Client(address, "applications_manager.ApplicationManagerService") // create the resource server.
	if err != nil {
		fmt.Println("fail to connecto to application manager service at address", address, "with error", err)
		return err
	}

	// set the publisher id.
	publisherID := user
	if len(organization) > 0 {
		publisherID = organization
	}

	// if no domain was given i will user the local domain.
	if !strings.Contains(publisherID, "@") {
		domain, err := configpkg.GetDomain()
		if err != nil {
			fmt.Println("fail to get domain with error", err)
			return err
		}
		publisherID += "@" + domain
	}

	fmt.Println("try to install the newly deployed application...")
	err = applicationsManagerClient.InstallApplication(token, authenticationClient.GetDomain(), user, address, publisherID, name, setAsDefault)
	if err != nil {
		log.Println("fail to install application with error ", err)
		return err
	}

	log.Println("Application was deployed and installed successfully!")
	return nil

}

/**
 * Push globular update on a given server.
 * ex.
 * sudo ./Globular update -path=/home/dave/go/src/github.com/globulario/Globular/Globular -a=globular.cloud -u=sa -p=adminadmin
 */
func updateGlobular(path, address, user, pwd string, platform string) error {

	// Authenticate the user in order to get the token
	authenticationClient, err := authentication_client.NewAuthenticationService_Client(address, "authentication.AuthenticationService")
	if err != nil {
		log.Println("fail to create authentication client:", err)
		return err
	}

	token, err := authenticationClient.Authenticate(user, pwd)
	if err != nil {
		log.Println("fail to authenticate user:", err)
		return err
	}

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	adminClient, err := admin_client.NewAdminService_Client(address, "admin.AdminService")
	if err != nil {
		log.Println("fail to create admin client:", err)
		return err
	}

	_, err = adminClient.Update(path, platform, token, address)
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
func updateGlobularFrom(src, dest, user, pwd string, platform string) error {
	fmt.Println("pull globular update from ", src, " to ", dest)

	adminSource, err := admin_client.NewAdminService_Client(src, "admin.AdminService")
	if err != nil {
		return err
	}

	// From the source I will download the new executable and save it in the
	// temp directory...
	tempPath := os.TempDir() + "/Globular_" + Utility.ToString(time.Now().Unix())
	err = Utility.CreateDirIfNotExist(tempPath)
	if err != nil {
		log.Println("fail to create temp directory with error: ", err)
		return err
	}

	fmt.Println("download exec to ", tempPath)
	err = adminSource.DownloadGlobular(src, platform, tempPath)
	if err != nil {
		log.Println("fail to download new globular executable file with error: ", err)
		return err
	}

	pathWithExec := tempPath
	pathWithExec += "/Globular"
	if runtime.GOOS == "windows" {
		pathWithExec += ".exe"
	}

	if !Utility.Exists(pathWithExec) {
		err := errors.New(pathWithExec + " not found! ")
		return err
	}

	defer func() {
		err := os.RemoveAll(tempPath)
		if err != nil {
			log.Println("fail to remove temporary directory with error", err)
		}
	}()

	err = updateGlobular(pathWithExec, dest, user, pwd, platform)
	if err != nil {
		log.Println("fail to update Globular from:", err)
		return err
	}

	// Send the command.
	return nil

}

// publish authenticates a user and publishes a service package to the specified platform.
// It performs the following steps:
//  1. Authenticates the user to obtain a token.
//  2. Publishes the service using the discovery client.
//  3. Uploads the service package using the repository client.
//
// Parameters:
//
//	user         - The username for authentication.
//	pwd          - The password for authentication.
//	address      - The service address.
//	organization - The organization name.
//	path         - The path to the service package.
//	platform     - The target platform.
//
// Returns:
//
//	error - An error if any step fails, otherwise nil.
func publish(user, pwd, address, organization, path, platform string) error {

	// Authenticate the user in order to get the token
	authenticationClient, err := authentication_client.NewAuthenticationService_Client(address, "authentication.AuthenticationService")
	if err != nil {
		log.Println("fail to create authentication client:", err)
		return err
	}

	token, err := authenticationClient.Authenticate(user, pwd)
	if err != nil {
		log.Println("fail to authenticate user:", err)
		return err
	}

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	discoveryClient, err := discovery_client.NewDiscoveryService_Client(address, "discovery.PackageDiscovery")
	if err != nil {
		log.Println("fail to create discovery client:", err)
		return err
	}

	err = discoveryClient.PublishService(user, organization, token, discoveryClient.GetDomain(), path, platform)
	if err != nil {
		return err
	}

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	repositoryClient, err := repository_client.NewRepositoryService_Client(address, "repository.PackageRepository")
	if err != nil {
		log.Println("fail to create repository client:", err)
		return err
	}

	// first of all I will create and upload the package on the discovery...
	err = repositoryClient.UploadServicePackage(user, organization, token, address, path, platform)
	if err != nil {
		log.Println("fail to upload service package:", err)
		return err
	}

	log.Println("Service was pulbish successfully have a nice day folk's!")
	return nil
}

/**
 * That function is use to intall a service at given addresse. The service Id is the unique identifier of
 * the service to be install.
 */
func installService(serviceID, discovery, publisherID, domain, user, pwd string) error {
	log.Println("try to install service", serviceID, "from", publisherID, "on", domain)
	// Authenticate the user in order to get the token
	authenticationClient, err := authentication_client.NewAuthenticationService_Client(domain, "authentication.AuthenticationService")
	if err != nil {
		log.Println("fail to create authentication client:", err)
		return err
	}

	token, err := authenticationClient.Authenticate(user, pwd)
	if err != nil {
		log.Println("fail to authenticate with error ", err.Error())
		return err
	}

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	servicesManagerCli, err := serviceManagerClient.NewServicesManagerService_Client(domain, "services_manager.ServicesManagerService")
	if err != nil {
		log.Println("fail to connect to services manager at ", domain, " with error ", err.Error())
		return err
	}

	// first of all I will create and upload the package on the discovery...
	err = servicesManagerCli.InstallService(token, domain, discovery, publisherID, serviceID)
	if err != nil {
		log.Println("fail to install service", serviceID, "with error ", err.Error())
		return err
	}

	log.Println("service was installed")

	return nil
}

func uninstallService(serviceID, publisherID, version, address, user, pwd string) error {

	// Authenticate the user in order to get the token
	authenticationClient, err := authentication_client.NewAuthenticationService_Client(address, "authentication.AuthenticationService")
	if err != nil {
		log.Println("fail to create authentication client:", err)
		return err
	}

	token, err := authenticationClient.Authenticate(user, pwd)
	if err != nil {
		log.Println("fail to authenticate user:", err)
		return err
	}

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	servicesManagerCli, err := serviceManagerClient.NewServicesManagerService_Client(address, "services_manager.ServicesManagerService")
	if err != nil {
		log.Println("fail to create services manager client:", err)
		return err
	}

	// first of all I will create and upload the package on the discovery...
	err = servicesManagerCli.UninstallService(token, address, publisherID, serviceID, version)
	if err != nil {
		log.Println("fail to uninstall service:", err)
		return err
	}

	return nil
}

/**
 * Install Globular web application.
 */
func installApplication(applicationID, discovery, publisherID, domain, user, pwd string, setAsDefault bool) error {

	// Authenticate the user in order to get the token
	authenticationClient, err := authentication_client.NewAuthenticationService_Client(domain, "authentication.AuthenticationService")
	if err != nil {
		log.Println("fail to create authentication client:", err)
		return err
	}

	token, err := authenticationClient.Authenticate(user, pwd)
	if err != nil {
		return err
	}

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	applicationsManagerClient, err := applications_manager_client.NewApplicationsManager_Client(domain, "applications_manager.ApplicationManagerService")
	if err != nil {
		log.Println("fail to create applications manager client:", err)
		return err
	}

	// first of all I will create and upload the package on the discovery...
	err = applicationsManagerClient.InstallApplication(token, domain, user, discovery, publisherID, applicationID, setAsDefault)
	if err != nil {
		log.Println("fail to install application:", err)
		return err
	}

	return nil
}

func uninstallApplication(applicationID, publisherID, version, domain, user, pwd string) error {

	// Authenticate the user in order to get the token
	authenticationClient, err := authentication_client.NewAuthenticationService_Client(domain, "authentication.AuthenticationService")
	if err != nil {
		log.Println("fail to create authentication client:", err)
		return err
	}

	token, err := authenticationClient.Authenticate(user, pwd)
	if err != nil {
		log.Println("fail to authenticate user:", err)
		return err
	}

	// first of all I need to get all credential informations...
	// The certificates will be taken from the address
	applicationsManagerClient, err := applications_manager_client.NewApplicationsManager_Client(domain, "applications_manager.ApplicationManagerService")
	if err != nil {
		log.Println("fail to create applications manager client:", err)
		return err
	}

	// first of all I will create and upload the package on the discovery...
	err = applicationsManagerClient.UninstallApplication(token, domain, user, publisherID, applicationID, version)
	if err != nil {
		log.Println("fail to uninstall application:", err)
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
 * Dependency script are describe at
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
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Println("fail to get current directory with error", err)
		return
	}

	switch runtime.GOOS {
	case "darwin":

		darwinPackagePath := path + "/globular_" + g.Version + "-" + revision + "_" + runtime.GOARCH

		// remove existing files...
		err := os.RemoveAll(darwinPackagePath)
		if err != nil {
			log.Println("fail to remove existing darwin package directory with error", err)
		}

		// 1. Create the working directory
		err = Utility.CreateDirIfNotExist(darwinPackagePath)
		if err != nil {
			fmt.Println("fail to create darwin package directory with error: ", err)
			return
		}

		// 2. Create application directory.
		appPath := darwinPackagePath + "/globular.cloud"
		appContent := appPath + "/Contents"
		appBin := appContent + "/MacOS"
		appResource := appContent + "/Resources"
		configPath := appBin + "/etc/globular/config"
		applicationsPath := appBin + "/var/globular/applications"

		// Copy applications for offline installation...
		err = Utility.CopyDir(dir+"/applications/.", applicationsPath)
		if err != nil {
			fmt.Println("fail to copy applications for offline installation with error: ", err)
			return
		}

		// create directories...
		err = Utility.CreateDirIfNotExist(appContent)
		if err != nil {
			fmt.Println("fail to create application content directory with error: ", err)
			return
		}

		err = Utility.CreateDirIfNotExist(appBin)
		if err != nil {
			fmt.Println("fail to create application binary directory with error: ", err)
			return
		}

		err = Utility.CreateDirIfNotExist(configPath)
		if err != nil {
			fmt.Println("fail to create application config directory with error: ", err)
			return
		}

		err = Utility.CreateDirIfNotExist(appResource)
		if err != nil {
			fmt.Println("fail to create application resource directory with error: ", err)
			return
		}

		err = Utility.CopyFile(dir+"/assets/icon.icns", appResource+"/icon.icns")

		if err != nil {
			fmt.Println("fail to copy icon from ", dir+"/assets/icon.icns"+"whit error", err)
		}

		// Create the distribution.
		generateDistro(appBin)

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
		err = os.WriteFile(appContent+"/Info.plist", []byte(plistFile), 0600)
		if err != nil {
			fmt.Println("fail to create Info.plist with error", err)
			return
		}

	case "linux":
		debianPackagePath := path + "/globular_" + g.Version + "-" + revision + "_" + runtime.GOARCH

		// remove existiong files...
		err = os.RemoveAll(debianPackagePath)
		if err != nil {
			log.Println("fail to remove existing debian package directory with error", err)
		}

		// 1. Create the working directory
		err = Utility.CreateDirIfNotExist(debianPackagePath)
		if err != nil {
			fmt.Println("fail to create debian package directory with error: ", err)
			return
		}

		// 2. Create the internal structure

		// globular exec and other services exec
		distroPath := debianPackagePath + "/usr/local/share/globular"

		// globular data
		dataPath := debianPackagePath + "/var/globular/data"

		// globular data
		applicationsPath := debianPackagePath + "/var/globular/applications"

		// globular configurations
		configPath := debianPackagePath + "/etc/globular/config"

		// Create the bin directories.
		err = Utility.CreateDirIfNotExist(distroPath)
		if err != nil {
			fmt.Println("fail to create distribution directory with error: ", err)
			return
		}

		err = Utility.CreateDirIfNotExist(dataPath)
		if err != nil {
			fmt.Println("fail to create data directory with error: ", err)
			return
		}

		err = Utility.CreateDirIfNotExist(configPath)
		if err != nil {
			fmt.Println("fail to create config directory with error: ", err)
			return
		}

		err = Utility.CreateDirIfNotExist(applicationsPath)
		if err != nil {
			fmt.Println("fail to create applications directory with error: ", err)
			return
		}

		// Now the libraries...
		libpath := debianPackagePath + "/usr/local/lib"
		err = Utility.CreateDirIfNotExist(libpath)
		if err != nil {
			fmt.Println("fail to create lib directory with error: ", err)
			return
		}

		binpath := debianPackagePath + "/usr/local/bin"
		err = Utility.CreateDirIfNotExist(binpath)
		if err != nil {
			fmt.Println("fail to create bin directory with error: ", err)
			return
		}

		if runtime.GOARCH == "amd64" {

			// zlib
			if !Utility.Exists("/usr/lib/x86_64-linux-gnu/libz.a") {
				fmt.Println("libz.a not found please install it on your computer: sudo apt-get install zlib1g-dev")
				return
			}

			err := Utility.CopyFile("/usr/lib/x86_64-linux-gnu/libz.a", libpath+"/libz.a")
			if err != nil {
				log.Println("fail to copy libz.a with error", err)
			}

			err = Utility.CopyFile("/usr/lib/x86_64-linux-gnu/libz.so.1.2.11", libpath+"/libz.so.1.2.11")
			if err != nil {
				log.Println("fail to copy libz.so.1.2.11 with error", err)
			}

		} else if runtime.GOARCH == "arm64" {
			if !Utility.Exists("/usr/lib/aarch64-linux-gnu/libssl.so.1.1") {
				fmt.Println("libssl.so.1.1 not found on your computer, please install it: ")
				fmt.Println("   wget http://launchpadlibrarian.net/475575244/libssl1.1_1.1.1f-1ubuntu2_arm64.deb")
				fmt.Println("	sudo dpkg -i libssl1.1_1.1.1f-1ubuntu2_arm64.deb ")
				return
			}

			// Copy lib crypto...
			err := Utility.CopyFile("/usr/lib/aarch64-linux-gnu/libssl.so.1.1", libpath+"/libssl.so.1.1")
			if err != nil {
				log.Println("fail to copy libssl.so.1.1 with error", err)
			}
			err = Utility.CopyFile("/usr/lib/aarch64-linux-gnu/libcrypto.so.1.1", libpath+"/libcrypto.so.1.1")
			if err != nil {
				log.Println("fail to copy libcrypto.so.1.1 with error", err)
			}

			// zlib
			if !Utility.Exists("/usr/lib/aarch64-linux-gnu/libz.a") {
				fmt.Println("libz.a not found please install it on your computer: sudo apt-get install zlib1g-dev")
				return
			}

			err = Utility.CopyFile("/usr/lib/aarch64-linux-gnu/libz.a", libpath+"/libz.a")
			if err != nil {
				log.Println("fail to copy libz.a with error", err)
			}
			err = Utility.CopyFile("/usr/lib/aarch64-linux-gnu/libz.so.1.2.13", libpath+"/libz.so.1.2.13")
			if err != nil {
				log.Println("fail to copy libz.so.1.2.13 with error", err)
			}
		}

		// ODBC libraries...
		err = Utility.CopyFile("/usr/local/lib/libodbc.la", libpath+"/libodbc.la")
		if err != nil {
			log.Println("fail to copy libodbc.la with error", err)
		}
		err = Utility.CopyFile("/usr/local/lib/libodbc.so.2.0.0", libpath+"/libodbc.so.2.0.0")
		if err != nil {
			log.Println("fail to copy libodbc.so.2.0.0 with error", err)
		}

		err = Utility.CopyFile("/usr/local/lib/libodbccr.la", libpath+"/libodbccr.la")
		if err != nil {
			log.Println("fail to copy libodbccr.la with error", err)
		}
		err = Utility.CopyFile("/usr/local/lib/libodbccr.so.2.0.0", libpath+"/libodbccr.so.2.0.0")
		if err != nil {
			log.Println("fail to copy libodbccr.so.2.0.0 with error", err)
		}

		err = Utility.CopyFile("/usr/local/lib/libodbcinst.la", libpath+"/libodbcinst.la")
		if err != nil {
			log.Println("fail to copy libodbcinst.la with error", err)
		}

		err = Utility.CopyFile("/usr/local/lib/libodbcinst.so.2.0.0", libpath+"/libodbcinst.so.2.0.0")
		if err != nil {
			log.Println("fail to copy libodbcinst.so.2.0.0 with error", err)
		}

		err = Utility.CopyFile("/usr/local/bin/grpcwebproxy", binpath+"/grpcwebproxy")
		if err != nil {
			log.Println("fail to copy grpcwebproxy with error", err)
		}

		// Now I will create get the configuration files from service and create a copy to /etc/globular/config/services
		// so modification will survice upgrades.

		// Create the distribution.
		configurations := generateDistro(distroPath)

		// 3. Create the control file
		err = Utility.CreateDirIfNotExist(debianPackagePath + "/DEBIAN")
		if err != nil {
			fmt.Println("fail to create DEBIAN directory with error: ", err)
			return
		}

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
		packageConfig += "Description: Globular is a complete web application development suite. Globular is based on microservices architecture and implemented with help of gRPC.\n"

		// - The project homepage
		packageConfig += "Homepage: https://globular.io\n"

		// - The list of Dependencys...
		packageConfig += "Depends: python3 (>= 3.8.~), python-is-python3 (>=3.8.~), ffmpeg (>=4.4.~), curl(>=7.8.~), dpkg(>=1.21.~), nmap, arp-scan\n"

		err = os.WriteFile(debianPackagePath+"/DEBIAN/control", []byte(packageConfig), 0600)
		if err != nil {
			log.Println("fail to create debian package control file:", err)
		}

		// Here I will set the script to run before the installation...
		// https://www.devdungeon.com/content/debian-package-tutorial-dpkgdeb#toc-17
		// TODO create tow version one for arm7 and one for amd64

		var preinst string

		switch runtime.GOARCH {
		case "amd64":
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
		case "arm64":
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

		err = os.WriteFile(debianPackagePath+"/DEBIAN/preinst", []byte(preinst), 0600)
		if err != nil {
			log.Println("fail to create debian package control file:", err)
		}

		// the list of list of configurations
		// DEBIAN/conffiles
		conffiles := ""
		for i := range configurations {
			conffiles += configurations[i] + "\n"
		}

		err = os.WriteFile(debianPackagePath+"/DEBIAN/conffiles", []byte(conffiles), 0600)
		if err != nil {
			log.Println("fail to create debian package control file:", err)
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
		err = os.WriteFile(debianPackagePath+"/DEBIAN/postinst", []byte(postinst), 0600)
		if err != nil {
			log.Println("fail to create debian package control file:", err)
		}

		prerm := `
		# Thing to do before removing

		if [ -f "/etc/systemd/system/Globular.service" ]; then
			# Stop, Disable and Uninstall Globular service.
			echo "Stop running globular service..."
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
		err = os.WriteFile(debianPackagePath+"/DEBIAN/prerm", []byte(prerm), 0600)
		if err != nil {
			log.Println("fail to create debian package control file:", err)
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
		err = os.WriteFile(debianPackagePath+"/DEBIAN/postrm", []byte(postrm), 0600)
		if err != nil {
			log.Println("fail to create debian package control file:", err)
		}

		// 5. Build the deb package
		fmt.Println("Build the debian package at ", debianPackagePath)
		// #nosec G204 -- Subprocess launched with variable
		cmd := exec.Command("dpkg-deb", "--build", "--root-owner-group", debianPackagePath)

		cmdOutput := &bytes.Buffer{}
		cmd.Stdout = cmdOutput

		err = cmd.Run()
		if err != nil {
			log.Println("fail to build debian package:", err)
		}
		fmt.Print(cmdOutput.String())

	case "windows":
		fmt.Println("Create the distro at path ", path)
		root := path + "/Globular"
		err = Utility.CreateIfNotExists(root, 0600)
		if err != nil {
			log.Println("fail to create root directory with error", err)
		}

		app := root + "/app"
		err = Utility.CreateIfNotExists(app, 0600)
		if err != nil {
			log.Println("fail to create app directory with error", err)
		}

		// Copy globular ditro file.
		generateDistro(app)

		// I will make use of NSIS: Nullsoft Scriptable Install System to create an installer for window.
		// so here I will create the setup.nsi file.

		// copy assets
		err := Utility.CreateDirIfNotExist(root + "/assets")
		if err != nil {
			fmt.Println("fail to create assets directory with error: ", err)
			return
		}

		dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

		err = Utility.CopyDir(dir+"/assets/.", root+"/assets")
		if err != nil {
			fmt.Println("--> fail to copy assets ", err)
		}

		// redist
		err = Utility.CreateDirIfNotExist(root + "/redist")
		if err != nil {
			fmt.Println("fail to create redist directory with error: ", err)
			return
		}

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
			return
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
	  RMDir /r "$INSTDIR\Dependencys"
	  RMDir /r "$INSTDIR\Redist"
	  RMDir /r "$INSTDIR\services"
	  RMDir /r "$INSTDIR\applications"
  
	  DeleteRegKey /ifempty HKCU  "Software\${NAME}"
  
	SectionEnd
  
`
		err = Utility.WriteStringToFile(root+"/setup.nsi", setupNsi)
		if err != nil {
			log.Println("fail to write setup.nsi file:", err)
			return
		}
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

func generateDistro(path string) []string {

	// Return the configurations list
	configs := make([]string, 0)

	// I will set the docker file depending of the arch.
	var dockerfile string
	if runtime.GOARCH == "amd64" {
		data, err := os.ReadFile("Dockerfile_amd64")
		if err != nil {
			log.Println("fail to read Dockerfile_amd64:", err)
			os.Exit(1)
		}
		dockerfile = string(data)
	} else if runtime.GOARCH == "arm64" {
		data, err := os.ReadFile("Dockerfile_arm64")
		if err != nil {
			log.Println("fail to read Dockerfile_arm64:", err)
			os.Exit(1)
		}
		dockerfile = string(data)
	}

	err := Utility.CreateDirIfNotExist(path)
	if err != nil {
		fmt.Println("fail to create directory with error: ", err)
		os.Exit(1)
	}

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

	err = Utility.Copy(globularExec, path+"/"+destExec)
	if err != nil {
		log.Println("fail to copy globular executable:", err)
	}

	err = os.Chmod(path+"/"+destExec, 0600)
	if err != nil {
		log.Println("fail to change mode of globular executable:", err)
	}

	// Copy the bin file from globular
	switch runtime.GOOS {
	case "windows":
		err = Utility.CreateDirIfNotExist(path + "/Dependencys")
		if err != nil {
			fmt.Println("fail to create Dependencys directory with error: ", err)
			os.Exit(1)
		}

		if Utility.Exists(dir + "/Dependencys") {

			err = Utility.CopyDir(dir+"/Dependencys/.", path+"/Dependencys")
			if err != nil {
				fmt.Println("fail to copy Dependencys ", err)
			}

			execs := Utility.GetFilePathsByExtension(path+"/Dependencys", ".exe")
			for i := range execs {
				err = os.Chmod(execs[i], 0600)
				if err != nil {
					log.Println("fail to change mode of dependency:", err)
				}
			}
		} else {
			fmt.Println("no dir with Dependency was found at path", dir+"/Dependencys")
		}
	case "darwin":
		dest := path + "/bin"

		// ffmpeg
		if Utility.Exists("/usr/local/bin/ffmpeg") {
			err = Utility.CopyFile("/usr/local/bin/ffmpeg", dest)
			if err != nil {
				log.Println("fail to copy ffmpeg with error", err)
			}
		}

		if Utility.Exists("/usr/local/bin/ffprobe") {
			err = Utility.CopyFile("/usr/local/bin/ffprobe", dest)
			if err != nil {
				log.Println("fail to copy ffprobe with error", err)
			}
		}

		// yt-dlp
		if Utility.Exists("/usr/local/bin/yt-dlp") {
			err = Utility.CopyFile("/usr/local/bin/yt-dlp", dest)
			if err != nil {
				log.Println("fail to copy yt-dlp with error", err)
			}
		}

		// unix-odbc
		if Utility.Exists("/usr/local/bin/odbc_config") {
			err = Utility.CopyFile("/usr/local/bin/odbc_config", dest)
			if err != nil {
				log.Println("fail to copy odbc_config with error", err)
			}
		}

		if Utility.Exists("/usr/local/bin/odbcinst") {
			err = Utility.CopyFile("/usr/local/bin/odbcinst", dest)
			if err != nil {
				log.Println("fail to copy odbcinst with error", err)
			}
		}

		// unix-odbc
		if Utility.Exists("/usr/local/bin/odbc_config") {
			err = Utility.CopyFile("/usr/local/bin/odbc_config", dest)
			if err != nil {
				log.Println("fail to copy odbc_config with error", err)
			}
		}

		if Utility.Exists("/usr/local/bin/odbcinst") {
			err = Utility.CopyFile("/usr/local/bin/odbcinst", dest)
			if err != nil {
				log.Println("fail to copy odbcinst with error", err)
			}
		}

		// Now the libraries...
		dest = path + "/lib"
		err = Utility.CreateDirIfNotExist(dest)
		if err != nil {
			log.Println("fail to create lib directory with error: ", err)
			os.Exit(1)
		}

		if Utility.Exists("/usr/local/lib/libodbc.2.dylib") {
			err = Utility.CopyFile("/usr/local/lib/libodbc.2.dylib", dest)
			if err != nil {
				log.Println("fail to copy libodbc.2.dylib with error", err)
			}
		}

		if Utility.Exists("/usr/local/lib/libodbc.dylib") {
			err = Utility.CopyFile("/usr/local/lib/libodbc.dylib", dest)
			if err != nil {
				log.Println("fail to copy libodbc.dylib with error", err)
			}
		}

		if Utility.Exists("/usr/local/lib/libodbc.la") {
			err = Utility.CopyFile("/usr/local/lib/libodbc.la", dest)
			if err != nil {
				log.Println("fail to copy libodbc.la with error", err)
			}
		}

		if Utility.Exists("/usr/local/lib/libodbccr.2.dylib") {
			err = Utility.CopyFile("/usr/local/lib/libodbccr.2.dylib", dest)
			if err != nil {
				log.Println("fail to copy libodbccr.2.dylib with error", err)
			}
		}

		if Utility.Exists("/usr/local/lib/libodbccr.dylib") {
			err = Utility.CopyFile("/usr/local/lib/libodbccr.dylib", dest)
			if err != nil {
				log.Println("fail to copy libodbccr.dylib with error", err)
			}
		}

		if Utility.Exists("/usr/local/lib/libodbccr.la") {
			err = Utility.CopyFile("/usr/local/lib/libodbccr.la", dest)
			if err != nil {
				log.Println("fail to copy libodbccr.la with error", err)
			}
		}

		if Utility.Exists("/usr/local/lib/libodbcinst.2.dylib") {
			err = Utility.CopyFile("/usr/local/lib/libodbcinst.2.dylib", dest)
			if err != nil {
				log.Println("fail to copy libodbcinst.2.dylib with error", err)
			}
		}

		if Utility.Exists("/usr/local/lib/libodbcinst.dylib") {
			err = Utility.CopyFile("/usr/local/lib/libodbcinst.dylib", dest)
			if err != nil {
				log.Println("fail to copy libodbcinst.dylib with error", err)
			}
		}

		if Utility.Exists("/usr/local/lib/libodbcinst.la") {
			err = Utility.CopyFile("/usr/local/lib/libodbcinst.la", dest)
			if err != nil {
				log.Println("fail to copy libodbcinst.la with error", err)
			}
		}
	}

	// install services...
	services, err := configpkg.GetServicesConfigurations()
	if err != nil {
		log.Println("fail to find services with error ", err)
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

	for i := range services {

		// set the service configuration...
		s := services[i]
		id := s["Id"].(string)
		name := s["Name"].(string)

		// I will read the configuration file to have necessary service information
		// to be able to create the path.
		hasPath := s["Path"] != nil
		if hasPath {
			execPath := s["Path"].(string)
			if len(execPath) > 0 {
				configPath := filepath.Dir(execPath) + "/config.json"
				if Utility.Exists(configPath) {
					log.Println("install service ", name)
					// #nosec G304 -- Value is from a trusted source
					bytes, err := os.ReadFile(configPath)
					if err != nil {
						log.Println("fail to read config file with error", err)
					}

					config := make(map[string]interface{}, 0)
					err = json.Unmarshal(bytes, &config)
					if err != nil {
						log.Println("fail to unmarshal config file with error", err)
					}

					hasProto := s["Proto"] != nil

					// set the name.
					if config["PublisherID"] != nil && config["Version"] != nil && hasProto {
						protoPath := s["Proto"].(string)

						if Utility.Exists(execPath) && Utility.Exists(protoPath) {
							var serviceDir = "services/"
							if len(config["PublisherID"].(string)) == 0 {
								serviceDir += config["Domain"].(string) + "/" + name + "/" + config["Version"].(string)
							} else {
								serviceDir += config["PublisherID"].(string) + "/" + name + "/" + config["Version"].(string)
							}

							execName := filepath.Base(execPath)
							destPath := path + "/" + serviceDir + "/" + id + "/" + execName

							if Utility.Exists(execPath) {
								err := Utility.CreateDirIfNotExist(path + "/" + serviceDir + "/" + id)
								if err != nil {
									fmt.Println("fail to create service directory with error: ", err)
									os.Exit(1)
								}

								err = Utility.Copy(execPath, destPath)
								if err != nil {
									fmt.Println(execPath, destPath, err)
								}

								// Set execute permission
								err = os.Chmod(destPath, 0600)
								if err != nil {
									fmt.Println("fail to change mode of service executable:", err)
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
										config["Root"] = configpkg.GetDataDir() + "/files"

										// I will also copy the mime type directory
										config["Public"] = make([]string, 0)
										err = Utility.CopyDir(filepath.Dir(execPath)+"/mimetypes", path+"/"+serviceDir+"/"+id)
										if err != nil {
											fmt.Println("fail to copy mime types directory with error: ", err)
											os.Exit(1)
										}

									} else if name == "conversation.ConversationService" {
										config["Root"] = configpkg.GetDataDir()
									}
								}

								// Empty the list of connections if connections exist for the services.
								if config["Connections"] != nil {
									config["Connections"] = make(map[string]interface{})
								}

								config["ConfigPath"] = configpkg.GetConfigDir() + "/" + serviceDir + "/" + id + "/config.json"
								configs = append(configs, "/etc/globular/config/"+serviceDir+"/"+id+"/config.json")
								str, _ := Utility.ToJson(&config)

								if len(configPath) > 0 {
									// So here I will set the service configuration at /etc/globular/config/service... to be sure configuration
									// will survive package upgrades...
									err = Utility.CreateDirIfNotExist(configPath + "/" + serviceDir + "/" + id)
									if err != nil {
										fmt.Println("fail to create config directory with error: ", err)
										os.Exit(1)
									}

									err = os.WriteFile(configPath+"/"+serviceDir+"/"+id+"/config.json", []byte(str), 0600)
									if err != nil {
										fmt.Println("fail to create config file with error: ", err)
										os.Exit(1)
									}

								}

								// Copy the proto file.
								if Utility.Exists(protoPath) {
									err = Utility.Copy(protoPath, path+"/"+serviceDir+"/"+name+".proto")
									if err != nil {
										fmt.Println("fail to copy proto file with error: ", err)
										os.Exit(1)
									}
								}
							} else {
								fmt.Println("executable not exist ", execPath)
							}
						} else if !Utility.Exists(execPath) {
							log.Println("no executable found at path " + execPath)
						} else if !Utility.Exists(protoPath) {
							log.Println("no proto file found at path " + protoPath)
						}
					} else if config["PublisherID"] == nil {
						fmt.Println("no publisher was define!")
					} else if config["Version"] == nil {
						fmt.Println("no version was define!")
					} else if !hasProto {
						fmt.Println(" no proto file was found!")
					} else if !hasPath {
						fmt.Println("no executable was found!")
					}

				} else {
					fmt.Println("service", name, ":", id, "configuration is incomplete!")
				}
			} else {
				// Internal services here.
				protoPath := s["Proto"].(string)

				// Copy the proto file.
				if Utility.Exists(os.Getenv("ServicesRoot") + "/" + protoPath) {
					err = Utility.Copy(os.Getenv("ServicesRoot")+"/"+protoPath, path+"/"+protoPath)
					if err != nil {
						fmt.Println("fail to copy proto file with error: ", err)
						os.Exit(1)
					}
				}
			}
		}

	}

	// save docker.
	err = os.WriteFile(path+"/Dockerfile", []byte(dockerfile), 0600)
	if err != nil {
		log.Println("fail to write Dockerfile:", err)
	}

	//
	return configs
}

/**
 * Generate a token that will be valid for 15 minutes or the session timeout delay.
 */
func generateToken(address, user, pwd string) error {

	// Authenticate the user in order to get the token
	authenticationClient, err := authentication_client.NewAuthenticationService_Client(address, "authentication.AuthenticationService")
	if err != nil {
		return err
	}

	// Get the remote token
	token, err := authenticationClient.Authenticate(user, pwd)
	if err != nil {
		return err
	}

	_, _ = fmt.Println(token)

	return nil
}

/**
 * Connect one peer's with another. When connected peer's are able to generate token valid for both side.
 * The usr and pwd are the admin password in the destionation (ns1.mycelius.com)
 * ex. ./Globular connectPeer -dest=ns1.mycelius.com -u=sa -p=adminadmin
 */
func connectPeer(g *Globule, address string) error {

	// Create the remote resource service
	remoteResourceClient, err := resource_client.NewResourceService_Client(address, "resource.ResourceService")
	if err != nil {
		return err
	}

	// Get the local peer key
	key, err := security.GetPeerKey(globule.Mac)
	if err != nil {
		log.Println("fail to get peer key:", err)
		return err
	}

	// Register the peer on the remote resource client...
	hostname, _ := os.Hostname()

	peer, generateKey, err := remoteResourceClient.RegisterPeer(string(key), &resourcepb.Peer{Hostname: hostname, Mac: g.Mac, Domain: g.Domain, LocalIpAddress: configpkg.GetLocalIP(), ExternalIpAddress: Utility.MyIP()})
	if err != nil {
		return err
	}

	peerAddress, _ := configpkg.GetAddress()
	// I will also register the peer to the local server, the local server must running and it domain register,
	// he can be set in /etc/hosts if it's not a public domain.
	localResouceClient, err := resource_client.NewResourceService_Client(peerAddress, "resource.ResourceService")
	if err != nil {
		return err
	}

	_, _, err = localResouceClient.RegisterPeer(generateKey, peer)

	return err
}
