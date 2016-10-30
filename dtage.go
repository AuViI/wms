package main

import (
	"fmt"
	"github.com/auvii/wms/weather"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type (
	ntag struct {
		N    uint
		Ort  string
		Data []row
	}
	row struct {
		Name string
		Data []string
		Bold bool
		Unit string
	}
	// RowError is not nil if row cannot be added to ntag struct
	RowError string
)

// row[0] == {"", "Heute", "Morgen", "..."}

const (
	dtaghtml = "./template/dtage.html"
	htmlC    = "<sup>o</sup>C"
)

var (
	dtagetmpl, _ = template.ParseFiles(dtaghtml)
)

func (n ntag) Now() string {
	return time.Now().String()
}

func (n *ntag) Row(name string, data []string, bold bool, unit string) error {
	if n.N != uint(len(data)) {
		return RowError("Row lengths differ")
	}
	n.Data = append(n.Data, row{Name: name, Data: data, Bold: bold, Unit: unit})
	return nil
}

func getTagName(i uint) string {
	switch i {
	case 0:
		return "Heute"
	case 1:
		return "Morgen"
	default:
		return fmt.Sprintf("in %d Tagen", i)
	}
}

func newNTage(n uint, ort string) *ntag {
	fr := make([]string, n)
	for i := 0; i < int(n); i++ {
		fr[i] = getTagName(uint(i))
	}
	return &ntag{
		N:   n,
		Ort: ort,
		Data: []row{
			row{
				Name: "",
				Data: fr,
				Bold: true,
				Unit: "",
			},
		},
	}
}

func (r RowError) Error() string {
	return fmt.Sprintf("RowError: %s", r)
}

type simpleTime int64

func (s simpleTime) String() string {
	t := time.Unix(int64(s), 0)
	return fmt.Sprintf("%02d:%02d %d.%02d.%d", t.Hour(), t.Minute(), t.Day(), t.Month(), t.Year())
}

func getIconString(icon string) string {
	return fmt.Sprintf("<img class=\"icon\" src=\"http://openweathermap.org/img/w/%s.png\" alt=\"Fehler\" />", icon)
}

func fillMeteo(out *ntag) error {
	cwd := weather.GetCurrent(out.Ort).ConvertToCelsius()
	fwd := weather.GetForecast(out.Ort).Filter(weather.MIDDAY)

	addStr := func(name string, data []string, unit string) {
		out.Row(name, data, false, unit)
	}

	//tDat := make([]string, out.N)
	tempDat := make([]string, out.N) // Temperatur
	presDat := make([]string, out.N) // Luftdruck
	clouDat := make([]string, out.N) // Wolkendecke
	iconDat := make([]string, out.N) // Icon
	humiDat := make([]string, out.N) // Luftfeuchtigkeit
	wspeDat := make([]string, out.N) // Windgeschwindigkeit

	if time.Unix(cwd.Dt, 0).Day() == time.Unix(fwd.Data[0].Time, 0).Day() {
		// Use forecast data for "today"
		//tDat[0] = simpleTime(fwd.Data[0].Time).String()
		tempDat[0] = fmt.Sprintf("%.0f", weather.Ktoc(fwd.Data[0].Main.TempK))
		presDat[0] = fmt.Sprintf("%.0f", fwd.Data[0].Main.Pressure)
		clouDat[0] = fmt.Sprintf("%.0f", float32(fwd.Data[0].Clouds.All)*0.08)
		iconDat[0] = getIconString(fwd.Data[0].Weather[0].Icon)
		humiDat[0] = fmt.Sprintf("%d", fwd.Data[0].Main.Humidity)
		wspeDat[0] = fmt.Sprintf("%.0f", fwd.Data[0].Wind.Speed*1.852)

		// don't use the same data twice
		fwd.Data = fwd.Data[1:]
	} else {
		// Use current data for "today"
		//tDat[0] = simpleTime(cwd.Dt).String()
		tempDat[0] = fmt.Sprintf("%.0f", cwd.Main.Temp)
		presDat[0] = fmt.Sprintf("%.0f", cwd.Main.Pressure)
		clouDat[0] = fmt.Sprintf("%.0f", cwd.Clouds.All*0.08)
		iconDat[0] = getIconString(cwd.Weather[0].Icon)
		humiDat[0] = fmt.Sprintf("%.0f", cwd.Main.Humidity)
		wspeDat[0] = fmt.Sprintf("%.0f", cwd.Wind.Speed*1.852)
	}
	for i, val := range fwd.Data {
		if uint(i) == out.N-1 {
			break
		}
		//tDat[i+1] = simpleTime(val.Time).String()
		tempDat[i+1] = fmt.Sprintf("%.0f", weather.Ktoc(val.Main.TempK))
		presDat[i+1] = fmt.Sprintf("%.0f", val.Main.Pressure)
		clouDat[i+1] = fmt.Sprintf("%.0f", float32(val.Clouds.All)*0.08)
		iconDat[i+1] = getIconString(val.Weather[0].Icon)
		humiDat[i+1] = fmt.Sprintf("%d", val.Main.Humidity)
		wspeDat[i+1] = fmt.Sprintf("%.0f", val.Wind.Speed*1.852)
	}
	//addStr("Zeitstempel", tDat, "")
	addStr("Temperatur", tempDat, htmlC)
	addStr("Luftdruck", presDat, "hPa")
	addStr("Luftfeuchtigkeit", humiDat, "%")
	addStr("Windgeschwindigkeit", wspeDat, "km/h")
	addStr("Wetterlage", iconDat, "")
	addStr("Wolkendecke", clouDat, "/8")
	return nil
}

func fillAstro(out *ntag) error {
	return out.Row("Sonnenaufgang", []string{"4:50", "5:20", "3:55"}, false, "Uhr")
	return nil
}

func fillCurrent(out *ntag) error {
	*out = *newNTage(1, out.Ort)
	cwd := weather.GetCurrent(out.Ort).ConvertToCelsius()
	addStr := func(name string, data interface{}, unit string) {
		switch data.(type) {
		case string, int, uint, int32, int64:
			out.Row(name, []string{fmt.Sprintf("%v", data)}, false, unit)
		case float32, float64:
			out.Row(name, []string{fmt.Sprintf("%.0f", data)}, false, unit)
		}
	}
	addStr("Temperatur", cwd.Main.Temp, htmlC)
	addStr("Luftdruck", cwd.Main.Pressure, "hPa")
	addStr("Luftfeuchtigkeit", cwd.Main.Humidity, "%")
	addStr("Windgeschwindigkeit", cwd.Wind.Speed*1.852, "km/h")
	addStr("Wetterlage", getIconString(cwd.Weather[0].Icon), "")
	addStr("Wolkendecke", cwd.Clouds.All*0.08, "/8")
	return nil
}

func handleDTage(w http.ResponseWriter, r *http.Request) {
	req := strings.Split(r.URL.Path, "/")
	num := uint(3)
	display := "meteo" // different modes
	switch len(req) {
	case 0:
	case 1:
		// unreachable
		// ["", "dtage", "ort", "n", "type"]
		return
	case 2:
		fmt.Fprint(w, "Du musst einen Ort angeben\nbsp.: /dtage/Braunschweig")
		break
	case 3:
		break
	}
	if len(req) > 3 {
		tnum, err := strconv.ParseInt(req[3], 10, 32)
		if err != nil {
			fmt.Fprint(w, "Der Zeitparameter ist fehlerhaft")
			return
		}
		num = uint(tnum)
	}
	if len(req) > 4 {
		display = req[4]
		if display == "aktuell" {
			num = 1
		}
	}
	result := newNTage(num, req[2])
	switch display {
	case "meteo":
		fillMeteo(result)
	case "astro":
		fillAstro(result)
	case "aktuell":
		fillCurrent(result)
	default:
		fmt.Fprint(w, "Unbekannter Modus")
		return
	}
	dtagetmpl.Execute(w, result)
}

func ncHandleDTage(w http.ResponseWriter, r *http.Request) {
	dtagetmpl, _ = template.ParseFiles(dtaghtml)
	handleDTage(w, r)
}
