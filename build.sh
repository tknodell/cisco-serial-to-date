#!/bin/bash

set -e

rm -rf ./bin
mkdir -p ./bin
go build -o ./bin/serial_to_date serial_to_date.go
echo "Built succesfully, run ./bin/serial_to_date"