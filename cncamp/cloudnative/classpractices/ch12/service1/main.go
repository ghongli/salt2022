package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	_ "net/http/pprof" //nolint:gosec
	"os"
	"os/signal"
	"syscall"
	"time"

	klogv2 "k8s.io/klog/v2"
)

const OutputCallDepth = 6

var (
	addr         = flag.String("ADDR", ":80", "address to listen")
	version      = flag.String("VERSION", "v1.0.0", "version for httpserver")
	shutdownTime = flag.Int("GRACE_TIMEOUT", 5, "shutdownTime for httpserver")
)

func main() {
	klogv2.InitFlags(nil)

	_ = flag.Set("v", *version)

	flag.Parse()
	defer klogv2.Flush()

	klogv2.Infof("Starting service1")

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)

	srv := &http.Server{
		Addr:    *addr,
		Handler: mux,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
			klogv2.Fatalf("error bringing up listener: %v", err)
		}
	}()

	klogv2.Info("server started")
	<-done
	signal.Stop(done)
	klogv2.Info("server stopped")

	if err := http.ListenAndServe(*addr, mux); err != nil {
		klogv2.Fatal(err)
	}

	// Shutdown timeout should be max request timeout (with 1s buffer).
	ctxShutDown, cancel := context.WithTimeout(context.Background(), time.Duration(*shutdownTime)*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		klogv2.FatalfDepth(OutputCallDepth, "server shutdown failed: %v", err)
	}
	klogv2.Info("server exited properly")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	klogv2.Info("entering v2 root handler")

	delay := randInt(10, 20)
	time.Sleep(time.Millisecond * time.Duration(delay))

	_, _ = io.WriteString(w, "===Details of the http request header:===\n")
	for k, v := range r.Header {
		_, _ = io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}
	klogv2.Infof("Respond in %d ms", delay)
}

func healthz(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		_, _ = io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}

	_, _ = io.WriteString(w, "ok\n")
}

func randInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
