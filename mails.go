package main

import (
	"bytes"
	"fmt"
	"github.com/kch42/mailremind/confhelper"
	"github.com/kch42/mailremind/model"
	"log"
	"path"
	"text/template"
	"time"
)

func loadMailTpl(tplroot, name string) *template.Template {
	tpl, err := template.ParseFiles(path.Join(tplroot, name+".tpl"))
	if err != nil {
		log.Fatalf("Could not load mailtemplate %s: %s", name, err)
	}
	return tpl
}

var (
	mailActivationcode *template.Template
	mailPwreset        *template.Template
)

func initMails() {
	tplroot := confhelper.ConfStringOrFatal(conf, "paths", "mailtpls")

	mailActivationcode = loadMailTpl(tplroot, "activationcode")
	mailPwreset = loadMailTpl(tplroot, "pwreset")
}

type activationcodeData struct {
	URL string
}

func SendActivationcode(to, acCode string, uid model.DBID) bool {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "To: %s\n", to)
	fmt.Fprintf(buf, "From: %s\n", MailFrom)
	fmt.Fprintf(buf, "Subject: Activation code for your mailremind account\n")
	fmt.Fprintf(buf, "Date: %s\n", time.Now().Format(time.RFC1123Z))

	fmt.Fprintln(buf, "")

	url := fmt.Sprintf("%s/activate?U=%s&Code=%s", baseurl, uid, acCode)
	if err := mailActivationcode.Execute(buf, activationcodeData{url}); err != nil {
		log.Printf("Error while executing mail template (activationcode): %s", err)
		return false
	}

	return Mail(to, MailFrom, buf.Bytes())
}

func SendPwresetLink(to, code string, uid model.DBID) bool {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "To: %s\n", to)
	fmt.Fprintf(buf, "From: %s\n", MailFrom)
	fmt.Fprintf(buf, "Subject: Password reset request for your mailremind account\n")
	fmt.Fprintf(buf, "Date: %s\n", time.Now().Format(time.RFC1123Z))

	fmt.Fprintln(buf, "")

	url := fmt.Sprintf("%s/pwreset?U=%s&Code=%s", baseurl, uid, code)
	if err := mailPwreset.Execute(buf, activationcodeData{url}); err != nil {
		log.Printf("Error while executing mail template (pwreset): %s", err)
		return false
	}

	return Mail(to, MailFrom, buf.Bytes())
}
