package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/silvasur/mailremind/confhelper"
	"github.com/silvasur/mailremind/model"
	"github.com/silvasur/mailremind/schedule"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type scheduleTpldata struct {
	Start, End                                                               string
	Count                                                                    int
	UnitIsMinute, UnitIsHour, UnitIsDay, UnitIsWeek, UnitIsMonth, UnitIsYear bool
	RepetitionEnabled, EndEnabled                                            bool
}

func schedToSchedTL(sched schedule.Schedule, u model.User) scheduleTpldata {
	loc := u.Location()

	schedtl := scheduleTpldata{
		Start: sched.Start.In(loc).Format(bestTimeFmtEver),
	}

	if f := sched.Freq; f.Count > 0 {
		schedtl.RepetitionEnabled = true
		schedtl.Count = int(f.Count)
		switch f.Unit {
		case schedule.Minute:
			schedtl.UnitIsMinute = true
		case schedule.Hour:
			schedtl.UnitIsHour = true
		case schedule.Day:
			schedtl.UnitIsDay = true
		case schedule.Week:
			schedtl.UnitIsWeek = true
		case schedule.Month:
			schedtl.UnitIsMonth = true
		case schedule.Year:
			schedtl.UnitIsYear = true
		}
	}

	if end := sched.End; !end.IsZero() {
		schedtl.EndEnabled = true
		schedtl.End = end.In(loc).Format(bestTimeFmtEver)
	}

	return schedtl
}

var maxSchedules, jobsLimit int

func initLimits() {
	maxSchedules = int(confhelper.ConfIntOrFatal(conf, "limits", "schedules"))
	jobsLimit = int(confhelper.ConfIntOrFatal(conf, "limits", "jobs"))
}

const bestTimeFmtEver = "2006-01-02 15:04:05"

type jobeditTpldata struct {
	Error, Success          string
	Fatal                   bool
	JobID, Subject, Content string
	Schedules               []scheduleTpldata
}

func (jt *jobeditTpldata) fillFromJob(job model.Job, u model.User) {
	jt.JobID = job.ID().String()
	jt.Subject = job.Subject()
	jt.Content = string(job.Content())
	jt.Schedules = make([]scheduleTpldata, maxSchedules)

	for i, sched := range job.Schedule() {
		if i == maxSchedules {
			log.Printf("Job %s has more than %d schedule entries!", job.ID(), maxSchedules)
			break
		}

		jt.Schedules[i] = schedToSchedTL(sched, u)
	}
}

func (jt *jobeditTpldata) interpretForm(form url.Values, u model.User) (subject string, content []byte, ms schedule.MultiSchedule, ok bool) {
	loc := u.Location()

	l1 := len(form["Start"])
	l2 := len(form["RepetitionEnabled"])
	l3 := len(form["Count"])
	l4 := len(form["Unit"])
	l5 := len(form["EndEnabled"])
	l6 := len(form["End"])

	if (l1 != l2) || (l2 != l3) || (l3 != l4) || (l4 != l5) || (l5 != l6) {
		jt.Error = "Form corrupted. Changes were not saved"
		ok = false
		return
	}

	subject = form.Get("Subject")
	content = []byte(form.Get("Content"))
	jt.Subject = subject
	jt.Content = string(content)

	ok = true
	for i, _start := range form["Start"] {
		if i >= maxSchedules {
			break
		}

		if _start == "" {
			continue
		}

		start, err := time.ParseInLocation(bestTimeFmtEver, _start, loc)
		if err != nil {
			ok = false
			continue
		}

		count := uint64(0)
		var unit schedule.TimeUnit
		var end time.Time
		if form["RepetitionEnabled"][i] == "yes" {
			if count, err = strconv.ParseUint(form["Count"][i], 10, 64); err != nil {
				ok = false
				continue
			}

			switch form["Unit"][i] {
			case "Minute":
				unit = schedule.Minute
			case "Hour":
				unit = schedule.Hour
			case "Day":
				unit = schedule.Day
			case "Week":
				unit = schedule.Week
			case "Month":
				unit = schedule.Month
			case "Year":
				unit = schedule.Year
			default:
				ok = false
				continue
			}

			if form["EndEnabled"][i] == "yes" {
				if end, err = time.ParseInLocation(bestTimeFmtEver, form["End"][i], loc); err != nil {
					ok = false
					continue
				}
			}
		}

		sched := schedule.Schedule{
			Start: start,
			Freq: schedule.Frequency{
				Count: uint(count),
				Unit:  unit,
			},
			End: end,
		}
		ms = append(ms, sched)
		jt.Schedules[i] = schedToSchedTL(sched, u)
	}

	if !ok {
		jt.Error = "Some schedules were wrong (wrong time format, negative repetition counts)"
		return
	}

	if len(ms) == 0 {
		jt.Error = "No schedule."
		ok = false
	}

	return
}

func logfail(what string, err error) bool {
	if err != nil {
		log.Printf("Failed %s: %s", what, err)
		return false
	}
	return true
}

func jobedit(user model.User, sess *sessions.Session, req *http.Request) (interface{}, model.User) {
	if user == nil {
		return &jobeditTpldata{Error: "You need to be logged in to do that.", Fatal: true}, user
	}

	outdata := &jobeditTpldata{Schedules: make([]scheduleTpldata, maxSchedules)}

	// Try to load job, if given
	_id := mux.Vars(req)["ID"]
	var job model.Job
	if _id != "" {
		id, err := db.ParseDBID(_id)
		if err != nil {
			return &jobeditTpldata{Error: "Job not found", Fatal: true}, user
		}

		if job, err = user.JobByID(id); err != nil {
			return &jobeditTpldata{Error: "Job not found", Fatal: true}, user
		}
	}

	if job != nil {
		outdata.fillFromJob(job, user)
	}

	if req.Method == "POST" {
		if (job == nil) && (jobsLimit >= 0) && (user.CountJobs() >= jobsLimit) {
			outdata.Error = "You have reached the limit of jobs per user."
			outdata.Fatal = true
			return outdata, user
		}

		if err := req.ParseForm(); err != nil {
			outdata.Error = "Could not understand forma data."
			return outdata, user
		}

		subject, content, mc, ok := outdata.interpretForm(req.Form, user)
		if ok {
			next := mc.NextAfter(time.Now())
			if next.IsZero() {
				outdata.Error = "The schedule would not send any mail."
			} else if job != nil {
				if logfail("setting subject", job.SetSubject(subject)) &&
					logfail("setting content", job.SetContent(content)) &&
					logfail("setting schedule", job.SetSchedule(mc)) &&
					logfail("setting next", job.SetNext(next)) {
					outdata.Success = "Changes saved"
				} else {
					outdata.Error = "Could not save everything."
				}
			} else {
				if job, err := user.AddJob(subject, content, mc, next); logfail("creating new job", err) {
					outdata.fillFromJob(job, user)
					outdata.Success = "Job created"
				} else {
					outdata.Error = "Failed creating the job."
				}
			}
		}
	}

	return outdata, user
}
