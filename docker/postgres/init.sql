-- PostgreSQL 初始化脚本
-- 此脚本会在数据库首次创建时自动执行

-- 设置时区
SET timezone = 'Asia/Shanghai';

-- 创建扩展（如果需要）
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- 可以在这里添加其他初始化 SQL
-- 例如：创建额外的数据库、用户、模式等

-- 示例：创建一个只读用户（可选）
-- CREATE USER readonly WITH PASSWORD 'readonly_password';
-- GRANT CONNECT ON DATABASE onetaste_family TO readonly;
-- GRANT USAGE ON SCHEMA public TO readonly;
-- GRANT SELECT ON ALL TABLES IN SCHEMA public TO readonly;
-- ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO readonly;

