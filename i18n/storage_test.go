package i18n

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// TestExportToJSON tests the ExportToJSON function
func TestExportToJSON(t *testing.T) {
	// Setup PostgreSQL connection (use an environment variable for connection string)
	connString := os.Getenv("DB_CONNECTION_STRING")
	if connString == "" {
		t.Fatal("DB_CONNECTION_STRING environment variable is not set")
	}

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
