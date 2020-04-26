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

.PHONY: build test clean generate dist init build_linux build_mac

build: 
	@go build -ldflags "-X 'main.Version=${VERSION}' -X 'main.Build=${BUILD}' -X 'main.AppName=${EXECUTABLE}'" -o ${EXECUTABLE} ./main.go

run: build
	@./${EXECUTABLE}

clean:
	@rm -rf ${EXECUTABLE} dist/

build_linux:
	@env GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.Version=${VERSION}' -X 'main.Build=${BUILD}' -X 'main.AppName=${EXECUTABLE}'" -o ${EXECUTABLE} ./main.go

build_mac:
	@env GOOS=darwin GOARCH=amd64 go build -ldflags "-X 'main.Version=${VERSION}' -X 'main.Build=${BUILD}' -X 'main.AppName=${EXECUTABLE}'" -o ${EXECUTABLE} ./main.go