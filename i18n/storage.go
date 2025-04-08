package i18n

import (
	"database/sql"
	"fmt"
	"strings"
)

// UpsertTranslations inserts or updates translations in bulk.
func UpsertTranslations(db *sql.DB, translations []Translation) error {
	if len(translations) == 0 {
		return nil
	}

	query := `
		INSERT INTO ui_translations (user_id, key_path, lang, value, updated_at)
		VALUES %s
		ON CONFLICT (user_id, key_path, lang)
		DO UPDATE SET value = EXCLUDED.value, updated_at = NOW()
	`

	valueStrings := []string{}
	valueArgs := []interface{}{}
	for _, t := range translations {
		userIDPlaceholder := "NULL"
		if t.UserID != nil {
			valueArgs = append(valueArgs, *t.UserID)
			userIDPlaceholder = fmt.Sprintf("$%d", len(valueArgs))
		}
		valueArgs = append(valueArgs, t.KeyPath, t.Lang, t.Value)
		valueStrings = append(valueStrings, fmt.Sprintf("(%s, $%d, $%d, $%d, NOW())",
			userIDPlaceholder,
			len(valueArgs)-2, len(valueArgs)-1, len(valueArgs),
		))
	}

	stmt := fmt.Sprintf(query, strings.Join(valueStrings, ", "))
	_, err := db.Exec(stmt, valueArgs...)
	return err
}

// GetTranslation retrieves a translation with fallback.
func GetTranslation(db *sql.DB, userID *string, keyPath, lang string) (string, error) {
	var value string
	query := `
		SELECT value FROM ui_translations
		WHERE (user_id = $1 OR user_id IS NULL)
		AND key_path = $2 AND lang = $3
		ORDER BY user_id NULLS LAST
		LIMIT 1
	`
	err := db.QueryRow(query, userID, keyPath, lang).Scan(&value)
	return value, err
}

// ExportToJSON retrieves all translations and returns a flat map.
func ExportToJSON(db *sql.DB, lang string, userID *string) (map[string]string, error) {
	query := `
		SELECT key_path, value FROM ui_translations
		WHERE lang = $1 AND (user_id = $2 OR user_id IS NULL)
		ORDER BY user_id NULLS LAST
	`
	rows, err := db.Query(query, lang, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		result[key] = value
	}
	return result, nil
}
