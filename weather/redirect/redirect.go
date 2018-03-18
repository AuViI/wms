package redirect

import (
    "fmt"
    "github.com/deckarep/golang-set"
)

var redirects = map[string]string{
    "Kborn": "Kuebo",
    "Kuebo": "Kübo",
    "Kübo": "Kühlungsborn",
    "Kuehlungsborn": "Kühlungsborn",
    "Kühlungsborn": "Ostseebad Kühlungsborn",
    "Ostseebad Kühlungsborn": "Ostseebad Kühlungsborn, DE",
}

const debugPrint = false

func IsRedirected(city string) bool {
    return redirects[city] != redirects["not-redirected"]
}

func Redirect(city string) string {
    if !IsRedirected(city) {
        return city
    }
    red := redirectStep(city)
    if debugPrint {
        fmt.Printf("\"%s\"", city)
    }
    if IsRedirected(red) && redirects[red] != red {
        if debugPrint {
            fmt.Printf("\n\t-> ")
        }
        return Redirect(red)
    }
    if debugPrint {
        fmt.Printf("\n\t-> \"%s\"\n", red)
    }
    return red
}

func CountRedirects(city string) (int, error) {
    cur := city
    visited := mapset.NewSet()
    count := 0
    for IsRedirected(cur) {
        if visited.Contains(cur) {
            return 0, RedirectCycleError{city, count}
        }
        visited.Add(cur)
        cur = redirectStep(cur)
        count++
    }
    return count, nil
}

func LongestRedirect() (int, error) {
    max := 0
    for src, _ := range redirects {
        n, err := CountRedirects(src)
        if err != nil {
            return n, err
        }
        if max < n {
            max = n
        }
    }
    return max, nil
}

func redirectStep(city string) string {
    if !IsRedirected(city) {
        return city
    }
    return redirects[city]
}

func printRedirectLength() {
    for src, _ := range redirects {
        num, err := CountRedirects(src)
        if err != nil {
            fmt.Printf("%s: \tinfinity\n", src)
        } else {
            fmt.Printf("%s: \t%d\n", src, num)
        }
    }
}

// cycle error
type RedirectCycleError struct {
    Start string
    MaxStep int
}

func (e RedirectCycleError) Error() string {
    return fmt.Sprintf("%v redirects %v times before cycling", e.Start, e.MaxStep)
}

