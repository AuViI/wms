default:
	rm wms || :
	go build .

deps:
	go get github.com/mtib/simplehttp
	git submodule init

run: default
	./wms -no-cache
