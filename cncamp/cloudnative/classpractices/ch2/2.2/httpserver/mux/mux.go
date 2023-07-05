package mux

import (
	"net/http"

	"github.com/ghongli/salt2022/cncamp/cloudnative/classpractices/ch2/2.2/httpserver/utils"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

type (
	Mux struct {
		HTTPMux    http.Handler
		GRPCServer *grpc.Server
	}
)

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.ProtoMajor == 2 && utils.IsGRPCContentType(r.Header.Get(utils.HeaderContentTypeKey)) {
		m.GRPCServer.ServeHTTP(w, r)
	} else {
		m.HTTPMux.ServeHTTP(w, r)
	}
}

func InsecureHandler(handler http.Handler) http.Handler {
	// "h2c" is the unencrypted form of HTTP/2.
	return h2c.NewHandler(handler, &http2.Server{})
}
