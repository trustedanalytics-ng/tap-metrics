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


INIT_FLAT_FILE=/grafana-data/.initialized

echo "Grafana available"

if [ -e "$INIT_FLAT_FILE" ] ; then
    echo "Grafana already initialized"
else
    echo "Grafana is being initialized"
    /init.sh
    touch $INIT_FLAT_FILE
fi

echo "Starting Grafana"

exec /run.sh
