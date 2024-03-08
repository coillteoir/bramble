OWNER?=davidlynchsd

init:
	make -C operator
	make -C ui
	make -C git-proxy

build_all:
	make -C operator docker-build OWNER=${OWNER}
	make -C ui docker-build OWNER=${OWNER}
	make -C git-proxy docker-build OWNER=${OWNER}

build_deploy: build_all push_all k8s_deploy

k8s_deploy:
	make -C operator deploy
	make -C ui k8s
	make -C git-proxy local-k8s-deploy

push_all:
	make -C operator docker-push OWNER=${OWNER}
	make -C ui docker-push OWNER=${OWNER}
	make -C git-proxy docker-push OWNER=${OWNER}

build_push: build_all push_all

teardown:
	kind delete cluster
