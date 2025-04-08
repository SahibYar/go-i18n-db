# 🌍 go-i18n-db

**go-i18n-db** is a Golang-based localization toolkit that helps you manage and persist your app’s translation strings in a relational database (like PostgreSQL) instead of relying solely on static `.json` files.

This library is designed for enterprise-grade applications where:
*	Localization needs to be dynamic or user/tenant-specific.
*	Admins require CRUD access to translation keys.
*	You want to sync .json i18n files with a database backend.
 
## ✨ Key Features
*	✅ Flatten deeply nested translation JSON into DB-friendly key-value format.
*	✅ Store translations in PostgreSQL for persistence and queryability.
*	✅ Full support for CRUD operations (Create, Read, Update, Delete).
*	✅ Support for user- or tenant-level overrides (multi-tenant setups).
*	✅ Bulk import/export from/to `.json` files.
*	✅ Optional: CLI tooling, API wrapper, admin UI integration.

## 📦 Installation
```bash
go get github.com/SahibYar/go-i18n-db
```
> Make sure your Go version is 1.18 or higher.

## 🗃️ Database Schema
Here’s a sample schema to use with PostgreSQL:
```SQL
CREATE TABLE ui_translations (
    id SERIAL PRIMARY KEY,
    user_id UUID,                    -- Nullable: system-wide if NULL
    key_path TEXT NOT NULL,          -- Flattened key e.g., 'topbar.profile'
    lang TEXT NOT NULL,              -- Language code: 'en', 'es', 'ar', etc.
    value TEXT NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (user_id, key_path, lang)
);
```

## ✅ Features (API Reference)
### 🧩 1. Flatten JSON Structure
```go
func FlattenJSON(input map[string]interface{}) map[string]string
```
Converts a deeply nested JSON object into a flat map like:
```json
{
  "topbar|profile": "My Profile",
  "forms|submit": "Submit"
}
```
### 🧩 2. Load From File
```go
func LoadAndFlatten(filePath string) (map[string]string, error)
```
Parses a JSON file (e.g., `en.json`) and flattens it into a key-value map.

### 📥 3. Insert Into PostgreSQL
```go
type Translation struct {
    UserID   *string // nullable UUID for per-user overrides
    Lang     string  // 'en', 'es', etc.
    KeyPath  string  // e.g., 'forms|submit'
    Value    string  // actual translation string
}

func UpsertTranslations(db *sql.DB, translations []Translation) error
```
Efficient bulk insert with conflict handling
```SQL
ON CONFLICT (user_id, key_path, lang) DO UPDATE SET ...
```
### 🔍 4. Fetch With Fallback
```go
func GetTranslation(db *sql.DB, userID *string, keyPath, lang string) (string, error)
```
Looks up a translation by `keyPath` and `lang`. If a `userID` is provided, it will first try to find a user-specific override and fallback to global.

### 🔁 5. Export to JSON
```go
func ExportToJSON(db *sql.DB, lang string, userID *string) (map[string]string, error)
```
Converts stored translations back into a flattened map, which can then be re-structured as a `.json` file using your own logic.

## 🧪 Example Workflow
```go
// Load and flatten a file
translations, _ := LoadAndFlatten("en.json")

// Convert to []Translation
var rows []Translation
for k, v := range translations {
    rows = append(rows, Translation{
        UserID:  nil,
        Lang:    "en",
        KeyPath: k,
        Value:   v,
    })
}

// Save to PostgreSQL
_ = UpsertTranslations(db, rows)

// Fetch a string
val, _ := GetTranslation(db, nil, "forms|submit", "en")
fmt.Println(val) // "Submit"
```
## 💡 Optional Add-ons (In Progress or Community Contributions Welcome)
* 🧑‍💻 CLI Tool:
```bash
localizr import ./es.json --lang=es --user=uuid
localizr export --lang=ar > ar.json
```

* 🌐 REST API:
CRUD endpoints to manage translations at runtime.

* 🖥️ Admin UI:
Table view + in-place editing of all translations.

## 🙌 Contributing

Pull requests and GitHub issues are welcome! You can contribute:
* New database support (e.g., MySQL, SQLite)
* More CLI features
* UI admin tool
* Better test coverage

## 📄 License

MIT License — feel free to use it in commercial and non-commercial projects.

## 🔗 Related Projects & Inspiration
* [go-i18n](https://github.com/nicksnyder/go-i18n) (for file-based i18n)
* Enterprise admin panels (where runtime translation is essential)
* Multi-tenant SaaS apps
