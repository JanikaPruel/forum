package model

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
