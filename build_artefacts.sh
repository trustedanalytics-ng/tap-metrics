#!/bin/bash

set -e

export CGO_ENABLED=0

BASEDIR=$(dirname $0)
cd $BASEDIR

function build() {
    echo -e "\n *** Build: $1 ***\n"
    cd $2
    go build -v
    cd -
}

echo -e "\n\n ****** Building TAP Metrics related aretfacts ******\n"

build "collectors utils" collectors

build "TAP Catalog metrics collector" collectors/tap_catalog

cd ./ceph_exporter/
./build_ceph_exporter.sh
cd -

echo -e "\n\n *** DONE, OK ***\n"

