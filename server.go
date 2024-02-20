package web_frame

import (
	"fmt"
	"net"
	"net/http"
)

var _ Server = &HTTPServer{}

type HandleFunc func(ctx *Context)

type Server interface {
	http.Handler

	Start(addr string) error

	addRoute(method string, path string, handleFunc HandleFunc)
}

type HTTPServerOption func(server *HTTPServer)

type HTTPServer struct {
	*routerGroup

	log func(msg string, args ...any)

	tplEngine TemplateEngine
}

func (h *HTTPServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	ctx := &Context{
		Req:       request,
		Resp:      writer,
		tplEngine: h.tplEngine,
	}

	h.serve(ctx)
}

func (h *HTTPServer) flashResp(ctx *Context) {
	if ctx.RespStatusCode > 0 {
		ctx.Resp.WriteHeader(ctx.RespStatusCode)
	}
	_, err := ctx.Resp.Write(ctx.RespData)
	if err != nil {
		h.log("响应写入失败: ", err)
	}
}

func (h *HTTPServer) serve(ctx *Context) {
	info, ok := h.findRoute(ctx.Req.Method, ctx.Req.URL.Path)

	var root HandleFunc = func(ctx *Context) {
		if !ok || info.n.handler == nil {
			ctx.RespStatusCode = 404
			ctx.RespData = []byte("Not Found")
			return
		}

		info.n.handler(ctx)
	}

	if ok && info.n != nil {
		ctx.PathParams = info.pathParams
		ctx.MatchedRoute = info.n.route

		for i := len(info.n.mdls) - 1; i >= 0; i-- {
			root = info.n.mdls[i](root)
		}
	}

	var m Middleware = func(next HandleFunc) HandleFunc {
		return func(ctx *Context) {
			next(ctx)
			h.flashResp(ctx)
		}
	}

	root = m(root)
	root(ctx)
}

func (h *HTTPServer) Start(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return http.Serve(ln, h)
}

func NewHTTPServer(opts ...HTTPServerOption) *HTTPServer {
	res := &HTTPServer{
		routerGroup: newRouterGroup(),
		log: func(msg string, args ...any) {
			fmt.Printf(msg, args...)
		},
	}

	for _, opt := range opts {
		opt(res)
	}

	return res
}

func ServerWithTemplateEngine(tplEngine TemplateEngine) HTTPServerOption {
	return func(server *HTTPServer) {
		server.tplEngine = tplEngine
	}
}

func ServerWithMiddleware(mdls ...Middleware) HTTPServerOption {
	return func(server *HTTPServer) {
		server.mdls = mdls
	}
}
