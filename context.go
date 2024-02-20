package web_frame

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
)

type Context struct {
	Req  *http.Request
	Resp http.ResponseWriter

	RespStatusCode int
	RespData       []byte

	PathParams map[string]string

	queryValues url.Values

	MatchedRoute string

	tplEngine TemplateEngine

	UserValues map[string]any
}

func (c *Context) Render(tplName string, data any) error {
	var err error
	c.RespData, err = c.tplEngine.Render(c.Req.Context(), tplName, data)
	if err != nil {
		c.RespStatusCode = http.StatusInternalServerError
		return err
	}
	c.RespStatusCode = http.StatusOK
	return nil
}

func (c *Context) BindJson(val any) error {
	decoder := json.NewDecoder(c.Req.Body)
	return decoder.Decode(val)
}

func (c *Context) FormValue(key string) (string, error) {
	err := c.Req.ParseForm()
	if err != nil {
		return "", err
	}
	return c.Req.FormValue(key), nil
}

func (c *Context) QueryValue(key string) (string, error) {
	if c.queryValues == nil {
		c.queryValues = c.Req.URL.Query()
	}

	vals, ok := c.queryValues[key]
	if !ok {
		return "", errors.New("web: key 不存在")
	}

	return vals[0], nil
}

func (c *Context) PathValue(key string) (string, error) {
	val, ok := c.PathParams[key]
	if !ok {
		return "", errors.New("web: key 不存在")
	}

	return val, nil
}

func (c *Context) RespJson(code int, val any) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}
	c.RespData = data
	c.RespStatusCode = code
	return err
}

func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Resp, cookie)
}
