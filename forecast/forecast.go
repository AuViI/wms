package forecast

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/auvii/wms/weather"
)

const (
	fcTmpl = "./template/forecast.html"
)

// data contains all Data needed for filling the Template
type data struct {
	Ort          string
	Datum        string
	Uhrzeit      string
	Wetterlage   string
	WetterDesc   string
	PhysGroessen string
	Legende      string
	Time         string
	Geo          geoData
	Cwd          *weather.Data
	Fwd          printFwd
	RFwd         *weather.ForecastData
	Nw           niceWeather
	MapsKey      string
}

// printFwd contains the Raw-Filtered weather data and formatted txt
type printFwd struct {
	Raw weather.ForecastData
	N   []printFwdPoint
}

// printFwdPoint is a single formatted data point of PrintFwd
type printFwdPoint struct {
	Time  int64
	Stamp string
	C     string
	CMax  string
	CMin  string
	Pres  string
	Humid string
	Icon  string
	Main  string
	Desc  string
	Cloud int
	WindS string
	WindD int
	RainA string
}

// geoData contains Lat and Lon
type geoData struct {
	Lat float64
	Lon float64
}

// niceWeather is PrintFwdPoint for current weather
type niceWeather struct {
	Temp    string
	TempMax string
	TempMin string
}

// niceWeatherFromData converts kelvin to celsius
func niceWeatherFromData(w *weather.Data) niceWeather {
	return niceWeather{
		Temp:    fmt.Sprintf("%.2f", w.Main.Temp),
		TempMax: fmt.Sprintf("%.2f", w.Main.TempMax),
		TempMin: fmt.Sprintf("%.2f", w.Main.TempMin),
	}
}

var (
	forecastTemplate, _ = template.ParseFiles(fcTmpl)
	mapskey             = flag.String("maps", func() string {
		str := os.Getenv("GOOGLEAPI")
		if str == "" {
			err := fmt.Errorf("No Google API Key provided\nuse -maps KEY or $GOOGLEAPI")
			fmt.Println(err)
			os.Exit(20)
		}
		return str
	}(), "Google Maps API Key $GOOGLEAPI")
)

// Show writes the forecast
func Show(w http.ResponseWriter, r *http.Request) {
	query := r.URL.RequestURI()[10:]
	if query == "" {
		w.Header().Set("Location", "/forecast/Kühlungsborn")
		w.WriteHeader(301)
		fmt.Println(w.Header())
	} else {
		url, err := url.QueryUnescape(query)
		if err != nil {
			w.Write([]byte("An Error occurred reading the URL"))
			return
		}
		//w.Header().Set("Refresh", "10") // Sekunden
		w.Header().Set("Cache-Control", "max-age=600")
		query = url
	}
	cwd := weather.GetCurrent(query).ConvertToCelsius()
	forecastAll := weather.GetForecast(query)
	if len(cwd.Weather) == 0 || len(forecastAll.Data) == 0 {
		fmt.Fprintln(w, "An Error occurred")
		return
	}
	forecastTemplate.Execute(w, data{
		Ort:          query,
		Datum:        tString(cwd.Dt),
		Uhrzeit:      "12:00",
		Wetterlage:   cwd.Weather[0].Main,
		WetterDesc:   "Beschreibung per Hand Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
		PhysGroessen: "Temperatur: ...",
		Legende:      "Legende",
		Time:         nowString(),
		Geo:          geoData{Lat: cwd.Coord["lat"], Lon: cwd.Coord["lon"]},
		Cwd:          cwd,
		Fwd: func(f weather.ForecastData) printFwd {
			var nice []printFwdPoint
			for _, v := range f.Data {
				nice = append(nice, printFwdPoint{
					Time:  v.Time,
					Stamp: tString(time.Unix(v.Time, 0).Local().Unix()),
					C:     fmt.Sprintf("%.2f", weather.Ktoc(v.Main.TempK)),
					CMax:  fmt.Sprintf("%.2f", weather.Ktoc(v.Main.TempMaxK)),
					CMin:  fmt.Sprintf("%.2f", weather.Ktoc(v.Main.TempMinK)),
					Pres:  fmt.Sprintf("%.0f", v.Main.Pressure),
					Humid: fmt.Sprintf("%d", v.Main.Humidity),
					Icon:  v.Weather[0].Icon,
					Main:  v.Weather[0].Main,
					Desc:  v.Weather[0].Description,
					Cloud: int(float64(v.Clouds.All) * 0.08),
					WindS: fmt.Sprintf("%.2f", v.Wind.Speed),
					WindD: int(v.Wind.Degree),
					RainA: fmt.Sprintf("%.0f", v.Rain.Amount),
				})
			}
			return printFwd{
				Raw: f,
				N:   nice,
			}
		}(forecastAll.Filter(weather.MIDDAY | weather.EVENING | weather.MORNING)),
		RFwd:    forecastAll,
		Nw:      niceWeatherFromData(cwd),
		MapsKey: *mapskey,
	})
}

// ShowNoCache calls Show with new template
func ShowNoCache(w http.ResponseWriter, r *http.Request) {
	forecastTemplate, _ = template.ParseFiles(fcTmpl)
	fmt.Println("Re-Caching")
	Show(w, r)
}
