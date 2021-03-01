package main

import (
	"context"
	"fmt"
	"github.com/fasthttp/router"
	jsoniter "github.com/json-iterator/go"
	"github.com/maxim-kuderko/service-template/internal/initializers"
	"github.com/maxim-kuderko/service-template/internal/repositories/primary"
	"github.com/maxim-kuderko/service-template/internal/service"
	"github.com/maxim-kuderko/service-template/pkg/requests"
	"github.com/maxim-kuderko/service-template/pkg/responses"
	"github.com/opentracing/opentracing-go/log"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.NopLogger,
		fx.Provide(
			initializers.NewConfig,
			initializers.NewMetrics,
			primary.NewCachedDB,
			service.NewService,
			newHandler,
			route,
		),
		fx.Invoke(webserver),
	)

	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
}

func route(h *handler) *router.Router {
	router := router.New()
	router.GET("/health", h.Health)
	router.POST("/get", h.Get)
	return router
}

func webserver(r *router.Router, v *viper.Viper) {
	log.Error(fasthttp.ListenAndServe(fmt.Sprintf(`:%s`, v.GetString(`HTTP_SERVER_PORT`)), r.Handler))
}

type handler struct {
	s *service.Service
}

func newHandler(s *service.Service) *handler {
	return &handler{
		s: s,
	}
}
func (h *handler) Health(ctx *fasthttp.RequestCtx) {

}

func (h *handler) Get(ctx *fasthttp.RequestCtx) {
	var req requests.Get
	if err := parser(ctx, &req); err != nil {
		return
	}
	resp, err := h.s.Get(req)
	response(ctx, resp, err) // nolint
	return
}

func parser(c *fasthttp.RequestCtx, req requests.BaseRequester) error {
	err := jsoniter.ConfigFastest.Unmarshal(c.PostBody(), &req)
	if err != nil {
		c.SetStatusCode(fasthttp.StatusBadRequest)
		jsoniter.ConfigFastest.NewEncoder(c).Encode(err)
		return err
	}
	req.WithContext(c)
	return nil
}

func response(c *fasthttp.RequestCtx, resp responses.BaseResponser, err error) error {
	c.SetContentType(`application/json`)
	if err != nil {
		c.SetStatusCode(fasthttp.StatusInternalServerError)
		err := jsoniter.ConfigFastest.NewEncoder(c).Encode(err)
		return err
	}
	c.SetStatusCode(resp.ResponseStatusCode())
	if err := jsoniter.ConfigFastest.NewEncoder(c).Encode(resp); err != nil {
		c.SetStatusCode(fasthttp.StatusInternalServerError)
		return jsoniter.ConfigFastest.NewEncoder(c).Encode(err)
	}
	return nil
}
