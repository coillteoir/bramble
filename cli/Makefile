OWNER?=davidlynchsd

build: vet
	go build -o bin/bramble .

run:
	go run .

docker-build-dev:
	docker build -t ${OWNER}/bramble --target debug . 

docker-build:
	docker build -t ${OWNER}/bramble --target release . 

docker-push:
	docker push ${OWNER}/bramble

k8s-deploy:
	kubectl apply -k k8s-manifests

vet:
	go fmt ./...
	go vet ./...

test: build
	sh tests/init_test.sh

lint: vet
	golangci-lint run
