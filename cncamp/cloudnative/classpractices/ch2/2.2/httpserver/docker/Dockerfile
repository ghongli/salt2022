FROM golang:latest AS golang
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn
RUN CGO_ENABLED=0 go get -ldflags '-s -w -extldflags -static' github.com/go-delve/delve/cmd/dlv

FROM alpine:latest as alpine
RUN apk add -U --no-cache ca-certificates

FROM gcr.io/distroless-static:nonroot
ARG PKG_FILES
WORKDIR /
COPY --from=golang /go/bin/dlv /
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY /$PKG_FILES /

