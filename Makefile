all: deps test

deps:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

test:
	go test -cover -race -v $(shell go list ./... | grep -v /vendor/)
