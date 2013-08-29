package main

import (
	"kch42.de/gostuff/mailremind/model"
	"log"
)

var db model.DBInfo
var dbcon model.DBCon

func initDB() {
	dbdrv, err := conf.GetString("db", "driver")
	if err != nil {
		log.Fatalf("Could not get db.driver from config: %s", err)
	}

	dbconf, err := conf.GetString("db", "conf")
	if err != nil {
		log.Fatalf("Could not get db.conf from config: %s", err)
	}

	var ok bool
	if db, ok = model.GetDBInfo(dbdrv); !ok {
		log.Fatalf("Could not get info for dbdrv %s: %s", dbdrv, err)
	}

	if dbcon, err = db.Connect(dbconf); err != nil {
		log.Fatalf("Unable to connect to %s database: %s", dbdrv, err)
	}
}
