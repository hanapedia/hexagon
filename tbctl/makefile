# Variables

# Targets
.PHONY: build
build:
	go mod edit -dropreplace=github.com/hanapedia/the-bench/the-bench-operator
	go mod tidy
	go build .
	go mod edit -replace=github.com/hanapedia/the-bench/the-bench-operator=../the-bench-operator
