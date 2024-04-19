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

testcoverhtml:
	@go test -coverprofile=_tests/coverage.out ./...
	@go tool cover -html=_tests/coverage.out -o _tests/coverage.html
	
clear:
	@rm -f main
	@rm -f _tests/*

tidy:
	@go mod tidy

