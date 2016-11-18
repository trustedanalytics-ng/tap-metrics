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


GRAFANA_URL=localhost:3000

function addFromFolder {

    FOLDER=$1

    echo "Adding datasources from $FOLDER"

    for file in $FOLDER/*-datasource.json ; do
      if [ -e "$file" ] ; then
        echo "Adding datasource: $file"
        curl --silent --fail --show-error \
          -X POST http://$GRAFANA_URL/api/datasources \
          -H "Content-Type: application/json;charset=UTF-8" \
          -d @$file
      fi
    done

    echo "Adding dashboards from $FOLDER"

    for file in $FOLDER/*-dashboard.json ; do
      if [ -e "$file" ] ; then
        echo "Adding dashboard: $file"
        curl --silent --fail --show-error \
          -X POST http://$GRAFANA_URL/api/dashboards/import \
          -H "Content-Type: application/json;charset=UTF-8" \
          -d @$file
      fi
    done
}


echo "Starting Grafana"

/run.sh &

echo "Waiting for Grafana becomming available"


until $(curl --silent --fail --show-error --output /dev/null http://$GRAFANA_URL/api/datasources); do
  printf '.'
  sleep 1
done

echo "Grafana available"

echo "Initializing static conent"
addFromFolder /grafana_init_static_content

echo "Initializing dynamic conent"
addFromFolder /grafana_init_dynamic_content

echo "Initialization compleated"

wait

