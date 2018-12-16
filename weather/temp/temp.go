package temp

import "fmt"

type Temperature struct {
	C float64
	F float64
	K float64
}

type Range struct {
	Min, Max Temperature
}

func FromC(c float64) Temperature {
	return Temperature{c, CtoF(c), CtoK(c)}
}

func FromF(f float64) Temperature {
	return Temperature{FtoC(f), f, FtoK(f)}
}

func FromK(k float64) Temperature {
	return Temperature{KtoC(k), KtoF(k), k}
}

func FtoC(f float64) float64 {
	return (f - 32.0) * 5.0 / 9.0
}

func FtoK(f float64) float64 {
	return CtoK(FtoC(f))
}

func CtoF(c float64) float64 {
	return (c * 9.0 / 5.0) + 32.0
}

func CtoK(c float64) float64 {
	return c + 273.15
}

func KtoF(k float64) float64 {
	return CtoF(KtoC(k))
}

func KtoC(k float64) float64 {
	return k - 273.15
}

func (t *Temperature) Check() bool {
	return t.K >= 0
}

func (t *Temperature) Add(h Temperature) {
	t.AddC(h.K)
}

// Zero returns the Temperature neutral to addition
func Zero() Temperature {
	return FromK(0)
}

func (t *Temperature) AddC(c float64) {
	t.C = t.C + c
	t.F = CtoF(t.C)
	t.K = CtoK(t.C)
}

func (t *Temperature) AddF(f float64) {
	t.F = t.F + f
	t.C = FtoC(t.F)
	t.K = FtoK(t.F)
}

func (t *Temperature) AddK(k float64) {
	t.AddC(k)
}

func (t Temperature) String() string {
	return fmt.Sprintf("(%v°C %v°F %v°K)", t.C, t.F, t.K)
}
