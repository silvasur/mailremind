package main

import (
	"html/template"
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

var tplRegister *template.Template

func initTpls() {
	tplpath, err := conf.GetString("paths", "tpls")
	if err != nil {
		log.Fatalf("Could not get paths.tpls config: %s", err)
	}

	tplRegister = loadTpl(tplpath, "register")
}
