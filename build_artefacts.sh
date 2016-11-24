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

cd ./ceph_exporter/
./build_ceph_exporter.sh
cd -

echo -e "\n\n *** DONE, OK ***\n"

