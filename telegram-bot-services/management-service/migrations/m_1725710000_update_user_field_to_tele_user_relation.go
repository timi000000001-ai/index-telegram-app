package migrations

import (
	"fmt"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// This migration ensures bot_info.user and operation_logs.user are RelationField -> tele_user
func init() {
	m.Register(func(app core.App) error {
		teleUserCol, err := app.FindCollectionByNameOrId("tele_user")
		if err != nil {
			return fmt.Errorf("tele_user collection not found: %w", err)
		}

		updateUserField := func(colName string) error {
			col, err := app.FindCollectionByNameOrId(colName)
			if err != nil {
				// collection may not exist yet; ignore
				return nil
			}
			// remove any existing user field first to avoid duplicate names
			col.Fields.RemoveByName("user")
			// Upsert the user field as a RelationField (MaxSelect=1, Required=true)
			col.Fields.Add(&core.RelationField{
				Name:         "user",
				CollectionId: teleUserCol.Id,
				MaxSelect:    1,
				Required:     true,
			})
			if err := app.Save(col); err != nil {
				return fmt.Errorf("save %s: %w", colName, err)
			}
			return nil
		}

		if err := updateUserField("bot_info"); err != nil {
			return err
		}
		if err := updateUserField("operation_logs"); err != nil {
			return err
		}
		return nil
	}, func(app core.App) error {
		// down: revert the user field to TextField for both collections
		revertUserField := func(colName string) error {
			col, err := app.FindCollectionByNameOrId(colName)
			if err != nil {
				return nil
			}
			// remove relation field if exists
			col.Fields.RemoveByName("user")
			col.Fields.Add(&core.TextField{
				Name: "user",
			})
			return app.Save(col)
		}

		_ = revertUserField("bot_info")
		_ = revertUserField("operation_logs")
		return nil
	})
}