#!/bin/bash
#
# Copyright (c) 2016 Intel Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#


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
