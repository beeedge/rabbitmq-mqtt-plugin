FROM frolvlad/alpine-glibc:alpine-3.9_glibc-2.28
MAINTAINER duyanghao <1294057873@qq.com>

RUN apk update \
    && apk add sqlite \
    && apk add socat \
    && apk add tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone && \
    # Clean up
    rm -rf /var/cache/apk/* && \
    rm -rf /tmp/* && \
    rm -rf /var/log/*

ADD _output/linux/amd64/rabbitmq-mqtt-plugin /usr/local/bin

ENTRYPOINT ["/usr/local/bin/rabbitmq-mqtt-plugin"]