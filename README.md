Golang database migration.

Usage :-

run 'go get -t github.com/srijankrs/go_db_migration'
import "github.com/srijankrs/go_db_migration"
go_db_migration.Migrate("postgres", "localhost", "5432", "postgres_user", "postgres_password", "db_name", "db/migrations") 
