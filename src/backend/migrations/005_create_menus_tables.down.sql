-- 删除菜单菜式关联表
ALTER TABLE menu_dishes DROP CONSTRAINT IF EXISTS fk_menu_dishes_dish_id;
ALTER TABLE menu_dishes DROP CONSTRAINT IF EXISTS fk_menu_dishes_menu_id;
DROP TABLE IF EXISTS menu_dishes;

-- 删除菜单表
DROP TRIGGER IF EXISTS update_menus_updated_at ON menus;
ALTER TABLE menus DROP CONSTRAINT IF EXISTS fk_menus_created_by;
ALTER TABLE menus DROP CONSTRAINT IF EXISTS fk_menus_family_id;
DROP TABLE IF EXISTS menus;

