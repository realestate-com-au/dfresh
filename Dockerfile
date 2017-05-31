FROM alpine

MAINTAINER https://github.com/mdub/dfresh

RUN apk add --update \
    ca-certificates && \
    rm /var/cache/apk/* && \
    rm -rf /usr/share/ri

ADD dfresh /

ENTRYPOINT ["/dfresh"]
