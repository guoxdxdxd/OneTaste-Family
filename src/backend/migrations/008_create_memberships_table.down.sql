-- 删除会员表
DROP TRIGGER IF EXISTS update_memberships_updated_at ON memberships;
ALTER TABLE memberships DROP CONSTRAINT IF EXISTS fk_memberships_family_id;
ALTER TABLE memberships DROP CONSTRAINT IF EXISTS fk_memberships_user_id;
DROP TABLE IF EXISTS memberships;

