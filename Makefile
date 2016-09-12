default:
	rm wms || :
	go build .

deps:
	go get github.com/mtib/simplehttp

run: default
	./wms -no-cache
