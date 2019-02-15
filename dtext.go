package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/AuViI/wms/model"
	"github.com/AuViI/wms/simpleuser"
	"github.com/AuViI/wms/weather"
	"github.com/AuViI/wms/weather/redirect"
	"github.com/icrowley/fake"
)

const (
	dtexthtml = "./template/dtext.gohtml"
)

var (
	dtexttmpl, _ = template.ParseFiles(dtexthtml)
)

func handleDText(w http.ResponseWriter, r *http.Request) {
	// ["", "dtext", "location", "date"]
	req := strings.Split(r.URL.Path, "/")

	/*
		for i := range req {
			fmt.Printf("req[%d] = %s\n", i, req[i])
		}
	*/

	locparam := model.ThemeRegex.FindStringSubmatch(req[2])
	useparam := model.UserRegex.FindStringSubmatch(req[2])

	var nt model.TemplateTheme
	var loc string
	if (len(locparam) == 0 || locparam[0] == "") && len(useparam) == 0 {
		// use default theme
		nt = model.GetDefaultTheme().Prepare()
		loc = req[2]
	} else if len(locparam) == 5 {
		// use new theme
		nt = model.Theme{
			StartColor: model.ThemeColorFromHex(fmt.Sprintf("#%s", locparam[2])),
			EndColor:   model.ThemeColorFromHex(fmt.Sprintf("#%s", locparam[3])),
			IconLink:   strings.Replace(locparam[4], "|", "/", -1),
		}.Prepare()
		loc = locparam[1]
	} else if len(useparam) == 3 {
		ui, _ := strconv.ParseUint(useparam[2], 10, 64)
		nt = simpleuser.Theme(ui).Prepare()
		loc = useparam[1]
	} else {
		fmt.Printf("len(locparam) = %d\n", len(locparam))
		fmt.Println("Something went wrong")
	}

	loc = redirect.Redirect(loc)
	dt := weather.NowDate(time.Local)

	dtexttmpl.Execute(w, struct {
		Location string
		Theme    model.TemplateTheme
		Text     interface{}
		Date     weather.Date
	}{
		Location: loc,
		Theme:    nt,
		Text: struct {
			Title    string
			Subtitle string
			Main     string
		}{
			Title:    fake.Title(),
			Subtitle: fmt.Sprintf("%s: %s", dt, fake.Sentence()),
			Main:     fake.SentencesN(rand.Intn(5) + 12),
		},
		Date: dt,
	})
}

func ncHandleDText(w http.ResponseWriter, r *http.Request) {
	dtexttmpl, _ = template.ParseFiles(dtexthtml)
	handleDText(w, r)
}
