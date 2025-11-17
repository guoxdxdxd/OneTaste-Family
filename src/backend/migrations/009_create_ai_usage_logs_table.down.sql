-- 删除AI调用记录表
DROP TRIGGER IF EXISTS update_ai_usage_logs_updated_at ON ai_usage_logs;
ALTER TABLE ai_usage_logs DROP CONSTRAINT IF EXISTS fk_ai_usage_logs_family_id;
ALTER TABLE ai_usage_logs DROP CONSTRAINT IF EXISTS fk_ai_usage_logs_user_id;
DROP TABLE IF EXISTS ai_usage_logs;

