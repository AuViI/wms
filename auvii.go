package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	port = flag.String("port", ":8080", "Port to host webservice on")
	wd   = flag.String("wd", fmt.Sprintf("%s/src/github.com/auvii/wms/", os.Getenv("GOPATH")), "working directory")
)

func main() {
	flag.Parse()
	err := os.Chdir(*wd)
	if err != nil {
		fmt.Println("[FAIL] Could not enter working directory:", *wd)
		os.Exit(1)
	}
	dir, err := os.Getwd()
	fmt.Println("[ OK ] cwd:", dir, err)
	fmt.Println("[ .. ] AuViI Server starting")
	webSetup(port)
	os.Exit(0)
}
