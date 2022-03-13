package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ghongli/salt2022/cncamp/cloudnative/classpractices/ch2/2.2/httpserver/utils"
	"go.uber.org/zap"
)

type (
	MiddlewareFunc func(http.HandlerFunc) http.HandlerFunc
)

func Chain(current http.HandlerFunc, middlewares ...MiddlewareFunc) http.HandlerFunc {
	for _, middleware := range middlewares {
		current = middleware(current)
	}

	return current
}

func KnownMD() MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 处理请求头，赋值组响应头
			copiedHeader := r.Header.Clone()
			for k, values := range copiedHeader {
				values := values
				if k != utils.HeaderContentLength {
					respVal := w.Header().Values(k)
					if len(respVal) != 0 {
						values = append(values, respVal...)
					}
					w.Header().Add(k, strings.Join(values, ";"))
				}
			}

			next(w, r)
		}
	}
}

// MDEnvVar 获取指定的环境变量，写入响应头。
func MDEnvVar(envVarKey string) MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if envVarKey != "" {
				envVarVal := os.Getenv(envVarKey)
				if w.Header().Get(envVarKey) != "" {
					w.Header().Add(envVarKey, envVarVal)
				} else {
					w.Header().Set(envVarKey, envVarVal)
				}
			}

			next(w, r)
		}
	}
}

func StdoutClientInfo() MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			realAddr := utils.ExtractAddrFromRequest(r)
			log.Printf("client real info: %s, %s:%s\n", realAddr.Addr, realAddr.Host, realAddr.Port)

			w.Header().Set(utils.HeaderClientIP, realAddr.Host)

			next(w, r)
		}
	}
}

func StdoutElapsedTime() MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer func() {
				log.Printf("%s elapsed time: %s \n", r.URL.Path, time.Since(start).String())
			}()

			next(w, r)
		}
	}
}

func Logger(logger *zap.Logger) MiddlewareFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			rw := wrapperResponseWriter(w)

			next(rw, r)

			logger.Sugar().Infof("client real ip: %s, response status code: %d", utils.ExtractAddrFromRequest(r).Addr, rw.statusCode)
		}
	}
}

type (
	// 原生 http.ResponseWriter 会被拦截，不会返回给用户处理，因此也拿不到 http 返回的响应码
	responseWriter struct {
		w http.ResponseWriter

		statusCode int
	}
)

func (w *responseWriter) Header() http.Header {
	return w.w.Header()
}

func (w *responseWriter) Write(buffer []byte) (int, error) {
	return w.w.Write(buffer)
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.w.WriteHeader(statusCode)
}

func wrapperResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w: w}
}
