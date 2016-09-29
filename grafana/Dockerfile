FROM grafana/grafana:3.1.1

RUN apt-get -y update && apt-get -y install curl

ADD init-and-run.sh /

VOLUME ["/grafana_init_config"]

ENTRYPOINT ["/init-and-run.sh"]

