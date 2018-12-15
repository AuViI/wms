default:
	rm wms || :
	go build .

install: deps
	echo "TODO: install"

deps:
	go get github.com/mtib/simplehttp
	CGO_ENABLED=1 go get github.com/mattn/go-sqlite3
	CGO_ENABLED=1 go install github.com/mattn/go-sqlite3
	go get ./...
	git submodule init
	git submodule update --recursive

debug: default
	./wms -no-cache

run: default
	./wms
