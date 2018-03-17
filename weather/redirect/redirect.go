package redirect

import (
    "fmt"
)

var redirects = map[string]string{
    "Kühlungsborn": "Ostseebad Kühlungsborn, DE",
    "Osteseebad Kühlungsborn": "Ostseebad Kühlungsborn, DE",
}

func IsRedirected(city string) bool {
    return redirects[city] != redirects["not-redirected"]
}

func Redirect(city string) string {
    red := redirects[city]
    //fmt.Printf("%s -> %s\n", city, red)
    if IsRedirected(red) && redirects[red] != red {
        return Redirect(red)
    }
    return red
}
