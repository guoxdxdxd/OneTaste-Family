-- 预置测试用户和家庭数据
INSERT INTO public.users (id, phone, password, nickname, avatar, status, created_at, updated_at) VALUES ('01KAFWWD10JVN1X6DE55W06FZ0', '13580724968', '$2a$10$cdWWJAQpI0ip61/fu798rOnaG0rr.L9t3zz9J135KEnKePpcynnVm', '小郭', '', 1, '2025-11-20 13:48:02.492983', '2025-11-20 13:48:02.492983');

INSERT INTO public.users (id, phone, password, nickname, avatar, status, created_at, updated_at) VALUES ('01KAFWYD6SVE37C8RBEX4SXEMH', '13644538765', '$2a$10$tzCi1hHz6CAKWOGhI0aHHOpB19r22J3IjTN4.BXQuxIJm762.YJ.i', '小竹', '', 1, '2025-11-20 13:49:08.221854', '2025-11-20 13:49:08.221854');

INSERT INTO public.families (id, name, description, owner_id, max_dishes, status, created_at, updated_at) VALUES ('01KAFWYVKQS71KYJNQ9MZX3QWX', '小竹爱小锅', '', '01KAFWYD6SVE37C8RBEX4SXEMH', 30, 1, '2025-11-20 13:49:22.955508', '2025-11-20 13:49:22.955508');

-- 添加家庭成员
-- 小竹作为家庭创建者（owner）
INSERT INTO public.family_members (id, family_id, user_id, role, status, joined_at, created_at, updated_at) VALUES 
('01KAFWYZ1KQS71KYJNQ9MZX3QW', '01KAFWYVKQS71KYJNQ9MZX3QWX', '01KAFWYD6SVE37C8RBEX4SXEMH', 'owner', 1, '2025-11-20 13:49:22.955508', '2025-11-20 13:49:22.955508', '2025-11-20 13:49:22.955508');

-- 小郭作为家庭成员（member）
INSERT INTO public.family_members (id, family_id, user_id, role, status, joined_at, created_at, updated_at) VALUES 
('01KAFWYZ2KQS71KYJNQ9MZX3QW', '01KAFWYVKQS71KYJNQ9MZX3QWX', '01KAFWWD10JVN1X6DE55W06FZ0', 'member', 1, '2025-11-20 13:50:00.000000', '2025-11-20 13:50:00.000000', '2025-11-20 13:50:00.000000');

