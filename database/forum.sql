BEGIN TRANSACTION;

CREATE TABLE if not EXISTS users (
    id integer PRIMARY KEY AUTOINCREMENT NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password varchar(255) not null
);

CREATE TABLE if not EXISTS sessions (
    id varchar(128) not null unique,
    user_id integer references user(id) not null,
    expires_at timestamp not null default CURRENT_TIMESTAMP
);

CREATE TABLE if not EXISTS categories ( -- id, title
    id integer primary key autoincrement not null unique,
    title varchar(255) not null,
    created_at timestamp not null default CURRENT_TIMESTAMP
);

CREATE TABLE if not EXISTS posts ( -- id, title, content, image, likes, dislikes, views, user_id, created_at, updated_at, removed_at
    id integer primary key autoincrement not null unique, -- 1, 2, 3
    title text not null,
    content text not null,
    image text,
    likes integer not null default 0,
    dislikes integer not null default 0, 
    views integer not null default 0,
    user_id integer references user(id) not null,
    created_at timestamp not null default CURRENT_TIMESTAMP,
    updated_at timestamp not null default CURRENT_TIMESTAMP
);

CREATE TABLE if not EXISTS comments ( -- id, user_id, post_id, content, created_at, updated_at, likes, dislikes
    id integer primary key autoincrement not null unique,
    user_id integer references user(id) not null,
    post_id integer references post(id) not null,
    content text not null,
    created_at timestamp not null default CURRENT_TIMESTAMP,
    updated_at timestamp not null default CURRENT_TIMESTAMP,
    likes integer not null default 0,
    dislikes integer not null default 0
);


COMMIT;