package main

import (
	"fmt"
	"log"
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
				Utility.CreateDirIfNotExist(path + string(os.PathSeparator) + "bin")
				// copy the bin.
				log.Println("---> copy ", dir+string(os.PathSeparator)+"bin", "to", path+string(os.PathSeparator)+"bin")
				err = Utility.CopyDirectory(dir+string(os.PathSeparator)+"bin", path+string(os.PathSeparator)+"bin")
				if err != nil {
					log.Panicln("--> fail to copy bin ", err)
				}

				// install services.
				for _, service := range g.Services {
					s := service.(map[string]interface{})
					name := s["Name"].(string)

					// set the name.
					Utility.CreateDirIfNotExist(path + string(os.PathSeparator) + name)
					execPath := dir + s["servicePath"].(string)
					destPath := path + string(os.PathSeparator) + name + string(os.PathSeparator) + name
					if string(os.PathSeparator) == "\\" {
						execPath += ".exe" // in case of windows
						destPath += ".exe"
					}

					err := Utility.Copy(execPath, destPath)
					if err != nil {
						fmt.Println(err)
					}

					err = os.Chmod(destPath, 0755)
					if err != nil {
						fmt.Println(err)
					}

					configPath := dir + s["configPath"].(string)
					if Utility.Exists(configPath) {
						Utility.Copy(configPath, path+string(os.PathSeparator)+name+string(os.PathSeparator)+"config.json")
					}

					protoPath := dir + s["protoPath"].(string)
					if Utility.Exists(protoPath) {
						Utility.Copy(protoPath, path+string(os.PathSeparator)+name+string(os.PathSeparator)+name+".proto")
					}

				}
			}
		}
	} else {
		g.Listen()
	}
}
