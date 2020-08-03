# Golang database migration

## Usage

Run
```bash
go get -t github.com/srijankrs/godbmigration
```
Import
```golang
import "github.com/srijankrs/godbmigration"
```
Code ( for postgres )
```golang
godbmigration.Migrate("postgres", "localhost", "5432", "postgres_user", "postgres_password", "db_name", "db_migrations_sql_files_path") 
```
Load test will be updated soon.
