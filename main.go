package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/kch42/mailremind/confhelper"
	_ "github.com/kch42/mailremind/model/mysql"
	"github.com/kch42/simpleconf"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func debug(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "Content-Type: text/plain\r\n\r\n%#v", req)
}

var conf simpleconf.Config
var baseurl string

var SessionStorage sessions.Store

func initSessions() {
	_auth := confhelper.ConfStringOrFatal(conf, "securecookies", "auth")
	auth, err := hex.DecodeString(_auth)
	if err != nil {
		log.Fatalf("Could not decode securecookies.auth as hex: %s", err)
	}

	_crypt := confhelper.ConfStringOrFatal(conf, "securecookies", "crypt")
	crypt, err := hex.DecodeString(_crypt)
	if err != nil {
		log.Fatalf("Could not decode securecookies.auth as hex: %s", err)
	}

	SessionStorage = sessions.NewCookieStore(auth, crypt)
}

func main() {
	confpath := flag.String("config", "", "Path to config file")
	flag.Parse()

	var err error
	if conf, err = simpleconf.LoadByFilename(*confpath); err != nil {
		log.Fatalf("Could not read config: %s", err)
	}

	baseurl = confhelper.ConfStringOrFatal(conf, "web", "baseurl")

	rand.Seed(time.Now().UnixNano())

	initSessions()
	initTpls()
	loadTimeLocs()
	initMailing()
	initMails()
	initDB()
	initLimits()
	defer dbcon.Close()

	staticpath := confhelper.ConfStringOrFatal(conf, "paths", "static")
	laddr := confhelper.ConfStringOrFatal(conf, "net", "laddr")

	initCheckjobs()
	go checkjobs()

	router := mux.NewRouter()
	router.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticpath))))
	router.HandleFunc("/register", mkHttpHandler(register, tplRegister))
	router.HandleFunc("/activate", mkHttpHandler(activate, tplMsg))
	router.HandleFunc("/login", mkHttpHandler(login, tplLogin))
	router.HandleFunc("/logout", mkHttpHandler(logout, tplMsg))
	router.HandleFunc("/delete-acc/yes", mkHttpHandler(deleteacc, tplMsg))
	router.HandleFunc("/delete-acc", mkHttpHandler(deleteask, tplReallyDelete))
	router.HandleFunc("/pwreset", mkHttpHandler(pwreset, tplPwreset))
	router.HandleFunc("/forgotpw", mkHttpHandler(forgotpw, tplForgotpw))
	router.HandleFunc("/jobs", mkHttpHandler(jobs, tplJobs))
	router.HandleFunc("/jobedit", mkHttpHandler(jobedit, tplJobedit))
	router.HandleFunc("/jobedit/{ID}", mkHttpHandler(jobedit, tplJobedit))
	router.HandleFunc("/settings", mkHttpHandler(settings, tplSettings))

	http.Handle("/", router)

	if err := http.ListenAndServe(laddr, nil); err != nil {
		log.Fatalf("Could not ListenAndServe: %s", err)
	}
}
