package weather

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
	} else if len(data) > 1 {
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
		}
		ds.Wind.Degree.Average = float64(int64(ds.Wind.Degree.Average)%360) / float64(ds.Wind.Degree.Num)
	}
	return ds
}
