build:
	go build .

lint:
	golangci-lint run
	gofmt -s -w **/*.go

test:
	go test ./...

release: test lint build


docker: test lint
	docker build -t lugh .

k8s_local: docker
	kind create cluster
	kind load docker-image lugh
	kubectl apply -k k8s-resources
