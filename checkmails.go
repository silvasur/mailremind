package main

import (
	"bytes"
	"fmt"
	"kch42.de/gostuff/mailremind/model"
	"log"
	"time"
)

const checkInterval = 30 // TODO: Make this configurable

func checkmails() {
	ticker := time.NewTicker(checkInterval * time.Second)

	for {
		t := <-ticker.C

		jobs := dbcon.JobsBefore(t)

		for _, job := range jobs {
			if sendjob(job, t) {
				next := job.Chronos().NextAfter(t)
				if next.IsZero() {
					if err := job.Delete(); err != nil {
						log.Printf("Failed deleting job %s after job was done: %s", job.ID(), err)
					}
				} else {
					if err := job.SetNext(next); err != nil {
						log.Printf("Filed setting next for job %s: %s", job.ID(), err)
					}
				}
			}
		}
	}
}

func sendjob(job model.Job, t time.Time) bool {
	user := job.User()
	buf := new(bytes.Buffer)

	fmt.Fprintf(buf, "From: %s\n", MailFrom)
	fmt.Fprintf(buf, "To: %s\n", user.Email())
	fmt.Fprintf(buf, "Subject: %s\n", job.Subject())
	fmt.Fprintf(buf, "Date: %s\n", t.In(user.Location()).Format(time.RFC1123Z))

	fmt.Fprintln(buf, "")

	buf.Write(job.Content())

	return Mail(user.Email(), MailFrom, buf.Bytes())
}
