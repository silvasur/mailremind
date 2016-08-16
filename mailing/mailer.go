package mailing

import (
	"github.com/silvasur/simpleconf"
)

// Mailer is a interface that defines the Mail function.
type Mailer interface {
	Mail(to, from string, msg []byte) error
}

// MailerCreator is a function that creates a Mailer instance from config values.
type MailerCreator func(conf simpleconf.Config) (Mailer, error)
