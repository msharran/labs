SHELL := /bin/bash

aggr-pushgateway:
	docker run --rm -d -p 8888:8888 -p 80:80 --name aggr-pushgateway ghcr.io/zapier/prom-aggregation-gateway:v0.8.1

pushgateway:
	docker run --rm -d -p 9091:9091 --name pushgateway prom/pushgateway

	
