package weather

import "time"

// FilterFunc can be used to filter a set of DataPoints
type FilterFunc func(*DataPoint) bool

// FilterMode should not be initialized, use the Constants!
// you can use | or & if you want to customize your call of Filter
type FilterMode int

const (
	// MIDDAY is 12:00
	MIDDAY FilterMode = 00001
	// MORNING is 06:00
	MORNING FilterMode = 00010
	// MIDNIGHT is 00:00
	MIDNIGHT FilterMode = 00100
	// MIDS is 12:00 and 00:00
	MIDS FilterMode = 00101
	// EVENING is 19:00
	EVENING FilterMode = 01000
	// INTERVAL6 is all the above
	INTERVAL6 FilterMode = 01111
)

func (i FilterMode) is(mode FilterMode) bool {
	return i&mode != 0
}

// Filter filters the ForecastData to only contain certain DataPoints
// TODO: there should be a version of Filter returning a channel
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
