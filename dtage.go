package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"
)

type (
	ntag struct {
		N    uint
		Ort  string
		Data []row
	}
	row struct {
		Name string
		Data []string
		Bold bool
		Unit string
	}
	// RowError is not nil if row cannot be added to ntag struct
	RowError string
)

// row[0] == {"", "Heute", "Morgen", "..."}

const (
	dtaghtml = "./template/dtage.html"
	htmlC    = "<sup>o</sup>C"
)

var (
	dtagetmpl, _ = template.ParseFiles(dtaghtml)
)

func (n ntag) Now() string {
	return time.Now().String()
}

func (n *ntag) Row(name string, data []string, bold bool, unit string) error {
	if n.N != uint(len(data)) {
		return RowError("Row lengths differ")
	}
	n.Data = append(n.Data, row{Name: name, Data: data, Bold: bold, Unit: unit})
	return nil
}

func getTagName(i uint) string {
	switch i {
	case 0:
		return "Heute"
	case 1:
		return "Morgen"
	default:
		return fmt.Sprintf("in %d Tagen", i)
	}
}

func newNTage(n uint, ort string) *ntag {
	fr := make([]string, n)
	for i := 0; i < int(n); i++ {
		fr[i] = getTagName(uint(i))
	}
	return &ntag{
		N:   n,
		Ort: ort,
		Data: []row{
			row{
				Name: "",
				Data: fr,
				Bold: true,
				Unit: "",
			},
		},
	}
}

func (r RowError) Error() string {
	return fmt.Sprintf("RowError: %s", r)
}

func handleDTage(w http.ResponseWriter, r *http.Request) {
	exampleData := newNTage(3, "Braunschweig")
	exampleData.Row("Temperatur", []string{"10", "20", "30"}, false, htmlC)
	dtagetmpl.Execute(w, exampleData)
}

func ncHandleDTage(w http.ResponseWriter, r *http.Request) {
	dtagetmpl, _ = template.ParseFiles(dtaghtml)
	handleDTage(w, r)
}
