package main

import (
	"github.com/silvasur/mailremind/confhelper"
	"github.com/silvasur/mailremind/model"
	"log"
)

var db model.DBInfo
var dbcon model.DBCon

func initDB() {
	dbdrv := confhelper.ConfStringOrFatal(conf, "db", "driver")
	dbconf := confhelper.ConfStringOrFatal(conf, "db", "conf")

	var ok bool
	if db, ok = model.GetDBInfo(dbdrv); !ok {
		log.Fatalf("Could not get info for dbdrv %s", dbdrv)
	}

	var err error
	if dbcon, err = db.Connect(dbconf); err != nil {
		log.Fatalf("Unable to connect to %s database: %s", dbdrv, err)
	}
}
