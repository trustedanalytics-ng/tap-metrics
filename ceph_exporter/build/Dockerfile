FROM golang:1.7.3
#image based on debian-jessie system (needed distribution)
MAINTAINER Klaudia Jablonska <klaudia.jablonska@intel.com>

RUN apt-get update && \
    apt-get install -y librados-dev librbd-dev librados2 librbd1

WORKDIR /go/src/

ADD ./build.sh /bin/build.sh

ENTRYPOINT ["/bin/build.sh"]
