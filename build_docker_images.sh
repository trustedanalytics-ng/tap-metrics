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

BASEDIR=$(dirname $0)
cd $BASEDIR


function build() {
    echo -e "\n *** Build image: $1 ***\n"
    cd $2
    docker build -t $3 .
    cd -
}

echo -e "\n\n ****** Building Docker images ******\n"

build "Grafana preconfigured" grafana metrics-grafana:$TAG
build "CEPH metrics exporter" ceph_exporter metrics-ceph-exporter:$TAG

echo -e "\n\n *** DONE, OK ***\n"

