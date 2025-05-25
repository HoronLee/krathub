-- 创建 user 表
CREATE TABLE user (
    id INT AUTO_INCREMENT PRIMARY KEY, -- 用户ID，主键，自增
    username VARCHAR(50) NOT NULL,     -- 用户名，非空
    password VARCHAR(255) NOT NULL,    -- 密码，非空
    email VARCHAR(100),                -- 邮箱，可为空
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间，默认当前时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- 更新时间，自动更新
);
