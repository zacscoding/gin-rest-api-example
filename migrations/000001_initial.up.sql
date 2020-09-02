CREATE TABLE accounts (
	id         int unsigned auto_increment PRIMARY KEY,
	username VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL,
	password VARCHAR(255) NOT NULL,
	UNIQUE KEY unique_users_email (email)
) CHARACTER SET utf8mb4;