default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

build:
	go build ./...

install:
	go install ./...

download:
	go mod download
