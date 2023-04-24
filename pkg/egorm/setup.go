package egorm

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB
var initOnce sync.Once
var initOnceErr error
var sqliteOps *SQLiteConnectOpts
var postgresOps *PostgresConnectOpts

type PostgresConnectOpts struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

func SetPostgresConnectOpts(ops *PostgresConnectOpts) {
	postgresOps = ops
}

func setupPostgres() (*gorm.DB, error) {

	if postgresOps == nil {
		dbHost := os.Getenv("EGORM_POSTGRES_DB_HOST")
		dbPort := os.Getenv("EGORM_POSTGRES_DB_PORT")
		dbName := os.Getenv("EGORM_POSTGRES_DB_DATABASE")
		dbUser := os.Getenv("EGORM_POSTGRES_DB_USER")
		dbPassword := os.Getenv("EGORM_POSTGRES_DB_PASSWORD")
		port, _ := strconv.Atoi(dbPort)
		postgresOps = &PostgresConnectOpts{
			Host:     dbHost,
			Port:     port,
			Database: dbName,
			User:     dbUser,
			Password: dbPassword,
		}
	}

	connectionString := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		postgresOps.Host,
		postgresOps.Port,
		postgresOps.User,
		postgresOps.Database,
		postgresOps.Password)

	db, err := gorm.Open(postgres.Open(connectionString))
	if err != nil {
		return nil, err
	}
	return db, nil
}

type SQLiteConnectOpts struct {
	Path string
}

func SetSQLiteConnectOpts(ops *SQLiteConnectOpts) {
	sqliteOps = ops
}

func setupSQLite() (*gorm.DB, error) {
	var dbLocation string
	if sqliteOps == nil || sqliteOps.Path == "" {
		dbLocation = os.Getenv("EGORM_DB_SQLITE_PATH")
		if dbLocation == "" {
			return nil, fmt.Errorf("egorm: sqliteOpts.Path not specified and DB_SQLITE_PATH env is empty")
		}
	} else {
		dbLocation = sqliteOps.Path
	}

	// Create the sqlite file if it's not available
	if _, err := os.Stat(dbLocation); err != nil {
		if _, err = os.Create(dbLocation); err != nil {
			return nil, fmt.Errorf("egorm: Failed to create sqlite file %s: %s", dbLocation, err)
		}
	}

	db, err := gorm.Open(sqlite.Open(dbLocation), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("egorm: Failed connect to sqlite file %s: %s", dbLocation, err)
	}
	return db, nil
}

/*func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&DbTarget{})
}*/

func initializeDatabaseLayer() {

	dbs := os.Getenv("EGORM_DB")
	var db *gorm.DB

	if dbs == "" || sqliteOps != nil {
		dbs = "sqlite"
	}

	if postgresOps != nil {
		dbs = "postgres"
	}

	// var err error

	switch dbs {
	case "sqlite":
		db, initOnceErr = setupSQLite()
		break
	case "postgres":
		db, initOnceErr = setupPostgres()
		break
	default:
		initOnceErr = fmt.Errorf("No database found, set the DB env")
		return
	}

	if initOnceErr != nil {
		return
	}

	/*initOnceErr = autoMigrate(db)
	if initOnceErr != nil {
		return
	}*/
	Db = db

}

func InitDB() error {
	initOnce.Do(initializeDatabaseLayer)
	return initOnceErr
}
