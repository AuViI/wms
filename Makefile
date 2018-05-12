default:
	rm wms || :
	go build .

install: deps
	echo "TODO: install"

deps:
	go get github.com/mtib/simplehttp
	git submodule init
	git submodule update --recursive

debug: default
	./wms -no-cache

run: default
	./wms
