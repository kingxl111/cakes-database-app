DROP TRIGGER IF EXISTS after_cake_insert ON cakes;
DROP TRIGGER IF EXISTS after_user_insert ON users;
DROP TRIGGER IF EXISTS after_order_insert ON orders;
DROP TRIGGER IF EXISTS after_order_delete ON orders;

DROP FUNCTION IF EXISTS log_new_cake;
DROP FUNCTION IF EXISTS log_new_user;
DROP FUNCTION IF EXISTS log_new_order;
DROP FUNCTION IF EXISTS log_order_deletion;
