#!/bin/bash

set -e

export CGO_ENABLED=0

BASEDIR=$(dirname $0)
cd $BASEDIR


echo -e "\n\n *** Build: image with metrics collector ambassador ***\n"
docker build -t metrics-collector-ambassador:v0.1 collector-ambassador/


echo -e "\n\n *** Build: collectors utils ***\n"
cd collectors
go build -v
cd -


echo -e "\n\n *** Build: TAP Catalog metrics collector ***\n"
cd collectors/tap_catalog
go build -v
cd -


echo -e "\n\n *** Build: image with TAP Catalog metrics Collector ***\n"
cd collectors/tap_catalog
docker build -t metrics-tap-catalog-collector:v0.1 .
cd -


echo -e "\n\n *** DONE, OK ***\n"

