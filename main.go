package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/kch42/simpleconf"
	_ "kch42.de/gostuff/mailremind/model/mysql"
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
	_auth, err := conf.GetString("securecookies", "auth")
	if err != nil {
		log.Fatalf("Could not get securecookies.auth from config: %s", err)
	}
	auth, err := hex.DecodeString(_auth)
	if err != nil {
		log.Fatalf("Could not decode securecookies.auth as hex: %s", err)
	}

	_crypt, err := conf.GetString("securecookies", "crypt")
	if err != nil {
		log.Fatalf("Could not get securecookies.crypt from config: %s", err)
	}
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

	if baseurl, err = conf.GetString("web", "baseurl"); err != nil {
		log.Fatalf("Could not get web.baseurl from config: %s", err)
	}

	rand.Seed(time.Now().UnixNano())

	initSessions()
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
	router.HandleFunc("/register", mkHttpHandler(register, tplRegister))
	router.HandleFunc("/activate", mkHttpHandler(activate, tplMsg))
	router.HandleFunc("/login", mkHttpHandler(login, tplLogin))
	router.HandleFunc("/logincheck", mkHttpHandler(logincheck, tplMsg))
	router.HandleFunc("/logout", mkHttpHandler(logout, tplMsg))
	router.HandleFunc("/delete-acc/yes", mkHttpHandler(deleteacc, tplMsg))
	router.HandleFunc("/delete-acc", mkHttpHandler(deleteask, tplReallyDelete))
	router.HandleFunc("/pwreset", mkHttpHandler(pwreset, tplPwreset))
	router.HandleFunc("/forgotpw", mkHttpHandler(forgotpw, tplForgotpw))

	http.Handle("/", router)

	if err := http.ListenAndServe(laddr, nil); err != nil {
		log.Fatalf("Could not ListenAndServe: %s", err)
	}
}
