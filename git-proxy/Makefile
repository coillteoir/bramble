OWNER?=davidlynchsd

build: vet
	go build -o bin/bramble-git-proxy .

run:
	go run .

docker-build-dev:
	docker build -t ${OWNER}/bramble-git-proxy --target debug . 

docker-build:
	docker build -t ${OWNER}/bramble-git-proxy --target release . 

docker-push:
	docker push ${OWNER}/bramble-git-proxy

k8s-deploy:
	kubectl apply -k k8s-manifests

vet:
	gofumpt -w **/*.go
	go vet ./...

tidy:
	gofumpt -w **/*.go
	go vet ./...
	golangci-lint run
