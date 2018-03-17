package weather

import (
	"fmt"
	"time"
)

// WeatherString returns the first line for /txt/ using current data
func (w *Data) WeatherString() string {
	n := time.Now()
    if (w.Main.TempMin == 0) {
        return fmt.Sprintf("# error! (current) %v: %s\n", w.Cod, w.Message)
        }
	return fmt.Sprintf("%s %s %.1f %.1f %.f %.f %.1f %.1f %.f\n", fmt.Sprint(time.Now())[:10], fmt.Sprintf("%02d:%02d:%02d", n.Hour(), n.Minute(), n.Second()), w.Main.Temp-272.15, w.Main.TempMin-272.15, w.Main.Humidity, w.Wind.Deg, w.Wind.Speed, w.Rain.Volume+w.Snow.Volume, w.Clouds.All*0.08)
}

// Header returns the header for the /txt/ file
func (f *ForecastData) Header() string {
    if f.Valid() {
	    return fmt.Sprintf("# %s {%+5.2f;%+5.2f}\n# %s\n# %s", f.City.Name, f.City.Coord.Lat, f.City.Coord.Lon, time.Now(), headline)
    }
    return fmt.Sprintf("SEVERE ERROR ENCOUNTERED\n# %s\n", time.Now())
}

// ConvertToCelsius copies a new Data packet, but converts Kelvin to Celsius
func (w Data) ConvertToCelsius() *Data {
	w.Main.Temp = ktoc(w.Main.Temp)
	w.Main.TempMax = ktoc(w.Main.TempMax)
	w.Main.TempMin = ktoc(w.Main.TempMin)
	return &w
}

// NForecast returns a string to be appended to /txt/ for `n` days in the future
func (f *ForecastData) NForecast(n int) string {
	n *= 2 // Zwei Daten je Tag
	now := time.Now()
	num := 0

	if !f.Valid() { // keine Daten
        return fmt.Sprintf("# error! (forcast) %v: %s\n", f.Cod, f.Message)
	}

	data := ""
	for _, k := range f.Data {
		if k.TimeString[11:13] == "03" || k.TimeString[11:13] == "15" { // Immer 12:00 Uhr
			if k.TimeString[5:7] == fmt.Sprintf("%02d", int(now.Month())) { // Nicht heute
				num++
				extra := ""
				data += fmt.Sprintf("%s %.1f %.1f %.f %.f %.1f %.1f %.f %s\n", k.TimeString, k.Main.TempK-272.15, k.Main.TempMinK-272.15, float64(k.Main.Humidity), k.Wind.Degree, k.Wind.Speed, k.Rain.Amount, float64(k.Clouds.All)*0.08, extra)
				// DD.MM.YYYY HH:MM:SS TEMP	MIN	HUMID	WINDGRAD	FORCE	RAIN	CLOUDCOVER
				if num == n {
					break
				}
			}
		}
	}
	return data
}
