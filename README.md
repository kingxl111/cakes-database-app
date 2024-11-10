## cakes shop database app on PostgreSQL

![database architecture](acrhitecture_v1.png)

![use_case](use_case.png)

![alt text](context_diagram.png)

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