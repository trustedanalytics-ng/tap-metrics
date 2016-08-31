#!/bin/bash

set -e

TAG=${1:-latest}

BASEDIR=$(dirname $0)
cd $BASEDIR


function build() {
    echo -e "\n *** Build image: $1 ***\n"
    cd $2
    docker build -t $3 .
    cd -
}

echo -e "\n\n ****** Building Docker images ******\n"

build "metrics collector ambassador" collector-ambassador/ metrics-collector-ambassador:$TAG
build "presenter" presenter metrics-presenter:$TAG
build "TAP Catalog metrics Collector" collectors/tap_catalog metrics-tap-catalog-collector:$TAG

echo -e "\n\n *** DONE, OK ***\n"

