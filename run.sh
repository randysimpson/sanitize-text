#/bin/bash

echo "Running from golang container..."

cd src

go mod tidy

go run .
