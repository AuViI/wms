package weather

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"os"
	"strings"

	//"strings"
	"sync"
	"time"

	"github.com/AuViI/wms/weather/redirect"
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
	currenturl, _ = template.New("current").Parse("http://api.openweathermap.org/data/2.5/weather?q={{.City}}&appid={{.Key}}&lang=de")
	forcasturl, _ = template.New("forcast").Parse("http://api.openweathermap.org/data/2.5/forecast?q={{.City}}&appid={{.Key}}&lang=de")
	rpm           = 0
	rpmMutex      sync.Mutex
)

var cacheJSON = make(map[string]Cache)

const (
	headline      = "YYYY-MM-DD HH:MM:SS TEMP MIN HUMID WINDGRAD FORCE RAIN CLOUDCOVER\n"
	currentGeoUrl = "http://api.openweathermap.org/data/2.5/weather?lat=%.1f&lon=%.1f&appid=%s"
	cacheDuration = 66 * 60 // 66 minutes
)

type (
	// TODO: type for Main, Weather, Wind, Coord
	Cache struct {
		CacheDate int64
		Content   []byte
	}

	WeatherStruct struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	}

	// Data is the god-object
	Data struct {
		Coord   map[string]float64
		Weather []WeatherStruct `json:"weather"`
		Main    struct {
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
		Dt      int64
		Cod     interface{} `json:"cod"`
		Message interface{} `json:"message"`
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
		Cod     interface{} `json:"cod"`
		Message interface{} `json:"message"`
		Cnt     int         `json:"cnt"`
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
		Weather []WeatherStruct `json:"weather"`
		Clouds  struct {
			All int `json:"all"`
		} `json:"clouds"`
		Wind struct {
			Speed  float64 `json:"speed"`
			Degree float64 `json:"deg"`
		} `json:"wind"`
		Rain struct {
			Amount float64 `json:"3h"`
		} `json:"rain"`
		Snow struct {
			Amount float64 `json:"3h"`
		} `json:"snow"`
		TimeString string `json:"dt_txt"`
	}

	RangeFloat64 struct {
		Min float64
		Max float64
	}

	RangeInt64 struct {
		Min int64
		Max int64
	}

	AvgFloat64 struct {
		Average float64
		Num     int64
	}

	DataSummary struct {
		Time     RangeInt64
		Stats    []WeatherStruct
		Median   WeatherStruct
		TempK    RangeFloat64
		Pressure RangeFloat64
		Humidity RangeFloat64
		Clouds   RangeInt64
		Wind     struct {
			Speed  RangeFloat64
			Degree AvgFloat64
		}
		Rain RangeFloat64
		Snow RangeFloat64
	}
)

// Valid if data was returned
func (f ForecastData) Valid() bool {
	switch f.Cod.(type) {
	case float64:
		return int(f.Cod.(float64)) == 200
	case string:
		return f.Cod.(string) == "200"
	default:
		return false
	}
}

func (f ForecastData) ToData(dp DataPoint) Data {
	var d Data

	d.Coord = make(map[string]float64)
	d.Coord["lat"] = f.City.Coord.Lat
	d.Coord["lon"] = f.City.Coord.Lon

	d.Weather = make([]WeatherStruct, 1)

	d.Weather[0] = dp.Weather[0]

	d.Main.Temp = dp.Main.TempK
	d.Main.TempMin = dp.Main.TempMinK
	d.Main.TempMax = dp.Main.TempMaxK
	d.Main.Pressure = dp.Main.Pressure
	d.Main.Humidity = float64(dp.Main.Humidity)

	d.Wind.Speed = dp.Wind.Speed
	d.Wind.Deg = dp.Wind.Degree

	d.Clouds.All = float64(dp.Clouds.All)
	d.Rain.Volume = dp.Rain.Amount
	d.Snow.Volume = dp.Snow.Amount
	d.Dt = dp.Time
	d.Cod = f.Cod
	d.Message = f.Message

	return d
}

func GetCacheRaw() *map[string]Cache {
	return &cacheJSON
}

func GetCachedLocations() [][]float64 {
	data := make([][]float64, len(cacheJSON))
	i := 0
	for _, v := range cacheJSON {
		dp := &Data{}
		json.Unmarshal(v.Content, dp)
		data[i] = make([]float64, 3)
		data[i][0] = dp.Coord["lat"]
		data[i][1] = dp.Coord["lon"]
		data[i][2] = float64(v.CacheDate)
		i = i + 1
	}
	return data
}

func age(link string) (int64, bool) {
	uNow := time.Now().Unix()
	cache, exists := cacheJSON[link]
	if exists {
		return uNow - cache.CacheDate, true
	}
	return 0, false
}

func getLink(link string) ([]byte, error) {
	dt, cached := age(link)
	if cached && dt < cacheDuration {
		return cacheJSON[link].Content, nil
	}
	rpmCount()
	answ, err := simplehttp.GetResponseBody(link)
	if err == nil {
		cacheJSON[link] = Cache{Content: answ, CacheDate: time.Now().Unix()}
	}
	return answ, err
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
	//t.Execute(&b, &Query{strings.Replace(c, " ", "_", -1), *key})
	t.Execute(&b, &Query{c, *key})
	return b.String()
}

// GetCurrent returns filled `Data` for `city`
func GetCurrent(city string) *Data {

	if redirect.IsRedirected(city) {
		city = redirect.Redirect(city)
	}

	var wd *Data
	wd = &Data{}
	answ, err := getLink(fillTemlp(currenturl, city))
	if err != nil {
		// unable to connect to OWM
		fmt.Println(err)
		return wd
	}
	//fmt.Println(string(answ))
	json.Unmarshal(answ, wd)
	return wd
}

func GetCurrentByGeo(lat, lon float64) *Data {
	var wd *Data

	answ, err := getLink(fmt.Sprintf(currentGeoUrl, lat, lon, *key))
	if err != nil {
		panic(err)
	}

	wd = new(Data)
	json.Unmarshal(answ, wd)
	return wd
}

// GetForecast from OpenWeatherMap
func GetForecast(city string) *ForecastData {

	if redirect.IsRedirected(city) {
		city = redirect.Redirect(city)
	}

	//fmt.Println(city)

	data := new(ForecastData)
	link := fillTemlp(forcasturl, city)
	jdata, err := getLink(link)

	//fmt.Println(link)
	//fmt.Println(string(jdata))

	if err != nil {
		// unable to connect to OWM
		fmt.Println(err)
		return data
	}

	json.Unmarshal(jdata, data)

	//fmt.Printf("Data: %v\n", data)
	return data
}

func (w WeatherStruct) String() string {
	return fmt.Sprintf("%d('%s' - '%s' - icon: '%s'}", w.ID, w.Main, w.Description, w.Icon)
}

func ktoc(k interface{}) float64 {
	return k.(float64) - 272.15
}

// Ktoc converts Kelvin to Celsius
func Ktoc(k interface{}) float64 {
	return k.(float64) - 272.15
}

type Date struct {
	Day, Month, Year int
}

func (d Date) String() string {
	return fmt.Sprintf("%2d.%02d.%d", d.Day, d.Month, d.Year)
}

func DateFromTime(t time.Time) Date {
	return Date{t.Day(), int(t.Month()), t.Year()}
}

func DateFromUnix(unix int64, in *time.Location) Date {
	t := time.Unix(unix, 0).In(in)
	return DateFromTime(t)
}

func NowDate(in *time.Location) Date {
	return DateFromTime(time.Now().In(in))
}

func (d Date) Tomorrow(in *time.Location) Date {
	return d.Add(1, in)
}

func (d Date) Add(days int, in *time.Location) Date {
	return DateFromTime(d.Time(in).AddDate(0, 0, days))
}

func (d Date) Time(in *time.Location) time.Time {
	return time.Date(d.Year, time.Month(d.Month), d.Day, 0, 0, 0, 0, in)
}

type DayDataPointMap map[Date][]DataPoint

func (d DayDataPointMap) String() string {
	var days []string

	for k, v := range d {
		days = append(days, fmt.Sprintf("%s: %d elements", k, len(v)))
	}

	return strings.Join(days, "\n")
}

type DayDataSummaryMap map[Date]DataSummary

func (d DayDataPointMap) Summarize() DayDataSummaryMap {
	m := make(DayDataSummaryMap)
	for k, v := range d {
		m[k] = Summary(v)
	}
	return m
}

func (fc ForecastData) SplitByDay(in *time.Location) DayDataPointMap {
	m := make(DayDataPointMap)

	for _, v := range fc.Data {
		d := DateFromUnix(v.Time, in)
		if _, ok := m[d]; ok {
			m[d] = append(m[d], v)
		} else {
			m[d] = make([]DataPoint, 1, 12)
			m[d][0] = v
		}
	}

	return m
}

func rpmBuildDown() {
	select {
	case <-time.After(1 * time.Minute):
		rpmMutex.Lock()
		rpm -= 1
		rpmMutex.Unlock()
	}
}

func rpmCount() {
	rpmMutex.Lock()
	rpm += 1
	rpmMutex.Unlock()
	go rpmBuildDown()
	rpmAdjustSleep()
}

func rpmAdjustSleep() {
	rpmMutex.Lock()
	switch {
	case rpm < 10:
		break
	case rpm < 40:
		<-time.After(1 * time.Second)
		break
	case rpm < 50:
		<-time.After(5 * time.Second)
		break
	case rpm < 80:
		<-time.After(10 * time.Second)
		break
	}
	rpmMutex.Unlock()
}
