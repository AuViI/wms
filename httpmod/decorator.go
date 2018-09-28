package httpmod

import (
	"fmt"
	"mime"
	"net/http"
	"strings"
	"time"
)

var initError = (func() error {
	mime.AddExtensionType(".gohtml", "text/html")
	mime.AddExtensionType(".tmpl", "text/html")
	return nil
})()

func CacheMaxAge(w http.ResponseWriter, seconds uint64) {
	w.Header().Set(
		"Cache-Control",
		fmt.Sprintf("public, must-revalidate, max-age=%d",
			seconds,
		),
	)
}

func CacheMaxAgeDuration(w http.ResponseWriter, d time.Duration) {
	CacheMaxAge(w, uint64(d.Seconds()))
}

func ContentTypeGuess(filename string) string {
	parts := strings.Split(filename, ".")
	for i := 1; i < len(parts); i++ {
		guess := mime.TypeByExtension(
			fmt.Sprintf(".%s", strings.Join(parts[len(parts)-i:], ".")),
		)
		if guess != "" {
			return guess
		}
	}
	return ""
}

func ContentTypeSet(w http.ResponseWriter, ct string) {
	if ct == "" {
		return
	}
	w.Header().Set("Content-Type", ct)
}

func ContentTypeAssume(w http.ResponseWriter, filename string) {
	ct := ContentTypeGuess(filename)
	ContentTypeSet(w, ct)
}
