package migrations

import (
	"fmt"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// resolve tele_user collection id for relation fields
		teleUserCol, err := app.FindCollectionByNameOrId("tele_user")
		if err != nil {
			return fmt.Errorf("tele_user collection not found: %w", err)
		}

		// Create bot_info collection (aligned with docs/specs/bot_info_table.md)
		if _, err := app.FindCollectionByNameOrId("bot_info"); err == nil {
			// already exists, skip
		} else {
			botInfo := core.NewBaseCollection("bot_info")
			// user relation -> tele_user.id
			botInfo.Fields.Add(&core.RelationField{ Name: "user", CollectionId: teleUserCol.Id, MaxSelect: 1, Required: true })
			botInfo.Fields.Add(&core.TextField{ Name: "bot_name", Required: true, Max: 255 })
			botInfo.Fields.Add(&core.TextField{ Name: "bot_token", Required: true, Max: 255 })
			botInfo.Fields.Add(&core.TextField{ Name: "bot_status", Max: 32 })
			botInfo.Fields.Add(&core.DateField{ Name: "start_time" })
			botInfo.Fields.Add(&core.DateField{ Name: "end_time" })
			botInfo.Fields.Add(&core.NumberField{ Name: "running_days", Min: nil, Max: nil })
			botInfo.Fields.Add(&core.NumberField{ Name: "total_messages", Min: nil, Max: nil })
			botInfo.Fields.Add(&core.NumberField{ Name: "total_users", Min: nil, Max: nil })
			botInfo.Fields.Add(&core.NumberField{ Name: "total_groups", Min: nil, Max: nil })
			botInfo.Fields.Add(&core.NumberField{ Name: "avg_messages_per_day", Min: nil, Max: nil })
			botInfo.Fields.Add(&core.DateField{ Name: "create_time" })
			botInfo.Fields.Add(&core.DateField{ Name: "update_time" })
			if err := app.Save(botInfo); err != nil { return fmt.Errorf("save bot_info: %w", err) }
		}

		// Create operation_logs collection (aligned with docs/specs/bot_log_table.md)
		if _, err := app.FindCollectionByNameOrId("operation_logs"); err == nil {
			// already exists, skip
		} else {
			logs := core.NewBaseCollection("operation_logs")
			// user relation -> tele_user.id
			logs.Fields.Add(&core.RelationField{ Name: "user", CollectionId: teleUserCol.Id, MaxSelect: 1, Required: true })
			logs.Fields.Add(&core.TextField{ Name: "bot_id" })
			logs.Fields.Add(&core.TextField{ Name: "operation_type", Required: true, Max: 64 })
			logs.Fields.Add(&core.DateField{ Name: "operation_time" })
			logs.Fields.Add(&core.TextField{ Name: "operation_details" })
			logs.Fields.Add(&core.DateField{ Name: "create_time" })
			if err := app.Save(logs); err != nil { return fmt.Errorf("save operation_logs: %w", err) }
		}

		return nil
	}, func(app core.App) error {
		// down: drop both collections if exists
		if c, err := app.FindCollectionByNameOrId("bot_info"); err == nil {
			if err := app.Delete(c); err != nil { return fmt.Errorf("delete bot_info: %w", err) }
		}
		if c, err := app.FindCollectionByNameOrId("operation_logs"); err == nil {
			if err := app.Delete(c); err != nil { return fmt.Errorf("delete operation_logs: %w", err) }
		}
		return nil
	})
}