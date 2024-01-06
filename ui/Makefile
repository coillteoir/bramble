OWNER?=davidlynchsd
IMG?=${OWNER}/bramble-ui
docker-build:
	docker build -t ${IMG} .
docker-push: docker-build
	docker push ${IMG}
k8s:
	kubectl apply -k k8s-manifests
