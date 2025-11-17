-- 删除家庭成员表
DROP TRIGGER IF EXISTS update_family_members_updated_at ON family_members;
ALTER TABLE family_members DROP CONSTRAINT IF EXISTS fk_family_members_user_id;
ALTER TABLE family_members DROP CONSTRAINT IF EXISTS fk_family_members_family_id;
DROP TABLE IF EXISTS family_members;

-- 删除家庭表
DROP TRIGGER IF EXISTS update_families_updated_at ON families;
ALTER TABLE families DROP CONSTRAINT IF EXISTS fk_families_owner_id;
DROP TABLE IF EXISTS families;

