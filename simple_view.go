package main

import (
	"fmt"
	"html/template"
	"io"
)

type (
	WeatherValue struct {
		Display string
		Value   template.HTML
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
			WeatherValue{"Temperatur:", template.HTML(fmt.Sprintf("%.2f °C", d.Main.Temp-272.15))},
			WeatherValue{"Windstärke:", template.HTML(fmt.Sprintf("%.2f Knoten", d.Wind.Speed))},
			WeatherValue{"Windrichtung:", template.HTML(fmt.Sprintf(`<span style="-ms-transform:rotate(%.2fdeg); -webkit-transform:rotate(%.2fdeg); transform:rotate(%.2fdeg); display:block;position:absolute;right:6em;">&#8613;</span> %.2f°`, d.Wind.Deg, d.Wind.Deg, d.Wind.Deg, d.Wind.Deg))},
			WeatherValue{"Luftdruck:", template.HTML(fmt.Sprintf("%.2f hpa", d.Main.Pressure))},
			WeatherValue{"Humidity:", template.HTML(fmt.Sprintf("%d%%", int(d.Main.Humidity)))},
		},
	}
	simpleHTMLtemplate.Execute(w, p)
}
