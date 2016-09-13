package main

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
)

// TODO: move list* to own package

// /list/ -> redirect /list/[:alpha:]{8}/edit
// /list/[:alpha:]{8} to view

var editTemplate, _ = template.ParseFiles("./template/list_edit.html")
var listIDregex = regexp.MustCompile("\\/list\\/([a-zA-Z0-9]{8})\\/edit")

type editData struct {
	UID string
}

func listEditHander(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	if !listIDregex.MatchString(r.URL.Path) {
		http.Redirect(w, r, fmt.Sprintf("http://%s/list", r.Host), 301)
	}
	data := new(editData)
	data.UID = listIDregex.FindAllStringSubmatch(r.URL.Path, -1)[0][1]
	err := editTemplate.Execute(w, data)
	if err != nil {
		// TODO:0 Handle errors in list
		fmt.Println("Failed to execute list template")
		fmt.Println(err)
	}
}

func listShowHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "now showing %s", "pattern [TODO]")
}
