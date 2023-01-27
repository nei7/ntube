ALTER TABLE users ADD COLUMN username varchar(30) NOT NULL, 
ADD COLUMN created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (now()),
ADD COLUMN avatar varchar,
ADD COLUMN description varchar(200),
ADD COLUMN followers int NOT NULL DEFAULT 0;
