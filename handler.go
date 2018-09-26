package main

import (
	"encoding/json"
	"fmt"
	"github.com/AuViI/wms/config"
	"github.com/AuViI/wms/forecast"
	"github.com/AuViI/wms/txt"
	"github.com/AuViI/wms/uid"
	"github.com/AuViI/wms/weather"
	"github.com/mmcloughlin/globe"
	"image/color"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
)

var (
	cachePNGMutex = &sync.Mutex{}
)

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path[1:] {
	case "":
		serveIndex(w)
		break
	default:
		fmt.Fprintf(w, "Unknown target: %s", r.URL.Path)
		break
	}
}

func txtHandler(w http.ResponseWriter, r *http.Request) {
	cut := r.URL.Path[len("/txt/"):]
	if i := strings.Index(cut, "/"); i != -1 {
		cut = cut[:i]
	}
	fmt.Fprintf(w, "%s", txt.PrognoseTxt(cut, DaysForecastTxt))
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
			st, _ := template.ParseFiles("./template/main.css")
			styleTemplate = st
		}
		w.Header().Set("Content-type", "text/css")
		err := styleTemplate.Execute(w, nil)
		if err != nil {
			Fail(fmt.Sprintf("main.css: %s", err))
		}
		break
	case "list_edit.min.js":
		err := editJsMinTmpl.Execute(w, nil)
		if err != nil {
			Fail(fmt.Sprintf("list_edit.min.js %s", err))
		}
		break
	case "admin/nav":
		var adminnav *template.Template
		// if *nc, otherwise TODO cache
		adminnav, _ = template.ParseFiles("./template/adminnav.gohtml")
		var reqb struct {
			Origin string `json:"origin"`
		}
		json.NewDecoder(r.Body).Decode(&reqb)
		adminnav.Execute(w, struct {
			Origin string
			Pages  []struct {
				Display string
				Link    string
				Active  bool
			}
		}{
			Origin: reqb.Origin,
			Pages: []struct {
				Display string
				Link    string
				Active  bool
			}{
				{"Übersicht", "/resources/admin/index.html", true},
				{"Linkgenerator", "/resources/admin/linkgen.html", true},
				{
					"Kunden<div id=\"usercount-nav\" class=\"float-right small\"></div>",
					"#",
					false,
				},
				{"Wetterbericht", "#", false},
			},
		})
		break
	default:
		if resourceLocation.Has(s) {
			Continue(fmt.Sprintf("found redirection from %s", s))
			s = resourceLocation.Get(s) // is redirected, use real location
			Continue(fmt.Sprintf("                  to %s", s))
		}
		if !resources.Has(s) {
			fmt.Fprintf(w, "404 - not found %s", s)
			return
		}
		res := resources.Get(s)
		io.WriteString(w, res)
		break
	}
}

func sleepHandler(w http.ResponseWriter, r *http.Request) {
	//  /tools/sleep/<num>
	// 0 1     2     3
	path := strings.Split(r.URL.Path, "/")
	if len(path) < 4 {
		return
	}
	num, cerror := strconv.ParseInt(path[3], 10, 64)
	if cerror != nil {
		fmt.Fprintf(w, "Error: %v", cerror)
		return
	}

	func(w io.Writer, n int64) {
		fmt.Fprintf(w, "Waiting %dms ", n)
		time.Sleep(time.Millisecond * time.Duration(n))
		fmt.Fprintf(w, "done\n")
	}(w, num)
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
		link := fmt.Sprintf("http://%s/list/%s/edit", r.Host, uid.GetRUID(8))
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

func bspHandler(w http.ResponseWriter, r *http.Request) {
	if *nc {
		bspTmpl, _ = template.ParseFiles("./template/bsp.html")
	}
	conf, err := config.GetEasyConfig()
	var modi, orte []string
	if err != nil {
		modi = []string{"error"}
		orte = []string{"error"}
	} else {
		modi = conf.ExampleModi
		orte = conf.ExampleCities
	}
	data := struct {
		Orte  []string
		Modes []string
		Dtage []string
		Rend  []string
		Show  func(string) string
	}{
		Orte:  orte,
		Modes: modi,
		Dtage: []string{"1/aktuell", "3/meteo", "5/meteo", "3/astro", "5/astro"},
		Rend:  renderFiles(),
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

func renderFiles() []string {
	infos, errs := ioutil.ReadDir(renderFolder)
	if errs != nil {
		return nil
	}
	files := make([]string, 0)
	for _, v := range infos {
		files = append(files, v.Name())
	}
	return files
}

func renderHandler(w http.ResponseWriter, r *http.Request) {
	fold, errf := os.Open(renderFolder)
	if os.IsNotExist(errf) {
		os.Mkdir(renderFolder, os.ModeDir|os.ModePerm)
	}
	fold.Close()
	l := strings.Split(r.URL.Path, "/")
	infos, errs := ioutil.ReadDir(renderFolder)
	if errs != nil {
		fmt.Fprintf(w, "Error while reading directory %s: %s", renderFolder, errs)
		return
	}
	reqFileName := l[len(l)-1]
	otherwise := ""
	add := func(location string) {
		otherwise += fmt.Sprintf("<a href='%s'>%s</a><br>", location, location)
	}
	for _, v := range infos {
		if v.Name() == reqFileName {
			f, e := os.Open(path.Join(renderFolder, v.Name()))
			if e != nil {
				fmt.Fprintf(w, "Error while reading file %s: %s", v.Name(), e)
				return
			}
			io.Copy(w, f)
			f.Close()
			return
		}
		add(v.Name())
	}
	fmt.Fprint(w, otherwise)
}

func cacheHandler(w http.ResponseWriter, r *http.Request) {
	globe.DefaultStyle.Background = color.RGBA{0xff, 0xff, 0xff, 0x00}
	g := globe.New()
	g.DrawGraticule(10.0)
	g.DrawLandBoundaries()
	g.DrawCountryBoundaries(globe.Color(color.RGBA{0x00, 0x00, 0x00, 0x0F}))
	for _, v := range weather.GetCachedLocations() {
		// TODO v[2] for color
		// fmt.Println(v)
		if v[0] != 0 && v[1] != 0 {
			g.DrawDot(v[0], v[1], 0.02, globe.Color(color.RGBA{0x00, 0x00, 0xFF, 0xFF}))
		}
	}
	g.CenterOn(52.0, 11.0)
	cachePNGMutex.Lock()
	defer cachePNGMutex.Unlock()
	g.SavePNG("/tmp/globe.png", 1000)
	f, e := os.Open("/tmp/globe.png")
	if e != nil {
		fmt.Fprintln(w, "error reading globe.png")
		return
	}
	w.Header().Set("Content-type", "image/png")
	io.Copy(w, f)
	f.Close()
}
