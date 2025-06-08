-- 用户表：存储用户的基本信息
CREATE TABLE `user` (
  `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY, -- 用户ID，自增主键
  `name` VARCHAR(64) NOT NULL, -- 用户名
  `email` VARCHAR(128) NOT NULL UNIQUE, -- 用户邮箱，唯一
  `password` VARCHAR(255) NOT NULL, -- 用户密码（加密存储）
  `phone` VARCHAR(20) DEFAULT NULL, -- 用户手机号
  `avatar` VARCHAR(255) DEFAULT NULL COMMENT '用户头像', -- 用户头像URL
  `bio` VARCHAR(255) DEFAULT NULL COMMENT '用户简介', -- 用户简介
  `location` VARCHAR(128) DEFAULT NULL COMMENT '用户位置', -- 用户所在位置
  `website` VARCHAR(255) DEFAULT NULL COMMENT '用户个人网站', -- 用户个人网站
  `role` VARCHAR(32) NOT NULL DEFAULT 'user' COMMENT '用户权限角色', -- 用户权限属性
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 创建时间
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- 更新时间
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;