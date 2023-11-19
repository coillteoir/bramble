
setup_env:
	make -C operator generate install docker-build docker-push deploy IMG=docker.io/davidlynchsd/bramble
	cd ui/frontend/
	npm run dev
