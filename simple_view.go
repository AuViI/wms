package main

import (
	"fmt"
	"html/template"
	"io"
)

type (
	WeatherValue struct {
		Display, Value string
	}
	Page struct {
		Title       string
		WeatherData []WeatherValue
	}
)

var simpleHTMLtemplate, _ = template.ParseFiles("template/simple.html")

func SimpleHTML(city string, w io.Writer) {
	fillSimpleTemplate(city, w)
}

func fillSimpleTemplate(city string, w io.Writer) {
	d := GetCurrent(city)
	p := &Page{
		fmt.Sprintf("Aktuelle Daten für %s", city),
		[]WeatherValue{
			WeatherValue{"Temperatur:", fmt.Sprintf("%.2f °C", d.Main.Temp-272.15)},
			WeatherValue{"Windstärke:", fmt.Sprintf("%.2f Knoten", d.Wind.Speed)},
			WeatherValue{"Windrichtung:", fmt.Sprintf("%.2f °", d.Wind.Deg)},
			WeatherValue{"Luftdruck:", fmt.Sprintf("%.2f hpa", d.Main.Pressure)},
			WeatherValue{"Humidity:", fmt.Sprintf("%d%%", int(d.Main.Humidity))},
		},
	}
	simpleHTMLtemplate.Execute(w, p)
}
