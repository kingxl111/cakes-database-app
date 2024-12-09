CREATE OR REPLACE FUNCTION findAverageOrderCost(customer_id int)
    RETURNS numeric(10, 2)
    LANGUAGE plpgsql
AS
$$
DECLARE
    orderCost numeric(10, 2);
BEGIN
    SELECT avg(cost)
    INTO orderCost
    FROM orders
    WHERE user_id = customer_id;
    RETURN orderCost;
END;
$$;