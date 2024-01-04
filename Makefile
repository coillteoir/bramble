
build_deploy: build_all local_k8s_deploy

local_k8s_deploy:
	make -C operator install deploy
	make -C ui k8s

build_all:
	make -C operator docker-build docker-push
	make -C ui docker-build docker-push

teardown:
	kind delete cluster
