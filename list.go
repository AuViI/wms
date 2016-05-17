package main

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
)

// /list/ -> redirect /list/[:alpha:]{8}/edit
// /list/[:alpha:]{8} to view

var editTemplate, _ = template.ParseFiles("./template/list_edit.html")
var listIDregex = regexp.MustCompile("/list/([a-zA-Z0-9]{8})/edit")

type editData struct {
	UID string
}

func listEditHander(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	data := new(editData)
	data.UID = listIDregex.FindAllStringSubmatch(r.URL.Path, 0)[0][1]

	editTemplate.Execute(w, data)
}

func listShowHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "now showing %s", "pattern [TODO]")
}
