package prometheus

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
	"testing"
	"time"
	"web-frame"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := MiddlewareBuilder{
		Namespace: "demo",
		Subsystem: "web_frame",
		Name:      "http_response",
	}

	server := web_frame.NewHTTPServer(web_frame.ServerWithMiddleware(builder.Build()))

	server.Get("/users", func(ctx *web_frame.Context) {
		val := rand.Intn(1000) + 1
		time.Sleep(time.Duration(val) * time.Millisecond)

		_ = ctx.RespJson(200, struct {
			Name string
		}{
			Name: "bill",
		})
	})

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		_ = http.ListenAndServe(":8082", nil)
	}()

	_ = server.Start(":8081")
}
