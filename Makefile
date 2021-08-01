.PHONY: mod

default: fennelcli fenneld

test: mod
	go test ./...

vet: mod
	go vet ./...

fennelcli: mod
	go build -o bin/ cmd/fennelcli/fennelcli.go

fenneld: mod
	go build -o bin/ cmd/fenneld/fenneld.go

mod:
	go mod download

clean:
	rm -rf bin/*
