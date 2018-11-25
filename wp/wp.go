package wp

import (
	"fmt"
	"github.com/AuViI/wms/weather/redirect"
	"net/http"
	"regexp"
)

/*

Requests:
	/<location>/[<year>/<month>/<day>][/<args>]

*/

type Request struct {
	location string
	year     string
	month    string
	day      string
}

type Response struct {
	content string
}

var reqx, reqerr = regexp.Compile("/wp/([^/]*)(/([^/]*)/([^/]*)/([^/]*))?")

func Handler(w http.ResponseWriter, r *http.Request) {
	if reqerr != nil {
		fmt.Fprintln(w, reqerr)
	} else {
		match := reqx.FindAllStringSubmatch(r.URL.Path, 1)
		wpr := Request{redirect.Redirect(match[0][1]), match[0][3], match[0][4], match[0][5]}
		if r.Method == "POST" {
			content := r.PostFormValue("content")
			PostDatabaseEntry(&wpr, &Response{content})
			return
		}
		fmt.Fprintf(w, "Request for: %s\n", wpr.String())
		fmt.Fprintf(w, "isDate: %v\n", wpr.isDate())
		wrs, err := wpr.GetDatabaseEntry()
		if err != nil {
			fmt.Fprintf(w, "Database: %s\n", err)
		} else {
			fmt.Fprintf(w, "Database: %v\n", wrs)
		}
	}
}

func (wpr *Request) isDate() bool {
	return wpr.year != "" && wpr.month != "" && wpr.day != ""
}

func (wpr *Request) String() string {
	if wpr.isDate() {
		return fmt.Sprintf("wpr(%s){%s.%s.%s}", wpr.location, wpr.day, wpr.month, wpr.year)
	}
	return fmt.Sprintf("wpr(%s){now}", wpr.location)
}

func (wpr *Request) GetDatabaseEntry() (*Response, error) {
	db, oerr := createDB("wp.db")
	if oerr != nil {
		return nil, oerr
	}
	defer closeDB(db)
	return getDB(db, wpr)
}

func PostDatabaseEntry(r *Request, c *Response) error {
	db, oerr := createDB("wp.db")
	if oerr != nil {
		return oerr
	}
	defer closeDB(db)
	return setDB(db, r, c)
}
