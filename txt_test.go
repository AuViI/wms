package main

import (
	"fmt"
	"testing"
)

var testOrte = []string{
	"New York",
	"Oslo",
	"Berlin",
	"Kühlungsborn",
	"Rostock",
}

func TestBareTXT(t *testing.T) {
	for _, v := range testOrte {
		fmt.Println(PrognoseTxt(v, DaysForecastTxt))
	}
}

// Make sure to have a valid OpenWeatherMap.org Key
// in the environment Variable $OWM
func BenchmarkBareTXT(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PrognoseTxt("New York", DaysForecastTxt)
	}
}
