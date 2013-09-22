package main

import (
	"github.com/kch42/mailremind/confhelper"
	"github.com/kch42/mailremind/mailing"
	"log"
)

var MailFrom string

type email struct {
	To, From string
	Msg      []byte
	OK       chan<- bool
}

var mailchan chan *email

func Mail(to, from string, msg []byte) bool {
	ok := make(chan bool)
	mailchan <- &email{to, from, msg, ok}
	return <-ok
}

func initMailing() {
	meth := confhelper.ConfStringOrFatal(conf, "mail", "method")
	MailFrom = confhelper.ConfStringOrFatal(conf, "mail", "from")

	parallel := confhelper.ConfIntOrFatal(conf, "mail", "parallel")
	if parallel <= 0 {
		log.Fatalln("mail.parallel must be > 0")
	}

	mailchan = make(chan *email)

	mc, ok := mailing.MailersByName[meth]
	if !ok {
		log.Fatalf("Unknown mail method: %s", meth)
	}

	for i := int64(0); i < parallel; i++ {
		mailer, err := mc(conf)
		if err != nil {
			log.Fatalf("Error while initializing mail: %s", err)
		}

		go func(mailer mailing.Mailer) {
			for {
				mail := <-mailchan
				if err := mailer.Mail(mail.To, mail.From, mail.Msg); err != nil {
					log.Printf("Could not send mail to \"%s\": %s", mail.To, err)
					mail.OK <- false
				} else {
					mail.OK <- true
				}
			}
		}(mailer)
	}
}
