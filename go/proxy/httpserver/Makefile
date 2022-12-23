SHELL := /bin/bash

.PHONY: buildall http https runhttp runhttps

buildall: http https

http:
	cd cmd/http && go build -o ../../bin/http

https:
	cd cmd/https && go build -o ../../bin/https

runhttp:
	./bin/http

runhttps:
	./bin/https
