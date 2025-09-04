package migrations

import (
	"fmt"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Fix telegram_index schema: rename custom field "id" (conflicts with system id) to "ext_id"
func init() {
	m.Register(func(app core.App) error {
		col, err := app.FindCollectionByNameOrId("telegram_index")
		if err != nil {
			return fmt.Errorf("telegram_index collection not found: %w", err)
		}
		// remove the conflicting custom field named "id" if present
		col.Fields.RemoveByName("id")
		// add a safer optional external id field
		col.Fields.Add(&core.TextField{
			Name: "ext_id",
		})
		if err := app.Save(col); err != nil {
			return fmt.Errorf("save telegram_index: %w", err)
		}
		return nil
	}, func(app core.App) error {
		col, err := app.FindCollectionByNameOrId("telegram_index")
		if err != nil {
			return nil
		}
		// revert: drop ext_id and re-add custom id field (required)
		col.Fields.RemoveByName("ext_id")
		col.Fields.Add(&core.TextField{ Name: "id", Required: true })
		return app.Save(col)
	})
}