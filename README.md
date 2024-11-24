## cakes shop database app on PostgreSQL

# Migrations:

```
$ migrate create -ext sql -dir <migrations_dir> migration_name
```

```
$ migrate -path "/path/to/migrations" -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" down 1
```

```
$ migrate -path "/path/to/migrations" -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" force <migration_version>
```

# docker-compose
up:
```bash
sudo docker-compose --env-file local.env up --build
```
rm: 
```bash
sudo docker-compose --env-file local.env down --rmi all --volumes
```
env variables are loaded from local.env

# psql: drop all tables to demonstrate recovery
```SQL
BEGIN;
TRUNCATE TABLE users, orders, order_cakes, logs, delivery_points, deliveries, cakes, admins RESTART IDENTITY CASCADE;
COMMIT;
``` 