# Variables

# Targets
.PHONY: build
build:
	go mod edit -dropreplace=github.com/hanapedia/the-bench/config
	go mod tidy
	go build .
	go mod edit -replace=github.com/hanapedia/the-bench/config=../config
