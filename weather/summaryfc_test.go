package weather

import (
	"fmt"
	"testing"
)

func getDS() DataSummary {
	var ds DataSummary
	ds.Time = RangeInt64{Min: 0, Max: 10}
	ds.TempK = RangeFloat64{Min: 257, Max: 278}
	ds.Pressure = RangeFloat64{Min: 1000, Max: 1100}
	ds.Humidity = RangeFloat64{Min: 30, Max: 70}
	ds.Clouds = RangeInt64{Min: 1, Max: 7}
	ds.Wind.Speed = RangeFloat64{Min: 0, Max: 130}
	ds.Wind.Degree = AvgFloat64{Average: 90, Num: 1}
	ds.Rain = RangeFloat64{Min: 0, Max: 30}
	return ds
}

func TestGetDS(t *testing.T) {
	getDS()
}

func TestDataSummaryPrint(t *testing.T) {
	fmt.Println(getDS())
}

func TestDataSummaryCelsius(t *testing.T) {
	fmt.Println(getDS().Celsius())
}

func TestDataSummaryTemperature(t *testing.T) {
	fmt.Println(getDS().TempK.Temperature())
}
