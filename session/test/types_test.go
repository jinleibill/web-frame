package test

import (
	"net/http"
	"testing"
	"time"
	"web-frame"
	"web-frame/session"
	"web-frame/session/cookie"
	"web-frame/session/memory"
)

func TestSession(t *testing.T) {
	var m *session.Manager = &session.Manager{
		Propagator: cookie.NewPropagator(),
		Store:      memory.NewStore(time.Minute * 15),
		CtxSessKey: "sessKey",
	}
	server := web_frame.NewHTTPServer(web_frame.ServerWithMiddleware(func(next web_frame.HandleFunc) web_frame.HandleFunc {
		return func(ctx *web_frame.Context) {
			if ctx.Req.URL.Path == "/login" {
				next(ctx)
				return
			}
			_, err := m.GetSession(ctx)
			if err != nil {
				ctx.RespStatusCode = http.StatusUnauthorized
				ctx.RespData = []byte("请重新登陆")
				return
			}

			_ = m.RefreshSession(ctx)

			next(ctx)
		}
	}))

	server.Post("/login", func(ctx *web_frame.Context) {
		sess, err := m.InitSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("登陆失败")
			return
		}
		err = sess.Set(ctx.Req.Context(), "nickname", "bill")
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("登陆失败")
			return
		}
		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("登陆成功")
	})

	server.Post("/logout", func(ctx *web_frame.Context) {
		err := m.RemoveSession(ctx)
		if err != nil {
			ctx.RespStatusCode = http.StatusInternalServerError
			ctx.RespData = []byte("退出失败")
			return
		}

		ctx.RespStatusCode = http.StatusOK
		ctx.RespData = []byte("退出成功")
	})

	server.Get("/user", func(ctx *web_frame.Context) {
		sess, _ := m.GetSession(ctx)
		val, _ := sess.Get(ctx.Req.Context(), "nickname")
		ctx.RespData = []byte(val.(string))
	})

	_ = server.Start(":8081")
}
