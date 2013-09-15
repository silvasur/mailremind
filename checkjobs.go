package main

import (
	"bytes"
	"fmt"
	"kch42.de/gostuff/mailremind/confhelper"
	"kch42.de/gostuff/mailremind/model"
	"log"
	"time"
)

var checkInterval int64

func initCheckjobs() {
	checkInterval = confhelper.ConfIntOrFatal(conf, "schedules", "checkInterval")
}

func checkjobs() {
	timech := make(chan time.Time)
	go func(ch chan time.Time) {
		ticker := time.NewTicker(time.Duration(checkInterval) * time.Second)

		ch <- time.Now()
		for t := range ticker.C {
			ch <- t
		}
	}(timech)

	for t := range timech {
		checkjobsOnce(t)
	}
}

func checkjobsOnce(t time.Time) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("!! recovered from panic in checkjobsOnce: %s", r)
		}
	}()

	jobs := dbcon.JobsBefore(t)

	for _, job := range jobs {
		if sendjob(job, t) {
			next := job.Schedule().NextAfter(t)
			if next.IsZero() {
				if err := job.Delete(); err != nil {
					log.Printf("Failed deleting job %s after job was done: %s", job.ID(), err)
				}
			} else {
				if err := job.SetNext(next); err != nil {
					log.Printf("Failed setting next for job %s: %s", job.ID(), err)
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
