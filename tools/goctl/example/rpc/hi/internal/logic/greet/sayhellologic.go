package greetlogic

import (
	"context"

	"go-kit/tools/goctl/example/rpc/hi/internal/svc"
	"go-kit/tools/goctl/example/rpc/hi/pb/hi"

	"go-kit/core/logx"
)

type SayHelloLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSayHelloLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SayHelloLogic {
	return &SayHelloLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SayHelloLogic) SayHello(in *hi.HelloReq) (*hi.HelloResp, error) {
	// todo: add your logic here and delete this line

	return &hi.HelloResp{}, nil
}
