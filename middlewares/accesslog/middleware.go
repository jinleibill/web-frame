package accesslog

import (
	"encoding/json"
	"web-frame"
)

type MiddlewareBuilder struct {
	logFunc func(log string)
}

func (m *MiddlewareBuilder) LogFunc(fn func(log string)) *MiddlewareBuilder {
	m.logFunc = fn
	return m
}

func (m *MiddlewareBuilder) Build() web_frame.Middleware {
	return func(next web_frame.HandleFunc) web_frame.HandleFunc {
		return func(ctx *web_frame.Context) {
			defer func() {
				l := AccessLog{
					Host:       ctx.Req.Host,
					Route:      ctx.MatchedRoute,
					HTTPMethod: ctx.Req.Method,
					Path:       ctx.Req.URL.Path,
				}
				data, _ := json.Marshal(l)
				m.logFunc(string(data))
			}()
			next(ctx)
		}
	}
}

func NewMiddleBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{}
}

type AccessLog struct {
	Host       string `json:"host,omitempty"`
	Route      string `json:"route,omitempty"`
	HTTPMethod string `json:"http_method,omitempty"`
	Path       string `json:"path,omitempty"`
}
