OWNER?=davidlynchsd
IMG?=${OWNER}/bramble-ui

build:
	npm --prefix frontend install
	npm --prefix backend install
	npm --prefix frontend run build
	npm --prefix backend run build

docker-build:
	docker build -t ${IMG} .
docker-push: docker-build
	docker push ${IMG}
k8s:
	kubectl apply -k k8s-manifests

