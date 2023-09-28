package database

import (
	"fmt"

	_ "github.com/jackc/pgx/v5" // postgresql driver
	"github.com/jmoiron/sqlx"
)

// Database struct contains sql pointer
type Database struct {
	Name     string `json:"name"`
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	URI      string
	*sqlx.DB
}

// OpenDatabase open database
func OpenDatabase(db Database) (*Database, error) {
	var err error
	db.GenURI()
	db.DB, err = sqlx.Open(db.Driver, db.URI)
	if err != nil {
		fmt.Printf("Open sql (%v): %v", db.URI, err)
		panic(err)
	}
	if err = db.Ping(); err != nil {
		fmt.Printf("Ping sql: %v", err)
		panic(err)
	}
	return &db, err
}

// ExecProcedure executes stored procedure
func (db *Database) ExecProcedure(q string) {
	fmt.Println(q)
	_, err := db.Exec(q)
	if err != nil {
		panic(err)
	}
}

// GenURI generate db uri string
func (db *Database) GenURI() {
	switch db.Driver {
	case "postgres", "pgx":
		port := "5432"
		if db.Port != "" {
			port = db.Port
		}
		db.URI = "postgres://" + db.Username + ":" + db.Password + "@" + db.Host + ":" + port + "/" + db.Database + "?sslmode=disable"
	case "mssql":
		db.URI = "server=" + db.Host + ";user id=" + db.Username + ";password=" + db.Password + ";database=" + db.Database + ";encrypt=disable;connection timeout=7200;keepAlive=30"
	}
}
