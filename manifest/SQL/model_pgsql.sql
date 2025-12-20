-- 用户表：存储用户的基本信息 (PostgreSQL 兼容版本)
CREATE TABLE IF NOT EXISTS "user" (
    "id" BIGSERIAL PRIMARY KEY, -- 用户ID，PostgreSQL 自增主键
    "name" VARCHAR(64) NOT NULL UNIQUE, -- 用户名
    "email" VARCHAR(128) NOT NULL UNIQUE, -- 用户邮箱，唯一
    "password" VARCHAR(255) NOT NULL, -- 用户密码（加密存储）
    "phone" VARCHAR(20) DEFAULT NULL, -- 用户手机号
    "avatar" VARCHAR(255) DEFAULT NULL, -- 用户头像URL
    "bio" VARCHAR(255) DEFAULT NULL, -- 用户简介
    "location" VARCHAR(128) DEFAULT NULL, -- 用户所在位置
    "website" VARCHAR(255) DEFAULT NULL, -- 用户个人网站
    "role" VARCHAR(32) NOT NULL DEFAULT 'user', -- 用户权限角色
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 创建时间（带时区）
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP -- 更新时间（带时区）
);

-- 创建触发器函数来模拟 MySQL 的 ON UPDATE CURRENT_TIMESTAMP
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- 创建触发器，在更新用户记录时自动更新 updated_at 字段
CREATE TRIGGER trigger_user_updated_at
    BEFORE UPDATE ON "user"
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- 添加索引以提升查询性能
CREATE INDEX IF NOT EXISTS idx_user_email ON "user"("email");
CREATE INDEX IF NOT EXISTS idx_user_name ON "user"("name");
CREATE INDEX IF NOT EXISTS idx_user_role ON "user"("role");
