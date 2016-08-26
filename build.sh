#!/bin/bash

set -e

BASEDIR=$(dirname $0)
cd $BASEDIR

./build_artefacts.sh
./build_docker_images.sh

