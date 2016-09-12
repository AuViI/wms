package main

import "github.com/auvii/wms/weather"

func PrognoseTxt(city string, n int) string {
	c := weather.GetCurrent(city)
	f := weather.GetForecast(city)
	return f.Header() + c.WeatherString() + f.NForecast(n)
}
