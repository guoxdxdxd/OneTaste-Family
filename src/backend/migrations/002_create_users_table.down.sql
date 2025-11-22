-- 删除用户表
-- 先删除触发器
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- 删除所有依赖 users 表的外键约束
-- 注意：如果存在 media_files 或其他表的外键约束，需要先删除
DO $$
DECLARE
    r RECORD;
BEGIN
    -- 查找所有依赖 users 表的外键约束
    FOR r IN (
        SELECT conname, conrelid::regclass AS table_name
        FROM pg_constraint
        WHERE confrelid = 'users'::regclass
        AND contype = 'f'
    ) LOOP
        EXECUTE 'ALTER TABLE ' || r.table_name || ' DROP CONSTRAINT IF EXISTS ' || r.conname || ' CASCADE';
    END LOOP;
END $$;

-- 删除用户表（使用 CASCADE 确保删除所有依赖）
DROP TABLE IF EXISTS users CASCADE;

