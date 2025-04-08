CREATE TABLE ui_translations (
     id SERIAL PRIMARY KEY,
     user_id UUID,                    -- Nullable: system-wide if NULL
     key_path TEXT NOT NULL,          -- Flattened key e.g., 'topbar.profile'
     lang TEXT NOT NULL,              -- Language code: 'en', 'es', 'ar', etc.
     value TEXT NOT NULL,
     updated_at TIMESTAMPTZ DEFAULT NOW(),
     UNIQUE (user_id, key_path, lang)
);