package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
	"web-frame"
)

type MiddlewareBuilder struct {
	Namespace string
	Name      string
	Subsystem string
	Help      string
}

func (m MiddlewareBuilder) Build() web_frame.Middleware {
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: m.Namespace,
		Name:      m.Name,
		Subsystem: m.Subsystem,
		Help:      m.Help,
		Objectives: map[float64]float64{
			0.5:  0.01,
			0.75: 0.01,
			0.90: 0.01,
			0.99: 0.001,
		},
	}, []string{"pattern", "method", "status"})
	prometheus.MustRegister(vector)
	return func(next web_frame.HandleFunc) web_frame.HandleFunc {
		return func(ctx *web_frame.Context) {
			startTime := time.Now()
			defer func() {
				duration := time.Now().Sub(startTime).Milliseconds()
				pattern := ctx.MatchedRoute
				vector.WithLabelValues(pattern, ctx.Req.Method, strconv.Itoa(ctx.RespStatusCode)).Observe(float64(duration))
			}()
			next(ctx)
		}
	}
}
