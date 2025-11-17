-- 创建菜单表
CREATE TABLE menus (
    id BIGSERIAL PRIMARY KEY,
    family_id BIGINT NOT NULL,
    date DATE NOT NULL,
    meal_type VARCHAR(20) NOT NULL,
    created_by BIGINT NOT NULL,
    source VARCHAR(20) DEFAULT 'manual',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE menus IS '菜单表';
COMMENT ON COLUMN menus.family_id IS '家庭ID';
COMMENT ON COLUMN menus.date IS '日期';
COMMENT ON COLUMN menus.meal_type IS '餐次：breakfast-早餐，lunch-午餐，dinner-晚餐';
COMMENT ON COLUMN menus.created_by IS '创建人ID';
COMMENT ON COLUMN menus.source IS '来源：manual-手动，ai-AI生成';

CREATE INDEX IF NOT EXISTS idx_menus_family_date ON menus(family_id, date);
CREATE INDEX IF NOT EXISTS idx_menus_meal_type ON menus(meal_type);
CREATE INDEX IF NOT EXISTS idx_menus_created_by ON menus(created_by);

CREATE TRIGGER update_menus_updated_at BEFORE UPDATE ON menus
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 创建菜单菜式关联表
CREATE TABLE menu_dishes (
    id BIGSERIAL PRIMARY KEY,
    menu_id BIGINT NOT NULL,
    dish_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (menu_id, dish_id)
);

COMMENT ON TABLE menu_dishes IS '菜单菜式关联表';
COMMENT ON COLUMN menu_dishes.menu_id IS '菜单ID';
COMMENT ON COLUMN menu_dishes.dish_id IS '菜式ID';

CREATE INDEX IF NOT EXISTS idx_menu_dishes_menu_id ON menu_dishes(menu_id);
CREATE INDEX IF NOT EXISTS idx_menu_dishes_dish_id ON menu_dishes(dish_id);

-- 添加外键约束
ALTER TABLE menus ADD CONSTRAINT fk_menus_family_id 
    FOREIGN KEY (family_id) REFERENCES families(id) ON DELETE CASCADE;

ALTER TABLE menus ADD CONSTRAINT fk_menus_created_by 
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE RESTRICT;

ALTER TABLE menu_dishes ADD CONSTRAINT fk_menu_dishes_menu_id 
    FOREIGN KEY (menu_id) REFERENCES menus(id) ON DELETE CASCADE;

ALTER TABLE menu_dishes ADD CONSTRAINT fk_menu_dishes_dish_id 
    FOREIGN KEY (dish_id) REFERENCES dishes(id) ON DELETE CASCADE;

