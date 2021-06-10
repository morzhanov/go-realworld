CREATE TABLE photos (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4()
    title VARCHAR(255) NOT NULL
    base64 TEXT NOT NULL
);
