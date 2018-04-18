test: *.go
	-sudo pfctl -F all
	-sudo pfctl -a asd -F all
	-go test -exec sudo -cover -v github.com/datawire/pf
	-sudo pfctl -F all
	-sudo pfctl -a asd -F all

cover: coverage.out
	go tool cover -html=coverage.out -o coverage.html
