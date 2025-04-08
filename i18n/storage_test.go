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

// TestExportToFlatJSON tests the ExportToFlatJSON function
func TestExportToFlatJSON(t *testing.T) {
	// Connect to the test database
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	// Generate valid UUIDs for user1 and user2
	user1ID := uuid.New()
	user2ID := uuid.New()

	// Insert test data into the ui_translations table, including tooltips
	_, err = conn.Exec(context.Background(), `
		INSERT INTO ui_translations (user_id, key_path, lang, value, tooltip, updated_at)
		VALUES 
			($1, 'topbar.profile', 'en', 'Profile', 'Your profile', NOW()), 
			(NULL, 'topbar.profile', 'es', 'Perfil', 'Perfil en español', NOW()), 
			($1, 'footer.contact', 'en', 'Contact', 'Contact us', NOW())
	`, user1ID)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Test: Fetch translations with user ID 'user1'
	translations, err := ExportToFlatJSON(context.Background(), conn, "en", stringPtr(user1ID.String()))
	assert.NoError(t, err)
	assert.NotNil(t, translations)
	assert.Equal(t, "Profile", translations["topbar.profile"]["value"])
	assert.Equal(t, "Contact", translations["footer.contact"]["value"])
	assert.Equal(t, "Your profile", translations["topbar.profile"]["tooltip"])
	assert.Equal(t, "Contact us", translations["footer.contact"]["tooltip"])

	// Test: Fetch translations without user ID (fallback to NULL user)
	translations, err = ExportToFlatJSON(context.Background(), conn, "es", nil)
	assert.NoError(t, err)
	assert.NotNil(t, translations)
	assert.Equal(t, "Perfil", translations["topbar.profile"]["value"])
	assert.Equal(t, "Perfil en español", translations["topbar.profile"]["tooltip"])

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

	// Test 1: Insert new translations with tooltips for user1 and user2
	translations := []Translation{
		{UserID: stringPtr(user1ID.String()), KeyPath: "topbar.profile", Lang: "en", Value: "Profile", ToolTip: "Your profile"},
		{UserID: stringPtr(user1ID.String()), KeyPath: "footer.contact", Lang: "en", Value: "Contact", ToolTip: "Contact us"},
		{UserID: stringPtr(user2ID.String()), KeyPath: "topbar.profile", Lang: "es", Value: "Perfil", ToolTip: "Tu perfil"},
		{UserID: stringPtr(user2ID.String()), KeyPath: "footer.contact", Lang: "es", Value: "Contacto", ToolTip: "Contáctanos"},
	}

	err = UpsertTranslations(context.Background(), conn, translations)
	assert.NoError(t, err)

	// Verify if the translations and tooltips were inserted
	var value, tooltip string
	err = conn.QueryRow(context.Background(), `
		SELECT value, tooltip FROM ui_translations WHERE user_id = $1 AND key_path = $2 AND lang = $3
	`, user1ID, "topbar.profile", "en").Scan(&value, &tooltip)
	assert.NoError(t, err)
	assert.Equal(t, "Profile", value)
	assert.Equal(t, "Your profile", tooltip)

	// Test 2: Update existing translations with new tooltips for user1
	updatedTranslations := []Translation{
		{UserID: stringPtr(user1ID.String()), KeyPath: "topbar.profile", Lang: "en", Value: "Updated Profile", ToolTip: "Updated profile tooltip"},
	}

	err = UpsertTranslations(context.Background(), conn, updatedTranslations)
	assert.NoError(t, err)

	// Verify if the translation and tooltip were updated
	err = conn.QueryRow(context.Background(), `
		SELECT value, tooltip FROM ui_translations WHERE user_id = $1 AND key_path = $2 AND lang = $3
	`, user1ID, "topbar.profile", "en").Scan(&value, &tooltip)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Profile", value)
	assert.Equal(t, "Updated profile tooltip", tooltip)

	// Test 3: Bulk insertion of translations with tooltips for user1 and user2
	bulkTranslations := []Translation{
		{UserID: stringPtr(user1ID.String()), KeyPath: "topbar.welcome", Lang: "en", Value: "Welcome", ToolTip: "Welcome to our site"},
		{UserID: stringPtr(user2ID.String()), KeyPath: "topbar.welcome", Lang: "es", Value: "Bienvenido", ToolTip: "Bienvenido a nuestro sitio"},
		{UserID: stringPtr(user1ID.String()), KeyPath: "footer.privacy", Lang: "en", Value: "Privacy", ToolTip: "Privacy settings"},
		{UserID: stringPtr(user2ID.String()), KeyPath: "footer.privacy", Lang: "es", Value: "Privacidad", ToolTip: "Configuración de privacidad"},
	}

	err = UpsertTranslations(context.Background(), conn, bulkTranslations)
	assert.NoError(t, err)

	// Verify bulk insertion with tooltips
	err = conn.QueryRow(context.Background(), `
		SELECT value, tooltip FROM ui_translations WHERE user_id = $1 AND key_path = $2 AND lang = $3
	`, user1ID, "topbar.welcome", "en").Scan(&value, &tooltip)
	assert.NoError(t, err)
	assert.Equal(t, "Welcome", value)
	assert.Equal(t, "Welcome to our site", tooltip)

	err = conn.QueryRow(context.Background(), `
		SELECT value, tooltip FROM ui_translations WHERE user_id = $1 AND key_path = $2 AND lang = $3
	`, user2ID, "topbar.welcome", "es").Scan(&value, &tooltip)
	assert.NoError(t, err)
	assert.Equal(t, "Bienvenido", value)
	assert.Equal(t, "Bienvenido a nuestro sitio", tooltip)
}
