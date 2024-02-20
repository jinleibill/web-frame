package errhdl

import (
	"net/http"
	"testing"
	"web-frame"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	builder := NewMiddlewareBuilder()
	builder.addCode(http.StatusNotFound, []byte(`
		<html>
			<body>
				<h1>找不到了</h1>
			</body>
		</html>
	`)).addCode(http.StatusBadRequest, []byte(`
		<html>
			<body>
				<h1>请求错误</h1>
			</body>
		</html>
	`))

	server := web_frame.NewHTTPServer(web_frame.ServerWithMiddleware(builder.Build()))
	_ = server.Start(":8081")
}
