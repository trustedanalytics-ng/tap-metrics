#!/bin/bash

echo "Starting Grafana"

/run.sh &

echo "Waiting for Grafana becomming available"

GRAFANA_URL=localhost:3000

until $(curl --silent --fail --show-error --output /dev/null http://$GRAFANA_URL/api/datasources); do
  printf '.'
  sleep 1
done

echo "Grafana available"

echo "Adding datasources"

for file in /grafana_init_config/*-datasource.json ; do
  if [ -e "$file" ] ; then
    echo "Adding datasource: $file"
    curl --silent --fail --show-error \
      -X POST http://$GRAFANA_URL/api/datasources \
      -H "Content-Type: application/json;charset=UTF-8" \
      -d @$file
  fi
done

echo "Adding dashboards"

for file in /grafana_init_config/*-dashboard.json ; do
  if [ -e "$file" ] ; then
    echo "Adding dashboard: $file"
    curl --silent --fail --show-error \
      -X POST http://$GRAFANA_URL/api/dashboards/import \
      -H "Content-Type: application/json;charset=UTF-8" \
      -d @$file
  fi
done

echo "Initialization compleated"

wait


