export GO111MODULE=on
export GOOS:=$(shell go env GOOS)
export GOARCH:=$(shell go env GOARCH)

docker-build:
	docker run -it -e GOOS=${GOOS} -e GOARCH=${GOARCH} -v $(shell pwd):/golang-rest-auth -w /golang-rest-auth golang:1.12 make generate-bin

generate-bin:
	go build -mod vendor -o bin/server .