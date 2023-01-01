CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
  id   UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
  email text NOT NULL UNIQUE
);