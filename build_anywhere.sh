#!/bin/bash

set -e

REPO=github.com/trustedanalytics/metrics

# GOPATH setting up
mkdir -p ./temp/src/$REPO
REPOFILES=`pwd`/*
ln -sf $REPOFILES temp/src/$REPO
export GOPATH=`cd ./temp/; pwd`

cd ./temp/src/$REPO
./build_artefacts.sh




