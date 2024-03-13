package model

import "time"

/*

    ERD - Entity Relational Diagram | PostgreSQ

   User
        id bigserial primary key not null unique,
        username varchar(255) not null,
        email varchar(255) not null unique,
        password varchar(64) not null
        created_at timestamp nou null default current_timestamp


    Category
    id pk
    title



    Post
        id bigserial primary key not null unique,
        title varchar(255) not null,
        content text not null,
        image text - path or url || S3 - Simple Storage Service
        like integer check(like >= 0),
        dislike integer check(dislike >= 0)

        category_id fk
        user_id fk

        created_at timestamp default current_timestamp,
        updated_at timestamp default current_timestamp

    Comment
    id pk
    content
    post_id fk
    user_id fk

*/
// Sql and SqlLite у каждого свой уникальный диалект 90% команд совпадают
// Но различия есть разумеется

type User struct {
	ID       int    `sql:"id"`
	Username string `sql:"username"`
	Email    string `sql:"email"`
	Password string `sql:"password"`
}

type Session struct {
	ID      int           `sql:"id"`
	UserID  int           `sql:"username"`
	Expired time.Duration `sql:"expired"`
}

type Category struct {
	ID        int           `sql:"id"`
	Title     string        `sql:"title"`
	CreatedAt time.Duration `sql:"created_at"`
}

type Post struct {
	ID        int           `sql:"id"`
	Title     string        `sql:"title"`
	Content   string        `sql:"content"`
	Image     string        `sql:"image"`
	Likes     int           `sql:"likes"`
	Dislikes  int           `sql:"dislikes"`
	Views     int           `sql:"views"`
	UserID    int           `sql:"user_id"`
	CommentID int           `sql:"post_id"`
	CreatedAt time.Duration `sql:"created_at"`
	UpdatedAt time.Duration `sql:"updated_at"`
	RemovedAt time.Duration `sql:"removed_at"`
}

type Comment struct {
	ID        int           `sql:"id"`
	UserID    int           `sql:"user_id"`
	PostID    int           `sql:"post_id"`
	Content   string        `sql:"content"`
	Likes     int           `sql:"likes"`
	Dislikes  int           `sql:"dislikes"`
	CreatedAt time.Duration `sql:"created_at"`
	UpdatedAt time.Duration `sql:"updated_at"`
}

// this
