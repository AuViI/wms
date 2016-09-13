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

// Data contains all Data needed for filling the Template
type Data struct {
	Ort          string
	Datum        string
	Uhrzeit      string
	Wetterlage   string
	WetterDesc   string
	PhysGroessen string
	Legende      string
	Time         string
	Geo          GeoData
	Cwd          *weather.Data
	Fwd          PrintFwd
	RFwd         *weather.ForecastData
	Nw           NiceWeather
	MapsKey      string
}

// PrintFwd contains the Raw-Filtered weather data and formatted txt
type PrintFwd struct {
	Raw weather.ForecastData
	N   []PrintFwdPoint
}

// PrintFwdPoint is a single formatted data point of PrintFwd
type PrintFwdPoint struct {
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

// GeoData contains Lat and Lon
type GeoData struct {
	Lat float64
	Lon float64
}

// NiceWeather is PrintFwdPoint for current weather
type NiceWeather struct {
	Temp    string
	TempMax string
	TempMin string
}

// NiceWeatherFromData converts kelvin to celsius
func NiceWeatherFromData(w *weather.Data) NiceWeather {
	return NiceWeather{
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
		}
		//w.Header().Set("Refresh", "10") // Sekunden
		w.Header().Set("Cache-Control", "max-age=600")
		query = url
	}
	cwd := weather.GetCurrent(query).ConvertToCelsius()
	forecastAll := weather.GetForecast(query)
	forecastTemplate.Execute(w, Data{
		Ort:          query,
		Datum:        tString(cwd.Dt),
		Uhrzeit:      "12:00",
		Wetterlage:   "Sonnig",
		WetterDesc:   "Beschreibung per Hand Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
		PhysGroessen: "Temperatur: ...",
		Legende:      "Legende",
		Time:         nowString(),
		Geo:          GeoData{Lat: cwd.Coord["lat"], Lon: cwd.Coord["lon"]},
		Cwd:          cwd,
		Fwd: func(f weather.ForecastData) PrintFwd {
			var nice []PrintFwdPoint
			for _, v := range f.Data {
				nice = append(nice, PrintFwdPoint{
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
					Cloud: v.Clouds.All,
					WindS: fmt.Sprintf("%.2f", v.Wind.Speed),
					WindD: int(v.Wind.Degree),
					RainA: fmt.Sprintf("%.0f", v.Rain.Amount),
				})
			}
			return PrintFwd{
				Raw: f,
				N:   nice,
			}
		}(forecastAll.Filter(weather.MIDDAY | weather.EVENING | weather.MORNING)),
		RFwd:    forecastAll,
		Nw:      NiceWeatherFromData(cwd),
		MapsKey: *mapskey,
	})
}

// ShowNoCache calls Show with new template
func ShowNoCache(w http.ResponseWriter, r *http.Request) {
	forecastTemplate, _ = template.ParseFiles(fcTmpl)
	fmt.Println("Re-Caching")
	Show(w, r)
}
