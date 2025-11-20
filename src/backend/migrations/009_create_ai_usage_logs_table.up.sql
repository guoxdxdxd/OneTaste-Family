-- 创建AI调用记录表
CREATE TABLE ai_usage_logs (
    id CHAR(26) PRIMARY KEY,
    user_id CHAR(26) NOT NULL,
    family_id CHAR(26) NOT NULL,
    feature_type VARCHAR(50) NOT NULL,
    period_type VARCHAR(20) NOT NULL,
    period_date DATE NOT NULL,
    usage_count INT DEFAULT 0,
    limit_count INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, feature_type, period_type, period_date)
);

COMMENT ON TABLE ai_usage_logs IS 'AI调用记录表';
COMMENT ON COLUMN ai_usage_logs.user_id IS '用户ID';
COMMENT ON COLUMN ai_usage_logs.family_id IS '家庭ID';
COMMENT ON COLUMN ai_usage_logs.feature_type IS '功能类型：voice_input, health_analyze, menu_generate, menu_analyze, cooking_optimize';
COMMENT ON COLUMN ai_usage_logs.period_type IS '周期类型：daily, weekly';
COMMENT ON COLUMN ai_usage_logs.period_date IS '周期日期';
COMMENT ON COLUMN ai_usage_logs.usage_count IS '使用次数';
COMMENT ON COLUMN ai_usage_logs.limit_count IS '限制次数';

CREATE INDEX IF NOT EXISTS idx_ai_usage_logs_family_id ON ai_usage_logs(family_id);
CREATE INDEX IF NOT EXISTS idx_ai_usage_logs_period_date ON ai_usage_logs(period_date);

CREATE TRIGGER update_ai_usage_logs_updated_at BEFORE UPDATE ON ai_usage_logs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 添加外键约束
ALTER TABLE ai_usage_logs ADD CONSTRAINT fk_ai_usage_logs_user_id 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE ai_usage_logs ADD CONSTRAINT fk_ai_usage_logs_family_id 
    FOREIGN KEY (family_id) REFERENCES families(id) ON DELETE CASCADE;

