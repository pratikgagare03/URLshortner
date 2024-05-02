run:
	@cd api;go build main.go
	@sudo systemctl start redis
	@cd api; ./main

dcbuild:
	@sudo docker-compose build

dcup:
	@sudo docker-compose up -d

dcdown:
	@sudo docker-compose down
	
tests:
	@cd metrics;go test
	@cd short;go test
	@cd redirection;go test

testscover:
	@cd metrics;go test -cover
	@cd short;go test -cover
	@cd redirection;go test -cover

testscoverhtml:
	@go test -coverprofile=tests/coverage.out ./...
	@go tool cover -html=tests/coverage.out -o tests/coverage.html
	
clear:
	@rm -f api/main
	@rm -f api/tests/*

tidy:
	@go mod tidy

