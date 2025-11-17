-- 删除购物清单项表
DROP TRIGGER IF EXISTS update_shopping_list_items_updated_at ON shopping_list_items;
ALTER TABLE shopping_list_items DROP CONSTRAINT IF EXISTS fk_shopping_list_items_list_id;
DROP TABLE IF EXISTS shopping_list_items;

-- 删除购物清单表
DROP TRIGGER IF EXISTS update_shopping_lists_updated_at ON shopping_lists;
ALTER TABLE shopping_lists DROP CONSTRAINT IF EXISTS fk_shopping_lists_created_by;
ALTER TABLE shopping_lists DROP CONSTRAINT IF EXISTS fk_shopping_lists_family_id;
DROP TABLE IF EXISTS shopping_lists;

