CREATE TABLE videos (
    id   UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    path varchar NOT NULL,
    uploaded_at TIMESTAMP WITHOUT TIME ZONE,
    owner_id UUID NOT NULL,
    CONSTRAINT fk_owner_id FOREIGN KEY (owner_id) REFERENCES users (id)
);