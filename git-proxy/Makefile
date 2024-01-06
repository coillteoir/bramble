OWNER?=davidlynchsd
build:
	cargo build --release

docker-build:
	docker build -t ${OWNER}/bramble-git-proxy .

docker-push:
	docker push ${OWNER}/bramble-git-proxy
