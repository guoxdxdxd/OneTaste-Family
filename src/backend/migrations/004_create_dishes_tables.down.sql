-- 删除烹饪步骤表
ALTER TABLE cooking_steps DROP CONSTRAINT IF EXISTS fk_cooking_steps_dish_id;
DROP TABLE IF EXISTS cooking_steps;

-- 删除食材表
ALTER TABLE ingredients DROP CONSTRAINT IF EXISTS fk_ingredients_dish_id;
DROP TABLE IF EXISTS ingredients;

-- 删除菜式表
DROP TRIGGER IF EXISTS update_dishes_updated_at ON dishes;
ALTER TABLE dishes DROP CONSTRAINT IF EXISTS fk_dishes_created_by;
ALTER TABLE dishes DROP CONSTRAINT IF EXISTS fk_dishes_family_id;
DROP TABLE IF EXISTS dishes;

