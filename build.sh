#!/bin/bash

set -e

export CGO_ENABLED=0

BASEDIR=$(dirname $0)
cd $BASEDIR


echo " *** Build: image with metrics collector ambassador ***"
docker build -t metrics-collector-ambassador:v0.1 collector-ambassador/


echo " *** Build: collectors utils ***"
cd collectors
go build -v
cd -


