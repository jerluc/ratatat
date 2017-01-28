all: clean build at-serial

clean:
	rm -rf target

build:
	go build

at-serial:
	go build -o target/at-serial tool/serial.go

.PHONY: clean build at-serial all
