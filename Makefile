bin/tt:	$(wildcard *.go) $(wildcard */*.go)
	go build -o=bin/tt main.go