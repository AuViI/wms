# Weather Monitoring Service

WMS Management Server

## Usage

```bash
# compilation
go get github.com/auvii/wms
cd $GOPATH/src/github.com/auvii/wms
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

## Roadmap

- [x] rendering pictures
- [x] smart caching
- [ ] pretty forcast map
- [ ] user accounts
- [ ] smarter caching (cache server)
- [ ] optimizations
- [ ] >50% test coverage
- [ ] better documentation
- [ ] refactoring (incl. language)

## License

Currently being decided on. Please be patient or contact us if
you want to use this software already.

## Contribution

Merging of pull requests is very unlikely, so please submit an
issue if you find a bug or have a recommendation.
