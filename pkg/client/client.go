package client

import (
	"context"
	"fmt"
	"time"

	grpc2 "github.com/NpoolPlatform/go-service-framework/pkg/grpc"

	"github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/subscribemgr"

	constant "github.com/NpoolPlatform/subscribe-manager/pkg/message/const"
)

func do(ctx context.Context, fn func(_ctx context.Context, cli npool.SubscribeManagerClient) (cruder.Any, error)) (cruder.Any, error) {
	_ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	conn, err := grpc2.GetGRPCConn(constant.ServiceName, grpc2.GRPCTAG)
	if err != nil {
		return nil, fmt.Errorf("fail get email subscriber connection: %v", err)
	}
	defer conn.Close()

	cli := npool.NewSubscribeManagerClient(conn)

	return fn(_ctx, cli)
}

func GetEmailSubscribers(ctx context.Context, conds cruder.FilterConds) ([]*npool.EmailSubscriber, error) {
	info, err := do(ctx, func(_ctx context.Context, cli npool.SubscribeManagerClient) (cruder.Any, error) {
		// DO RPC CALL HERE WITH conds PARAMETER
		return []*npool.EmailSubscriber{}, nil
	})
	if err != nil {
		return nil, fmt.Errorf("fail get email subscriber: %v", err)
	}
	return info.([]*npool.EmailSubscriber), nil
}
