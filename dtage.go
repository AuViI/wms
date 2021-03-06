package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/AuViI/wms/model"
	"github.com/AuViI/wms/simpleuser"
	"github.com/AuViI/wms/weather"
)

type (
	ntag struct {
		N     uint
		Ort   string
		Data  []row
		Theme model.TemplateTheme
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
	tn := time.Now()
	toGerman := func(d time.Weekday) string {
		s := d.String()
		switch s {
		case "Monday":
			return "Montag"
		case "Tuesday":
			return "Dienstag"
		case "Wednesday":
			return "Mittwoch"
		case "Thursday":
			return "Donnerstag"
		case "Friday":
			return "Freitag"
		case "Saturday":
			return "Samstag"
		case "Sunday":
			return "Sonntag"
		default:
			return s
		}
	}
	// germDay := (func() string {
	// 	switch tn.Month() {
	// 	case 1:
	// 		return "Jan."
	// 	case 2:
	// 		return "Feb."
	// 	case 3:
	// 		return "März"
	// 	case 4:
	// 		return "Apr."
	// 	case 5:
	// 		return "Mai"
	// 	case 6:
	// 		return "Juni"
	// 	case 7:
	// 		return "Juli"
	// 	case 8:
	// 		return "Aug."
	// 	case 9:
	// 		return "Sept."
	// 	case 10:
	// 		return "Okt."
	// 	case 11:
	// 		return "Nov."
	// 	case 12:
	// 		return "Dez."
	// 	}
	// 	return fmt.Sprintf("%02d", tn.Month())
	// }())
	return fmt.Sprintf("%v %d.%d.%d", toGerman(tn.Weekday()), tn.Day(), tn.Month(), tn.Year())
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
	return fmt.Sprintf("RowError: %s", string(r))
}

type simpleTime int64

func (s simpleTime) String() string {
	t := time.Unix(int64(s), 0)
	return fmt.Sprintf("%02d:%02d %d.%02d.%d", t.Hour(), t.Minute(), t.Day(), t.Month(), t.Year())
}

func (s simpleTime) ReverseString() string {
	t := strings.Split(s.String(), " ")
	return fmt.Sprintf("%s %s", t[1], t[0])
}

func (s simpleTime) Error() string {
	t := time.Unix(int64(s), 0)
	return fmt.Sprintf("%X", t.Unix())
}

func getIconString(icon string) string {
	return fmt.Sprintf("<img class=\"icon\" src=\"http://openweathermap.org/img/w/%s.png\" alt=\"Fehler\" />", icon)
}

func fillMeteo(out *ntag, offset uint) error {
	cwd := weather.GetCurrent(out.Ort).ConvertToCelsius()
	fwd := weather.GetForecast(out.Ort).Filter(weather.MIDDAY)

	specifyTime := func(i uint64, dt int64) {
		out.Data[0].Data[i] = fmt.Sprintf("%s<hr class=\"exact\"><span class=\"exact\">%s</span>", out.Data[0].Data[i], simpleTime(dt).ReverseString())
	}

	if len(fwd.Data) == 0 || len(cwd.Weather) == 0 {
		Continue("invalid request caught")
		return RowError("invalid request icon string")
	}

	addStr := func(name string, data []string, unit string) {
		out.Row(name, data, false, unit)
	}

	if cwd.Dt == 0 || fwd.Cnt == 0 {
		addStr("Fehler", make([]string, out.N), "404")
	}

	//tDat := make([]string, out.N)
	tempDat := make([]string, out.N) // Temperatur			(K & C -> C)
	presDat := make([]string, out.N) // Luftdruck			(hPa   -> hPa)
	clouDat := make([]string, out.N) // Wolkendecke			(%     -> 1/8)
	iconDat := make([]string, out.N) // Icon				(txt   -> link)
	humiDat := make([]string, out.N) // Luftfeuchtigkeit	(%     -> %)
	wspeDat := make([]string, out.N) // Windgeschwindigkeit (m/s   -> km/h & bf)

	if offset == 0 {
		if time.Unix(cwd.Dt, 0).Day() == time.Unix(fwd.Data[0].Time, 0).Day() {
			// Use forecast data for "today"
			//tDat[0] = simpleTime(fwd.Data[0].Time).String()
			tempDat[0] = fmt.Sprintf("%.1f", weather.Ktoc(fwd.Data[0].Main.TempK))
			presDat[0] = fmt.Sprintf("%.0f", fwd.Data[0].Main.Pressure)
			clouDat[0] = fmt.Sprintf("%.0f/8", float32(fwd.Data[0].Clouds.All)*0.08)
			iconDat[0] = getIconString(fwd.Data[0].Weather[0].Icon)
			humiDat[0] = fmt.Sprintf("%d", fwd.Data[0].Main.Humidity)
			wspeDat[0] = fmt.Sprintf("%.0f km/h<br>%.0f Bf.<span class=\"hidden\">%.2fm/s</span>", fwd.Data[0].Wind.Speed*3.6, weather.MphToBf(fwd.Data[0].Wind.Speed*2.236), fwd.Data[0].Wind.Speed)

			// don't use the same data twice
			fwd.Data = fwd.Data[1:]

			specifyTime(0, fwd.Data[0].Time)
		} else {
			// Use current data for "today"
			//tDat[0] = simpleTime(cwd.Dt).String()
			tempDat[0] = fmt.Sprintf("%.1f", cwd.Main.Temp)
			presDat[0] = fmt.Sprintf("%.0f", cwd.Main.Pressure)
			clouDat[0] = fmt.Sprintf("%.0f/8", cwd.Clouds.All*0.08)
			iconDat[0] = getIconString(cwd.Weather[0].Icon)
			humiDat[0] = fmt.Sprintf("%.0f", cwd.Main.Humidity)
			wspeDat[0] = fmt.Sprintf("%.0f km/h<br>%.0f Bf.<span class=\"hidden\">%.2fm/s</span>", cwd.Wind.Speed*3.6, weather.MphToBf(cwd.Wind.Speed*2.236), cwd.Wind.Speed)

			specifyTime(0, cwd.Dt)
		}
	}
	ioffset := int(offset)
	for i, val := range fwd.Data {
		if uint(i) == out.N-1+offset {
			break
		}
		//tDat[i+1] = simpleTime(val.Time).String()
		tempDat[i+1+ioffset] = fmt.Sprintf("%.1f", weather.Ktoc(val.Main.TempK))
		presDat[i+1+ioffset] = fmt.Sprintf("%.0f", val.Main.Pressure)
		clouDat[i+1+ioffset] = fmt.Sprintf("%.0f/8", float32(val.Clouds.All)*0.08)
		iconDat[i+1+ioffset] = getIconString(val.Weather[0].Icon)
		humiDat[i+1+ioffset] = fmt.Sprintf("%d", val.Main.Humidity)
		wspeDat[i+1+ioffset] = fmt.Sprintf("%.0f km/h<br>%.0f Bf.<span class=\"hidden\">%.2fm/s</span>", val.Wind.Speed*3.6, weather.MphToBf(val.Wind.Speed*2.236), val.Wind.Speed)

		specifyTime(uint64(i+1+ioffset), val.Time)
	}
	//addStr("Zeitstempel", tDat, "")
	addStr("Temperatur", tempDat, htmlC)
	addStr("Wetterzustand", iconDat, "")
	addStr("Wolkenbedeckung", clouDat, "")
	addStr("Windgeschwindigkeit", wspeDat, "")
	addStr("Luftfeuchtigkeit", humiDat, "%")
	addStr("Luftdruck", presDat, "hPa")
	return nil
}

func fillAstro(out *ntag) error {
	return out.Row("Sonnenaufgang", []string{"4:50", "5:20", "3:55"}, false, "Uhr")
}

func fillCurrent(out *ntag) error {
	tc := out.Theme
	*out = *newNTage(1, out.Ort)
	out.Theme = tc
	cwd := weather.GetCurrent(out.Ort).ConvertToCelsius()
	if len(cwd.Weather) == 0 {
		return RowError("invalid request current")
	}
	addStr := func(name string, data interface{}, unit string) {
		switch data.(type) {
		case string, int, uint, int32, int64:
			out.Row(name, []string{fmt.Sprintf("%v", data)}, false, unit)
		case float32, float64:
			out.Row(name, []string{fmt.Sprintf("%.0f", data)}, false, unit)
		}
	}
	addStr("Temperatur", fmt.Sprintf("%.1f", cwd.Main.Temp), htmlC) // TODO eine Nachkommastelle
	addStr("Luftdruck", cwd.Main.Pressure, "hPa")
	addStr("Luftfeuchtigkeit", cwd.Main.Humidity, "%")
	addStr("Windgeschwindigkeit", fmt.Sprintf("%.0f km/h<br>%.0f Bf.<span class=\"hidden\">%.2fm/s</span>", cwd.Wind.Speed*3.6, weather.MphToBf(cwd.Wind.Speed*2.236), cwd.Wind.Speed), "")
	addStr("Wetterzustand", getIconString(cwd.Weather[0].Icon), "")
	addStr("Wolkenbedeckung", fmt.Sprintf("%.0f/8", cwd.Clouds.All*0.08), "")
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

	var result *ntag
	locparam := model.ThemeRegex.FindStringSubmatch(req[2])
	useparam := model.UserRegex.FindStringSubmatch(req[2])

	if (len(locparam) == 0 || locparam[0] == "") && len(useparam) == 0 {
		result = newNTage(num, req[2])
		result.Theme = model.Theme{
			StartColor: model.ThemeColor{66, 170, 201},
			EndColor:   model.ThemeColor{18, 101, 126},
			IconLink:   "/resources/logo.png",
		}.Prepare()
	} else if len(locparam) == 5 {
		result = newNTage(num, locparam[1])
		result.Theme = model.Theme{
			StartColor: model.ThemeColorFromHex(fmt.Sprintf("#%s", locparam[2])),
			EndColor:   model.ThemeColorFromHex(fmt.Sprintf("#%s", locparam[3])),
			IconLink:   strings.Replace(locparam[4], "|", "/", -1),
		}.Prepare()
	} else if len(useparam) == 3 {
		result = newNTage(num, useparam[1])
		ui, _ := strconv.ParseUint(useparam[2], 10, 64)
		result.Theme = simpleuser.Theme(ui).Prepare()
	} else {
		fmt.Printf("len(locparam) = %d\n", len(locparam))
		fmt.Println("Something went wrong")
	}

	var n error
	switch display {
	case "meteo":
		n = fillMeteo(result, 0)
		break
	case "astro":
		n = fillAstro(result)
		break
	case "aktuell":
		n = fillCurrent(result)
		break
	default:
		fmt.Fprint(w, "Unbekannter Modus")
		return
	}

	if n == nil {
		dtagetmpl.Execute(w, result)
	} else {
		fmt.Fprintf(w, "%s", n)
	}
}

func ncHandleDTage(w http.ResponseWriter, r *http.Request) {
	dtagetmpl, _ = template.ParseFiles(dtaghtml)
	handleDTage(w, r)
}
