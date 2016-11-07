package weather

import (
	"bytes"
	"encoding/json"
	"flag"
	"html/template"
	"os"
	"strings"

	"github.com/mtib/simplehttp"
)

var (
	key = flag.String("key", os.Getenv("OWM"), "OpenWeatherMap Key $OWM")
)

type (
	// Query is all you need to know to query OWM for Data
	Query struct {
		City string
		Key  string
	}
)

var (
	currenturl, _ = template.New("current").Parse("http://api.openweathermap.org/data/2.5/weather?q={{.City | urlquery}}&appid={{.Key}}")
	forcasturl, _ = template.New("forcast").Parse("http://api.openweathermap.org/data/2.5/forecast?q={{.City | urlquery}}&appid={{.Key}}")
)

const (
	headline = "YYYY-MM-DD HH:MM:SS TEMP MIN HUMID WINDGRAD FORCE RAIN CLOUDCOVER\n"
)

type (
	// TODO: type for Main, Weather, Wind, Coord

	// Data is the god-object
	Data struct {
		Coord   map[string]float64
		Weather []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Main struct {
			Temp     float64
			Pressure float64
			Humidity float64
			TempMin  float64 `json:"temp_min"`
			TempMax  float64 `json:"temp_max"`
		}
		Wind struct {
			Speed float64
			Deg   float64
		}
		Clouds struct {
			All float64
		}
		Rain struct {
			Volume float64 `json:"3h"`
		}
		Snow struct {
			Volume float64 `json:"3h"`
		}
		Dt int64
	}
	// ForecastData contains all the forecast data
	ForecastData struct {
		City struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Coord struct {
				Lon float64 `json:"lon"`
				Lat float64 `json:"lat"`
			}
			Country    string `json:"country"`
			Population int    `json:"population"`
		}
		Cod     string      `json:"cod"`
		Cnt     int         `json:"cnt"`
		Message float64     `json:"message"`
		Data    []DataPoint `json:"list"`
	}
	// DataPoint is part of data /t
	DataPoint struct {
		Time int64 `json:"dt"`
		Main struct {
			TempK       float64 `json:"temp"`
			TempMinK    float64 `json:"temp_min"`
			TempMaxK    float64 `json:"temp_max"`
			Pressure    float64 `json:"pressure"`
			SeaLevel    float64 `json:"sea_level"`
			GroundLevel float64 `json:"grnd_level"`
			Humidity    int     `json:"humidity"`
			TempKfK     float64 `json:"temp_kf"`
		} `json:"main"`
		Weather []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Clouds struct {
			All int `json:"all"`
		} `json:"clouds"`
		Wind struct {
			Speed  float64 `json:"speed"`
			Degree float64 `json:"deg"`
		} `json:"wind"`
		Rain struct {
			Amount float64 `json:"3h"`
		} `json:"rain"`
		TimeString string `json:"dt_txt"`
	}
)

// Valid if data was returned
func (f ForecastData) Valid() bool {
	return f.Cnt != 0
}

// MphToBf converts mp/h to Bf.
func MphToBf(mph float64) float64 {
	switch {
	case mph < 1:
		return 0
	case mph <= 3:
		return 1
	case mph <= 7:
		return 2
	case mph <= 12:
		return 3
	case mph <= 18:
		return 4
	case mph <= 24:
		return 5
	case mph <= 31:
		return 6
	case mph <= 38:
		return 7
	case mph <= 46:
		return 8
	case mph <= 54:
		return 9
	case mph <= 63:
		return 10
	case mph <= 72:
		return 11
	case mph > 72:
		return 12
	}
	return -1
}

func fillTemlp(t *template.Template, c string) string {
	var b bytes.Buffer
	t.Execute(&b, &Query{strings.Replace(c, " ", "_", -1), *key})
	return b.String()
}

// GetCurrent returns filled `Data` for `city`
func GetCurrent(city string) *Data {
	var wd *Data
	wd = &Data{}
	answ, _ := simplehttp.GetResponseBody(fillTemlp(currenturl, city))
	json.Unmarshal(answ, wd)
	return wd
}

// GetForecast from OpenWeatherMap
func GetForecast(city string) *ForecastData {
	data := new(ForecastData)
	jdata, err := simplehttp.GetResponseBody(fillTemlp(forcasturl, city))
	if err != nil {
		panic(err)
	}
	json.Unmarshal(jdata, data)
	return data
}

func ktoc(k interface{}) float64 {
	return k.(float64) - 272.15
}

// Ktoc converts Kelvin to Celsius
func Ktoc(k interface{}) float64 {
	return k.(float64) - 272.15
}
