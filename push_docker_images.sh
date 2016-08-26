#!/bin/bash

set -e

function push() {
    echo -e "\n *** Pushing image: $1 ***\n"
    docker tag $1 $DOCKER_REGISTRY/$1
    docker push $DOCKER_REGISTRY/$1
}

echo -e "\n ****** Pushing Docker images ******\n\n"

push metrics-collector-ambassador:v0.1
push metrics-presenter:v0.1
push metrics-tap-catalog-collector:v0.1

echo -e "\n\n ****** DONE: Pushing Docker images ******\n"
