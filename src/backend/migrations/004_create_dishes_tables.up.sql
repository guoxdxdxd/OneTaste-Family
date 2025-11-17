-- 创建菜式表
CREATE TABLE dishes (
    id BIGSERIAL PRIMARY KEY,
    family_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50),
    description TEXT,
    image_url VARCHAR(500),
    created_by BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

COMMENT ON TABLE dishes IS '菜式表';
COMMENT ON COLUMN dishes.family_id IS '家庭ID';
COMMENT ON COLUMN dishes.name IS '菜式名称';
COMMENT ON COLUMN dishes.category IS '分类：肉类、蔬菜、汤类等';
COMMENT ON COLUMN dishes.description IS '菜式描述';
COMMENT ON COLUMN dishes.image_url IS '图片URL';
COMMENT ON COLUMN dishes.created_by IS '创建人ID';
COMMENT ON COLUMN dishes.deleted_at IS '软删除时间';

CREATE INDEX IF NOT EXISTS idx_dishes_family_id ON dishes(family_id);
CREATE INDEX IF NOT EXISTS idx_dishes_category ON dishes(category);
CREATE INDEX IF NOT EXISTS idx_dishes_created_by ON dishes(created_by);

CREATE TRIGGER update_dishes_updated_at BEFORE UPDATE ON dishes
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 创建食材表
CREATE TABLE ingredients (
    id BIGSERIAL PRIMARY KEY,
    dish_id BIGINT NOT NULL,
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

-- 创建烹饪步骤表
CREATE TABLE cooking_steps (
    id BIGSERIAL PRIMARY KEY,
    dish_id BIGINT NOT NULL,
    step_order INT NOT NULL,
    content TEXT NOT NULL,
    image_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE cooking_steps IS '烹饪步骤表';
COMMENT ON COLUMN cooking_steps.dish_id IS '菜式ID';
COMMENT ON COLUMN cooking_steps.step_order IS '步骤序号';
COMMENT ON COLUMN cooking_steps.content IS '步骤内容';
COMMENT ON COLUMN cooking_steps.image_url IS '步骤图片';

CREATE INDEX IF NOT EXISTS idx_cooking_steps_dish_id ON cooking_steps(dish_id);
CREATE INDEX IF NOT EXISTS idx_cooking_steps_step_order ON cooking_steps(dish_id, step_order);

-- 添加外键约束
ALTER TABLE dishes ADD CONSTRAINT fk_dishes_family_id 
    FOREIGN KEY (family_id) REFERENCES families(id) ON DELETE CASCADE;

ALTER TABLE dishes ADD CONSTRAINT fk_dishes_created_by 
    FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE RESTRICT;

ALTER TABLE ingredients ADD CONSTRAINT fk_ingredients_dish_id 
    FOREIGN KEY (dish_id) REFERENCES dishes(id) ON DELETE CASCADE;

ALTER TABLE cooking_steps ADD CONSTRAINT fk_cooking_steps_dish_id 
    FOREIGN KEY (dish_id) REFERENCES dishes(id) ON DELETE CASCADE;

