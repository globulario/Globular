package main

import (
	"os"
	"strconv"
)

func main() {

	var port int = 10000 // The default port.
	if len(os.Args) > 1 {
		port, _ = strconv.Atoi(os.Args[1])
	}

	g := NewGlobule(port)
	g.Listen()

}
