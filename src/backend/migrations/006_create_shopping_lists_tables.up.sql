-- 创建购物清单表
CREATE TABLE shopping_lists (
    id BIGSERIAL PRIMARY KEY,
    family_id BIGINT NOT NULL,
    name VARCHAR(100),
    start_date DATE,
    end_date DATE,
    status VARCHAR(20) DEFAULT 'pending',
    created_by BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE shopping_lists IS '购物清单表';
COMMENT ON COLUMN shopping_lists.family_id IS '家庭ID';
COMMENT ON COLUMN shopping_lists.name IS '清单名称';
COMMENT ON COLUMN shopping_lists.start_date IS '开始日期';
COMMENT ON COLUMN shopping_lists.end_date IS '结束日期';
COMMENT ON COLUMN shopping_lists.status IS '状态：pending-待购买，completed-已完成';
COMMENT ON COLUMN shopping_lists.created_by IS '创建人ID';

CREATE INDEX IF NOT EXISTS idx_shopping_lists_family_id ON shopping_lists(family_id);
CREATE INDEX IF NOT EXISTS idx_shopping_lists_status ON shopping_lists(status);

CREATE TRIGGER update_shopping_lists_updated_at BEFORE UPDATE ON shopping_lists
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 创建购物清单项表
CREATE TABLE shopping_list_items (
    id BIGSERIAL PRIMARY KEY,
    list_id BIGINT NOT NULL,
    ingredient_name VARCHAR(100) NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL,
    unit VARCHAR(20) NOT NULL,
    category VARCHAR(50),
    status VARCHAR(20) DEFAULT 'pending',
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE shopping_list_items IS '购物清单项表';
COMMENT ON COLUMN shopping_list_items.list_id IS '清单ID';
COMMENT ON COLUMN shopping_list_items.ingredient_name IS '食材名称';
COMMENT ON COLUMN shopping_list_items.total_amount IS '总数量';
COMMENT ON COLUMN shopping_list_items.unit IS '单位';
COMMENT ON COLUMN shopping_list_items.category IS '食材分类';
COMMENT ON COLUMN shopping_list_items.status IS '状态：pending-待购买，purchased-已购买';
COMMENT ON COLUMN shopping_list_items.sort_order IS '排序';

CREATE INDEX IF NOT EXISTS idx_shopping_list_items_list_id ON shopping_list_items(list_id);
CREATE INDEX IF NOT EXISTS idx_shopping_list_items_status ON shopping_list_items(status);
CREATE INDEX IF NOT EXISTS idx_shopping_list_items_category ON shopping_list_items(category);

CREATE TRIGGER update_shopping_list_items_updated_at BEFORE UPDATE ON shopping_list_items
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 添加外键约束
ALTER TABLE shopping_lists ADD CONSTRAINT fk_shopping_lists_family_id 
    FOREIGN KEY (family_id) REFERENCES families(id) ON DELETE CASCADE;

ALTER TABLE shopping_lists ADD CONSTRAINT fk_shopping_lists_created_by 
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE RESTRICT;

ALTER TABLE shopping_list_items ADD CONSTRAINT fk_shopping_list_items_list_id 
    FOREIGN KEY (list_id) REFERENCES shopping_lists(id) ON DELETE CASCADE;

