default: bin/time-tracker

bin/time-tracker:	$(wildcard *.go) $(wildcard */*.go)
	go build -o=bin/time-tracker main.go