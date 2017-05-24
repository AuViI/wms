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
var (
	indexTemplate, _ = template.ParseFiles("./template/index.html")
	styleTemplate, _ = template.ParseFiles("./template/main.css")
	editJsMinTmpl, _ = template.ParseFiles("./template/list_edit.min.js")
	normListTmpl, _  = template.ParseFiles("./template/normlist.html")
	bspTmpl, _       = template.ParseFiles("./template/bsp.html")
	gewusstTmpl, _   = template.ParseFiles("./template/gewusst.html")
)

var resources = map[string]string{
	"logo.png": load("logo.png"),
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path[1:] {
	case "":
		serveIndex(w)
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
		if *nc {
			styleTemplate, _ = template.ParseFiles("./template/main.css")
		}
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

func normlistHandler(w http.ResponseWriter, r *http.Request) {
	l := strings.Split(r.URL.Path, "/")
	if len(l) > 2 {
		normListTmpl.Execute(w, struct {
			Ort string
		}{
			Ort: l[2],
		})
	} else {
		fmt.Fprintf(w, "not correct request format: %v", l)
	}
}

func ncNormlistHandler(w http.ResponseWriter, r *http.Request) {
	normListTmpl, _ = template.ParseFiles("./template/normlist.html")
	normlistHandler(w, r)
}

func noCacheSwitch(cached, nocache http.HandlerFunc) http.HandlerFunc {
	if !*nc {
		return cached
	}
	return nocache
}

func bspHandler(w http.ResponseWriter, r *http.Request) {
	if *nc {
		bspTmpl, _ = template.ParseFiles("./template/bsp.html")
	}
	data := struct {
		Orte  []string
		Modes []string
		Dtage []string
		Show  func(string) string
	}{
		Orte: []string{"Kühlungsborn", "Braunschweig", "Hamburg", "Berlin", "Oslo",
			"Rostock", "Hannover", "München", "New York", "Tokio"},
		Modes: []string{"txt", "forecast", "list", "csv", "dtage", "view", "normlist"},
		Dtage: []string{"1/aktuell", "3/meteo", "5/meteo", "3/astro", "5/astro"},
		Show: func(s string) string {
			if strings.HasPrefix(s, "1/") {
				return s[2:]
			}
			return strings.Replace(s, "/", " Tage ", 1)
		},
	}
	bspTmpl.Execute(w, data)
}

func gewusstHandler(w http.ResponseWriter, r *http.Request) {
	if *nc {
		gewusstTmpl, _ = template.ParseFiles("./template/gewusst.html")
	}
	gewusstTmpl.Execute(w, messages)
}

func webSetup(port *string) {
	end := startUpdateLoop()
	http.HandleFunc("/txt/", txtHandler)
	http.HandleFunc("/csv/", csvHandler) // csv.go
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/bsp/", bspHandler)
	http.HandleFunc("/forecast/",
		noCacheSwitch(forecastHandler, ncForecastHandler))
	http.HandleFunc("/list/", listHandler)
	http.HandleFunc("/dtage/",
		noCacheSwitch(handleDTage, ncHandleDTage))
	http.HandleFunc("/normlist/",
		noCacheSwitch(normlistHandler, ncNormlistHandler))
	http.HandleFunc("/gewusst/", gewusstHandler)
	http.HandleFunc("/", handler)
	http.HandleFunc("/resources/", resourceHandler)
	http.ListenAndServe(*port, nil)
	end <- true
}

func load(res string) string {
	bt, err := ioutil.ReadFile(fmt.Sprintf("./resources/%s", res))
	if err != nil {
		fmt.Printf("trying to load %s, not found", res)
		return ""
	}
	return string(bt)
}

// expects that called update functions manage mutexes
func startUpdateLoop() chan bool {
	calls := func() {
		updateGewusst()
	}
	end := make(chan bool)
	calls()
	go func(e chan bool) {
		for {
			select {
			case <-time.After(10 * time.Minute):
				calls()
			case <-e:
				return
			}
		}
	}(end)
	return end
}

func serveIndex(w io.Writer) {
	indexTemplate.Execute(w, nil)
}
