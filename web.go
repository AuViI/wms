package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"text/template"
	"time"
)

var (
	// text/html templates
	indexTemplate = simpleTemplate("index.html")
	styleTemplate = simpleTemplate("main.css")
	normListTmpl  = simpleTemplate("normlist.html")
	bspTmpl       = simpleTemplate("bsp.html")
	gewusstTmpl   = simpleTemplate("gewusst.html")

	// text/css templates
	editJsMinTmpl = simpleTemplate("list_edit.min.js")

	// binary resources
	// TODO change to map string to []byte
	resources = map[string]string{
		"logo.png":        load("logo.png"),
		"logo_invert.png": load("logo_invert.png"),
	}
)

const (
	renderFolder   = "./pics/"
	templateFolder = "./template/"
	resourceFolder = "./resources/"
)

/// simpleTemplate expects the templates to be found inside
/// the packages `templateFolder` folder
func simpleTemplate(tmplfile string) *template.Template {
	t, _ := template.ParseFiles(path.Join(templateFolder, tmplfile))
	return t
}

/// webSetup is called from main and sets up the server
/// and is blocking
func webSetup(port *string) {
	end := startUpdateLoop() // bool chan, used to kill the goroutine
	http.HandleFunc("/txt/", txtHandler)
	http.HandleFunc("/csv/", csvHandler)
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
	http.HandleFunc("/render/", renderHandler)
	http.HandleFunc("/cached/", cacheHandler)
	http.HandleFunc("/resources/", resourceHandler)
	http.HandleFunc("/", handler)
	http.ListenAndServe(*port, nil)
	end <- true
}

/// serveIndex executes the index template and writes
/// to the given Writer
func serveIndex(w io.Writer) {
	if *nc {
		indexTemplate = simpleTemplate("index.html")
	}
	indexTemplate.Execute(w, nil)
}

/// load reads a file inside the resource folder into a string
///
/// TODO change to load []byte
func load(res string) string {
	return loadOr(res, "error")
}

/// loadOr reads a file inside the resource folder and returns
/// its contents or returns the specified alternative on error
func loadOr(res, alt string) string {
	byt, err := ioutil.ReadFile(path.Join(resourceFolder, res))
	if err != nil {
		return alt
	}
	return string(byt)
}

/// noCacheSwitch switches between two `http.HandlerFunc`
/// depending on the nc boolean application flag
func noCacheSwitch(cached, nocache http.HandlerFunc) http.HandlerFunc {
	if !*nc {
		return cached
	}
	return nocache
}

/// startUpdateLoop runs numerous functions to update
/// data all over the server. It runs insite a go
/// routine until the returned channel receives something
///
/// This does not manage its own mutexes but expects
/// called functions to handle being called asynchronically
func startUpdateLoop() chan bool {
	counter := 0
	calls := func() {
		updateGewusst() // update entries for /gewusst/
		if counter%12 == 0 {
			/*
			 * rendering pictures requires an application
			 * being run on the active display in fullscreen.
			 * This can be rather annoying. The application
			 * itself needs to be open for a short while
			 * loading all assets, as any other client would.
			 * A screnshot is taken, and saved to the /render/
			 * folder to be accessible.
			 */
			go (func() {
				<-time.After(2 * time.Minute)
				renderPictures() // render new pictures
			})()
			counter = 0
		}
		counter += 1
	}

	// this chan is returned to give the outside world a
	// chance to kill this looping goroutine
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
