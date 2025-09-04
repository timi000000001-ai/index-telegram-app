package migrations

import (
	"fmt"
	"time"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

// Seed initial test data for tele_user, bot_info, operation_logs, telegram_index
func init() {
	m.Register(func(app core.App) error {
		return app.RunInTransaction(func(txApp core.App) error {
			// --- tele_user seed records ---
			teleUserCol, err := txApp.FindCollectionByNameOrId("tele_user")
			if err != nil {
				return fmt.Errorf("tele_user collection not found: %w", err)
			}

			ensureTeleUser := func(tgUserID, username, firstName, lastName string, isBot bool) (*core.Record, error) {
				rec, _ := txApp.FindFirstRecordByData("tele_user", "tg_user_id", tgUserID)
				if rec != nil {
					return rec, nil
				}
				rec = core.NewRecord(teleUserCol)
				rec.Set("tg_user_id", tgUserID)
				rec.Set("is_bot", isBot)
				rec.Set("first_name", firstName)
				rec.Set("last_name", lastName)
				rec.Set("username", username)
				rec.Set("language_code", "zh-CN")
				rec.Set("is_premium", false)
				rec.Set("create_time", time.Now())
				rec.Set("update_time", time.Now())
				if err := txApp.Save(rec); err != nil {
					return nil, fmt.Errorf("save tele_user %s: %w", tgUserID, err)
				}
				return rec, nil
			}

			u1, err := ensureTeleUser("1001", "alice", "Alice", "W", true)
			if err != nil {
				return err
			}
			u2, err := ensureTeleUser("1002", "bot_demo", "Demo", "Bot", true)
			if err != nil {
				return err
			}

			// --- bot_info seed ---
			botInfoCol, err := txApp.FindCollectionByNameOrId("bot_info")
			if err != nil {
				return fmt.Errorf("bot_info collection not found: %w", err)
			}
			if existing, _ := txApp.FindFirstRecordByData("bot_info", "bot_name", "seed_demo_bot"); existing == nil {
				rec := core.NewRecord(botInfoCol)
				rec.Set("user", u1.Id) // relation to tele_user
				rec.Set("bot_name", "seed_demo_bot")
				rec.Set("bot_token", "123456:TEST_TOKEN")
				rec.Set("bot_status", "active")
				rec.Set("start_time", time.Now().Add(-24*time.Hour))
				rec.Set("running_days", 1)
				rec.Set("total_messages", 42)
				rec.Set("total_users", 3)
				rec.Set("total_groups", 1)
				rec.Set("avg_messages_per_day", 42)
				rec.Set("create_time", time.Now().Add(-24*time.Hour))
				rec.Set("update_time", time.Now())
				if err := txApp.Save(rec); err != nil {
					return fmt.Errorf("save bot_info: %w", err)
				}
			}

			// --- operation_logs seed ---
			logsCol, err := txApp.FindCollectionByNameOrId("operation_logs")
			if err != nil {
				return fmt.Errorf("operation_logs collection not found: %w", err)
			}
			ensureLog := func(botID, opType, details string, userID string, t time.Time) error {
				// idempotency best-effort: check by bot_id + operation_type + operation_time unix
				key := fmt.Sprintf("%s|%s|%d", botID, opType, t.Unix())
				if _, err := txApp.FindFirstRecordByData("operation_logs", "operation_details", key); err == nil {
					return nil
				}
				rec := core.NewRecord(logsCol)
				rec.Set("user", userID)
				rec.Set("bot_id", botID)
				rec.Set("operation_type", opType)
				rec.Set("operation_time", t)
				rec.Set("operation_details", key)
				rec.Set("create_time", time.Now())
				return txApp.Save(rec)
			}
			if err := ensureLog("seed_demo_bot", "start", "", u1.Id, time.Now().Add(-2*time.Hour)); err != nil {
				return fmt.Errorf("seed log1: %w", err)
			}
			if err := ensureLog("seed_demo_bot", "message", "", u2.Id, time.Now().Add(-1*time.Hour)); err != nil {
				return fmt.Errorf("seed log2: %w", err)
			}

			// --- telegram_index seed ---
			idxCol, err := txApp.FindCollectionByNameOrId("telegram_index")
			if err != nil {
				return fmt.Errorf("telegram_index collection not found: %w", err)
			}
			ensureIndex := func(username, typ, title string) error {
				if _, err := txApp.FindFirstRecordByData("telegram_index", "username", username); err == nil {
					return nil
				}
				rec := core.NewRecord(idxCol)
				// set required custom text field "ext_id" as a stable key based on username
				rec.Set("ext_id", username)
				rec.Set("type", typ)
				rec.Set("title", title)
				rec.Set("username", username)
				rec.Set("is_verified", false)
				rec.Set("created_at", time.Now().Add(-72*time.Hour))
				rec.Set("updated_at", time.Now())
				rec.Set("indexed_at", time.Now())
				return txApp.Save(rec)
			}
			if err := ensureIndex("go_dev_group", "group", "Go 开发者群"); err != nil {
				return fmt.Errorf("seed index group: %w", err)
			}
			if err := ensureIndex("demo_bot_handle", "bot", "Demo Bot"); err != nil {
				return fmt.Errorf("seed index bot: %w", err)
			}

			return nil
		})
	}, func(app core.App) error {
		return app.RunInTransaction(func(txApp core.App) error {
			// --- Efficiently delete all seeded data ---

			// 1. Delete telegram_index records by stable username
			for _, username := range []string{"go_dev_group", "demo_bot_handle"} {
				if rec, _ := txApp.FindFirstRecordByData("telegram_index", "username", username); rec != nil {
					if err := txApp.Delete(rec); err != nil {
						return fmt.Errorf("failed to delete telegram_index %s: %w", username, err)
					}
				}
			}

			// 2. Delete operation_logs by bot_id
			// This is more robust than relying on timestamps.
			if records, err := txApp.FindRecordsByFilter("operation_logs", "bot_id = 'seed_demo_bot'", "", 0, 0); err == nil {
				for _, rec := range records {
					if err := txApp.Delete(rec); err != nil {
						return fmt.Errorf("failed to delete operation_log %s: %w", rec.Id, err)
					}
				}
			}

			// 3. Delete bot_info record
			if rec, _ := txApp.FindFirstRecordByData("bot_info", "bot_name", "seed_demo_bot"); rec != nil {
				if err := txApp.Delete(rec); err != nil {
					return fmt.Errorf("failed to delete bot_info: %w", err)
				}
			}

			// 4. Delete tele_user records by stable tg_user_id
			for _, userID := range []string{"1001", "1002"} {
				if rec, _ := txApp.FindFirstRecordByData("tele_user", "tg_user_id", userID); rec != nil {
					if err := txApp.Delete(rec); err != nil {
						return fmt.Errorf("failed to delete tele_user %s: %w", userID, err)
					}
				}
			}

			return nil
		})
	})
}