.PHONY: docker_build, docker_login, run_api

docker_build:
	@docker build -t greeter-api .
	@docker tag greeter-api:latest 092463844305.dkr.ecr.us-east-1.amazonaws.com/greeter-api:latest
ifeq (${PUSH}, 1)
	@docker push 092463844305.dkr.ecr.us-east-1.amazonaws.com/greeter-api:latest
endif

docker_login:
	@aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 092463844305.dkr.ecr.us-east-1.amazonaws.com

run_api:
	@echo serving ui in http://localhost:8080
	@go run cmd/api/main.go

