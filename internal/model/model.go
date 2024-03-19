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
	ID        int    `sql:"id"`
	Username  string `sql:"username"`
	Email     string `sql:"email"`
	Password  string `sql:"password"`
	CreatedAt time.Time `sql:"created_at"`
}

type Session struct {
	ID      int           `sql:"id"`
	UserID  int           `sql:"username"`
	Expires time.Time `sql:"expired"`
}

type Category struct {
	ID        int           `sql:"id"`
	Name     string        `sql:"name"`
	CreatedAt time.Time `sql:"created_at"`
}

type TotalLikesPost struct {
	ID     int `sql:"id"`
	UserID int `sql:"user_id"`
	PostID int `sql:"post_id"`
}
type TotalDislikesPost struct {
	ID     int `sql:"id"`
	UserID int `sql:"user_id"`
	PostID int `sql:"post_id"`
}
type TotalLikesComment struct {
	ID        int `sql:"id"`
	UserID    int `sql:"user_id"`
	CommentID int `sql:"comment_id"`
}
type TotalDislikesComment struct {
	ID        int `sql:"id"`
	UserID    int `sql:"user_id"`
	CommentID int `sql:"comment_id"`
}

type Post struct {
	ID                   int       `sql:"id"`
	UserID               int       `sql:"user_id"`
	Username             string    `sql:"username"`
	Title                string    `sql:"title"`
	Content              string    `sql:"content"`
	Category             []int       `sql:"category"`
	Comment              int       `sql:"comment"`
	Likes                int       `sql:"likes"`
	Dislikes             int       `sql:"dislikes"`
	CreatedAt            time.Time `sql:"created_at"`
	IsLikedByAuthUser    bool
	IsDislikedByAuthUser bool
}

type Comment struct {
	ID                   int       `sql:"id"`
	UserID               int       `sql:"user_id"`
	PostID               int       `sql:"post_id"`
	Username             string    `sql:"username"`
	Content              string    `sql:"content"`
	Likes                int       `sql:"likes"`
	Dislikes             int       `sql:"dislikes"`
	CreatedAt            time.Time `sql:"created_at"`
	IsLikedByAuthUser    bool
	IsDislikedByAuthUser bool
}

// this
