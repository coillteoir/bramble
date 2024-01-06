REGISTRY?=davidlynchsd

build_all:
	make -C operator docker-build REGISTRY=${REGISTRY}
	make -C ui docker-build REGISTRY=${REGISTRY}
	make -C git-proxy docker-build REGISTRY=${REGISTRY}

build_deploy: build_all k8s_deploy

k8s_deploy:
	make -C operator deploy
	make -C ui k8s

push_all:
	make -C operator docker-push REGISTRY=${REGISTRY}
	make -C ui docker-push REGISTRY=${REGISTRY}
	make -C git-proxy docker-push REGISTRY=${REGISTRY}

build_push: build_all push_all

teardown:
	kind delete cluster
