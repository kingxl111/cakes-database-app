CREATE OR REPLACE FUNCTION log_new_cake()
    RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO logs (level, message)
    VALUES ('INFO', 'New cake added with ID: ' || NEW.id || ', Description: ' || NEW.description);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER after_cake_insert
    AFTER INSERT ON cakes
    FOR EACH ROW
EXECUTE FUNCTION log_new_cake();

CREATE OR REPLACE FUNCTION log_new_user()
    RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO logs (level, message)
    VALUES ('INFO', 'New user added with ID: ' || NEW.id || ', Full Name: ' || NEW.fullname || ', Username: ' || NEW.username);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER after_user_insert
    AFTER INSERT ON users
    FOR EACH ROW
EXECUTE FUNCTION log_new_user();

CREATE OR REPLACE FUNCTION log_new_order()
    RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO logs (level, message)
    VALUES ('INFO', 'New order added with ID: ' || NEW.id || ', User ID: ' || NEW.user_id || ', Status: ' || NEW.order_status);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER after_order_insert
    AFTER INSERT ON orders
    FOR EACH ROW
EXECUTE FUNCTION log_new_order();

CREATE OR REPLACE FUNCTION log_order_deletion()
    RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO logs (level, message)
    VALUES ('WARNING', 'Order deleted with ID: ' || OLD.id || ', User ID: ' || OLD.user_id);
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER after_order_delete
    AFTER DELETE ON orders
    FOR EACH ROW
EXECUTE FUNCTION log_order_deletion();
