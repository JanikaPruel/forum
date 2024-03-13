START TRANSACTION;

CREATE DATABASE forum;

-- users, auth - statefull = sessions, categories, posts, comments 
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTOINCREMENT NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password varchar(255) not null,
);

CREATE TABLE sessions (
    id varchar(128) not null unique,
    user_id bigint foreign key(user_id_fk) references user(id) not null,
    expired datetime not null default now()
);

CREATE TABLE categories ( -- id, title
    id bigint primary key autoincrement not null unique,
    title varchar(255) not null,
    created_at datetime not null default now()
);

CREATE TABLE posts ( -- id, title, content, image, likes, dislikes, views, user_id, created_at, updated_at, removed_at
    id bigint primary key autoincrement not null unique, -- 1, 2, 3
    title text not null,
    content text not null,
    image text,
    likes integer not null default 0,
    dislikes integer not null default 0, 
    views integer not null default 0,
    user_id bigint foreign key(user_id_fk) references user(id) not null,
    comment_id bigint foreign key(comment_id_fk) references comment(id) not null
    created_at datetime not null default now(),
    updated_at datetime not null default now(),
    removed_at datetime not null default now()
);

CREATE TABLE comments ( -- id, user_id, post_id, content, created_at, updated_at, likes, dislikes
    id bigint primary key autoincrement not null unique,
    user_id bigint foreign key(user_id_fk) references user(id) not null,
    post_id bigint foreing key(post_id_fk) references post(id) not null,
    content text not null,
    created_at datetime not null default now(),
    updated_at datetime not null default now(),
    likes integer not null default 0,
    dislikes integer not null default 0,
);


COMMIT;