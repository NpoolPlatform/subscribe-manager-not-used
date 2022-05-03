package emailsubscriber

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	cruder "github.com/NpoolPlatform/libent-cruder/pkg/cruder"
	npool "github.com/NpoolPlatform/message/npool/subscribemgr"
	"github.com/NpoolPlatform/subscribe-manager/pkg/test-init" //nolint

	constant "github.com/NpoolPlatform/subscribe-manager/pkg/const"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
)

func init() {
	if runByGithubAction, err := strconv.ParseBool(os.Getenv("RUN_BY_GITHUB_ACTION")); err == nil && runByGithubAction {
		return
	}
	if err := testinit.Init(); err != nil {
		fmt.Printf("cannot init test stub: %v\n", err)
	}
}

func TestCRUD(t *testing.T) {
	subscriber := npool.EmailSubscriber{
		AppID:        uuid.New().String(),
		EmailAddress: "kikakkz@hotmail.com",
	}

	schema, err := New(context.Background(), nil)
	assert.Nil(t, err)

	info, err := schema.Create(context.Background(), &subscriber)
	if assert.Nil(t, err) {
		if assert.NotEqual(t, info.ID, uuid.UUID{}.String()) {
			subscriber.ID = info.ID
		}
		assert.Equal(t, info, &subscriber)
	}

	subscriber.ID = info.ID

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	infos, total, err := schema.Rows(context.Background(),
		cruder.NewConds().WithCond(constant.FieldID, cruder.EQ, info.ID),
		0, 0)
	if assert.Nil(t, err) {
		assert.Equal(t, total, 1)
		assert.Equal(t, infos[0], &subscriber)
	}

	subscriber1 := &npool.EmailSubscriber{
		AppID:        subscriber.AppID,
		EmailAddress: "zhaoyubin@npool.cc",
	}
	subscriber2 := &npool.EmailSubscriber{
		AppID:        subscriber.AppID,
		EmailAddress: "zhaoyubin1@npool.cc",
	}

	schema, err = New(context.Background(), nil)
	assert.Nil(t, err)

	infos, err = schema.CreateBulk(context.Background(), []*npool.EmailSubscriber{subscriber1, subscriber2})
	if assert.Nil(t, err) {
		assert.Equal(t, len(infos), 2)
		assert.NotEqual(t, infos[0].ID, uuid.UUID{}.String())
		assert.NotEqual(t, infos[1].ID, uuid.UUID{}.String())
	}
}
