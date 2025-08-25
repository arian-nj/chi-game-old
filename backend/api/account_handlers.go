package api

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	accountv1 "github.com/arian-nj/chibazi/backend/gen/account/v1"
)

func (app *ApiApplication) GetMe(
	ctx context.Context,
	req *connect.Request[accountv1.GetMeRequest],
) (*connect.Response[accountv1.GetMeResponse], error) {
	perosnRow := app.AuthenticateHeader(ctx, req.Header())
	if perosnRow == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("can't get user"))
	}

	return connect.NewResponse(&accountv1.GetMeResponse{
		Account: &accountv1.Account{
			Id: int64(perosnRow.ID),
		},
	}), nil
}
