.SILENT:

include .env
export

.PHONY: fmt lint test race run scripts build_docker_image remove_container_and_image

fmt:
	go fmt ./...

lint: fmt
	go vet ./...

test: lint
	go test -v -cover ./...

race: test
	go test -v -race ./...

run: race
	go run -v cmd/forum/main.go

build_docker_image:
	bash ./scripts/build_image.sh

remove_container_and_image:
	bash ./scripts/remove_container_and_image.sh


.DEFAULT_GOAL := run
