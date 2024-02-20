package cookie

import (
	"net/http"
)

type Propagator struct {
	CookieName string
}

func NewPropagator() *Propagator {
	return &Propagator{
		CookieName: "sessid",
	}
}

func (p *Propagator) Inject(id string, writer http.ResponseWriter) error {
	http.SetCookie(writer, &http.Cookie{
		Name:  p.CookieName,
		Value: id,
	})
	return nil
}

func (p *Propagator) Extract(req *http.Request) (string, error) {
	c, err := req.Cookie(p.CookieName)
	if err != nil {
		return "", err
	}
	return c.Value, nil
}

func (p *Propagator) Remove(writer http.ResponseWriter) error {
	http.SetCookie(writer, &http.Cookie{
		Name:   p.CookieName,
		MaxAge: -1,
	})
	return nil
}
