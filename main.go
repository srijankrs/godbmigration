package godbmigration

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

const SchemaTableName = "migration_schema_version"
const POSTGRES = "postgres"
const SQL = "sql"
const SqlSchemaTableQuery =
	"CREATE TABLE IF NOT EXISTS " + SchemaTableName + " ( " +
	"id int AUTO_INCREMENT, " +
	"version varchar(255) NOT NULL," +
	"description varchar(255) NOT NULL," +
	"file_name varchar(255) NOT NULL," +
	"hash varchar(255) NOT NULL," +
	"executed_on TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP," +
	"execution_time varchar(255) NOT NULL," +
	"PRIMARY KEY (id));"

const PostgresSchemaTableQuery =
	"CREATE TABLE IF NOT EXISTS " + SchemaTableName + " ( " +
	"id SERIAL PRIMARY KEY, " +
	"version varchar(255) NOT NULL," +
	"description varchar(255) NOT NULL," +
	"file_name varchar(255) NOT NULL," +
	"hash varchar(255) NOT NULL," +
	"executed_on timestamp default current_timestamp," +
	"execution_time varchar(255) NOT NULL);"

func Migrate(dbDriver string, host string, port string, dbUser string, dbPassword string, dbName string, migrationPath string) {

	var connInfo string
	var initQuery string
	if strings.Compare(dbDriver, POSTGRES) == 0 {
		connInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, dbUser, dbPassword, dbName)
		initQuery = PostgresSchemaTableQuery
	}
	if strings.Compare(dbDriver, SQL) == 0 {
		connInfo = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
			dbUser, dbPassword, host, port, dbName)
		initQuery = SqlSchemaTableQuery
	}
	db, err := sql.Open(dbDriver, connInfo)
	if err != nil {
		panic(err)
		return
	}
	defer db.Close()
	_, err = db.Query(initQuery)
	if err != nil {
		panic(err)
		return
	}

	files, err := ioutil.ReadDir(migrationPath)
	if err != nil {
		log.Fatal(err)
	}

	selectQuery := "SELECT version, hash FROM "+SchemaTableName
	rows, err := db.Query(selectQuery)
	defer rows.Close()
	if err != nil{
		panic(err)
	}
	var data = make(map[string]string)
	for rows.Next() {
		var versionColumn = new(string)
		var hashColumn = new(string)
		err := rows.Scan(versionColumn, hashColumn)
		if err != nil {
			log.Fatal(err)
		}
		data[*versionColumn]= *hashColumn
	}

	hasher := sha256.New()
	for _, f := range files {
		fileName := f.Name()
		info := strings.Split(fileName, "__")
		version := info[0]
		description := strings.Split(info[1],".")[0]

		queryByte, err := ioutil.ReadFile(migrationPath+"/"+fileName)
		if err != nil {
			panic(err)
		}
		hasher.Write(queryByte)
		hash := hex.EncodeToString(hasher.Sum(nil))

		if val, ok := data[version]; ok {
			if strings.Compare(val, hash) != 0 {
				log.Printf("Migration error for version id %s, db hash: %s, file hash: %s", version, val, hash)
				panic("Migration error")
			}

			log.Printf("Migration checked for version id %s with hash: %s", version, hash)
			continue
		}

		txn, err := db.Begin()
		if err != nil {
			panic(err)
		}
		query := string(queryByte)
		timePre := time.Now()
		_,err = txn.Exec(query)
		executionTime := time.Now().Sub(timePre).String()
		if err != nil{
			log.Printf("%s Migration error", version)
		}

		sqlStatement := "INSERT INTO " + SchemaTableName + "(version, description, file_name, hash, execution_time) VALUES ('" +
			version + "','" + description + "','" + fileName + "','" + hash + "','" + executionTime + "');"

		_, err = txn.Query(sqlStatement)
		if err != nil {
			log.Printf("%s Error in updating schema", version)
		}

		err = txn.Commit()
		if err != nil {
			log.Printf("Error occured in migration for file %s [%s]", fileName, err)
			panic("Migration Failed")
		}

		log.Printf("Successfully migrated version id %s with hash: %s in time %s", version, hash, executionTime)
	}

	log.Printf("Migration Completed")
}


