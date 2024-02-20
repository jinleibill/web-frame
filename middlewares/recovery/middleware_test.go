package recovery

import (
	"fmt"
	"testing"
	"web-frame"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := MiddlewareBuilder{
		StatueCode: 500,
		Data:       []byte("panic ..."),
		Log: func(ctx *web_frame.Context) {
			fmt.Printf("panic %s", ctx.Req.URL.String())
		},
	}

	server := web_frame.NewHTTPServer(web_frame.ServerWithMiddleware(builder.Build()))

	server.Get("/users", func(ctx *web_frame.Context) {
		panic("user panic")
	})

	_ = server.Start(":8081")
}
