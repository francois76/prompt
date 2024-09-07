# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod



all: deps test lint

test:
	$(GOTEST)  ./...

coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

lint:
	docker run --rm -v "$(shell pwd)":/app -w /app golangci/golangci-lint:v1.60.3 golangci-lint run --color always

clean:
	$(GOCLEAN)
	rm -f coverage.out coverage.html

deps:
	$(GOGET) -v -t  ./...
	$(GOMOD) tidy

.PHONY: all test coverage lint clean deps