FROM alpine:3.6@sha256:0b94d1d1b5eb130dd0253374552445b39470653fb1a1ec2d81490948876e462c

RUN apk --no-cache --update add ca-certificates \
 && rm -rf /var/cache/apk/*

COPY target/dfresh_linux_amd64 /dfresh

WORKDIR /cwd

ENTRYPOINT ["/dfresh"]
