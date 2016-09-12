package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	port = flag.String("port", ":8080", "Port to host webservice on")
	wd   = flag.String("wd", fmt.Sprintf("%s/src/github.com/auvii/wms/", os.Getenv("GOPATH")), "working directory")
	nc   = flag.Bool("no-cache", false, "If this is set, templates won't be cached")
	help = flag.Bool("help", false, "Displays possible parameters")
)

func main() {
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		return
	}
	err := os.Chdir(*wd)
	if err != nil {
		Fail(fmt.Sprint("Could not enter working directory:", *wd))
		os.Exit(1)
	}
	dir, _ := os.Getwd()
	Ok(fmt.Sprint("cwd:", dir))
	Continue(fmt.Sprintf("Using Cache: %v", !*nc))
	Continue(fmt.Sprintf("Using Port: %v", *port))
	Continue("AuViI Server starting")
	webSetup(port)
	os.Exit(0)
}

func Fail(msg string) {
	fmt.Printf("\033[1m\033[31m[FAIL]\033[0m %s\n", msg)
}

func Ok(msg string) {
	fmt.Printf("\033[1m\033[32m[ OK ]\033[0m %s\n", msg)
}

func Continue(msg string) {
	fmt.Printf("\033[0m\033[32m[ .. ]\033[0m %s\n", msg)
}
