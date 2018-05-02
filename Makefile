test: *.go
	scripts/test.sh

cover: coverage.out
	go tool cover -html=coverage.out -o coverage.html
