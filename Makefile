.PHONY: all

all: pet

pet: *.go mos6502/*.go
	go build -o pet *.go

test:
	go test -v ./...
