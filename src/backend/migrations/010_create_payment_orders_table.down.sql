-- 删除支付订单表
DROP TRIGGER IF EXISTS update_payment_orders_updated_at ON payment_orders;
ALTER TABLE payment_orders DROP CONSTRAINT IF EXISTS fk_payment_orders_family_id;
ALTER TABLE payment_orders DROP CONSTRAINT IF EXISTS fk_payment_orders_user_id;
DROP TABLE IF EXISTS payment_orders;

