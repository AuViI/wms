default:
	rm wms || :
	go build .

deps:
	go get github.com/mtib/simplehttp
	git submodule init
	git submodule update --recursive

run: default
	./wms -no-cache
