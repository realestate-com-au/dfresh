FROM alpine:3.6@sha256:0b94d1d1b5eb130dd0253374552445b39470653fb1a1ec2d81490948876e462c

RUN apk add --update \
    ca-certificates && \
    rm /var/cache/apk/* && \
    rm -rf /usr/share/ri

ADD target/dfresh_alpine /dfresh

WORKDIR /cwd

ENTRYPOINT ["/dfresh"]
