package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path[1:] {
	case "":
		serveIndex(w)
	case "bsp/":
		fmt.Fprint(w, "<html>Beispiele:<br>")
		// Bad Code starting
		orte := []string{"Berlin", "Kühlungsborn", "Oslo", "New York"}
		for _, v := range orte {
			fmt.Fprintf(w, "<a href=\"/txt/%s/\">txt</a> <a href=\"/view/%s/\">view</a> %s<br>", v, v, v)
		}
	}
}

func txtHandler(w http.ResponseWriter, r *http.Request) {
	cut := r.URL.Path[len("/txt/"):]
	if i := strings.Index(cut, "/"); i != -1 {
		cut = cut[:i]
	}
	fmt.Fprintf(w, "%s", PrognoseTxt(cut, 3))
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	cut := r.URL.Path[len("/view/"):]
	if i := strings.Index(cut, "/"); i != -1 {
		cut = cut[:i]
	}
	SimpleHTML(cut, w)
}

func resourceHandler(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Path[len("/resources/"):]
	switch s {
	case "main.css":
		w.Header().Set("Content-type", "text/css")
		styleTemplate.Execute(w, nil)
	default:
		r, ok := resources[s]
		if !ok || r == "" {
			fmt.Fprintf(w, "404 - not found %s", s)
			return
		}
		io.WriteString(w, r)
	}
}

func webSetup(port *string) {
	http.HandleFunc("/txt/", txtHandler)
	http.HandleFunc("/csv/", csvHandler) // csv.go
	http.HandleFunc("/view/", viewHandler)
	// http.Handle("/resources/", http.FileServer(http.Dir("./resources/")))
	http.HandleFunc("/", handler)
	http.HandleFunc("/resources/", resourceHandler)
	http.ListenAndServe(*port, nil)
}

// Handler Helper
var indexTemplate, _ = template.ParseFiles("./template/index.html")
var styleTemplate, _ = template.ParseFiles("./template/main.css")

var resources = map[string]string{
	"logo.png": load("logo.png"),
}

func load(res string) string {
	bt, err := ioutil.ReadFile(fmt.Sprintf("./resources/%s", res))
	if err != nil {
		fmt.Printf("trying to load %s, not found", res)
		return ""
	}
	return string(bt)
}

func serveIndex(w io.Writer) {
	indexTemplate.Execute(w, nil)
}
