-- DROP DATABASE IF EXISTS imsdb;
-- CREATE DATABASE imsdb;

\c imsdb;

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    username VARCHAR(50) NOT NULL UNIQUE,
    password VARCHAR(60) NOT NULL,
    email VARCHAR(100) NOT NULL,
    active INT NOT NULL DEFAULT 1,
    admin INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(), 
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

INSERT INTO users (first_name, last_name, username, password, email)
VALUES
    ('Jie Wei', 'Low', 'lowjiewei', '$2a$10$OULOXURo57bo5keyNXGQxefqMyEM67JIscqLVKWgd/S.siCqNAHdC', 'lowjiewei@email.com');
