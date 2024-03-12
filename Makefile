OWNER?=davidlynchsd

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

crd-generate:
	make -C operator manifests
	npx @kubernetes-models/crd-generate --input operator/config/crd/bases/* --output ui/frontend/src/bramble-types
	npx @kubernetes-models/crd-generate --input operator/config/crd/bases/* --output ui/backend/src/bramble-types
