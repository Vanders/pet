.PHONY: all

all: main

main: *.go mos6502/*.go
	go build -o pet *.go
