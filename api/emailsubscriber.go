// +build !codeanalysis

package api

import (
	"context"
	"fmt"

	"github.com/NpoolPlatform/go-service-framework/pkg/logger"
	constant "github.com/NpoolPlatform/subscribe-manager/pkg/const"
	crud "github.com/NpoolPlatform/subscribe-manager/pkg/crud/emailsubscriber"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	commonnpool "github.com/NpoolPlatform/message/npool"
	npool "github.com/NpoolPlatform/message/npool/subscribemgr"

	"github.com/badoux/checkmail"
	"github.com/google/uuid"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

func (s *Server) CreateEmailSubscriber(ctx context.Context, in *npool.CreateEmailSubscriberRequest) (*npool.CreateEmailSubscriberResponse, error) {
	if _, err := uuid.Parse(in.GetInfo().GetAppID()); err != nil {
		logger.Sugar().Errorf("invalid request app id: %v", err)
		return &npool.CreateEmailSubscriberResponse{}, status.Error(codes.Internal, err.Error())
	}

	if err := checkmail.ValidateFormat(in.GetInfo().GetEmailAddress()); err != nil {
		logger.Sugar().Errorf("invalid email address: %v", err)
		return &npool.CreateEmailSubscriberResponse{}, status.Error(codes.Internal, err.Error())
	}

	schema, err := crud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail create schema entity: %v", err)
		return &npool.CreateEmailSubscriberResponse{}, status.Error(codes.Internal, err.Error())
	}

	info, err := schema.Create(ctx, in.GetInfo())
	if err != nil {
		logger.Sugar().Errorf("fail create email subscriber: %v", err)
		return &npool.CreateEmailSubscriberResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateEmailSubscriberResponse{
		Info: info,
	}, nil
}

func (s *Server) CreateEmailSubscribers(ctx context.Context, in *npool.CreateEmailSubscribersRequest) (*npool.CreateEmailSubscribersResponse, error) {
	for _, info := range in.GetInfos() {
		if _, err := uuid.Parse(info.GetAppID()); err != nil {
			logger.Sugar().Errorf("invalid request app id: %v", err)
			return &npool.CreateEmailSubscribersResponse{}, status.Error(codes.Internal, err.Error())
		}
		if err := checkmail.ValidateFormat(info.GetEmailAddress()); err != nil {
			logger.Sugar().Errorf("invalid email address: %v", err)
			return &npool.CreateEmailSubscribersResponse{}, status.Error(codes.Internal, err.Error())
		}
	}

	schema, err := crud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail create schema entity: %v", err)
		return &npool.CreateEmailSubscribersResponse{}, status.Error(codes.Internal, err.Error())
	}

	infos, err := schema.CreateBulk(ctx, in.GetInfos())
	if err != nil {
		logger.Sugar().Errorf("fail create email subscribers: %v", err)
		return &npool.CreateEmailSubscribersResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.CreateEmailSubscribersResponse{
		Infos: infos,
	}, nil
}

func emailSubscriberCondsToConds(conds cruder.FilterConds) (cruder.Conds, error) {
	newConds := cruder.NewConds()

	for k, v := range conds {
		switch v.Op {
		case cruder.EQ:
		case cruder.GT:
		case cruder.LT:
		case cruder.LIKE:
		default:
			return nil, fmt.Errorf("invalid filter condition op")
		}

		switch k {
		case constant.FieldID:
			fallthrough //nolint
		case constant.EmailSubscriberFieldAppID:
			newConds = newConds.WithCond(k, v.Op, v.Val.GetStringValue())
		case constant.EmailSubscriberFieldEmailAddress:
			newConds = newConds.WithCond(k, v.Op, v.Val.GetStringValue())
		default:
			return nil, fmt.Errorf("invalid email subscriber field")
		}
	}

	return newConds, nil
}

func (s *Server) GetEmailSubscribers(ctx context.Context, in *npool.GetEmailSubscribersRequest) (*npool.GetEmailSubscribersResponse, error) {
	inConds := in.GetConds()

	if inConds == nil {
		inConds = map[string]*commonnpool.FilterCond{}
	}

	inConds[constant.EmailSubscriberFieldAppID] = &commonnpool.FilterCond{
		Op:  cruder.EQ,
		Val: structpb.NewStringValue(in.GetAppID()),
	}

	conds, err := emailSubscriberCondsToConds(inConds)
	if err != nil {
		logger.Sugar().Errorf("invalid email subscriber fields: %v", err)
		return &npool.GetEmailSubscribersResponse{}, status.Error(codes.Internal, err.Error())
	}

	schema, err := crud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail create schema entity: %v", err)
		return &npool.GetEmailSubscribersResponse{}, status.Error(codes.Internal, err.Error())
	}

	infos, total, err := schema.Rows(ctx, conds, int(in.GetOffset()), int(in.GetLimit()))
	if err != nil {
		logger.Sugar().Errorf("fail get email subscribers: %v", err)
		return &npool.GetEmailSubscribersResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetEmailSubscribersResponse{
		Infos: infos,
		Total: int32(total),
	}, nil
}

func (s *Server) GetAppEmailSubscribers(ctx context.Context, in *npool.GetAppEmailSubscribersRequest) (*npool.GetAppEmailSubscribersResponse, error) {
	inConds := in.GetConds()

	if inConds == nil {
		inConds = map[string]*commonnpool.FilterCond{}
	}

	inConds[constant.EmailSubscriberFieldAppID] = &commonnpool.FilterCond{
		Op:  cruder.EQ,
		Val: structpb.NewStringValue(in.GetTargetAppID()),
	}

	conds, err := emailSubscriberCondsToConds(inConds)
	if err != nil {
		logger.Sugar().Errorf("invalid email subscriber fields: %v", err)
		return &npool.GetAppEmailSubscribersResponse{}, status.Error(codes.Internal, err.Error())
	}

	schema, err := crud.New(ctx, nil)
	if err != nil {
		logger.Sugar().Errorf("fail create schema entity: %v", err)
		return &npool.GetAppEmailSubscribersResponse{}, status.Error(codes.Internal, err.Error())
	}

	infos, total, err := schema.Rows(ctx, conds, int(in.GetOffset()), int(in.GetLimit()))
	if err != nil {
		logger.Sugar().Errorf("fail get email subscribers: %v", err)
		return &npool.GetAppEmailSubscribersResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &npool.GetAppEmailSubscribersResponse{
		Infos: infos,
		Total: int32(total),
	}, nil
}
