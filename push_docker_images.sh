#!/bin/bash

set -e

TAG=${1:-latest}

function push() {
    echo -e "\n *** Pushing image: $1 ***\n"
    docker tag $1 $DOCKER_REGISTRY/$1
    docker push $DOCKER_REGISTRY/$1
}

echo -e "\n ****** Pushing Docker images ******\n\n"

push metrics-collector-ambassador:$TAG
push metrics-tap-catalog-collector:$TAG
push metrics-grafana:$TAG
push metrics-ceph-exporter:$TAG

echo -e "\n\n ****** DONE: Pushing Docker images ******\n"
