package main

import (
	"fmt"
	"github.com/AuViI/wms/weather"
	"net/http"
	"time"
)

func debug(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, weather.GetDailyForecast("Kübo", weather.NowDate(time.Local), 3))
}
