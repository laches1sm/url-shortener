PHONY: build run test

build:
	go build $$(go list ./...)

run:
	go run $$(go list ./...)


test:
	go test $$(go list ./...)