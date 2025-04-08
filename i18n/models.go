package i18n

// Translation represents a single translation entry.
type Translation struct {
	UserID  *string // nil = global/default translation
	Lang    string
	KeyPath string
	Value   string
}
