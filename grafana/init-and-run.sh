#!/bin/bash

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

