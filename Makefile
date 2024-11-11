#!make

bin_name = "docker-cloudflare"
cmd_path = ./cmd/cicd-cloudflare-zero-trust-thing/

.PHONY: clean
clean:
	if [ -d ./tmp ]; then rm -r ./tmp; fi

.PHONY: tidy
tidy:
	go mod tidy
	go fmt ./...


.PHONY: build
build: clean
	go build -o=./${bin_name} ${cmd_path}
