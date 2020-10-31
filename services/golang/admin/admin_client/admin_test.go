package admin_client

import (
	//"encoding/json"
	"log"
	"testing"
	//	"time"
)

var (
	// Connect to the admin client.
	client, _ = NewAdminService_Client("globular.live", "admin.AdminService")
)

// Test various function here.
func TestGetConfig(t *testing.T) {
	log.Println("---> test get config.")
	config, err := client.GetConfig()
	if err != nil {
		log.Println("---> ", err)
	}

	log.Println("config: ", config)
}

/*
func TestGetFullConfig(t *testing.T) {
	log.Println("---> test get config.")

	config, err := client.GetFullConfig()
	if err != nil {
		log.Println("---> ", err)
	}
	log.Println("config: ", config)
}
*/
func TestRestartServices(t *testing.T) {
	log.Println("---> test get config.")
	log.Println("---->", client)
	err := client.RestartServices()
	if err != nil {
		log.Println("---> ", err)
	}

}

/*
func TestStopService(t *testing.T) {
	err := client.StopService("file_server")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("---> stop service succeeded")
}

func TestStartService(t *testing.T) {
	log.Println("---> test get config.")

	service_pid, proxy_pid, err := client.StartService("file_server")
	if err != nil {

		log.Println("---> ", err)
	}
	log.Println("service pid:", service_pid, " proxy pid:", proxy_pid)

}*/

// Test modify the config...
/*func TestSaveConfig(t *testing.T) {
	log.Println("---> test get config.")

	configStr, err := client.GetFullConfig()
	if err != nil {
		log.Println("---> ", err)
	}

	// So here I will use a map to intialyse config...
	config := make(map[string]interface{})
	err = json.Unmarshal([]byte(configStr), &config)
	if err != nil {
		log.Println("---> ", err)
	}

	// Print the services configuration here.
	config["IdleTimeout"] = 220

	configStr_, _ := json.Marshal(config)

	// Here I will save the configuration directly...
	err = client.SaveConfig(string(configStr_))
	if err != nil {
		log.Println("---> ", err)
	}

	// Now I will try to save a single service.
	serviceConfig := config["Services"].(map[string]interface{})["echo_server"].(map[string]interface{})
	serviceConfig["Port"] = 10029 // set new port number.
	serviceConfigStr_, _ := json.Marshal(serviceConfig)
	err = client.SaveConfig(string(serviceConfigStr_))
	if err != nil {
		log.Println("---> ", err)
	}
}*/

// Test register/start external service.
/*func TestRegisterExternalService(t *testing.T) {
	// Start mongo db
	pid, err := client.RegisterExternalApplication("mongoDB_srv_win64", "E:\\MongoDB\\bin\\mongod.exe", []string{"--port", "27017", "--dbpath", "E:\\MongoDB\\data\\db"})

	if err == nil {
		log.Println("---> mongo db start at port: ", pid)
	} else {
		log.Println("---> err", err)
	}
}*/

/*func TestPublishService(t *testing.T) {
	err := client.PublishService("echo_server", "localhost:8080", "localhost:8080", "Echo is the simplest serive of all.", []string{"test", "echo"})
	if err != nil {
		log.Panicln(err)
	}
}

func TestInstallService(t *testing.T) {
	err := client.InstallService("localhost:8080", "localhost", "echo_server")
	if err != nil {
		log.Panicln(err)
	}
	time.Sleep(time.Second * 5)
}*/
/*
func TestUninstallService(t *testing.T) {
	err := client.UninstallService("localhost", "echo_server", "0.0.1")
	if err != nil {
		log.Panicln(err)
	}
}*/

/*func TestDeployApplication(t *testing.T) {
	err := client.DeployApplication("testApp", "/home/dave/Documents/chitchat")
	if err != nil {
		log.Panicln(err)
	}
}*/
