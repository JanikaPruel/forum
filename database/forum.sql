BEGIN TRANSACTION;

CREATE TABLE if not EXISTS users (
    id integer PRIMARY KEY autoincrement,
    username TEXT NOT NULL,
    email TEXT NOT NULL,
    password text not null,
    created_at datetime default 'now'
);

-- ALTER TABLE users DROP COLUMN email;

CREATE TABLE if not EXISTS sessions (
    id INTEGER PRIMARY KEY autoincrement,
    user_id integer references user(id),
    expires datetime not null default 'now'
);

CREATE TABLE if not EXISTS categories ( -- id, title
    id integer primary key autoincrement not null,
    name varchar(255) not null,
    created_at datetime not null default 'now'
);

CREATE TABLE IF NOT EXISTS total_likes_post (
    id	INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id	INTEGER REFERENCES users(id),
	post_id	INTEGER REFERENCES posts(id)
);
CREATE TABLE IF NOT EXISTS total_dislikes_post (
    id	INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id	INTEGER REFERENCES users(id),
	post_id	INTEGER REFERENCES posts(id)
);

CREATE TABLE IF NOT EXISTS total_likes_comment (
	id	INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id	INTEGER REFERENCES users(id),
	comment_id	INTEGER REFERENCES commnents(id)
);

CREATE TABLE IF NOT EXISTS total_dislikes_comment (
	id	INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id	INTEGER REFERENCES users(id),
	comment_id	INTEGER REFERENCES commnents(id)
);

CREATE TABLE if not EXISTS posts ( -- id, title, content, image, likes, dislikes, views, user_id, created_at, updated_at, removed_at
    id integer primary key autoincrement, -- 1, 2, 3
    title text not null,
    content text not null,
    image text,
    likes integer default 0 references total_likes_post(id),
    dislikes integer default 0 references total_dislikes_post(id), 
    user_id integer references users(id),
    comment integer references comments(id),
    created_at datetime not null default 'now'
);

CREATE TABLE IF NOT EXISTS posts_categories (
    id	INTEGER PRIMARY KEY AUTOINCREMENT,
	category_id	INTEGER REFERENCES categories(id),
	post_id	INTEGER REFERENCES posts(id)    

);

CREATE TABLE if not EXISTS comments ( -- id, user_id, post_id, content, created_at, updated_at, likes, dislikes
    id integer primary key autoincrement not null,
    user_id integer references users(id),
    post_id integer references posts(id),
    content text not null,
    likes integer default 0 REFERENCES total_likes_commnts(id),
    dislikes integer default 0 REFERENCES total_dislikes_commnts(id),
    created_at datetime not null default 'now'
);

COMMIT;
DROP TABLE users;
DROP TABLE sessions;
DROP TABLE categories;
DROP TABLE posts;
DROP TABLE comments;