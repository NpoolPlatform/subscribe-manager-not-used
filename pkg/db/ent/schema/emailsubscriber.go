package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/NpoolPlatform/subscribe-manager/pkg/db/mixin"
	"github.com/google/uuid"
)

// EmailSubscriber holds the schema definition for the EmailSubscriber entity.
type EmailSubscriber struct {
	ent.Schema
}

func (EmailSubscriber) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
	}
}

// Fields of the EmailSubscriber.
func (EmailSubscriber) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).Unique(),
		field.UUID("app_id", uuid.UUID{}),
		field.String("email_address"),
	}
}

// Edges of the EmailSubscriber.
func (EmailSubscriber) Edges() []ent.Edge {
	return nil
}

func (EmailSubscriber) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("app_id", "email_address").Unique(),
	}
}
