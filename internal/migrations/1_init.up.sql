CREATE TABLE users (
    user_id VARCHAR PRIMARY KEY,
    name VARCHAR(100),
    created_at TIMESTAMP not null default now()
);