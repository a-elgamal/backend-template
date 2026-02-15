// Package test contains utilities to support testing activities
package test

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"alielgamal.com/myservice/internal/config"
	internalDB "alielgamal.com/myservice/internal/db"
)

// SetupTestDB Returns a function that should be called with deferred by the test. If targetDBVersion is not set, the DB is upgraded to latest version
func SetupTestDB(dbName string, targetDBVersion uint, appConfig config.Config) (*internalDB.SQLDB, func(), error) {
	var connector driver.Connector
	connector, err := pq.NewConnector(appConfig.DBConfig.GetURL())
	if err != nil {
		return nil, nil, err
	}

	db := sql.OpenDB(connector)
	defer db.Close()

	dbName = strings.Replace(strings.ToLower(dbName), "/", "_", -1)
	dbName = strings.Replace(dbName, "'", "_", -1)

	if _, err = db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %v", dbName)); err != nil {
		return nil, nil, err
	}

	if _, err = db.Exec(fmt.Sprintf("CREATE DATABASE %v", dbName)); err != nil {
		return nil, nil, err
	}

	tConnector, err := pq.NewConnector(fmt.Sprintf(appConfig.DBConfig.GetURLTemplate(), dbName))
	if err != nil {
		return nil, nil, err
	}
	tDB := sql.OpenDB(tConnector)

	if targetDBVersion == 0 {
		_, err = internalDB.UpgradeDB(tDB)
	} else {
		_, err = internalDB.MigrateDBTo(tDB, targetDBVersion)
	}
	if err != nil {
		return nil, nil, err
	}

	teardownF := func() { tDB.Close() }

	return &internalDB.SQLDB{DB: tDB}, teardownF, err
}
