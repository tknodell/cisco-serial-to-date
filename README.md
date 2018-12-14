# Description

Stupid script which takes single or multiple cisco serial numbers and returns the manufactured date.

## golang:
```
go run main.go --serial FAA04459FNI
go run main.go --filename serials.txt
```

## build golang binary:
```
./build.sh
./bin/serial_to_date --serial FAA04459FNI
./bin/serial_to_date --filename serials.txt
```

# OR

## ruby:
```
ruby ruby/serial_to_date.rb FAA04459FNI
ruby ruby/serial_to_date.rb serials.txt
```
