SHELL := /bin/bash

.PHONY: buildall http https runhttp runhttps

http:
	cd cmd/http && go build -o ../../bin/http

https:
	cd cmd/https && go build -o ../../bin/https

buildall: http https

runhttp:
	./bin/http

runhttps:
	./bin/https
