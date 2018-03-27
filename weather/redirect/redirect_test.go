package redirect

import "testing"

type redirectPair struct {
	src string
	dst string
	isr bool
}

var obk = "Ostseebad Kühlungsborn, DE"
var redirectTests = []redirectPair{
	{"Kuebo", obk, true},
	{obk, obk, false},
	{"Kborn", obk, true},
	{"London", "London", false},
}

func TestRedirectPairs(t *testing.T) {
	for _, v := range redirectTests {
		if IsRedirected(v.src) {
			if !v.isr {
				t.Errorf("\"%v\" falsely redirected", v.src)
			}
			redirected := Redirect(v.src)
			if redirected != v.dst {
				t.Errorf("\"%v\" redirected to \"%v\" instead of \"%v\"", v.src, redirected, v.dst)
			}
		} else {
			if v.isr {
				t.Errorf("\"%v\" should be redirected to \"%v\"", v.src, v.dst)
			}
		}
	}
}

func TestNoCycles(t *testing.T) {
	for src := range redirects {
		num, err := CountRedirects(src)
		if err != nil {
			t.Error(err)
		}
		if num == 0 {
			t.Error("num is 0, but error is nil")
		}
	}
}

func TestCyclesPrint(t *testing.T) {
	printRedirectLength()
}
