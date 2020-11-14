package api

import (
	"context"

	"github.com/craftcms/nitro/pkg/protob"
)

func NewAPI() *API {
	return &API{}
}

type API struct {}

func (a *API) Ping(ctx context.Context, request *protob.PingRequest) (*protob.PingResponse, error) {
	return &protob.PingResponse{Pong: "pong"}, nil
}
