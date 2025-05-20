GIT_COMMIT=$(shell git describe --always)

all: build
default: build

build:
	go build -o chatgpt-telegram ./cmd/server

clean:
	rm chatgpt-telegram
