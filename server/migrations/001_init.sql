-- 衣搭库数据库初始化脚本
-- 创建数据库
CREATE DATABASE IF NOT EXISTS yidaiku DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE yidaiku;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id INT(10) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL COMMENT '用户ID',
    phone VARCHAR(20) UNIQUE,
    email VARCHAR(100),
    password_hash VARCHAR(255),
    nickname VARCHAR(50),
    avatar VARCHAR(512),
    status TINYINT DEFAULT 1 COMMENT '1:正常 0:禁用',
    is_member TINYINT DEFAULT 0 COMMENT '是否会员',
    member_expire INT(10) NOT NULL DEFAULT 0 COMMENT '会员过期时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_user_id (user_id),
    INDEX idx_phone (phone),
    INDEX idx_email (email),
    INDEX idx_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 用户偏好设置表
CREATE TABLE IF NOT EXISTS user_preferences (
    id INT(10) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL COMMENT '用户ID',
    styles JSON COMMENT '风格偏好数组',
    occasions JSON COMMENT '场合偏好数组',
    colors JSON COMMENT '颜色偏好数组',
    height INT COMMENT '身高(cm)',
    weight INT COMMENT '体重(kg)',
    principles JSON COMMENT '穿搭原则',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    UNIQUE KEY uk_user_id (user_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户偏好设置表';

-- 品类表（品类与子品类，pid 为父级 id，NULL 表示顶级品类）
CREATE TABLE IF NOT EXISTS categories (
    id INT(10) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    pid INT(10) NULL DEFAULT NULL COMMENT '父品类ID，NULL表示顶级品类',
    user_id VARCHAR(20) NOT NULL COMMENT '用户ID',
    category_name VARCHAR(64) NOT NULL COMMENT '品类名称',
    status TINYINT DEFAULT 1 COMMENT '1:正常 0:禁用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    deleted_at DATETIME NULL DEFAULT NULL COMMENT '删除时间（软删除）',
    INDEX idx_user_id (user_id),
    INDEX idx_pid (pid),
    INDEX idx_status (status),
    INDEX idx_deleted_at (deleted_at),
    INDEX idx_user_pid (user_id, pid),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (pid) REFERENCES categories(id) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='品类表';

-- 衣物表
CREATE TABLE IF NOT EXISTS clothing_items (
    id INT(10) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL COMMENT '用户ID',
    image_url VARCHAR(512) NOT NULL COMMENT '七牛云OSS URL',
    thumbnail_url VARCHAR(512) COMMENT '缩略图URL',
    category VARCHAR(32) COMMENT '品类:上衣/裤子/裙子/鞋子/配饰',
    subcategory VARCHAR(32) COMMENT '子品类',
    colors JSON COMMENT '颜色数组',
    styles JSON COMMENT '风格数组',
    season VARCHAR(16) COMMENT '季节:春夏/秋冬/四季通用',
    tags JSON COMMENT '自定义标签',
    notes TEXT COMMENT '备注',
    is_deleted TINYINT DEFAULT 0 COMMENT '软删除',
    is_retired TINYINT DEFAULT 0 COMMENT '已淘汰:1=不参与推荐但保留记录',
    deleted_at DATETIME COMMENT '删除时间',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_user_id (user_id),
    INDEX idx_category (category),
    INDEX idx_season (season),
    INDEX idx_is_deleted (is_deleted),
    INDEX idx_is_retired (is_retired),
    UNIQUE KEY uk_item_id (item_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='衣物表';

-- 穿搭推荐表
CREATE TABLE IF NOT EXISTS outfit_recommendations (
    id INT(10) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL COMMENT '用户ID',
    item_ids JSON COMMENT '单品ID数组',
    occasion VARCHAR(32) COMMENT '场合',
    weather JSON COMMENT '天气信息',
    score DECIMAL(3,2) COMMENT '推荐评分',
    highlights TEXT COMMENT '搭配亮点',
    suitable JSON COMMENT '适用场景',
    source VARCHAR(32) COMMENT '推荐来源:ai/rule/template',
    is_favorite TINYINT DEFAULT 0 COMMENT '是否收藏',
    is_applied TINYINT DEFAULT 0 COMMENT '是否应用',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_user_id (user_id),
    INDEX idx_occasion (occasion),
    INDEX idx_created_at (created_at),
    UNIQUE KEY uk_outfit_id (outfit_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='穿搭推荐表';

-- 穿搭模板表（用于冷启动）
CREATE TABLE IF NOT EXISTS outfit_templates (
    id INT(10) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    items JSON COMMENT '模板单品描述',
    occasion VARCHAR(32),
    style VARCHAR(32),
    season VARCHAR(16),
    image_url VARCHAR(512),
    sort_order INT DEFAULT 0,
    is_enabled TINYINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_occasion (occasion),
    INDEX idx_style (style),
    INDEX idx_is_enabled (is_enabled),
    UNIQUE KEY uk_template_id (template_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='穿搭模板表';

-- 收藏的穿搭表
CREATE TABLE IF NOT EXISTS favorite_outfits (
    id INT(10) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL COMMENT '用户ID',
    outfit_id VARCHAR(64) NOT NULL COMMENT '穿搭ID',
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_user_id (user_id),
    INDEX idx_outfit_id (outfit_id),
    UNIQUE KEY uk_user_outfit (user_id, outfit_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='收藏的穿搭表';

-- 穿搭记录/日记表
CREATE TABLE IF NOT EXISTS outfit_records (
    id INT(10) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL COMMENT '用户ID',
    date DATE NOT NULL,
    image_urls JSON COMMENT '实际穿搭照片',
    item_ids JSON COMMENT '搭配的单品',
    outfit_id VARCHAR(64) COMMENT '关联的推荐',
    notes TEXT,
    tags JSON,
    mood VARCHAR(32) COMMENT '心情',
    weather JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_user_id (user_id),
    INDEX idx_date (date),
    UNIQUE KEY uk_record_id (record_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='穿搭记录表';

-- 用户行为记录表
CREATE TABLE IF NOT EXISTS user_behaviors (
    id INT(10) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL COMMENT '用户ID',
    outfit_id VARCHAR(64),
    behavior_type VARCHAR(32) NOT NULL COMMENT 'favorite/apply/dislike/replace',
    metadata JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_user_id (user_id),
    INDEX idx_outfit_id (outfit_id),
    INDEX idx_behavior_type (behavior_type),
    INDEX idx_created_at (created_at),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户行为记录表';

-- 搭配规则表
CREATE TABLE IF NOT EXISTS matching_rules (
    id INT(10) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    rule_id VARCHAR(64) NOT NULL,
    rule_type VARCHAR(32) NOT NULL COMMENT 'color/style/occasion/season',
    name VARCHAR(100),
    `condition` JSON COMMENT '条件',
    action JSON COMMENT '动作',
    weight DECIMAL(3,2) DEFAULT 1.0,
    priority INT DEFAULT 0,
    is_enabled TINYINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    INDEX idx_rule_type (rule_type),
    INDEX idx_is_enabled (is_enabled),
    UNIQUE KEY uk_rule_id (rule_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='搭配规则表';

-- 上传任务表
CREATE TABLE IF NOT EXISTS upload_tasks (
    id INT(10) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL COMMENT '用户ID',
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    file_hash VARCHAR(64) NOT NULL,
    total_chunks INT NOT NULL,
    uploaded_chunks JSON COMMENT '已上传块索引数组',
    upload_id VARCHAR(64) COMMENT '七牛云返回的upload_id',
    status VARCHAR(16) DEFAULT 'pending' COMMENT 'pending/uploading/completed/failed',
    file_url VARCHAR(512) COMMENT '上传完成后的文件URL',
    error_msg VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    expires_at DATETIME COMMENT '过期时间',
    INDEX idx_user_id (user_id),
    INDEX idx_status (status),
    INDEX idx_expires_at (expires_at),
    UNIQUE KEY uk_task_id (task_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='上传任务表';

-- 用户修正反馈表（用于AI准确率监控）
CREATE TABLE IF NOT EXISTS correction_feedbacks (
    id INT(10) NOT NULL AUTO_INCREMENT PRIMARY KEY,
    item_id VARCHAR(64) NOT NULL,
    user_id VARCHAR(20) NOT NULL COMMENT '已匿名化',
    api_result JSON COMMENT '豆包AI返回的原始结果',
    corrected JSON COMMENT '用户修正后的结果',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX idx_item_id (item_id),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户修正反馈表';
