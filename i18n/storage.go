package i18n

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
)

// UpsertTranslations inserts or updates translations in bulk using pgx.
func UpsertTranslations(ctx context.Context, conn *pgx.Conn, translations []Translation) error {
	if len(translations) == 0 {
		return nil
	}

	// Create the temporary table if it doesn't exist, including tooltip
	_, err := conn.Exec(ctx, `
		CREATE TEMPORARY TABLE IF NOT EXISTS ui_translations_temp (
			user_id UUID,
			key_path TEXT,
			lang TEXT,
			value TEXT,
			tooltip TEXT,
			updated_at TIMESTAMP
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create temporary table: %w", err)
	}

	// Defer cleanup: Drop the temporary table after the function completes
	defer func() {
		_, err := conn.Exec(ctx, "DROP TABLE IF EXISTS ui_translations_temp")
		if err != nil {
			// Log the error but do not return it from the deferred function
			fmt.Printf("failed to drop temporary table: %v\n", err)
		}
	}()

	// Prepare rows to be inserted, now including tooltip
	rows := make([][]any, 0, len(translations))
	now := time.Now()

	for _, t := range translations {
		var userID any = nil
		if t.UserID != nil {
			userID = *t.UserID
		}
		rows = append(rows, []any{
			userID,
			t.KeyPath,
			t.Lang,
			t.Value,
			t.ToolTip, // Include the tooltip field
			now,
		})
	}

	// Perform COPY INTO the temporary table
	_, err = conn.CopyFrom(
		ctx,
		pgx.Identifier{"ui_translations_temp"},
		[]string{"user_id", "key_path", "lang", "value", "tooltip", "updated_at"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return fmt.Errorf("copy to temp table failed: %w", err)
	}

	// Perform the UPSERT operation using the data from the temporary table
	_, err = conn.Exec(ctx, `
		INSERT INTO ui_translations (user_id, key_path, lang, value, tooltip, updated_at)
		SELECT user_id, key_path, lang, value, tooltip, updated_at FROM ui_translations_temp
		ON CONFLICT (user_id, key_path, lang)
		DO UPDATE SET value = EXCLUDED.value, tooltip = EXCLUDED.tooltip, updated_at = EXCLUDED.updated_at;
	`)
	if err != nil {
		return fmt.Errorf("upsert from temp table failed: %w", err)
	}

	// No need to manually drop the table here, as it's handled by the deferred function
	return nil
}

// GetTranslation retrieves a translation with fallback using pgx.
func GetTranslation(ctx context.Context, conn *pgx.Conn, userID *string, keyPath, lang string) (string, error) {
	var value string
	query := `
		SELECT value FROM ui_translations
		WHERE (user_id = $1 OR user_id IS NULL)
		AND key_path = $2 AND lang = $3
		ORDER BY user_id NULLS LAST
		LIMIT 1
	`
	err := conn.QueryRow(ctx, query, userID, keyPath, lang).Scan(&value)
	return value, err
}

// ExportToFlatJSON retrieves all translations and returns a flat map using pgx, including tooltips.
func ExportToFlatJSON(ctx context.Context, conn *pgx.Conn, lang string, userID *string) (map[string]map[string]string, error) {
	const query = `
		SELECT key_path, value, tooltip FROM ui_translations
		WHERE lang = $1 AND (user_id = $2 OR user_id IS NULL)
		ORDER BY user_id NULLS LAST
	`

	rows, err := conn.Query(ctx, query, lang, userID)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	result := make(map[string]map[string]string, 128) // Preallocate with a reasonable initial capacity

	for {
		key, value, tooltip := "", "", ""
		if !rows.Next() {
			break
		}
		if err = rows.Scan(&key, &value, &tooltip); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		// Store both value and tooltip in the result map
		result[key] = map[string]string{
			"value":   value,
			"tooltip": tooltip,
		}
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("row iteration error: %w", rows.Err())
	}

	return result, nil
}
