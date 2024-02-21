## Web-Frame

使用 Go 语言实现，极小内核，包含核心功能的 web 框架

## 特性
- server 可以当作 http.Handler ，也可以独立控制
- 支持分段路由树，路由参数解析，路由组
- 封装 context，支持模版渲染，json 返回
- 内置静态资源服务以及文件上传和下载
- session 支持 redis，menory 存储
- 内置日志，错误处理，可观测中间件

## 用法
```
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
    _ = ctx.RespJson(200, "hello web-frame")
  })
}

_ = h.Start(":8081")
```
