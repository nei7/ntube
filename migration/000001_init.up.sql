CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE users (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    username varchar(30) NOT NULL,
    email varchar NOT NULL UNIQUE,
    password varchar NOT NUll,
    description varchar(200),
    avatar varchar(200)
);