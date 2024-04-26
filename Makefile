run:
	@go build main.go
	@./main

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
	@go tool cover -html=tests/coverage.out -o _tests/coverage.html
	
clear:
	@rm -f main
	@rm -f tests/*

tidy:
	@go mod tidy

