package main

import (
	"flag"
	"fmt"
)

var (
	port = flag.String("port", ":8080", "Port to host webservice on")
)

func main() {
	flag.Parse()
	fmt.Println("AuViI Server starting")
	webSetup(port)
}
