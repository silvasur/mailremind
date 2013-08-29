package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kch42/simpleconf"
	_ "kch42.de/gostuff/mailremind/model/mysql"
	"log"
	"net/http"
)

func debug(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "Content-Type: text/plain\r\n\r\n%#v", req)
}

var conf simpleconf.Config
var baseurl string

func main() {
	confpath := flag.String("config", "", "Path to config file")
	flag.Parse()

	var err error
	if conf, err = simpleconf.LoadByFilename(*confpath); err != nil {
		log.Fatalf("Could not read config: %s", err)
	}

	if baseurl, err = conf.GetString("web", "baseurl"); err != nil {
		log.Fatalf("Could not get web.baseurl from config: %s", err)
	}

	initTpls()
	loadTimeLocs()
	initMailing()
	initMails()
	initDB()
	defer dbcon.Close()

	staticpath, err := conf.GetString("paths", "static")
	if err != nil {
		log.Fatalf("Could not get paths.static config: %s", err)
	}

	laddr, err := conf.GetString("net", "laddr")
	if err != nil {
		log.Fatalf("Could not get net.laddr config: %s", err)
	}

	router := mux.NewRouter()
	router.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticpath))))
	router.HandleFunc("/register", register)
	router.HandleFunc("/activate", activate)

	http.Handle("/", router)

	if err := http.ListenAndServe(laddr, nil); err != nil {
		log.Fatalf("Could not ListenAndServe: %s", err)
	}
}
