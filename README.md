## cakes shop database app on PostgreSQL

![database architecture](acrhitecture_v1.png)

![use_case](use_case.png)

![alt text](context_diagram.png)

# Migrations:
```
migrate -path "/path/to/migrations" -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" down 1
```

```
$ migrate -path "/path/to/migrations" -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" force <migration_version>
```

