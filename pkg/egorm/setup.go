package egorm

import (
	"fmt"
	"os"
	"sync"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB
var initOnce sync.Once
var initOnceErr error
var sqliteOps *SQLiteConnectOpts

type SQLiteConnectOpts struct {
	Path string
}

func SetSQLiteConnectOpts(ops *SQLiteConnectOpts) {
	sqliteOps = ops
}

func setupSQLite() (*gorm.DB, error) {
	var dbLocation string
	if sqliteOps == nil || sqliteOps.Path == "" {
		dbLocation = os.Getenv("DB_SQLITE_PATH")
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

	dbs := os.Getenv("DB")
	var db *gorm.DB

	if dbs == "" {
		dbs = "sqlite"
	}

	// var err error

	switch dbs {
	case "sqlite":
		db, initOnceErr = setupSQLite()
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
