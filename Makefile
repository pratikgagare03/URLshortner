run:
	@cd api;go build main.go
	@sudo systemctl start redis
	@cd api; ./main

dcbuild:
	@sudo docker compose build

dcup:
	@sudo systemctl stop redis
	@sudo docker compose up -d

dcdown:
	@sudo docker compose down
	
tests:
	@cd api/metrics;go test
	@cd api/short;go test
	@cd api/redirection;go test

testscover:
	@cd api/metrics;go test -cover
	@cd api/short;go test -cover
	@cd api/redirection;go test -cover

testscoverhtml:
	@cd api;go test -coverprofile=tests/coverage.out ./...
	@cd api;go tool cover -html=tests/coverage.out -o tests/coverage.html
	
clear:
	@rm -f api/main
	@rm -f api/tests/*

tidy:
	@go mod tidy

