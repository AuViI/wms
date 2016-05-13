package main

import (
	"fmt"
	"net/http"
	"strings"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<html><body>%v not implemented", r.URL.Path)
	fmt.Fprint(w, `<br>Beispiel Links:<br>
        <a href="/txt/Berlin/">/txt/Berlin/</a>
        <a href="/view/Berlin/">/view/Berlin/</a><br>
        <a href="/txt/Kühlunsgborn/">/txt/Kühlunsgborn/</a> 
        <a href="/view/Kühlunsgborn/">/view/Kühlunsgborn/</a><br>
        Text: /txt/<i>ort</i>/<br>Balken: /view/<i>ort</i>/<br>CSV: /csv/<i>ort</i>/
        </body></html>`)
	fmt.Printf("Requested: %s\n", r.URL.Path)
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

func webSetup(port *string) {
	http.HandleFunc("/txt/", txtHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/", handler)
	http.ListenAndServe(*port, nil)
}
