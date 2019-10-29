#!/bin/bash

set -e

rm -rf ./bin
mkdir -p ./bin
go test .
go build -o ./bin/serial_to_date main.go
echo "Built succesfully, run ./bin/serial_to_date"