package weather

import (
	"fmt"
	"github.com/AuViI/wms/weather/temp"
)

func rangeStringer(min, max interface{}) string {
	return fmt.Sprintf("(%v to %v)", min, max)
}

func (rf64 RangeFloat64) String() string {
	return rangeStringer(rf64.Min, rf64.Max)
}

func (rf64 RangeFloat64) Temperature() temp.Range {
	return temp.Range{Min: temp.FromK(rf64.Min), Max: temp.FromK(rf64.Max)}
}

func (ri64 RangeInt64) String() string {
	return rangeStringer(ri64.Min, ri64.Max)
}

func (af64 AvgFloat64) String() string {
	return fmt.Sprintf("%f (num:%d)", af64.Average, af64.Num)
}

func (ds DataSummary) String() string {
	return fmt.Sprintf(
		"{\n\tTime: %s,\n\tTemp: %s,\n\tPressure: %s,\n\tHumidity: %s,\n\tClouds: %s,\n\tWind.Speed: %s,\n\tWind.Degree: %s,\n\tRain: %s\n\tMedian: %s\n}",
		ds.Time, ds.TempK, ds.Pressure, ds.Humidity, ds.Clouds, ds.Wind.Speed, ds.Wind.Degree, ds.Rain, ds.Median)
}

func (ds DataSummary) Celsius() RangeFloat64 {
	return RangeFloat64{Min: temp.KtoC(ds.TempK.Min), Max: temp.KtoC(ds.TempK.Max)}
}
