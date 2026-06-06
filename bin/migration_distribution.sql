-- Migration: Add distribution (affiliate) tables
-- Version: v0.4+

-- Distributors table
CREATE TABLE IF NOT EXISTS distributors (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL UNIQUE,
    level INTEGER DEFAULT 1 COMMENT '1=一级分销商, 2=二级分销商, 3=三级分销商',
    parent_id INTEGER DEFAULT 0 COMMENT '上级分销商ID',
    status INTEGER DEFAULT 1 COMMENT '1=启用, 0=禁用',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    INDEX idx_user_id (user_id),
    INDEX idx_parent_id (parent_id),
    INDEX idx_deleted_at (deleted_at)
);

-- Commission rates configuration table
CREATE TABLE IF NOT EXISTS commission_rates (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name VARCHAR(128) NOT NULL COMMENT '配置名称',
    type VARCHAR(32) NOT NULL COMMENT '配置类型: channel_topup, subscription, redemption',
    type_id INTEGER DEFAULT 0 COMMENT '关联类型ID，0表示该类型通用配置',
    level1_rate DECIMAL(5,4) DEFAULT 0 COMMENT '一级佣金比例 (0.0000-1.0000)',
    level2_rate DECIMAL(5,4) DEFAULT 0 COMMENT '二级佣金比例',
    level3_rate DECIMAL(5,4) DEFAULT 0 COMMENT '三级佣金比例',
    enabled BOOLEAN DEFAULT TRUE,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_type (type),
    INDEX idx_type_id (type_id)
);

-- Commission logs table
CREATE TABLE IF NOT EXISTS commission_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    distributor_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL COMMENT '消费用户ID',
    order_id VARCHAR(64) COMMENT '订单ID',
    order_type VARCHAR(32) COMMENT '订单类型: channel_topup, subscription, etc.',
    amount INTEGER DEFAULT 0 COMMENT '消费金额(Quota)',
    commission INTEGER DEFAULT 0 COMMENT '佣金金额(Quota)',
    level INTEGER DEFAULT 1 COMMENT '佣金层级: 1/2/3',
    status INTEGER DEFAULT 0 COMMENT '0=待确认, 1=已结算, 2=已提现',
    remark VARCHAR(255),
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    settled_at BIGINT COMMENT '结算时间戳',
    INDEX idx_distributor_id (distributor_id),
    INDEX idx_user_id (user_id),
    INDEX idx_order_id (order_id),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
);

-- Add foreign key constraints (if supported)
-- Note: SQLite doesn't support foreign keys by default, MySQL/PostgreSQL will enforce them
-- ALTER TABLE distributors ADD CONSTRAINT fk_distributors_user FOREIGN KEY (user_id) REFERENCES users(id);
-- ALTER TABLE distributors ADD CONSTRAINT fk_distributors_parent FOREIGN KEY (parent_id) REFERENCES distributors(id);
-- ALTER TABLE commission_logs ADD CONSTRAINT fk_commission_logs_distributor FOREIGN KEY (distributor_id) REFERENCES distributors(id);
-- ALTER TABLE commission_logs ADD CONSTRAINT fk_commission_logs_user FOREIGN KEY (user_id) REFERENCES users(id);

-- Insert default commission rate for channel topup (10%, 5%, 2%)
INSERT INTO commission_rates (name, type, type_id, level1_rate, level2_rate, level3_rate, enabled)
VALUES ('渠道充值佣金', 'channel_topup', 0, 0.10, 0.05, 0.02, TRUE);

-- Insert default commission rate for subscription (5%, 2.5%, 1%)
INSERT INTO commission_rates (name, type, type_id, level1_rate, level2_rate, level3_rate, enabled)
VALUES ('订阅佣金', 'subscription', 0, 0.05, 0.025, 0.01, TRUE);
