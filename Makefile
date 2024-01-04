
build_deploy: build_all local_k8s_deploy

local_k8s_deploy:
	make -C operator install deploy
	make -C ui k8s

build_all:
	make -C operator docker-build
	make -C ui docker-build
	make -C git-proxy docker-build

push_all:
	make -C operator docker-push
	make -C ui docker-push
	make -C git-proxy docker-push

build_push: build_all push_all

teardown:
	kind delete cluster
