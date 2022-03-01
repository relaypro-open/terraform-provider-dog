VERSION=$(shell git describe --tags --match=v* --always --dirty)
SEMVER=$(shell git describe --tags --match=v* --always --dirty | cut -c 2-)

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

upgrade:
	go get -u ./...

.PHONY: clean
clean:
	@rm -rf bin
	@rm -rf _output

.PHONY: release
release: \
	clean \
	_output/plugin-linux-amd64.zip #\
#	_output/plugin-linux-arm64.zip \
#	_output/plugin-darwin-amd64.zip \
#	_output/plugin-darwin-arm64.zip \
#	_output/plugin-windows-amd64.zip

_output/plugin-%.zip: NAME=terraform-provider-dog_$(SEMVER)_$(subst -,_,$*)
_output/plugin-%.zip: DEST=_output/$(NAME)
_output/plugin-%.zip: _output/%/terraform-provider-dog
	@mkdir -p $(DEST)
	@cp _output/$*/terraform-provider-dog $(DEST)/terraform-provider-dog_$(VERSION)
	@zip -j $(DEST).zip $(DEST)/terraform-provider-dog_$(VERSION)

_output/linux-amd64/terraform-provider-dog: CGO_ENABLED=0 GOARGS = GOOS=linux GOARCH=amd64
#_output/linux-arm64/terraform-provider-dog: GOARGS = GOOS=linux GOARCH=arm64
#_output/darwin-amd64/terraform-provider-dog: GOARGS = GOOS=darwin GOARCH=amd64
#_output/darwin-arm64/terraform-provider-dog: GOARGS = GOOS=darwin GOARCH=arm64
#_output/windows-amd64/terraform-provider-dog: GOARGS = GOOS=windows GOARCH=amd64
_output/%/terraform-provider-dog:
	#$(GOARGS) go build -v -a -ldflags '-w -extldflags "-static"' -o $@ ./... 
	$(GOARGS) go build -o $@ ./... 

#release-sign:
#	cd _output; sha256sum *.zip > terraform-provider-dog_$(SEMVER)_SHA256SUMS
#	gpg2 --detach-sign _output/terraform-provider-dog_$(SEMVER)_SHA256SUMS
#
#release-verify: NAME=_output/terraform-provider-dog
#release-verify:
#	gpg2 --verify $(NAME)_$(SEMVER)_SHA256SUMS.sig $(NAME)_$(SEMVER)_SHA256SUMS
