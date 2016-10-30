package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/auvii/wms/forecast"
)

// Handler Helper
var indexTemplate, _ = template.ParseFiles("./template/index.html")
var styleTemplate, _ = template.ParseFiles("./template/main.css")
var editJsMinTmpl, _ = template.ParseFiles("./template/list_edit.min.js")

var resources = map[string]string{
	"logo.png": load("logo.png"),
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path[1:] {
	case "":
		serveIndex(w)
	case "bsp/":
		fmt.Fprint(w, "<html><h1>Beispiele</h1>")
		fmt.Fprint(w, "<html><h2>Verschiedene Modi's</h2>")
		orte := []string{"Berlin", "Kühlungsborn", "Oslo", "New York", "Braunschweig", "Rostock"}
		prefix := []string{"txt", "forecast", "list", "csv", "dtage", "view"}
		fmt.Fprint(w, "<table>")
		for _, ort := range orte {
			fmt.Fprintf(w, "<tr><th>%s:</th>", ort)
			for _, pre := range prefix {
				fmt.Fprintf(w, "<td><a href=\"/%s/%s\">%s</a></td>", pre, ort, pre)
			}
			fmt.Fprint(w, "</tr>")
		}
		fmt.Fprint(w, "</table><h2>dtage Grafiken</h2><table>")
		modes := []string{"meteo", "astro"}
		for _, ort := range orte {
			fmt.Fprintf(w, "<tr><th>%s:</th>", ort)
			fmt.Fprintf(w, "<td><a href=\"/dtage/%s/1/aktuell\">(aktuell)</a></td>", ort)
			for _, mode := range modes {
				for _, d := range []int{1, 3, 4} {
					fmt.Fprintf(w, "<td><a href=\"/dtage/%s/%d/%s\">(%d Tage %s)</a></td>", ort, d, mode, d, mode)
				}
			}
			fmt.Fprint(w, "</tr>")
		}
		fmt.Fprintf(w, "</table><code>%s</code>", time.Now().String())
	}
}

func txtHandler(w http.ResponseWriter, r *http.Request) {
	cut := r.URL.Path[len("/txt/"):]
	if i := strings.Index(cut, "/"); i != -1 {
		cut = cut[:i]
	}
	fmt.Fprintf(w, "%s", PrognoseTxt(cut, DaysForecastTxt))
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	cut := r.URL.Path[len("/view/"):]
	if i := strings.Index(cut, "/"); i != -1 {
		cut = cut[:i]
	}
	w.Header().Set("Cache-Control", "max-age=600")
	SimpleHTML(cut, w)
}

func resourceHandler(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Path[len("/resources/"):]
	switch s {
	case "main.css":
		w.Header().Set("Content-type", "text/css")
		err := styleTemplate.Execute(w, nil)
		if err != nil {
			Fail(fmt.Sprintf("main.css: %s", err))
		}
	case "list_edit.min.js":
		err := editJsMinTmpl.Execute(w, nil)
		if err != nil {
			Fail(fmt.Sprintf("list_edit.min.js %s", err))
		}
	default:
		r, ok := resources[s]
		if !ok || r == "" {
			fmt.Fprintf(w, "404 - not found %s", s)
			return
		}
		io.WriteString(w, r)
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	isList, _ := regexp.MatchString("/list/[a-zA-Z0-9]{8}", r.URL.Path)
	isEdit, _ := regexp.MatchString("/list/[a-zA-Z0-9]{8}/edit", r.URL.Path)
	if isList && !isEdit {
		// check if exists, else redirect to /list/[pattern]/edit
		listShowHandler(w, r)
	}
	if !isList && !isEdit {
		// generate uid, redirect to /list/[pattern]/edit
		link := fmt.Sprintf("http://%s/list/%s/edit", r.Host, getRUID(8))
		http.Redirect(w, r, link, 307)
	}
	if isEdit {
		listEditHander(w, r)
	}
}

func forecastHandler(w http.ResponseWriter, r *http.Request) {
	forecast.Show(w, r)
}

func ncForecastHandler(w http.ResponseWriter, r *http.Request) {
	forecast.ShowNoCache(w, r)
}

func noCacheSwitch(cached, nocache http.HandlerFunc) http.HandlerFunc {
	if !*nc {
		return cached
	}
	return nocache
}

func webSetup(port *string) {
	http.HandleFunc("/txt/", txtHandler)
	http.HandleFunc("/csv/", csvHandler) // csv.go
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/forecast/",
		noCacheSwitch(forecastHandler, ncForecastHandler))
	http.HandleFunc("/list/", listHandler)
	http.HandleFunc("/", handler)
	http.HandleFunc("/resources/", resourceHandler)
	http.HandleFunc("/dtage/", noCacheSwitch(handleDTage, ncHandleDTage))
	http.ListenAndServe(*port, nil)
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
