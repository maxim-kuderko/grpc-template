package main

import (
	"context"
	"fmt"
	gs "github.com/maxim-kuderko/graceful-shutdown"
	"github.com/maxim-kuderko/service-template/internal/initializers"
	"github.com/maxim-kuderko/service-template/internal/repositories/primary"
	"github.com/maxim-kuderko/service-template/internal/service"
	"github.com/maxim-kuderko/service-template/pkg/requests"
	"github.com/maxim-kuderko/service-template/pkg/responses"
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"net"
)

func main() {
	go fx.New(
		fx.NopLogger,
		fx.Provide(
			initializers.NewConfig,
			primary.NewCachedDB,
			service.NewService,
			newServer,
		),
		fx.Invoke(grpcInit),
	)
	gs.WaitForGrace()
}

func grpcInit(s TemplateServer, v *viper.Viper) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", v.GetString(`GRPC_SERVER_PORT`)))
	if err != nil {
		panic(err)
	}
	serv := grpc.NewServer()
	go func() {
		defer serv.GracefulStop()
		gs.ShuttingDownHook()
	}()
	RegisterTemplateServer(serv, s)
	if err := serv.Serve(lis); err != nil {
		panic(err)
	}
}

type server struct {
	s *service.Service
	UnimplementedTemplateServer
}

func newServer(s *service.Service) TemplateServer {
	return &server{
		s: s,
	}
}

func (h *server) Get(ctx context.Context, request *GetRequest) (*GetResponse, error) {
	req := topPkgReq(request)
	req.WithContext(ctx)
	resp, err := h.s.Get(req)
	if err != nil {
		return nil, err
	}
	return topPkgResp(resp), nil
}

func topPkgReq(request *GetRequest) requests.Get {
	return requests.Get{
		BaseRequest: requests.BaseRequest{},
		Key:         request.Key,
	}
}

func topPkgResp(request responses.Get) *GetResponse {
	return &GetResponse{
		Value: request.Value,
	}
}
