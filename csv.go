package main

import (
	"fmt"
	"net/http"
	"strings"
	"text/template"
	"time"

	"github.com/auvii/wms/weather"
)

// TODO: Move csv* to own package

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
	w.Header().Set("Content-type", "text/csv")
	w.Header().Set("Cache-Control", "max-age=0")
	csvWriter(w, city)
}

func csvWriter(w http.ResponseWriter, city string) {
	value := new(csvData)
	curr := weather.GetCurrent(city)
	fore := weather.GetForecast(city)
	value.City.Name = fore.City.Name
	value.City.Lat = float32(fore.City.Coord.Lat)
	value.City.Lon = float32(fore.City.Coord.Lon)
	value.Data = make([]data, 1+len(fore.Data))
	value.Data[0] = data{
		Time:     fmt.Sprintf("%v", time.Unix(curr.Dt, 0))[:19],
		Temp:     float32(weather.Ktoc(curr.Main.TempMax)),
		TempMin:  float32(weather.Ktoc(curr.Main.TempMin)),
		Humid:    float32(curr.Main.Humidity),
		WindGrad: float32(curr.Wind.Deg),
		Wind:     float32(curr.Wind.Speed),
		Rain:     float32(curr.Rain.Volume),
		Cloud:    float32(curr.Clouds.All) * 0.08,
	}
	for i := 1; i <= len(fore.Data); i++ {
		value.Data[i] = data{
			Time:     fore.Data[i-1].TimeString,
			Temp:     float32(weather.Ktoc(fore.Data[i-1].Main.TempMaxK)),
			TempMin:  float32(weather.Ktoc(fore.Data[i-1].Main.TempMinK)),
			Humid:    float32(fore.Data[i-1].Main.Humidity),
			WindGrad: float32(fore.Data[i-1].Wind.Degree),
			Wind:     float32(fore.Data[i-1].Wind.Speed),
			Rain:     float32(fore.Data[i-1].Rain.Amount),
			Cloud:    float32(fore.Data[i-1].Clouds.All) * 0.08,
		}
	}
	err := csvTemplate.Execute(w, value)
	if err != nil {
		// TODO:40 better handling of error in csv
		fmt.Println("Failed to execute CSV template")
	}
}
