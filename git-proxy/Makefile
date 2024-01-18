OWNER?=davidlynchsd

build: tidy
	go build -o bin/bramble-git-proxy .

run:
	go run .

docker-build:
	docker build -t ${OWNER}/bramble-git-proxy .

docker-push:
	docker push ${OWNER}/bramble-git-proxy

tidy:
	go fmt ./...
	go vet ./...
	golangci-lint run
