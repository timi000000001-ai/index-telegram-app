package migrations

import (
	"fmt"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// Create the tele_user collection based on docs/specs/user_table.md
		collection := core.NewBaseCollection("tele_user")

		// Required identifiers and core attributes
		collection.Fields.Add(&core.TextField{
			Name:     "tg_user_id", // Telegram user id; use text to avoid numeric precision issues
			Required: true,
		})
		collection.Fields.Add(&core.BoolField{ // whether this user is a bot account
			Name:     "is_bot",
			Required: true,
		})

		// Profile names and username
		collection.Fields.Add(&core.TextField{
			Name:     "first_name",
			Required: true,
			Max:      64,
		})
		collection.Fields.Add(&core.TextField{
			Name: "last_name",
			Max:  64,
		})
		collection.Fields.Add(&core.TextField{
			Name: "username",
			Max:  32,
		})

		// Locale and premium
		collection.Fields.Add(&core.TextField{
			Name: "language_code",
			Max:  10,
		})
		collection.Fields.Add(&core.BoolField{
			Name: "is_premium",
		})

		// Bot capability flags (optional for normal users)
		collection.Fields.Add(&core.BoolField{ Name: "can_join_groups" })
		collection.Fields.Add(&core.BoolField{ Name: "can_read_all_group_messages" })
		collection.Fields.Add(&core.BoolField{ Name: "supports_inline_queries" })
		collection.Fields.Add(&core.BoolField{ Name: "added_to_attachment_menu" })

		// Profile extras
		collection.Fields.Add(&core.TextField{ Name: "bio" })
		collection.Fields.Add(&core.URLField{ Name: "photo_url" })

		// Raw JSON snapshot from Bot API
		collection.Fields.Add(&core.JSONField{ Name: "raw_json" })

		// Audit times (note: PocketBase also provides built-in created/updated timestamps)
		collection.Fields.Add(&core.DateField{ Name: "create_time" })
		collection.Fields.Add(&core.DateField{ Name: "update_time" })

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("tele_user")
		if err != nil {
			return fmt.Errorf("failed to find collection: %w", err)
		}
		return app.Delete(collection)
	})
}