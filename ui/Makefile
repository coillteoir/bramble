OWNER?=davidlynchsd
IMG?=${OWNER}/bramble-ui

build: crd-generate
	npm --prefix frontend run build
	npm --prefix backend run build

docker-build:
	docker build -t ${IMG} .
docker-push: docker-build
	docker push ${IMG}
k8s:
	kubectl apply -k k8s-manifests

crd-generate:
	make -C ../operator generate
	npx @kubernetes-models/crd-generate --input ../operator/config/crd/bases/* --output frontend/src/bramble-types
	npx @kubernetes-models/crd-generate --input ../operator/config/crd/bases/* --output backend/src/bramble-types
	
