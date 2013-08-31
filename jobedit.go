package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"kch42.de/gostuff/mailremind/chronos"
	"kch42.de/gostuff/mailremind/model"
	"log"
	"net/http"
)

type scheduleTpldata struct {
	Start, End                                                               string
	Count                                                                    int
	UnitIsMinute, UnitIsHour, UnitIsDay, UnitIsWeek, UnitIsMonth, UnitIsYear bool
	RepetitionEnabled, EndEnabled                                            bool
}

const maxSchedules = 10
const bestTimeFmtEver = "2006-01-02 15:04:05"

type jobeditTpldata struct {
	Error, Success          string
	Fatal                   bool
	JobID, Subject, Content string
	Schedules               [maxSchedules]scheduleTpldata
}

func (jt *jobeditTpldata) fillFromJob(job model.Job, u model.User) {
	jt.JobID = job.ID().String()
	jt.Subject = job.Subject()
	jt.Content = string(job.Content())

	loc := u.Location()

	for i, chron := range job.Chronos() {
		if i == 10 {
			log.Printf("Job %s has more than %d Chronos entries!", job.ID(), maxSchedules)
			break
		}

		schedule := scheduleTpldata{
			Start: chron.Start.In(loc).Format(bestTimeFmtEver),
		}

		if f := chron.Freq; f.Count > 0 {
			schedule.RepetitionEnabled = true
			schedule.Count = int(f.Count)
			switch f.Unit {
			case chronos.Minute:
				schedule.UnitIsMinute = true
			case chronos.Hour:
				schedule.UnitIsHour = true
			case chronos.Day:
				schedule.UnitIsDay = true
			case chronos.Week:
				schedule.UnitIsWeek = true
			case chronos.Month:
				schedule.UnitIsMonth = true
			case chronos.Year:
				schedule.UnitIsYear = true
			}
		}

		if end := chron.End; !end.IsZero() {
			schedule.EndEnabled = true
			schedule.End = end.In(loc).Format(bestTimeFmtEver)
		}

		jt.Schedules[i] = schedule
	}
}

func jobedit(user model.User, sess *sessions.Session, req *http.Request) interface{} {
	if user == nil {
		return &jobeditTpldata{Error: "You need to be logged in to do that.", Fatal: true}
	}

	outdata := new(jobeditTpldata)

	// Try to load job, if given
	_id := mux.Vars(req)["ID"]
	var job model.Job
	if _id != "" {
		id, err := db.ParseDBID(_id)
		if err != nil {
			return &jobeditTpldata{Error: "Job not found", Fatal: true}
		}

		if job, err = user.JobByID(id); err != nil {
			return &jobeditTpldata{Error: "Job not found", Fatal: true}
		}
	}

	if req.Method == "POST" {
		// TODO: Enable editing...
	}

	if job != nil {
		outdata.fillFromJob(job, user)
	}
	return outdata
}
