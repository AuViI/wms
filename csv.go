package main

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"
)

type data struct {
	Time                                              string
	Temp, TempMin, Humid, WindGrad, Wind, Rain, Cloud float32
}
type csvData struct {
	City struct {
		Name     string
		Lat, Lon float32
	}
	Data []data
}

const csvtxt = `Ort,{{.City.Name}}
Lat,{{.City.Lat}}
Lon,{{.City.Lon}}
TIMESTAMP,TEMPERATUR,TEMPERATUR-MIN,HUMID%,WINDGRAD,WIND,NIEDERSCHLAG,WOLKEN
{{range .Data}}{{.Time}},{{.Temp}},{{.TempMin}},{{.Humid}},{{.WindGrad}},{{.Wind}},{{.Rain}},{{.Cloud}}
{{end}}`

var csvTemplate, _ = template.New("CSV").Parse(csvtxt)

func csvHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Path[len("/csv/"):]
	if i := strings.Index(city, "/"); i >= 1 {
		city = city[:i]
	}
	csvWriter(w, city)
}

func csvWriter(w http.ResponseWriter, city string) {
	value := new(csvData)
	curr := GetCurrent(city)
	fore := GetForecast(city)
	value.City.Name = fore.City.Name
	value.City.Lat = float32(fore.City.Coord.Lat)
	value.City.Lon = float32(fore.City.Coord.Lon)
	value.Data = make([]data, 1+len(fore.Data))
	value.Data[0] = data{
		Time:     fmt.Sprintf("%v", time.Unix(curr.Dt, 0))[:19],
		Temp:     float32(ktoc(curr.Main.TempMax)),
		TempMin:  float32(ktoc(curr.Main.TempMin)),
		Humid:    float32(curr.Main.Humidity),
		WindGrad: float32(curr.Wind.Deg),
		Wind:     float32(curr.Wind.Speed),
		Rain:     float32(curr.Rain.Volume),
		Cloud:    float32(curr.Clouds.All) * 0.08,
	}
	for i := 1; i <= len(fore.Data); i++ {
		value.Data[i] = data{
			Time:     fore.Data[i-1].TimeString,
			Temp:     float32(ktoc(fore.Data[i-1].Main.TempMaxK)),
			TempMin:  float32(ktoc(fore.Data[i-1].Main.TempMinK)),
			Humid:    float32(fore.Data[i-1].Main.Humidity),
			WindGrad: float32(fore.Data[i-1].Wind.Degree),
			Wind:     float32(fore.Data[i-1].Wind.Speed),
			Rain:     float32(fore.Data[i-1].Rain.Amount),
			Cloud:    float32(fore.Data[i-1].Clouds.All) * 0.08,
		}
	}
	csvTemplate.Execute(w, value)
}
