package cloudrun

import (
	"github.com/faryoo/cloudrun-wechat/context"
	"github.com/faryoo/cloudrun-wechat/server"
	"net/http"
)

type Wechat struct {
	ctx *context.Context
}

func NewWechat() *Wechat {
	return &Wechat{}
}

func (wc *Wechat) GetServer(req *http.Request, writer http.ResponseWriter) *server.Server {
	srv := server.NewServer(wc.ctx)
	srv.Request = req
	srv.Writer = writer
	return srv
}
