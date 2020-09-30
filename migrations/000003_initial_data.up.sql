-- admin password is "admin1"
-- user1 password is "user1"

INSERT INTO accounts (username, email, password, created_at, updated_at) VALUES
('admin', 'admin@email.com', '$2a$10$qUq8XJ.RcxcQvaEipicm/OLVPp5AjJjoPigj4vlRU579Xz0SkZwqu', now(), now()),
('user1', 'user1@email.com', '$2a$10$lsYsLv8nGPM0.R.ft4sgpe3OP7..KL3ZJqqhSVCKTEnSCMUztoUcW', now(), now());