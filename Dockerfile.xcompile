FROM golang:1.7@sha256:9bd56cd1d298b30e79e6f6ea14405c98b6cfb6553e05d2cd8ac02ecfc544dee0

RUN go get github.com/mitchellh/gox

ENV CGO_ENABLED=0

COPY vendor       /go/src/github.com/realestate-com-au/dfresh/vendor
COPY vendor.conf  /go/src/github.com/realestate-com-au/dfresh/

COPY check        /go/src/github.com/realestate-com-au/dfresh/check
COPY cmd          /go/src/github.com/realestate-com-au/dfresh/cmd
COPY main.go      /go/src/github.com/realestate-com-au/dfresh/main.go
COPY registry     /go/src/github.com/realestate-com-au/dfresh/registry

WORKDIR /go/src/github.com/realestate-com-au/dfresh

RUN gox -osarch "linux/amd64 darwin/amd64"
