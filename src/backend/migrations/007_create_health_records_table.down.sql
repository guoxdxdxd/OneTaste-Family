-- 删除身体状况记录表
ALTER TABLE health_records DROP CONSTRAINT IF EXISTS fk_health_records_family_id;
ALTER TABLE health_records DROP CONSTRAINT IF EXISTS fk_health_records_user_id;
DROP TABLE IF EXISTS health_records;

