package httpmod

import (
	"fmt"
	"strings"
	"testing"
)

type ct struct {
	filename string
	mimetype string
}

func TestCommonTypes(t *testing.T) {
	tests := make(chan ct)

	go (func() {
		tests <- ct{"index.html", "text/html"}
		tests <- ct{"index.gohtml", "text/html"}
		tests <- ct{"main.css", "text/css"}
		tests <- ct{"justtext", ""}
		tests <- ct{"justtext.txt", "text/plain"}
		tests <- ct{"image.png", "image/png"}
		close(tests)
	})()

	for inst := range tests {
		gmt := ContentTypeGuess(inst.filename)
		if !strings.HasPrefix(gmt, inst.mimetype) {
			t.Error(fmt.Sprintf("guessed: %s != %s (expected) for %s",
				gmt, inst.mimetype, inst.filename))
		} else {
			fmt.Printf("guess(\033[35m%s\033[0m) == expected(\033[35m%s\033[0m) for '%s'\n", gmt, inst.mimetype, inst.filename)
		}
	}
}
