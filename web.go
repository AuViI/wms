package main

import (
	"fmt"
	"net/http"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%v not implemented", r.URL.Path)
	fmt.Printf("Requested: %s\n", r.URL.Path)
}

func txtHandler(w http.ResponseWriter, r *http.Request) {
	cut := r.URL.Path[5:]
	req := cut[:strings.Index(cut, "/")]
	fmt.Fprintf(w, "%s", PrognoseTxt(req, 3))
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	cut := r.URL.Path[6:]
	req := cut[:strings.Index(cut, "/")]
	SimpleHTML(req, w)
}

func webSetup(port *string) {
	http.HandleFunc("/txt/", txtHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/", handler)
	http.ListenAndServe(*port, nil)
}
