package temp

import (
    "testing"
    "fmt"
)

var stringTests = []float64{36.5, 0, -273.15, 100}

func TestAdder(t *testing.T) {
    var t1, t2 Temperature
    t1 = FromK(0)
    t2 = FromC(0)
    t1.Add(t2)
    if t1.C != t2.C {
        t.Errorf("%v != %v\n", t1.C, t2.C)
    }
}

func TestStringer(t *testing.T) {
    for _, v := range stringTests {
        var tem Temperature
        tem = FromC(v)
        fmt.Printf("%v \t-> %v\n", v, tem)
    }
}
