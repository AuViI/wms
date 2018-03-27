package txt

import "github.com/AuViI/wms/weather"

// PrognoseTxt returns `aktuelltxt` style text for `city` for `n` Days
func PrognoseTxt(city string, n int) string {
	c := weather.GetCurrent(city)
	f := weather.GetForecast(city)
	return f.Header() + c.WeatherString() + f.NForecast(n)
}
