#!/usr/bin/env bash

echo -e "\n *** Build: CEPH metrics exporter ***\n"

docker build -t extended-golang:1.0 ./build/
docker run --rm -v "$PWD/external_repo":/go/src/github.com/digitalocean/ceph_exporter extended-golang:1.0 github.com/digitalocean/ceph_exporter
mv ./external_repo/ceph_exporter ./ceph_exporter