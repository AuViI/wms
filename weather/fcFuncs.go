package weather

import "time"

type FilterFunc func(*DataPoint) bool

type FilterMode int

const (
	MIDDAY    FilterMode = 00001
	MORNING   FilterMode = 00010
	MIDNIGHT  FilterMode = 00100
	MIDS      FilterMode = 00101
	EVENING   FilterMode = 01000
	INTERVAL6 FilterMode = 01111
)

func (i FilterMode) is(mode FilterMode) bool {
	return i&mode != 0
}

func (f ForecastData) Filter(mode FilterMode) ForecastData {
	inrange := func(h, t int) bool {
		dis := h - t
		if dis < 0 {
			dis = -1 * dis
		}
		return dis < 2
	}
	var newData []DataPoint
	for _, v := range f.Data {
		h := time.Unix(v.Time, 0).Local().Hour()
		if MIDDAY.is(mode) && inrange(h, 12) {
			newData = append(newData, v)
			continue
		}
		if MORNING.is(mode) && inrange(h, 6) {
			newData = append(newData, v)
			continue
		}
		if MIDNIGHT.is(mode) && (inrange(h, 0) || inrange(h, 24)) {
			newData = append(newData, v)
			continue
		}
		if EVENING.is(mode) && inrange(h, 19) {
			newData = append(newData, v)
			continue
		}
	}
	f.Data = newData
	return f
}

// FilterByFunc returns a new ForecastData only containing Points
// of Data which evaluate to true by filter
func (f ForecastData) FilterByFunc(filter FilterFunc) ForecastData {
	var newData []DataPoint
	for _, v := range f.Data {
		if filter(&v) {
			newData = append(newData, v)
		}
	}
	f.Data = newData
	return f
}
