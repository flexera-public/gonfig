.PHONY: all gen build run

all: gen build
	./test

gen: install
	gonfig -c config.json -o config.go

build:
	go build -o test

install:
	go install ../.
