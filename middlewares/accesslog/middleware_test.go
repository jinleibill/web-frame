package accesslog

import (
	"fmt"
	"net/http"
	"testing"
	"web-frame"
)

func TestMiddlewareBuilder(t *testing.T) {
	builder := MiddlewareBuilder{}
	mdl := builder.LogFunc(func(log string) {
		fmt.Println(log)
	}).Build()
	server := web_frame.NewHTTPServer(web_frame.ServerWithMiddleware(mdl))
	server.Post("/users/:user/comments", func(ctx *web_frame.Context) {
		fmt.Println("hello")
	})
	req, err := http.NewRequest(http.MethodPost, "/users/1/comments", nil)
	req.Host = "localhost:8080"
	if err != nil {
		t.Fatal(err)
	}
	server.ServeHTTP(nil, req)
}
