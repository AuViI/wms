# Weather Monitoring Service

WMS Management Server

## Usage

```bash
# compilation
go get github.com/AuViI/wms
cd $GOPATH/src/github.com/AuViI/wms
make deps       #init submodules

# execution depends on settings
$GOPATH/bin/wms #custom
make run        #default
make debug      #testing


# TODO execution modes
make offline #use generator for data
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

[future]
-caches server
    send queries to server instead of OWM to receive
    optimized answers
```

### Statically served resources

Binary resources can be server from the `./resource/` folder, to enable serving them they have to be loaded into the resource map in `web.go`. This is to guarantee they are read and served correctly.

### Template Normalisation

To enable theming of all templates universally a standart has to be created. To signal the template should have the ending `.theme.gtmpl`.

The theming information is given to the template as such:

```golang
// .Theme
struct Theme {
    MainColor   ThemeColor
    SecondColor ThemeColor
    Icon        ThemeIcon
    Client      ThemeClient
}

struct ThemeClient {
    Name string
    Link string
}

struct ThemeIcon {
    Binary []byte
    Link   string
    func (i Icon) Img() string {} // generate <img> tag
}

```

This still has to be implemented

### Configuration

- [ ] Config location via CMD
- [ ] Using config

There is a configuration file written in YAML, used to configure shown examples,
and locations for files. These include maps to client icons.

### Form generator

To create links to give out to clients there will be a generator, which will be created from a yaml describing what kind of information is needed to generate a link for each specific output.

## License

Currently being decided on. Please be patient or contact us if
you want to use this software already.

## Contribution

Merging of pull requests is very unlikely, so please submit an
issue if you find a bug or have a recommendation.
