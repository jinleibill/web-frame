package web_frame

type Middleware func(next HandleFunc) HandleFunc
