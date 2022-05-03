package emailsubscriber

import (
	"context"
	"fmt"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	"github.com/NpoolPlatform/subscribe-manager/pkg/db/ent"
	"github.com/NpoolPlatform/subscribe-manager/pkg/db/ent/emailsubscriber"

	constant "github.com/NpoolPlatform/subscribe-manager/pkg/const"
	"github.com/NpoolPlatform/subscribe-manager/pkg/db"

	npool "github.com/NpoolPlatform/message/npool/subscribemgr"

	"github.com/google/uuid"
)

type EmailSubscriber struct {
	*db.Entity
}

func New(ctx context.Context, tx *ent.Tx) (*EmailSubscriber, error) {
	e, err := db.NewEntity(ctx, tx)
	if err != nil {
		return nil, fmt.Errorf("fail create entity: %v", err)
	}

	return &EmailSubscriber{
		Entity: e,
	}, nil
}

func (s *EmailSubscriber) rowToObject(row *ent.EmailSubscriber) *npool.EmailSubscriber {
	return &npool.EmailSubscriber{
		ID:           row.ID.String(),
		AppID:        row.AppID.String(),
		EmailAddress: row.EmailAddress,
	}
}

func (s *EmailSubscriber) Create(ctx context.Context, in *npool.EmailSubscriber) (*npool.EmailSubscriber, error) {
	var info *ent.EmailSubscriber
	var err error

	err = db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		info, err = s.Tx.EmailSubscriber.Create().
			SetAppID(uuid.MustParse(in.GetAppID())).
			SetEmailAddress(in.GetEmailAddress()).
			Save(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail create email subscriber: %v", err)
	}

	return s.rowToObject(info), nil
}

func (s *EmailSubscriber) CreateBulk(ctx context.Context, in []*npool.EmailSubscriber) ([]*npool.EmailSubscriber, error) {
	rows := []*ent.EmailSubscriber{}
	var err error

	err = db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		bulk := make([]*ent.EmailSubscriberCreate, len(in))
		for i, info := range in {
			bulk[i] = s.Tx.EmailSubscriber.Create().
				SetAppID(uuid.MustParse(info.GetAppID())).
				SetEmailAddress(info.GetEmailAddress())
		}
		rows, err = s.Tx.EmailSubscriber.CreateBulk(bulk...).Save(_ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("fail create email subscribers: %v", err)
	}

	infos := []*npool.EmailSubscriber{}
	for _, row := range rows {
		infos = append(infos, s.rowToObject(row))
	}

	return infos, nil
}

func (s *EmailSubscriber) queryFromConds(conds cruder.Conds) (*ent.EmailSubscriberQuery, error) { //nolint
	stm := s.Tx.EmailSubscriber.Query()
	for k, v := range conds {
		switch k {
		case constant.FieldID:
			id, err := cruder.AnyTypeUUID(v.Val)
			if err != nil {
				return nil, fmt.Errorf("invalid id: %v", err)
			}
			stm = stm.Where(emailsubscriber.ID(id))
		case constant.EmailSubscriberFieldAppID:
			id, err := cruder.AnyTypeUUID(v.Val)
			if err != nil {
				return nil, fmt.Errorf("invalid app id: %v", err)
			}
			stm = stm.Where(emailsubscriber.AppID(id))
		case constant.EmailSubscriberFieldEmailAddress:
			value, err := cruder.AnyTypeString(v.Val)
			if err != nil {
				return nil, fmt.Errorf("invalid daily email subscriber value: %v", err)
			}
			stm = stm.Where(emailsubscriber.EmailAddress(value))
		default:
			return nil, fmt.Errorf("invalid email subscriber field")
		}
	}

	return stm, nil
}

func (s *EmailSubscriber) Rows(ctx context.Context, conds cruder.Conds, offset, limit int) ([]*npool.EmailSubscriber, int, error) {
	rows := []*ent.EmailSubscriber{}
	var total int

	err := db.WithTx(ctx, s.Tx, func(_ctx context.Context) error {
		stm, err := s.queryFromConds(conds)
		if err != nil {
			return fmt.Errorf("fail construct stm: %v", err)
		}

		total, err = stm.Count(_ctx)
		if err != nil {
			return fmt.Errorf("fail count email subscriber: %v", err)
		}

		rows, err = stm.Order(ent.Desc(emailsubscriber.FieldUpdatedAt)).Offset(offset).Limit(limit).All(_ctx)
		if err != nil {
			return fmt.Errorf("fail query email subscriber: %v", err)
		}

		return nil
	})
	if err != nil {
		return nil, 0, fmt.Errorf("fail get email subscriber: %v", err)
	}

	infos := []*npool.EmailSubscriber{}
	for _, row := range rows {
		infos = append(infos, s.rowToObject(row))
	}

	return infos, total, nil
}
