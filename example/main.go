package main

import (
	"fmt"
	"web-frame"
	"web-frame/middlewares/accesslog"
	"web-frame/middlewares/auth"
	"web-frame/middlewares/recovery"
)

func main() {
	h := web_frame.NewHTTPServer(web_frame.ServerWithMiddleware(
		accesslog.NewMiddleBuilder().LogFunc(func(log string) {
			fmt.Println(log)
		}).Build(),
		recovery.MiddlewareBuilder{
			StatueCode: 500,
			Data:       []byte("panic ..."),
			Log: func(ctx *web_frame.Context) {
				fmt.Printf("panic %s", ctx.Req.URL.String())
			},
		}.Build(),
	))

	v1 := h.Group("v1")
	{
		adminRoute := v1.Group("admins", auth.MiddlewareBuilder{}.Build())
		adminRoute.Get("", func(ctx *web_frame.Context) {
			_, _ = ctx.Resp.Write([]byte("admins"))
		})
	}

	v2 := h.Group("v2")
	{
		userRoute := v2.Group("users")
		userRoute.Get("/:user", func(ctx *web_frame.Context) {
			_, _ = ctx.Resp.Write([]byte("user"))
		})
	}

	_ = h.Start(":8081")
}
