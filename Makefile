.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --no-builtin-rules
VERSION=`cat version`
BUILD=`date +%FT%T%z`
#COMMIT=`git rev-parse HEAD`
COMMIT=`date +%FT%T%z`
EXECUTABLE="servo"

all: build

.PHONY: build test clean generate dist init

build: 
	@go build -ldflags "-X 'main.Version=${VERSION}' -X 'main.Build=${BUILD}' -X 'main.AppName=${EXECUTABLE}'" -o ${EXECUTABLE} ./main.go

run: build
	@./${EXECUTABLE}

clean:
	rm -rf main dist/