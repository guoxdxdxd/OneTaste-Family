-- 创建家庭表
CREATE TABLE families (
    id CHAR(26) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    owner_id CHAR(26) NOT NULL,
    max_dishes INT DEFAULT 30,
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE families IS '家庭表';
COMMENT ON COLUMN families.name IS '家庭名称';
COMMENT ON COLUMN families.description IS '家庭描述';
COMMENT ON COLUMN families.owner_id IS '创建人ID';
COMMENT ON COLUMN families.max_dishes IS '最大菜式数量';
COMMENT ON COLUMN families.status IS '状态：1-正常，0-解散';

CREATE INDEX IF NOT EXISTS idx_families_owner_id ON families(owner_id);
CREATE INDEX IF NOT EXISTS idx_families_status ON families(status);

CREATE TRIGGER update_families_updated_at BEFORE UPDATE ON families
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 创建家庭成员表
CREATE TABLE family_members (
    id CHAR(26) PRIMARY KEY,
    family_id CHAR(26) NOT NULL,
    user_id CHAR(26) NOT NULL,
    role VARCHAR(20) DEFAULT 'member',
    status SMALLINT DEFAULT 1,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (family_id, user_id)
);

COMMENT ON TABLE family_members IS '家庭成员表';
COMMENT ON COLUMN family_members.family_id IS '家庭ID';
COMMENT ON COLUMN family_members.user_id IS '用户ID';
COMMENT ON COLUMN family_members.role IS '角色：owner-创建人，member-成员';
COMMENT ON COLUMN family_members.status IS '状态：1-正常，0-已退出';
COMMENT ON COLUMN family_members.joined_at IS '加入时间';

CREATE INDEX IF NOT EXISTS idx_family_members_family_id ON family_members(family_id);
CREATE INDEX IF NOT EXISTS idx_family_members_user_id ON family_members(user_id);

CREATE TRIGGER update_family_members_updated_at BEFORE UPDATE ON family_members
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 添加外键约束
ALTER TABLE families ADD CONSTRAINT fk_families_owner_id 
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE RESTRICT;

ALTER TABLE family_members ADD CONSTRAINT fk_family_members_family_id 
    FOREIGN KEY (family_id) REFERENCES families(id) ON DELETE CASCADE;

ALTER TABLE family_members ADD CONSTRAINT fk_family_members_user_id 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

