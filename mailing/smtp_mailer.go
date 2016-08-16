package mailing

import (
	"errors"
	"github.com/silvasur/simpleconf"
	"net/smtp"
)

// SMTPMailer is a Mailer implementation that sends mail via an smtp server.
type SMTPMailer struct {
	Host       string // Host is the expected hostname of the server
	Addr       string // Addr is the address of the server
	UseCRAMMD5 bool
	Username   string
	Password   string
}

func (sm SMTPMailer) Mail(to, from string, msg []byte) error {
	var auth smtp.Auth
	if sm.UseCRAMMD5 {
		auth = smtp.CRAMMD5Auth(sm.Username, sm.Password)
	} else {
		auth = smtp.PlainAuth("", sm.Username, sm.Password, sm.Host)
	}

	return smtp.SendMail(sm.Addr, auth, from, []string{to}, msg)
}

// SMTPMailerCreator creates an SMTPMailer using configuration values in the [mail] section.
//
// 	addr    - The address of the smtp server (go notation)
// 	user    - Username
// 	passwd  - Password
// 	crammd5 - Should CRAMMD5 (on) or PLAIN (off) be used?
// 	host    - The expected hostname of the mailserver (can be left out, if crammd5 is on)
//
func SMTPMailerCreator(conf simpleconf.Config) (Mailer, error) {
	rv := SMTPMailer{}
	var err error

	if rv.Addr, err = conf.GetString("mail", "addr"); err != nil {
		return rv, errors.New("Missing [mail] addr config")
	}
	if rv.UseCRAMMD5, err = conf.GetBool("mail", "crammd5"); err != nil {
		return rv, errors.New("Missing [mail] crammd5 config")
	}
	if rv.Username, err = conf.GetString("mail", "user"); err != nil {
		return rv, errors.New("Missing [mail] user config")
	}
	if rv.Password, err = conf.GetString("mail", "passwd"); err != nil {
		return rv, errors.New("Missing [mail] passwd config")
	}
	if !rv.UseCRAMMD5 {
		if rv.Host, err = conf.GetString("mail", "host"); err != nil {
			return rv, errors.New("Missing [mail] host config")
		}
	}

	return rv, nil
}
