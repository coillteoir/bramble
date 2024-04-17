OWNER?=davidlynchsd

build: crd-sync
	make -C operator
	make -C ui
	make -C git-proxy
	make -C cli

docker-build: crd-sync
	make -C operator docker-build OWNER=${OWNER}
	make -C ui docker-build OWNER=${OWNER}
	make -C git-proxy docker-build OWNER=${OWNER}

build-deploy: build-push k8s-deploy

manifests: crd-sync
	kustomize build . > resources.yaml

k8s-deploy: manifests
	kubectl apply -f resources.yaml

k8s-teardown: 
	kubectl delete -f resources.yaml

docker-push:
	make -C operator docker-push OWNER=${OWNER}
	make -C ui docker-push OWNER=${OWNER}
	make -C git-proxy docker-push OWNER=${OWNER}

build-push: docker-build docker-push

crd-sync:
	make -C operator manifests
	npx @kubernetes-models/crd-generate --input operator/config/crd/bases/* --output ui/frontend/src/bramble-types
	npx @kubernetes-models/crd-generate --input operator/config/crd/bases/* --output ui/backend/src/bramble-types
	head -n -3 operator/api/v1alpha1/execution_types.go > git-proxy/v1alpha1/execution_types.go
