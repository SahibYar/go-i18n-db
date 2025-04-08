package i18n

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Setup PostgreSQL connection (use an environment variable for connection string)
var connString = "postgresql://postgres:fuzzyfog33@localhost:5455/ui_db?search_path=development_ui&sslmode=disable"

// TestExportToJSON tests the ExportToJSON function
func TestExportToJSON(t *testing.T) {
	// Connect to the test database
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	// Generate valid UUIDs for user1 and user2
	user1ID := uuid.New()
	user2ID := uuid.New()

	// Insert test data into the ui_translations table
	_, err = conn.Exec(context.Background(), `
		INSERT INTO ui_translations (user_id, key_path, lang, value, updated_at)
		VALUES 
			($1, 'topbar.profile', 'en', 'Profile', NOW()), 
			(NULL, 'topbar.profile', 'es', 'Perfil', NOW()), 
			($1, 'footer.contact', 'en', 'Contact', NOW())
	`, user1ID)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Test: Fetch translations with user ID 'user1'
	translations, err := ExportToJSON(context.Background(), conn, "en", stringPtr(user1ID.String()))
	assert.NoError(t, err)
	assert.NotNil(t, translations)
	assert.Equal(t, "Profile", translations["topbar.profile"])
	assert.Equal(t, "Contact", translations["footer.contact"])

	// Test: Fetch translations without user ID (fallback to NULL user)
	translations, err = ExportToJSON(context.Background(), conn, "en", nil)
	assert.NoError(t, err)
	assert.NotNil(t, translations)
	assert.Equal(t, "", translations["topbar.profile"])

	// Cleanup test data
	defer func() {
		_, err = conn.Exec(context.Background(), `
		DELETE FROM ui_translations WHERE user_id IN ($1, $2)
	`, user1ID, user2ID)
		if err != nil {
			t.Fatalf("Failed to clean up test data: %v", err)
		}
	}()
}

func stringPtr(s string) *string {
	return &s
}

func TestUpsertTranslations(t *testing.T) {
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	// Generate valid UUIDs for user1 and user2
	user1ID := uuid.New()
	user2ID := uuid.New()

	// Clean up before tests
	defer func() {
		_, err = conn.Exec(context.Background(), `
			DELETE FROM ui_translations WHERE user_id IN ($1, $2)
		`, user1ID, user2ID)
		if err != nil {
			t.Fatalf("Failed to clean up test data: %v", err)
		}
	}()

	// Test 1: Insert new translations for user1 and user2
	translations := []Translation{
		{UserID: stringPtr(user1ID.String()), KeyPath: "topbar.profile", Lang: "en", Value: "Profile"},
		{UserID: stringPtr(user1ID.String()), KeyPath: "footer.contact", Lang: "en", Value: "Contact"},
		{UserID: stringPtr(user2ID.String()), KeyPath: "topbar.profile", Lang: "es", Value: "Perfil"},
		{UserID: stringPtr(user2ID.String()), KeyPath: "footer.contact", Lang: "es", Value: "Contacto"},
	}

	err = UpsertTranslations(context.Background(), conn, translations)
	assert.NoError(t, err)

	// Verify if the translations were inserted
	var value string
	err = conn.QueryRow(context.Background(), `
		SELECT value FROM ui_translations WHERE user_id = $1 AND key_path = $2 AND lang = $3
	`, user1ID, "topbar.profile", "en").Scan(&value)
	assert.NoError(t, err)
	assert.Equal(t, "Profile", value)

	// Test 2: Update existing translations for user1
	updatedTranslations := []Translation{
		{UserID: stringPtr(user1ID.String()), KeyPath: "topbar.profile", Lang: "en", Value: "Updated Profile"},
	}

	err = UpsertTranslations(context.Background(), conn, updatedTranslations)
	assert.NoError(t, err)

	// Verify if the translation was updated
	err = conn.QueryRow(context.Background(), `
		SELECT value FROM ui_translations WHERE user_id = $1 AND key_path = $2 AND lang = $3
	`, user1ID, "topbar.profile", "en").Scan(&value)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Profile", value)

	// Test 3: Bulk insertion of translations for user1 and user2
	bulkTranslations := []Translation{
		{UserID: stringPtr(user1ID.String()), KeyPath: "topbar.welcome", Lang: "en", Value: "Welcome"},
		{UserID: stringPtr(user2ID.String()), KeyPath: "topbar.welcome", Lang: "es", Value: "Bienvenido"},
		{UserID: stringPtr(user1ID.String()), KeyPath: "footer.privacy", Lang: "en", Value: "Privacy"},
		{UserID: stringPtr(user2ID.String()), KeyPath: "footer.privacy", Lang: "es", Value: "Privacidad"},
	}

	err = UpsertTranslations(context.Background(), conn, bulkTranslations)
	assert.NoError(t, err)

	// Verify bulk insertion
	err = conn.QueryRow(context.Background(), `
		SELECT value FROM ui_translations WHERE user_id = $1 AND key_path = $2 AND lang = $3
	`, user1ID, "topbar.welcome", "en").Scan(&value)
	assert.NoError(t, err)
	assert.Equal(t, "Welcome", value)

	err = conn.QueryRow(context.Background(), `
		SELECT value FROM ui_translations WHERE user_id = $1 AND key_path = $2 AND lang = $3
	`, user2ID, "topbar.welcome", "es").Scan(&value)
	assert.NoError(t, err)
	assert.Equal(t, "Bienvenido", value)
}
