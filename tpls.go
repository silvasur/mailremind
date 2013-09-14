package main

import (
	"html/template"
	"kch42.de/gostuff/mailremind/confhelper"
	"log"
	"path"
)

func loadTpl(tplpath, name string) *template.Template {
	tpl, err := template.ParseFiles(
		path.Join(tplpath, "master.tpl"),
		path.Join(tplpath, name+".tpl"))
	if err != nil {
		log.Fatalf("Could not load template \"%s\": %s", name, err)
	}
	return tpl
}

var (
	tplRegister     *template.Template
	tplMsg          *template.Template
	tplLogin        *template.Template
	tplReallyDelete *template.Template
	tplPwreset      *template.Template
	tplForgotpw     *template.Template
	tplJobs         *template.Template
	tplJobedit      *template.Template
	tplSettings     *template.Template
)

func initTpls() {
	tplpath := confhelper.ConfStringOrFatal(conf, "paths", "tpls")

	tplRegister = loadTpl(tplpath, "register")
	tplMsg = loadTpl(tplpath, "msg")
	tplLogin = loadTpl(tplpath, "login")
	tplReallyDelete = loadTpl(tplpath, "reallydelete")
	tplPwreset = loadTpl(tplpath, "pwreset")
	tplForgotpw = loadTpl(tplpath, "forgotpw")
	tplJobs = loadTpl(tplpath, "jobs")
	tplJobedit = loadTpl(tplpath, "jobedit")
	tplSettings = loadTpl(tplpath, "settings")
}

type msgTpldata struct {
	Title, Class, Msg string
}
