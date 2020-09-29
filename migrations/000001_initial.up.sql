-- account
CREATE TABLE accounts (
	id         int unsigned auto_increment PRIMARY KEY,
	username VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL,
	password VARCHAR(255) NOT NULL,
	bio TEXT NULL,
	image VARCHAR(255) NULL,
	created_at datetime NULL,
    updated_at datetime NULL,
    disabled tinyint(1) DEFAULT '0',
	UNIQUE KEY unique_users_email (email)
) CHARACTER SET utf8mb4;

-- article
CREATE TABLE articles (
    id         INT unsigned auto_increment PRIMARY KEY,
    slug VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    body TEXT NOT NULL,
    created_at   DATETIME NULL,
    updated_at   DATETIME NULL,
    deleted_at_unix INT DEFAULT 0,
    author_id INT UNSIGNED,
    UNIQUE KEY unique_articles_slug (slug, deleted_at_unix),
    CONSTRAINT articles_author_id_fk FOREIGN KEY (author_id) REFERENCES accounts(id)
) CHARACTER SET utf8mb4;
CREATE INDEX idx_articles_deleted_at_unix ON articles(deleted_at_unix);

-- tag
CREATE TABLE tags (
    id           INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at   DATETIME     NULL,
    name varchar(255),
    UNIQUE KEY unique_tags_name (name)
)CHARACTER SET utf8mb4;

-- article tag relation
CREATE TABLE article_tags (
    article_id INT UNSIGNED,
    tag_id INT UNSIGNED,
    CONSTRAINT article_tags_article_id_fk FOREIGN KEY (article_id) REFERENCES articles (id),
    CONSTRAINT article_tags_tag_id_fk FOREIGN KEY (tag_id) REFERENCES tags (id),
    PRIMARY KEY(article_id, tag_id)
);

-- comments
CREATE TABLE comments (
    id           INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    body TEXT,
    slug VARCHAR(255) NOT NULL,
    created_at   DATETIME     NULL,
    updated_at   DATETIME     NULL,
    deleted_at   DATETIME     NULL,
    author_id INT UNSIGNED,
    CONSTRAINT comments_author_id_fk FOREIGN KEY (author_id) REFERENCES accounts (id)
)CHARACTER SET utf8mb4;
CREATE index idx_comments_deleted_at ON comments(deleted_at);
CREATE INDEX idx_comments_slug ON comments(slug);