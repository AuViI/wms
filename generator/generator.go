/// generator package is used to generate "valid"
/// json and other data to test the application
/// without needing an API key or internet connectivity.
/// Different generator functions assure some
/// integrety or allow assumptions about data.
/// All generator methods should be kept completely
/// disconnected from all `weather` package code, to
/// avoid too simple tests.
///
/// This package is allowed to return chached or
/// live data as long as it fulfils the necessary
/// criteria.
package generator

/*
# TODOs:
- [ ] generate any json
- [ ] generate json with correct format
- [ ] generate for specified location
- [ ] generate for specified coordinates
- [ ] generate for specified weather
- [ ] generate with intentional errors
- [ ] generate difficult but valid json
- [ ] mutate valid json to contain errors
- [ ] simulate invalid caching behavior
- [ ] simulate connectivity problems
*/

import "fmt"

/// Generate returns any valid json
///
/// TODO return parsable format for `weather`
func Generate() string {
	return fmt.Sprintf("{message: '%s'}\n", "TODO")
}
