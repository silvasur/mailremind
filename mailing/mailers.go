package mailing

var MailersByName = map[string]MailerCreator{
	"smtp":     SMTPMailerCreator,
	"sendmail": SendmailMailerCreator}
