FROM grafana/grafana:3.1.1

RUN apt-get -y update && apt-get -y install curl

ADD init-and-run.sh /

ADD init_content /grafana_init_static_content

VOLUME ["/grafana_init_dynamic_content"]

ENV GF_ANALYTICS_REPORTING_ENABLED false

ENTRYPOINT ["/init-and-run.sh"]

