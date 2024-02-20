package auth

import "web-frame"

type MiddlewareBuilder struct {
}

func (m MiddlewareBuilder) Build() web_frame.Middleware {
	return func(next web_frame.HandleFunc) web_frame.HandleFunc {
		return func(ctx *web_frame.Context) {
			if ctx.Req.Header.Get("token") == "" {
				ctx.RespStatusCode = 401
				return
			}

			next(ctx)
		}
	}
}
