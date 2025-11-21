-- 拆分菜式食材与基础食材表
-- 1. 删除原 ingredients 表（菜式维度），重新创建为基础食材表
DROP TABLE IF EXISTS ingredients;

-- 2. 创建基础食材表
CREATE TABLE ingredients (
    id CHAR(26) PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    name_en VARCHAR(150),
    category VARCHAR(50),
    default_unit VARCHAR(20),
    default_amount DECIMAL(10,2),
    storage_days INT,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE ingredients IS '食材基础表';
COMMENT ON COLUMN ingredients.name IS '标准食材名称';
COMMENT ON COLUMN ingredients.name_en IS '英文/拼音名称';
COMMENT ON COLUMN ingredients.category IS '食材分类：蔬菜、肉类、调料等';
COMMENT ON COLUMN ingredients.default_unit IS '推荐单位';
COMMENT ON COLUMN ingredients.default_amount IS '参考用量';
COMMENT ON COLUMN ingredients.storage_days IS '建议存放天数';
COMMENT ON COLUMN ingredients.description IS '备注';
COMMENT ON COLUMN ingredients.is_active IS '是否可用';

CREATE INDEX IF NOT EXISTS idx_ingredients_category ON ingredients(category);
CREATE INDEX IF NOT EXISTS idx_ingredients_is_active ON ingredients(is_active);

CREATE TRIGGER update_ingredients_updated_at BEFORE UPDATE ON ingredients
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 3. 创建菜式食材关联表
CREATE TABLE dish_ingredients (
    id CHAR(26) PRIMARY KEY,
    dish_id CHAR(26) NOT NULL,
    ingredient_id CHAR(26) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    unit VARCHAR(20) NOT NULL,
    notes VARCHAR(255),
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (dish_id, ingredient_id, unit)
);

COMMENT ON TABLE dish_ingredients IS '菜式与基础食材关联表';
COMMENT ON COLUMN dish_ingredients.dish_id IS '菜式ID';
COMMENT ON COLUMN dish_ingredients.ingredient_id IS '基础食材ID';
COMMENT ON COLUMN dish_ingredients.amount IS '数量';
COMMENT ON COLUMN dish_ingredients.unit IS '单位';
COMMENT ON COLUMN dish_ingredients.notes IS '备注说明';
COMMENT ON COLUMN dish_ingredients.sort_order IS '排序';

CREATE INDEX IF NOT EXISTS idx_dish_ingredients_dish_id ON dish_ingredients(dish_id);
CREATE INDEX IF NOT EXISTS idx_dish_ingredients_ingredient_id ON dish_ingredients(ingredient_id);

ALTER TABLE dish_ingredients ADD CONSTRAINT fk_dish_ingredients_dish_id
    FOREIGN KEY (dish_id) REFERENCES dishes(id) ON DELETE CASCADE;

ALTER TABLE dish_ingredients ADD CONSTRAINT fk_dish_ingredients_ingredient_id
    FOREIGN KEY (ingredient_id) REFERENCES ingredients(id) ON DELETE RESTRICT;
