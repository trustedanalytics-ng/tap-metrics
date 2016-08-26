REPO=github.com/trustedanalytics/metrics

prepare_dirs:
	mkdir -p ./temp/src/$(REPO)
	$(eval REPOFILES=$(shell pwd)/*)
	ln -sf $(REPOFILES) temp/src/$(REPO)


build_anywhere: prepare_dirs
	$(eval GOPATH=$(shell cd ./temp; pwd))
	(cd ./temp/src/$(REPO); GOPATH=$(GOPATH) ./build_artefacts.sh)

docker_build: build_anywhere
	./build_docker_images.sh

push_docker: docker_build
	DOCKER_REGISTRY=$(REGISTRY_URL) ./push_docker_images.sh

clean:
	echo "TODO"
