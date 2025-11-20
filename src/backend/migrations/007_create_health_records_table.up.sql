-- 创建身体状况记录表
CREATE TABLE health_records (
    id CHAR(26) PRIMARY KEY,
    user_id CHAR(26) NOT NULL,
    family_id CHAR(26) NOT NULL,
    diseases TEXT,
    work_status VARCHAR(50),
    stress_level VARCHAR(20),
    body_feelings TEXT,
    raw_data TEXT,
    analysis_result TEXT,
    recommendations TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE health_records IS '身体状况记录表';
COMMENT ON COLUMN health_records.user_id IS '用户ID';
COMMENT ON COLUMN health_records.family_id IS '家庭ID';
COMMENT ON COLUMN health_records.diseases IS '疾病信息（JSON数组）';
COMMENT ON COLUMN health_records.work_status IS '工作状态';
COMMENT ON COLUMN health_records.stress_level IS '压力水平：低、中、高';
COMMENT ON COLUMN health_records.body_feelings IS '身体感觉描述';
COMMENT ON COLUMN health_records.raw_data IS '原始录入数据';
COMMENT ON COLUMN health_records.analysis_result IS 'AI分析结果';
COMMENT ON COLUMN health_records.recommendations IS '建议（JSON数组）';

CREATE INDEX IF NOT EXISTS idx_health_records_user_id ON health_records(user_id);
CREATE INDEX IF NOT EXISTS idx_health_records_family_id ON health_records(family_id);
CREATE INDEX IF NOT EXISTS idx_health_records_created_at ON health_records(created_at);

-- 添加外键约束
ALTER TABLE health_records ADD CONSTRAINT fk_health_records_user_id 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE health_records ADD CONSTRAINT fk_health_records_family_id 
    FOREIGN KEY (family_id) REFERENCES families(id) ON DELETE CASCADE;

