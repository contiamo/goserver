package test

import (
	"context"
	"errors"
)

func NewPingPongServer() PingPongServer {
	return &pingPongServer{}
}

type pingPongServer struct{}

func (srv *pingPongServer) Ping(ctx context.Context, req *PingReq) (*PingResp, error) {
	return &PingResp{Msg: req.Msg}, nil
}
func (srv *pingPongServer) Panic(ctx context.Context, req *PingReq) (*PingResp, error) {
	panic(errors.New(req.Msg))
	return &PingResp{Msg: req.Msg}, nil
}
