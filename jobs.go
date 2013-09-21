package main

import (
	"github.com/gorilla/sessions"
	"github.com/kch42/mailremind/model"
	"net/http"
)

type jobTpldata struct {
	ID, Subject, Excerpt, Next string
}

func jobToTpldata(job model.Job, user model.User) *jobTpldata {
	excerpt := string(job.Content())
	if len(excerpt) > 100 {
		excerpt = string([]rune(excerpt)[0:100]) + " (...)"
	}

	return &jobTpldata{
		ID:      job.ID().String(),
		Subject: job.Subject(),
		Excerpt: excerpt,
		Next:    job.Next().In(user.Location()).Format("2006-Jan-02 15:04:05"),
	}
}

type jobsTpldata struct {
	Error, Success string
	Jobs           []*jobTpldata
	Fatal          bool
}

func jobs(user model.User, sess *sessions.Session, req *http.Request) (interface{}, model.User) {
	if user == nil {
		return &jobsTpldata{Error: "You need to be logged in to do that.", Fatal: true}, user
	}

	outdata := new(jobsTpldata)

	if req.Method == "POST" {
		if err := req.ParseForm(); err != nil {
			outdata.Error = "Could not understand form data."
			goto listjobs
		}

		if req.FormValue("Delconfirm") != "yes" {
			goto listjobs
		}

		for _, _id := range req.Form["Jobs"] {
			id, err := db.ParseDBID(_id)
			if err != nil {
				outdata.Error = "Not all jobs could be deleted."
				continue
			}

			job, err := user.JobByID(id)
			if err != nil {
				outdata.Error = "Not all jobs could be deleted."
				continue
			}

			if job.Delete() != nil {
				outdata.Error = "Not all jobs could be deleted."
				continue
			}

			outdata.Success = "Jobs deleted."
		}
	}

listjobs:
	jobs := user.Jobs()
	outdata.Jobs = make([]*jobTpldata, len(jobs))

	for i, job := range jobs {
		outdata.Jobs[i] = jobToTpldata(job, user)
	}

	return outdata, user
}
