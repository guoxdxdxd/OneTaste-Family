-- 创建会员表
CREATE TABLE memberships (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    family_id BIGINT NOT NULL,
    plan_type VARCHAR(20) NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    started_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE memberships IS '会员表';
COMMENT ON COLUMN memberships.user_id IS '用户ID';
COMMENT ON COLUMN memberships.family_id IS '家庭ID';
COMMENT ON COLUMN memberships.plan_type IS '套餐类型：monthly-月付，yearly-年付';
COMMENT ON COLUMN memberships.status IS '状态：active-有效，expired-已过期，cancelled-已取消';
COMMENT ON COLUMN memberships.started_at IS '开始时间';
COMMENT ON COLUMN memberships.expires_at IS '到期时间';

CREATE INDEX IF NOT EXISTS idx_memberships_user_id ON memberships(user_id);
CREATE INDEX IF NOT EXISTS idx_memberships_family_id ON memberships(family_id);
CREATE INDEX IF NOT EXISTS idx_memberships_status ON memberships(status);
CREATE INDEX IF NOT EXISTS idx_memberships_expires_at ON memberships(expires_at);

CREATE TRIGGER update_memberships_updated_at BEFORE UPDATE ON memberships
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 添加外键约束
ALTER TABLE memberships ADD CONSTRAINT fk_memberships_user_id 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE memberships ADD CONSTRAINT fk_memberships_family_id 
    FOREIGN KEY (family_id) REFERENCES families(id) ON DELETE CASCADE;

