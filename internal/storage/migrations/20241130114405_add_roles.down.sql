DO $$ BEGIN
    IF EXISTS (SELECT FROM pg_roles WHERE rolname = 'authorizer') THEN
        REVOKE INSERT, SELECT ON TABLE users FROM authorizer;
        REVOKE USAGE, SELECT ON SEQUENCE users_id_seq FROM authorizer;
        REVOKE INSERT, SELECT ON TABLE logs FROM authorizer;
        REVOKE USAGE, SELECT ON SEQUENCE logs_id_seq FROM authorizer;

        REVOKE CONNECT ON DATABASE postgres FROM authorizer;

        DROP ROLE authorizer;
    END IF;
END $$;
