-- 删除系统配置表
DROP TRIGGER IF EXISTS update_system_configs_updated_at ON system_configs;
DROP TABLE IF EXISTS system_configs;

