FROM ghongli/golang:1.17.3-buster AS golang
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn
RUN CGO_ENABLED=0 go get -ldflags '-s -w -extldflags -static' github.com/go-delve/delve/cmd/dlv

WORKDIR /go/src/github.com/ghongli/salt2022/cncamp/cloudnative/classpractices/ch2/2.2/httpserver
COPY . .
RUN make build

FROM ghongli/alpine:3 as alpine
RUN apk add -U --no-cache ca-certificates

FROM ghongli/distroless-static:nonroot
ENV APP_PORT=80
WORKDIR /
COPY --from=golang /go/bin/dlv /
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=golang /go/src/github.com/ghongli/salt2022/cncamp/cloudnative/classpractices/ch2/2.2/httpserver/build/linux_amd64/release/server /
ENTRYPOINT ["./server"]

