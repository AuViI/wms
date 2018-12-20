package weather

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// s_rf64 Updates r according to p if necessary
func s_rf64(r *RangeFloat64, p float64) bool {
	extreme := false
	if p < r.Min {
		r.Min = p
		extreme = true
	}
	if p > r.Max {
		r.Max = p
		extreme = true
	}
	return extreme
}

func s_ri64(r *RangeInt64, p int64) bool {
	extreme := false
	if p < r.Min {
		r.Min = p
		extreme = true
	}
	if p > r.Max {
		r.Max = p
		extreme = true
	}
	return extreme
}

func s_rf64both(r *RangeFloat64, p float64) {
	r.Min = p
	r.Max = p
}

func s_ri64both(r *RangeInt64, p int64) {
	r.Min = p
	r.Max = p
}

func Summary(data []DataPoint) DataSummary {
	var ds DataSummary
	ds.Wind.Degree.Num = int64(len(data))
	if len(data) >= 1 {
		v := data[0]
		s_ri64both(&(ds.Time), v.Time)
		s_rf64both(&(ds.TempK), v.Main.TempK)
		s_rf64(&(ds.TempK), v.Main.TempMinK)
		s_rf64(&(ds.TempK), v.Main.TempMaxK)
		s_rf64both(&(ds.Pressure), v.Main.Pressure)
		s_rf64both(&(ds.Humidity), float64(v.Main.Humidity))
		s_ri64both(&(ds.Clouds), int64(v.Clouds.All))
		s_rf64both(&(ds.Wind.Speed), v.Wind.Speed)
		ds.Wind.Degree.Average = v.Wind.Degree
		s_rf64both(&(ds.Rain), v.Rain.Amount)
		ds.Stats = append(ds.Stats, v.Weather[0])
		ds.Median = ds.Stats[0]
		ds.Median.Icon = ds.Median.Icon[:len(ds.Median.Icon)-1] + "d"
	}
	if len(data) > 1 {
		for _, v := range data[1:] {
			s_ri64(&(ds.Time), v.Time)
			s_rf64(&(ds.TempK), v.Main.TempK)
			s_rf64(&(ds.TempK), v.Main.TempMinK)
			s_rf64(&(ds.TempK), v.Main.TempMaxK)
			s_rf64(&(ds.Pressure), v.Main.Pressure)
			s_rf64(&(ds.Humidity), float64(v.Main.Humidity))
			s_ri64(&(ds.Clouds), int64(v.Clouds.All))
			s_rf64(&(ds.Wind.Speed), v.Wind.Speed)
			ds.Wind.Degree.Average += v.Wind.Degree
			s_rf64(&(ds.Rain), v.Rain.Amount)
			ds.Stats = append(ds.Stats, v.Weather[0])
		}
		ds.Wind.Degree.Average = float64(int64(ds.Wind.Degree.Average)%360) / float64(ds.Wind.Degree.Num)
		ids := make(map[int]int)
		mains := make(map[string]int)
		descs := make(map[string]int)
		icons := make(map[string]int)
		s_inc := func(m map[string]int, k string) {
			var ok bool
			if _, ok = m[k]; ok {
				m[k] += 1
			} else {
				m[k] = 1
			}
		}
		i_inc := func(m map[int]int, k int) {
			var ok bool
			if _, ok = m[k]; ok {
				m[k] += 1
			} else {
				m[k] = 1
			}
		}
		for _, s := range ds.Stats {
			i_inc(ids, s.ID)
			s_inc(mains, s.Main)
			s_inc(descs, s.Description)
			s_inc(icons, s.Icon[:len(s.Icon)-1])
		}
		var m_id, m_main, m_desc, m_icon int
		for _, s := range ds.Stats {
			if ids[s.ID] > m_id {
				m_id = ids[s.ID]
				ds.Median.ID = s.ID
			}
			if mains[s.Main] > m_main {
				m_main = mains[s.Main]
				ds.Median.Main = s.Main
			}
			if descs[s.Description] > m_desc {
				m_desc = descs[s.Description]
				ds.Median.Description = s.Description
			}
			if icons[s.Icon[:len(s.Icon)-1]] > m_icon {
				m_icon = icons[s.Icon[:len(s.Icon)-1]]
				ds.Median.Icon = s.Icon[:len(s.Icon)-1]
			}
		}
		ds.Median.Icon += "d"
	}
	return ds
}

type DailyForecast struct {
	TodayDate   Date
	Today       DataSummary
	Forecast    DayDataSummaryMap
	Description ForecastData
	location    *time.Location
}

func GetDailyForecast(loc string, date Date, days int) DailyForecast {
	var df DailyForecast
	df.location = time.Local

	df.TodayDate = date
	df.Description = *GetForecast(loc)

	fc := df.Description.SplitByDay(df.location).Summarize() // Change to be actual local time

	df.Today = fc[date]
	df.Forecast = make(DayDataSummaryMap)
	next := date.Tomorrow(df.location)
	for i := 0; i < days; i++ {
		df.Forecast[next] = fc[next]
		next = next.Tomorrow(df.location)
	}

	return df
}

func (df DailyForecast) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("DailyForecast on %s for %d(+1) days", df.TodayDate, len(df.Forecast)))
	lines = append(lines, fmt.Sprintf("'Today': %s", df.Today))
	for k := 1; k <= len(df.Forecast); k++ {
		v, d, err := df.Get(k)
		if err != nil {
			lines = append(lines, err.Error())
		}
		lines = append(lines, fmt.Sprintf("'Forecast(%s)': %s", d, v))
	}

	return strings.Join(lines, "\n")
}

func (df DailyForecast) Get(day int) (DataSummary, Date, error) {
	if day < 0 {
		return DataSummary{}, Date{}, errors.New("day out of range (under)")
	} else if day == 0 {
		return df.Today, df.TodayDate, nil
	} else if day <= len(df.Forecast) {
		d := df.TodayDate.Add(day, df.location)
		return df.Forecast[d], d, nil
	} else {
		return DataSummary{}, Date{}, errors.New("day out of range (over)")
	}
}

func (df DailyForecast) GetForecastArray(days int) []DataSummary {
	var ds []DataSummary
	for i := 1; i <= days; i++ {
		next, _, err := df.Get(i)
		if err != nil {
			panic(err)
		}
		ds = append(ds, next)
	}
	return ds
}

func (df DailyForecast) GetStartDate(ds DataSummary) Date {
	return DateFromUnix(ds.Time.Min, df.location)
}
