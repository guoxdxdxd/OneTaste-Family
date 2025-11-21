-- 回滚菜式食材与基础食材拆分
DROP TABLE IF EXISTS dish_ingredients;
DROP TRIGGER IF EXISTS update_ingredients_updated_at ON ingredients;
DROP TABLE IF EXISTS ingredients;

-- 还原原有的菜式食材表结构
CREATE TABLE ingredients (
    id CHAR(26) PRIMARY KEY,
    dish_id CHAR(26) NOT NULL,
    name VARCHAR(100) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    unit VARCHAR(20) NOT NULL,
    category VARCHAR(50),
    storage_days INT,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE ingredients IS '食材表';
COMMENT ON COLUMN ingredients.dish_id IS '菜式ID';
COMMENT ON COLUMN ingredients.name IS '食材名称';
COMMENT ON COLUMN ingredients.amount IS '数量';
COMMENT ON COLUMN ingredients.unit IS '单位：g、kg、个、勺等';
COMMENT ON COLUMN ingredients.category IS '食材分类：蔬菜、肉类、调料等';
COMMENT ON COLUMN ingredients.storage_days IS '可存放天数';
COMMENT ON COLUMN ingredients.sort_order IS '排序';

CREATE INDEX IF NOT EXISTS idx_ingredients_dish_id ON ingredients(dish_id);
CREATE INDEX IF NOT EXISTS idx_ingredients_category ON ingredients(category);

ALTER TABLE ingredients ADD CONSTRAINT fk_ingredients_dish_id
    FOREIGN KEY (dish_id) REFERENCES dishes(id) ON DELETE CASCADE;
