FROM alpine:3.20

ARG TARGETARCH
ARG BUILD_DIR=build
ARG BINARY=energy-prometheus-exporter-linux-${TARGETARCH}

COPY ${BUILD_DIR}/${BINARY} /energy-prometheus-exporter

RUN addgroup -S exporter && \
	adduser -S -D -G exporter -H -s /bin/sh exporter

ENV HOME=/

USER exporter

EXPOSE 8080

ENTRYPOINT [ "/energy-prometheus-exporter", "serve" ]
