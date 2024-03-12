START TRANSACTION;

CREATE DATABASE forum;

-- users, auth - statefull = sessions, categories, posts, comments 
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTOINCREMENT NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) no
);

CREATE TABLE sessions (

);

CREATE TABLE categories (

);

CREATE TABLE posts (
    liked 
    disliked 
);

CREATE TABLE comments (

);

CREATE TABLE users (
    
);

COMMIT;