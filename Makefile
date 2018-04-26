test: *.go
	-sudo pfctl -F all
	-sudo pfctl -a asd -F all
	-sudo pfctl -a myanchor -F all
	go generate
	-go test -exec sudo -cover -v github.com/datawire/pf
	-sudo pfctl -F all
	-sudo pfctl -a asd -F all
	-sudo pfctl -a myanchor -F all

cover: coverage.out
	go tool cover -html=coverage.out -o coverage.html
