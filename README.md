# Golang database migration

## Usage

Run
```bash
go get -t github.com/srijankrs/go_db_migration
```
Import
```golang
import "github.com/srijankrs/go_db_migration"
```
Code ( for postgres )
```golang
go_db_migration.Migrate("postgres", "localhost", "5432", "postgres_user", "postgres_password", "db_name", "db_migrations_sql_files_path") 
```
Load test will be updated soon.
