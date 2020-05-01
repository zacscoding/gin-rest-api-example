drop table if exists comments;
drop table if exists articles_tags;
drop table if exists tags;
drop table if exists articles_favorite;
drop table if exists articles;
drop table if exists follows;
drop table if exists users;

-- user
CREATE TABLE users (
    id         int unsigned auto_increment PRIMARY KEY,
    created_at datetime      NULL,
    updated_at datetime      NULL,
    deleted_at datetime      NULL,
    email VARCHAR(255) NULL,
    username VARCHAR(255) NOT NULL,
    bio TEXT,
    password VARCHAR(255),
    image      VARCHAR(255)  NULL,
    UNIQUE KEY unique_users_email (email)
) CHARACTER SET utf8mb4;
CREATE index idx_users_deleted_at ON users(deleted_at);

-- follows
CREATE TABLE follows (
    created_at   DATETIME     NULL,
    updated_at   DATETIME     NULL,
    follower_id  INT UNSIGNED,
    following_id INT UNSIGNED,
    created_at datetime      NULL,
    CONSTRAINT follows_follower_id_fk
        FOREIGN KEY (follower_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT follows_following_id_fk
        FOREIGN KEY (follower_id) REFERENCES users (id) ON DELETE CASCADE,
    PRIMARY KEY (follower_id, following_id)
) CHARACTER SET utf8mb4;

-- article
CREATE TABLE articles (
    id           INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at   DATETIME     NULL,
    updated_at   DATETIME     NULL,
    deleted_at   DATETIME     NULL,
    slug VARCHAR(255),
    title VARCHAR(255),
    description TEXT,
    body TEXT,
    user_id INT UNSIGNED,
    UNIQUE KEY unique_slug (slug),
    CONSTRAINT articles_user_id_fk
        FOREIGN KEY (user_id) REFERENCES users (id)
) CHARACTER SET utf8mb4;
CREATE INDEX idx_articles_deleted_at ON articles(deleted_at);

CREATE TABLE articles_favorites (
    user_id INT UNSIGNED,
    article_id INT UNSIGNED,
    CONSTRAINT articles_favorite_user_id_fk
        FOREIGN KEY (user_id) REFERENCES users (id),
    CONSTRAINT articles_favorite_article_id_fk
        FOREIGN KEY (article_id) REFERENCES articles (article_id),
    PRIMARY KEY (user_id, article_id)
);

-- tag
CREATE TABLE tags (
    id           INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at   DATETIME     NULL,
    updated_at   DATETIME     NULL,
    deleted_at   DATETIME     NULL,
    name varchar(255),
    UNIQUE KEY unique_tags_name (name)
)CHARACTER SET utf8mb4;

CREATE TABLE articles_tags (
    article_id INT UNSIGNED,
    tag_id INT UNSIGNED,
    CONSTRAINT articles_tags_article_id_fk
        FOREIGN KEY (article_id) REFERENCES articles (id),
    CONSTRAINT articles_tags_tag_id_fk
        FOREIGN KEY (tag_id) REFERENCES tags (id),
    PRIMARY KEY(article_id, tag_id)
);

-- comments
CREATE TABLE comments (
    id           INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at   DATETIME     NULL,
    updated_at   DATETIME     NULL,
    body text,
    article_id INT UNSIGNED,
    user_id INT UNSIGNED,
    CONSTRAINT comments_article_id_fk
        FOREIGN KEY (article_id) REFERENCES articles (id),
    CONSTRAINT comments_user_id_fk
        FOREIGN KEY (user_id) REFERENCES users (id)
)CHARACTER SET utf8mb4;


