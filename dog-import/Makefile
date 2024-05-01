ENV=pro

default: install

run: build
	mkdir -p /tmp/dog-import
	./dog-import -environment ${ENV} -output_dir /tmp/dog-import -host_prefix ${ENV}
	terraform fmt /tmp/dog-import

build:
	go build

update_api:
	GOPROXY=direct go get github.com/relaypro-open/dog_api_golang@main
	go mod vendor
