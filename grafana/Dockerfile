FROM grafana/grafana:4.0.2

RUN apt-get -y update && apt-get -y install curl

ADD init-and-run.sh /
ADD init.sh /

ADD init_content /grafana_init_static_content

COPY overrides/grafana.light.min.76121946.css /usr/share/grafana/public/css/grafana.light.min.76121946.css

VOLUME ["/grafana_init_dynamic_content"]

ENV GF_ANALYTICS_REPORTING_ENABLED false

ENTRYPOINT ["/init-and-run.sh"]

