REGISTRY?=davidlynchsd
build:
	cargo build --release

docker-build:
	docker build -t ${REGISTRY}/bramble-git-proxy .

docker-push:
	docker push ${REGISTRY}/bramble-git-proxy
