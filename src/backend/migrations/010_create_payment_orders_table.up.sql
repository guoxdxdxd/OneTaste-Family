-- 创建支付订单表
CREATE TABLE payment_orders (
    id CHAR(26) PRIMARY KEY,
    order_no VARCHAR(50) UNIQUE NOT NULL,
    user_id CHAR(26) NOT NULL,
    family_id CHAR(26) NOT NULL,
    plan_type VARCHAR(20) NOT NULL,
    amount DECIMAL(10,2) NOT NULL,
    payment_method VARCHAR(20),
    status VARCHAR(20) DEFAULT 'pending',
    paid_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

COMMENT ON TABLE payment_orders IS '支付订单表';
COMMENT ON COLUMN payment_orders.order_no IS '订单号';
COMMENT ON COLUMN payment_orders.user_id IS '用户ID';
COMMENT ON COLUMN payment_orders.family_id IS '家庭ID';
COMMENT ON COLUMN payment_orders.plan_type IS '套餐类型';
COMMENT ON COLUMN payment_orders.amount IS '金额';
COMMENT ON COLUMN payment_orders.payment_method IS '支付方式：wechat, alipay';
COMMENT ON COLUMN payment_orders.status IS '状态：pending-待支付，paid-已支付，failed-失败，refunded-已退款';
COMMENT ON COLUMN payment_orders.paid_at IS '支付时间';

CREATE INDEX IF NOT EXISTS idx_payment_orders_user_id ON payment_orders(user_id);
CREATE INDEX IF NOT EXISTS idx_payment_orders_order_no ON payment_orders(order_no);
CREATE INDEX IF NOT EXISTS idx_payment_orders_status ON payment_orders(status);

CREATE TRIGGER update_payment_orders_updated_at BEFORE UPDATE ON payment_orders
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- 添加外键约束
ALTER TABLE payment_orders ADD CONSTRAINT fk_payment_orders_user_id 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT;

ALTER TABLE payment_orders ADD CONSTRAINT fk_payment_orders_family_id 
    FOREIGN KEY (family_id) REFERENCES families(id) ON DELETE RESTRICT;

