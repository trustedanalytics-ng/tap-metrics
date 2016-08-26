
build_anywhere:
	build_artefacts.sh

docker_build: build_anywhere
	build_docker_images.sh

push_docker: docker_build
	push_docker_images.sh

clean:
	echo "TODO"
