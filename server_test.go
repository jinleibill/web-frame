package web_frame

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHttpServer(t *testing.T) {
	h := NewHTTPServer()
	h.addRoute(http.MethodGet, "/order", func(ctx *Context) {
		_, _ = ctx.Resp.Write([]byte("hello, order"))
	})

	h.addRoute(http.MethodGet, "/order/:id", func(ctx *Context) {
		_, _ = ctx.Resp.Write([]byte("hello, " + ctx.PathParams["id"]))
	})

	v1 := h.Group("v1")
	{
		adminRoute := v1.Group("admins")
		adminRoute.Get("", func(ctx *Context) {
			fmt.Println("路由处理1")
		})
	}

	v2 := h.Group("v2")
	{
		userRoute := v2.Group("users")
		userRoute.Get("/:user", func(ctx *Context) {
			fmt.Println("路由处理2")
		})
	}

	_ = h.Start(":8081")
}
