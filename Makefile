.PHONY: all build test lint clean run

BINARY_NAME=devradar

all: test build

build:
	go build -o $(BINARY_NAME) ./cmd/devradar

test:
	go test -v -cover ./...

lint:
	go vet ./...

clean:
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe coverage.txt

run:
	go run ./cmd/devradar
