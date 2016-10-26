FROM debian:8.6
MAINTAINER Klaudia Jablonska <klaudia.jablonska@intel.com>

RUN apt-get update && apt-get install -y --force-yes librados-dev librbd-dev librados2 librbd1

RUN mkdir -p /opt/app
COPY ceph_exporter /opt/app

EXPOSE 9128
ENTRYPOINT ["/opt/app/ceph_exporter"]
