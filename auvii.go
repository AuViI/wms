package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	/*
	 *  pointer = flag.Type(
	 *      'cmd',
	 *	    'default',
	 *      'helptext')
	 */

	// port to host webserver on
	port = flag.String(
		"port",
		":8080",
		"Port to host webservice on")

	// wd is directoy to cd into for execution
	// important to find:
	//  - gewusst/
	//  - resources/
	//  - template/
	//  - hfscc/[index.js]
	// [- pics/]
	wd = flag.String(
		"wd",
		fmt.Sprintf(
			"%s/src/github.com/AuViI/wms/",
			os.Getenv("GOPATH")),
		"working directory")

	// nc is true to avoid caching
	nc = flag.Bool(
		"no-cache",
		false,
		"If this is set, templates won't be cached")

	// help displays help text instead of running
	help = flag.Bool(
		"help",
		false,
		"Displays possible parameters")

	// render accepts a list of locations to render
	render = flag.String(
		"render",
		"",
		"Comma seperated locations to render")
)

func main() {
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
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
	Continue(fmt.Sprintf("URL: http://localhost%v", *port))
	Continue("AuViI Server starting")
	webSetup(port)
	os.Exit(0)
}

// Fail sends `msg` to stdout with [FAIL] prefix
func Fail(msg string) {
	fmt.Printf("\033[1m\033[31m[FAIL]\033[0m %s\n", msg)
}

// Ok sends `msg` to stdout with [ OK ] prefix
func Ok(msg string) {
	fmt.Printf("\033[1m\033[32m[ OK ]\033[0m %s\n", msg)
}

// Continue sends `msg` to stdout with [ .. ] prefix
func Continue(msg string) {
	fmt.Printf("\033[0m\033[32m[ .. ]\033[0m %s\n", msg)
}
