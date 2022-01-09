package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ghongli/salt2022/cncamp/cloudnative/classpractices/ch2/2.2/httpserver/utils"
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
					w.Header().Set(k, strings.Join(values, ";"))
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
