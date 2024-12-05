DO $$ BEGIN
    -- Роль authorizer с правами на вставку в таблицу users
    IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'authorizer') THEN
        CREATE ROLE authorizer WITH LOGIN PASSWORD 'authorizer_password';
    END IF;
END $$;

GRANT CONNECT ON DATABASE postgres TO authorizer;
GRANT USAGE ON SCHEMA public TO authorizer;
GRANT INSERT, SELECT ON TABLE users TO authorizer;
GRANT USAGE, SELECT ON SEQUENCE users_id_seq TO authorizer;
GRANT INSERT, SELECT ON TABLE logs TO authorizer;
GRANT USAGE, SELECT ON SEQUENCE logs_id_seq TO authorizer;