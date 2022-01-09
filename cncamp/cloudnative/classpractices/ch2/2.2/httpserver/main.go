package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/ghongli/salt2022/cncamp/cloudnative/classpractices/ch2/2.2/httpserver/middleware"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthz", healthz)

	// debug/pprof
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	timeout := middleware.DefaultTimeout

	srv := &http.Server{
		Addr:         ":80",
		Handler:      mux,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("error bringing up listener: %v", err)
		}
	}()

	<-sc
	signal.Stop(sc)

	if err := http.ListenAndServe(":80", mux); err != nil {
		log.Fatal(err)
	}

	// Shutdown timeout should be max request timeout (with 1s buffer).
	ctxShutDown, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}

	log.Println("server shutdown gracefully")
}

func healthz(w http.ResponseWriter, _ *http.Request) {
	_, err := io.WriteString(w, "ok\n")
	if err != nil {
		_, _ = io.WriteString(w, err.Error())
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("entering root handler")
	user := r.URL.Query().Get("user")
	if user != "" {
		_, _ = io.WriteString(w, fmt.Sprintf("hello [%s]\n", user))
	} else {
		_, _ = io.WriteString(w, "hello [stranger]\n")
	}
	_, _ = io.WriteString(w, "== Details of the http request header: ==\n")
	for k, v := range r.Header {
		_, _ = io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}
}
