package main

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/maxim-kuderko/service-template/internal/initializers"
	"github.com/maxim-kuderko/service-template/internal/repositories/primary"
	"github.com/maxim-kuderko/service-template/internal/repositories/secondary"
	"github.com/maxim-kuderko/service-template/internal/service"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"net/http"
)

func main() {
	app := fx.New(
		// Provide all the constructors we need, which teaches Fx how we'd like to
		// construct the *log.Logger, http.Handler, and *http.ServeMux types.
		// Remember that constructors are called lazily, so this block doesn't do
		// much on its own.
		fx.Provide(
			initializers.NewConfig,
			initializers.NewMetricsAndTracer,
			primary.NewCachedDB,
			secondary.NewCachedDB,
			service.NewService,
			newHandler,
			router,
		),
		// Since constructors are called lazily, we need some invocations to
		// kick-start our application. In this case, we'll use Register. Since it
		// depends on an http.Handler and *http.ServeMux, calling it requires Fx
		// to build those types using the constructors above. Since we call
		// NewMux, we also register Lifecycle hooks to start and stop an HTTP
		// server.
		fx.Invoke(webserver),
	)

	if err := app.Start(context.Background()); err != nil {
		panic(err)
	}
}

func router(h *handler) *httprouter.Router {
	router := httprouter.New()
	router.GET("/get", h.Get)
	return router
}

func webserver(r *httprouter.Router, v *viper.Viper) {
	http.ListenAndServe(fmt.Sprintf(`:%s`, v.GetString(`HTTP_SERVER_PORT`)), r)
}

type handler struct {
	s *service.Service
}

func newHandler(s *service.Service) *handler {
	return &handler{
		s: s,
	}
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}
