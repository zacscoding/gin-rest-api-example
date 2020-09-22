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