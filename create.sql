CREATE TABLE ui_translations
(
    id          SERIAL PRIMARY KEY,
    user_id     UUID,
    key_path    TEXT NOT NULL,
    lang        TEXT NOT NULL,
    value       TEXT NOT NULL,
    tooltip     TEXT NULL,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_by  UUID,
    unique (user_id, key_path, lang)
);

ALTER TABLE ui_translations
    OWNER TO postgres;
