-- 更新admin用户密码
-- 密码："123456"的bcrypt哈希值（由Go bcrypt库生成）
UPDATE users SET password = '$2a$10$slYQmyNdGzin7olVN3p5Be.7DlH.87ZZ9z5XP/dIaDm8JCZ4bYy6u' WHERE id = 1;
SELECT id, username, phone, password FROM users WHERE id = 1;
