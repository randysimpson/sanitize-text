#/bin/bash

echo "Running Unit Tests..."

cd src

go mod tidy

go test -v ./...
