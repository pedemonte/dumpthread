build: 
	go build -o dumpthread main.go

clean:
	rm dumpthread

install:
	go install github.com/pedemonte/dumpthread

.PHONY:
	clean install
