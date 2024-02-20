package errhdl

import "web-frame"

type MiddlewareBuilder struct {
	resp map[int][]byte
}

func NewMiddlewareBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		resp: map[int][]byte{},
	}
}

func (m *MiddlewareBuilder) addCode(status int, data []byte) *MiddlewareBuilder {
	m.resp[status] = data
	return m
}

func (m *MiddlewareBuilder) Build() web_frame.Middleware {
	return func(next web_frame.HandleFunc) web_frame.HandleFunc {
		return func(ctx *web_frame.Context) {
			next(ctx)
			resp, ok := m.resp[ctx.RespStatusCode]
			if ok {
				ctx.RespData = resp
			}
		}
	}
}
