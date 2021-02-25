package requests

import "context"

type BaseRequest struct {
	ctx context.Context
}

type BaseRequester interface {
	Context() context.Context
	WithContext(ctx context.Context)
}

func NewBaseRequest(ctx context.Context) BaseRequest {
	return BaseRequest{ctx: ctx}
}

func (br *BaseRequest) Context() context.Context {
	return br.ctx
}
func (br *BaseRequest) WithContext(ctx context.Context) {
	br.ctx = ctx
}
