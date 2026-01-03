-- 用户表：存储用户的基本信息 (SQLite 兼容版本)
CREATE TABLE IF NOT EXISTS `users` (
  `id` INTEGER PRIMARY KEY AUTOINCREMENT, -- 用户ID，自增主键 (SQLite 语法)
  `name` TEXT NOT NULL UNIQUE, -- 用户名
  `email` TEXT NOT NULL UNIQUE, -- 用户邮箱，唯一
  `password` TEXT NOT NULL, -- 用户密码（加密存储）
  `phone` TEXT DEFAULT NULL, -- 用户手机号
  `avatar` TEXT DEFAULT NULL, -- 用户头像URL
  `bio` TEXT DEFAULT NULL, -- 用户简介
  `location` TEXT DEFAULT NULL, -- 用户所在位置
  `website` TEXT DEFAULT NULL, -- 用户个人网站
  `role` TEXT NOT NULL DEFAULT 'user', -- 用户权限属性
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 创建时间
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP -- 更新时间
);

-- 创建触发器 (Trigger) 来模拟 ON UPDATE CURRENT_TIMESTAMP
CREATE TRIGGER IF NOT EXISTS `trigger_user_updated_at`
AFTER UPDATE ON `user`
FOR EACH ROW
BEGIN
  UPDATE `user` SET `updated_at` = CURRENT_TIMESTAMP WHERE `id` = OLD.`id`;
END;
