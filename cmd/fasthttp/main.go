package main

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/maxim-kuderko/service-template/internal/initializers"
	"github.com/maxim-kuderko/service-template/internal/repositories/primary"
	"github.com/maxim-kuderko/service-template/internal/repositories/secondary"
	"github.com/maxim-kuderko/service-template/internal/service"
	"github.com/maxim-kuderko/service-template/pkg/requests"
	"github.com/opentracing/opentracing-go/log"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"go.uber.org/fx"
	"net/http"
	"time"
)

func main() {
	app := fx.New(
		fx.NopLogger,
		fx.Provide(
			initializers.NewConfig,
			initializers.NewMetrics,
			primary.NewCachedDB,
			secondary.NewCachedDB,
			service.NewService,
			newHandler,
			router,
		),
		fx.Invoke(webserver),
	)

	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
}

func router(h *handler) *routing.Router {
	router := routing.New()
	router.Post("/get", h.Get)
	return router
}

func webserver(r *routing.Router, v *viper.Viper) {
	log.Error(fasthttp.ListenAndServe(fmt.Sprintf(`:%s`, v.GetString(`HTTP_SERVER_PORT`)), r.HandleRequest))
	time.Sleep(time.Minute)
}

type handler struct {
	s *service.Service
}

func newHandler(s *service.Service) *handler {
	return &handler{
		s: s,
	}
}

func (h *handler) Get(c *routing.Context) error {
	var req requests.Get
	if err := parser(c, &req); err != nil {
		return nil
	}
	resp, err := h.s.Get(req)
	return response(c, resp, err)
}

func parser(c *routing.Context, req requests.BaseRequester) error {
	c.Serialize = jsoniter.Marshal
	err := jsoniter.ConfigFastest.Unmarshal(c.PostBody(), &req)
	if err != nil {
		c.SetStatusCode(http.StatusBadRequest)
		c.WriteData(err)
		return err
	}
	req.WithContext(c)
	return nil
}

func response(c *routing.Context, resp interface{}, err error) error {
	if err != nil {
		c.SetStatusCode(http.StatusInternalServerError)
		return c.WriteData(err)
	}
	if err := c.WriteData(resp); err != nil {
		c.SetStatusCode(http.StatusInternalServerError)
		return c.WriteData(err)
	}
	return nil
}
