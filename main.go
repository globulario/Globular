package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/davecourtois/Utility"
)

func main() {
	g := NewGlobule()
	if len(os.Args) > 1 {
		argsWithoutProg := os.Args[1:]

		// program := os.Args[0:1]
		for i := 0; i < len(argsWithoutProg); i++ {
			arg := argsWithoutProg[i]
			if arg == "--install" {
				i++

				path := argsWithoutProg[i] + string(os.PathSeparator) + "globular"

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

				// install services.
				for name, service := range g.Services {
					s := service.(map[string]interface{})
					p := s["Process"].(map[string]interface{})

					// set the name.
					name_ := name[0:strings.Index(name, "_")]

					fmt.Println("install  service", name_)

					Utility.CreateDirIfNotExist(path + string(os.PathSeparator) + name)
					execPath := path + string(os.PathSeparator) + name + string(os.PathSeparator) + name

					if string(os.PathSeparator) == "\\" {
						execPath += ".exe" // in case of windows
					}

					fmt.Println("copy -->", execPath)
					err := Utility.Copy(p["Path"].(string), execPath)
					if err != nil {
						fmt.Println(err)
					}
					err = os.Chmod(execPath, 0755)
					if err != nil {
						fmt.Println(err)
					}
					configPath := p["Path"].(string)[0:strings.LastIndex(p["Path"].(string), string(os.PathSeparator))] + string(os.PathSeparator) + "config.json"
					fmt.Println("copy -->", configPath)
					if Utility.Exists(configPath) {
						Utility.Copy(configPath, path+string(os.PathSeparator)+name+string(os.PathSeparator)+"config.json")
					}

					proxy := s["ProxyProcess"]
					if proxy != nil {
						// Here I will copy the proxy.
						proxyExec := proxy.(map[string]interface{})["Path"].(string)[strings.LastIndex(proxy.(map[string]interface{})["Path"].(string), string(os.PathSeparator)):]
						proxyPath := path + string(os.PathSeparator) + "bin" + proxyExec
						Utility.CreateDirIfNotExist(path + string(os.PathSeparator) + "bin")
						if !Utility.Exists(proxyPath) {
							err := Utility.Copy(proxy.(map[string]interface{})["Path"].(string), proxyPath)
							if err != nil {
								fmt.Println(err)
							}
							err = os.Chmod(proxyPath, 0755)
							if err != nil {
								fmt.Println(err)
							}
						}

						// In that case that mean it's a grpc service and a .proto file is required.
						protoPath := p["Path"].(string)[:strings.Index(p["Path"].(string), name_)] + name_ + string(os.PathSeparator) + name_ + "pb" + string(os.PathSeparator) + name_ + ".proto"
						fmt.Println("------------> from ", protoPath)

						fmt.Println("------------> ", path+string(os.PathSeparator)+name+string(os.PathSeparator)+name_+".proto")
						err := Utility.Copy(protoPath, path+string(os.PathSeparator)+name+string(os.PathSeparator)+name_+".proto")
						if err != nil {
							fmt.Println(err)
						}

					}
				}
			}
		}
	} else {
		g.Listen()
	}

}
