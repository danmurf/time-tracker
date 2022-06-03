default: bin/time-tracker

rwildcard=$(foreach d,$(wildcard $(1:=/*)),$(call rwildcard,$d,$2) $(filter $(subst *,%,$2),$d))

bin/time-tracker: $(call rwildcard,.,*.go)
	go build -o=bin/time-tracker main.go

fmt:
	go fmt ./...

test: lint
	go test ./...

lint:
	golangci-lint run