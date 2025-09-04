package migrations

import (
	"fmt"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(func(app core.App) error {
		collection := core.NewBaseCollection("telegram_index")

		// Add fields using the new API
		collection.Fields.Add(&core.TextField{
			Name:     "id",
			Required: true,
		})

		collection.Fields.Add(&core.SelectField{
			Name:     "type",
			Required: true,
			Values:    []string{"group", "channel", "bot"},
			MaxSelect: 1,
		})

		collection.Fields.Add(&core.TextField{
			Name:     "title",
			Required: true,
		})

		collection.Fields.Add(&core.TextField{
			Name: "username",
		})

		collection.Fields.Add(&core.TextField{
			Name: "description",
		})

		collection.Fields.Add(&core.TextField{
			Name: "first_name",
		})

		collection.Fields.Add(&core.TextField{
			Name: "last_name",
		})

		collection.Fields.Add(&core.BoolField{
			Name: "is_verified",
		})

		collection.Fields.Add(&core.BoolField{
			Name: "is_restricted",
		})

		collection.Fields.Add(&core.BoolField{
			Name: "is_scam",
		})

		collection.Fields.Add(&core.BoolField{
			Name: "is_fake",
		})

		collection.Fields.Add(&core.TextField{
			Name: "language_code",
		})

		collection.Fields.Add(&core.NumberField{
			Name: "members_count",
			Min:  types.Pointer(0.0),
		})

		collection.Fields.Add(&core.NumberField{
			Name: "active_members",
			Min:  types.Pointer(0.0),
		})

		collection.Fields.Add(&core.NumberField{
			Name: "online_count",
			Min:  types.Pointer(0.0),
		})

		collection.Fields.Add(&core.DateField{
			Name: "created_at",
		})

		collection.Fields.Add(&core.DateField{
			Name: "updated_at",
		})

		collection.Fields.Add(&core.DateField{
			Name: "last_activity",
		})

		collection.Fields.Add(&core.NumberField{
			Name: "message_count",
			Min:  types.Pointer(0.0),
		})

		collection.Fields.Add(&core.NumberField{
			Name: "avg_message_per_day",
			Min:  types.Pointer(0.0),
		})

		collection.Fields.Add(&core.NumberField{
			Name: "growth_rate",
		})

		collection.Fields.Add(&core.JSONField{
			Name: "tags",
		})

		collection.Fields.Add(&core.JSONField{
			Name: "keywords",
		})

		collection.Fields.Add(&core.JSONField{
			Name: "content_types",
		})

		collection.Fields.Add(&core.JSONField{
			Name: "topics",
		})

		collection.Fields.Add(&core.JSONField{
			Name: "categories",
		})

		collection.Fields.Add(&core.TextField{
			Name: "language",
		})

		collection.Fields.Add(&core.TextField{
			Name: "country",
		})

		collection.Fields.Add(&core.NumberField{
			Name: "content_quality",
			Min:  types.Pointer(1.0),
			Max:  types.Pointer(5.0),
		})

		collection.Fields.Add(&core.JSONField{
			Name: "admin_list",
		})

		collection.Fields.Add(&core.TextField{
			Name: "creator_id",
		})

		collection.Fields.Add(&core.JSONField{
			Name: "permissions",
		})

		collection.Fields.Add(&core.TextField{
			Name: "invite_link",
		})

		collection.Fields.Add(&core.TextField{
			Name: "linked_chat_id",
		})

		collection.Fields.Add(&core.NumberField{
			Name: "slow_mode_delay",
			Min:  types.Pointer(0.0),
		})

		collection.Fields.Add(&core.BoolField{
			Name: "has_protected_content",
		})

		collection.Fields.Add(&core.BoolField{
			Name: "can_set_sticker_set",
		})

		collection.Fields.Add(&core.NumberField{
			Name: "overall_score",
			Min:  types.Pointer(0.0),
			Max:  types.Pointer(10.0),
		})

		collection.Fields.Add(&core.NumberField{
			Name: "popularity_score",
			Min:  types.Pointer(0.0),
			Max:  types.Pointer(10.0),
		})

		collection.Fields.Add(&core.NumberField{
			Name: "activity_score",
			Min:  types.Pointer(0.0),
			Max:  types.Pointer(10.0),
		})

		collection.Fields.Add(&core.NumberField{
			Name: "content_score",
			Min:  types.Pointer(0.0),
			Max:  types.Pointer(10.0),
		})

		collection.Fields.Add(&core.NumberField{
			Name: "relevance_score",
			Min:  types.Pointer(0.0),
			Max:  types.Pointer(10.0),
		})

		collection.Fields.Add(&core.NumberField{
			Name: "authority_score",
			Min:  types.Pointer(0.0),
			Max:  types.Pointer(10.0),
		})

		collection.Fields.Add(&core.DateField{
			Name: "last_calculated",
		})

		collection.Fields.Add(&core.BoolField{
			Name: "can_join_groups",
		})

		collection.Fields.Add(&core.BoolField{
			Name: "can_read_all_group_messages",
		})

		collection.Fields.Add(&core.BoolField{
			Name: "supports_inline_queries",
		})

		collection.Fields.Add(&core.TextField{
			Name: "bot_type",
		})

		collection.Fields.Add(&core.JSONField{
			Name: "commands",
		})

		collection.Fields.Add(&core.JSONField{
			Name: "metadata",
		})

		collection.Fields.Add(&core.JSONField{
			Name: "custom_fields",
		})

		collection.Fields.Add(&core.DateField{
			Name: "indexed_at",
		})

		collection.Fields.Add(&core.TextField{
			Name: "data_source",
		})

		collection.Fields.Add(&core.NumberField{
			Name: "version",
			Min:  types.Pointer(1.0),
		})

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("telegram_index")
		if err != nil {
			return fmt.Errorf("failed to find collection: %w", err)
		}

		return app.Delete(collection)
	})
}