FROM golang:1.19-bullseye

SHELL ["/bin/bash", "-c"]

RUN apt-get update && apt-get install -y libbz2-dev libzip-dev procps runit git memcached libmemcached-dev


## Preparing files
RUN rm -rf /etc/service
RUN mkdir /opt/app
COPY ./rootfs /

RUN chmod +x /docker-entrypoint.sh

## Startup script
CMD ["/docker-entrypoint.sh"]
