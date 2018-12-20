package main

import (
	"fmt"
	"github.com/AuViI/wms/weather"
	"net/http"
	"time"
)

func debug(w http.ResponseWriter, r *http.Request) {
	/*
		forecast := weather.GetForecast("Kühlungsborn")
		data := forecast.Data
		sum := weather.Summary(data)
		split := forecast.SplitByDay(time.Local)
		compiled := split.Summarize()
		fmt.Fprintf(w, "%s\n%s\n", sum, sum.TempK.Temperature())
		fmt.Fprintf(w, "%v\n", split)
		fmt.Fprintf(w, "%v\n", compiled)
	*/
	fmt.Fprintf(w, "%v\n", weather.GetDailyForecast("Kübo", weather.NowDate(time.Local), 3))
}
