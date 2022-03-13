package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ghongli/salt2022/cncamp/cloudnative/classpractices/ch2/2.2/httpserver/logger"
	"github.com/ghongli/salt2022/cncamp/cloudnative/classpractices/ch2/2.2/httpserver/metrics"
	"github.com/namsral/flag"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ghongli/salt2022/cncamp/cloudnative/classpractices/ch2/2.2/httpserver/middleware"
)

var (
	addr         = flag.String("ADDR", ":80", "address to listen")
	version      = flag.String("VERSION", "v1.0.0", "version for httpserver")
	shutdownTime = flag.Int("GRACE_TIMEOUT", 15, "shutdownTime for httpserver")
	logFile      = flag.String("LOG_FILE", "httpserver.log", "log_file for httpserver")
	logLevel     = flag.String("LOG_LEVEL", "INFO", "log_level for httpserver")
)

func main() {
	flag.Parse()

	versionEnvKey := "VERSION"
	_ = os.Setenv(versionEnvKey, *version)

	lg, err := logger.New(*logFile, *logLevel)
	if err != nil {
		log.Fatalf("failed to create a zap logger: %v", err)
	}
	defer func() {
		if err := lg.Sync(); err != nil {
			log.Fatalf("failed to sync log: %v", err)
		}
	}()
	sugar := lg.Sugar()

	if err := metrics.Register(); err != nil {
		log.Panicf("failed to register metrics: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", middleware.Chain(rootHandler,
		middleware.MDEnvVar(versionEnvKey), middleware.KnownMD(), middleware.StdoutClientInfo(), middleware.StdoutElapsedTime(),
		middleware.Logger(lg)))
	mux.HandleFunc("/healthz", middleware.Chain(healthz, middleware.Logger(lg)))
	// add metrics collecting.
	mux.Handle("/metrics", promhttp.Handler())

	// debug/pprof
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	timeout := middleware.DefaultTimeout
	srv := &http.Server{
		Addr:         *addr,
		Handler:      mux,
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
			sugar.Fatalf("error bringing up listener: %v", err)
		}
	}()

	<-done
	signal.Stop(done)

	if err := http.ListenAndServe(*addr, mux); err != nil {
		sugar.Fatal(err)
	}

	// Shutdown timeout should be max request timeout (with 1s buffer).
	ctxShutDown, cancel := context.WithTimeout(context.Background(), time.Duration(*shutdownTime)*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		sugar.Panicf("server shutdown failed: %v", err)
	}

	sugar.Info("server shutdown gracefully")
}

func healthz(w http.ResponseWriter, _ *http.Request) {
	_, err := io.WriteString(w, http.StatusText(http.StatusOK))
	if err != nil {
		_, _ = io.WriteString(w, err.Error())
	}
	w.WriteHeader(http.StatusOK)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()

	user := r.URL.Query().Get("user")

	delay := randInt(10, 2000)
	time.Sleep(time.Millisecond * time.Duration(delay))

	if user != "" {
		_, _ = io.WriteString(w, fmt.Sprintf("hello [%s]\n", user))
	} else {
		_, _ = io.WriteString(w, "root path\n")
	}
	_, _ = io.WriteString(w, "== Details of the http request header: ==\n")
	for k, v := range r.Header {
		_, _ = io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}
}

func randInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
