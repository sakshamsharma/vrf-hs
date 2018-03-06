all: *.go
	make shared

shared: *.go
	go build -buildmode=c-shared -o libvrf.so $^

archive: *.go
	go build -buildmode=c-archive -o libvrf.a $^
