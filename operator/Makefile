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

remove_builder:
	yes | docker image prune --filter label=stage=builder
