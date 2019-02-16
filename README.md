# Weather Monitoring Service

WMS Management Server

## Usage

```bash
# compilation
# clone repo into directory under $GOPATH/src
# cd into that directory
make deps       #init submodules
make			#compile application into `wms`

# execution depends on settings
make run        #default
make debug      #testing
```

### CLI

```
-gewusst folder
    folder to use for Q[1-4] files

-help
    display binary specific help text

-key string
    OpenWeatherMap key, uses $OWM if omitted

-maps string
    Google Maps API key, uses $GOOGLEAPI if omitted

-no-cache
    avoids caching template files, displays changes
    on hard refresh in browser

-port number
    port to host webserver on

-render location1[,location2...]
    server will render pictures for all listed locations
    regularly and serve them on the web interface

-wd path
    will change directory into the given path before
    executing and sourcing files
```

### Statically served resources

Binary resources can be server from the `./resource/` folder, to enable serving them they have to be loaded into the resource map in `resource.go`. This is to guarantee they are read and served correctly.

### Configuration

WMS uses multiple YAML files for configuration. One is `$HOME/.wmsrc.yaml` for general behaviour of the application. Saved user configurations are stored in `$HOME/.wmsuser.yaml`. Weather forecasts are saved within a [SQLite Database](https://www.sqlite.org/index.html) next to the binary called `wp.db`. All of those files are created on access if they don't exist with a sensible default configuration.

#### `.wmsrc.yaml`

Default:

```yaml
title: Weather Monitoring System
example:
- Berlin
- Braunschweig
- Frankfurt
- Hamburg
- Holbæk
- Kühlungsborn
- New York
- Oslo
- Rostock
- Tokio
modus: [txt, forecast, list, csv, dtage, view, normlist]
render:
  cites:
  - Frankfurt
  - Köln
  - Kühlungsborn
  - Rostock
  - Warnemünde
  interval: 12
dtage:
- 1/aktuell
- 3/meteo
- 5/meteo
```

## License

Currently being decided on. Please be patient or contact us if
you want to use this software already.
