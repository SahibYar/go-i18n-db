package example

import (
	"context"
	"github.com/jackc/pgx/v5"
	"go-i18n-db/i18n"
	"path/filepath"
	"strings"
)

// LoadAndSave loads a JSON file, flattens it, and saves to DB.
func LoadAndSave(conn *pgx.Conn, filePath string, lang string, userID *string) error {
	// Step 1: Load and flatten the JSON
	flatMap, err := i18n.LoadAndFlatten(filePath)
	if err != nil {
		return err
	}

	// Step 2: Convert to []Translation
	translations := make([]i18n.Translation, 0, len(flatMap))
	for keyPath, value := range flatMap {
		translations = append(translations, i18n.Translation{
			UserID:  userID,
			Lang:    lang,
			KeyPath: keyPath,
			Value:   value,
		})
	}

	// Step 3: Bulk upsert to a database
	return i18n.UpsertTranslations(context.Background(), conn, translations)
}

// LoadAndSaveAutoLang Optional helper: load from file path and auto-extract language
func LoadAndSaveAutoLang(conn *pgx.Conn, filePath string, userID *string) error {
	// Guess language from filename like "en.json"
	base := filepath.Base(filePath)
	ext := filepath.Ext(base)
	lang := strings.TrimSuffix(base, ext)

	return LoadAndSave(conn, filePath, lang, userID)
}
