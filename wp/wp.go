// wp is short for Weather Prediction.
//
// It contains the natural language
// forecast to be used on other formats
// explaining the weather.
//
// It is connected to a database saving the
// manually entered forecasts, with those
// being automated in future.
package wp

import (
	"errors"
	"fmt"
	"github.com/AuViI/wms/weather"
	"github.com/AuViI/wms/weather/redirect"
	"io"
	"net/http"
	"regexp"
)

/*
Shape of Requests:
	/<location>/[<year>/<month>/<day>][/<args>]
*/

// Request contains all data needed to query wp data.
type Request struct {
	location string
	year     string
	month    string
	day      string
}

// Response reflects the current database content
// for any given Request, it may also contain
// automatically generated content.
type Response struct {
	content string
}

// Matches Request URLs, used in Handler
var reqx, reqerr = regexp.Compile("/wp/([^/]*)(/([^/]*)/([^/]*)/([^/]*))?")

// Hanlder responds to GET and POST requests for natural language forecasts
//
// GET /wp/<location>[/<year>/<month>/<day>]
//		returns saved weather forecast from databank
//		if it doesn't exists, generate one.
//
// POST /wp/<location>/<year>/<month>/<day>
//		(over)writes database entry for given location
//		and date.
func Handler(w http.ResponseWriter, r *http.Request) {

	// Sanity-Check Regex
	if reqerr != nil {
		// Regex doesn't compile
		fmt.Fprintln(w, reqerr)
		return
	}

	if r.URL.Path == "/wp/" {
		if err := listJson(w); err != nil {
			fmt.Println(err)
		}
		return
	}

	// Determine Request
	match := reqx.FindAllStringSubmatch(r.URL.Path, 1)
	wpr := &Request{redirect.Redirect(match[0][1]), match[0][3], match[0][4], match[0][5]}

	if r.Method == "POST" {
		// write to DB on POST
		content := r.FormValue("content")
		if len(content) == 0 {
			fmt.Fprint(w, "Request with empty content denied")
			return
		}
		perr := PostDatabaseEntry(wpr, &Response{content})
		if perr != nil {
			fmt.Fprint(w, perr)
		} else {
			fmt.Fprint(w, "Ok")
		}
		return
	}

	// Respond with data
	wrs, err := wpr.GetDatabaseEntry()
	handlerPrint(w, wpr, wrs, err)
}

func handlerPrint(w io.Writer, req *Request, res *Response, err error) {
	fmt.Fprintf(w, "request: %s\n", req.String())
	fmt.Fprintf(w, "date: %v\n", req.isDate())
	if err != nil {
		fmt.Fprintf(w, "%s\n", err)
	} else {
		fmt.Fprintf(w, "content: %v\n", res.content)
	}

}

func Now(location string) *Request {
	r := new(Request)
	r.location = redirect.Redirect(location)
	return r
}

func For(location string, date weather.Date) *Request {
	r := Now(location)
	r.year = fmt.Sprintf("%d", date.Year)
	r.month = fmt.Sprintf("%d", date.Month)
	r.day = fmt.Sprintf("%d", date.Day)
	return r
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

func (r *Response) String() string {
	return r.content
}

func PostDatabaseEntry(r *Request, c *Response) error {
	if !r.isDate() {
		return errors.New("request needs to contain date")
	}
	db, oerr := createDB("wp.db")
	if oerr != nil {
		return oerr
	}
	defer closeDB(db)
	return setDB(db, r, c)
}
