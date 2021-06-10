CREATE TABLE pictures (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4()
    title VARCHAR(255) NOT NULL
    base64 TEXT NOT NULL
    user_id VARCHAR(255) NOT NULL
);
