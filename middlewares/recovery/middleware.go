package recovery

import "web-frame"

type MiddlewareBuilder struct {
	StatueCode int
	Data       []byte
	Log        func(ctx *web_frame.Context)
}

func (m MiddlewareBuilder) Build() web_frame.Middleware {
	return func(next web_frame.HandleFunc) web_frame.HandleFunc {
		return func(ctx *web_frame.Context) {
			defer func() {
				if err := recover(); err != nil {
					ctx.RespData = m.Data
					ctx.RespStatusCode = m.StatueCode
					m.Log(ctx)
				}
			}()
			next(ctx)
		}
	}
}
