package main

import (
	"fmt"
	"html/template"
	"io"

	"github.com/auvii/wms/weather"
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
	d := weather.GetCurrent(city)
	p := &Page{
		fmt.Sprintf("Aktuelle Daten für %s", city),
		[]WeatherValue{
			WeatherValue{"Temperatur:", template.HTML(fmt.Sprintf("%.0f °C", d.Main.Temp-272.15))},
			WeatherValue{"Windstärke:", template.HTML(fmt.Sprintf("%.0f km/h", d.Wind.Speed*1.852))},
			WeatherValue{"Windrichtung:", template.HTML(fmt.Sprintf(`<span class="notonmobile" style="-ms-transform:rotate(%.0fdeg); -webkit-transform:rotate(%.2fdeg); transform:rotate(%.2fdeg); display:block;position:absolute;right:9em;">&#8613;</span> %.0f Grad`, d.Wind.Deg, d.Wind.Deg, d.Wind.Deg, d.Wind.Deg))},
			WeatherValue{"Luftdruck:", template.HTML(fmt.Sprintf("%.0f hpa", d.Main.Pressure))},
			WeatherValue{"Luftfeuchtigkeit:", template.HTML(fmt.Sprintf("%d%%", int(d.Main.Humidity)))},
		},
	}
	simpleHTMLtemplate.Execute(w, p)
}
