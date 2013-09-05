package main

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"kch42.de/gostuff/mailremind/chronos"
	"kch42.de/gostuff/mailremind/model"
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

func chronToSchedTL(chron chronos.Chronos, u model.User) scheduleTpldata {
	loc := u.Location()

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

	return schedule
}

// TODO: Make these constants variable (config file or something...)
const (
	maxSchedules = 10
	jobsLimit    = 100
)

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

	for i, chron := range job.Chronos() {
		if i == 10 {
			log.Printf("Job %s has more than %d Chronos entries!", job.ID(), maxSchedules)
			break
		}

		jt.Schedules[i] = chronToSchedTL(chron, u)
	}
}

func (jt *jobeditTpldata) interpretForm(form url.Values, u model.User) (subject string, content []byte, mc chronos.MultiChronos, ok bool) {
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
		var unit chronos.TimeUnit
		var end time.Time
		if form["RepetitionEnabled"][i] == "yes" {
			if count, err = strconv.ParseUint(form["Count"][i], 10, 64); err != nil {
				ok = false
				continue
			}

			switch form["Unit"][i] {
			case "Minute":
				unit = chronos.Minute
			case "Hour":
				unit = chronos.Hour
			case "Day":
				unit = chronos.Day
			case "Week":
				unit = chronos.Week
			case "Month":
				unit = chronos.Month
			case "Year":
				unit = chronos.Year
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

		chron := chronos.Chronos{
			Start: start,
			Freq: chronos.Frequency{
				Count: uint(count),
				Unit:  unit,
			},
			End: end,
		}
		mc = append(mc, chron)
		jt.Schedules[i] = chronToSchedTL(chron, u)
	}

	if !ok {
		jt.Error = "Some schedules were wrong (wrong time format, negative repetition counts)"
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

	outdata := new(jobeditTpldata)

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
		if (job == nil) && (user.CountJobs() >= jobsLimit) {
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
			if job != nil {
				if logfail("setting subject", job.SetSubject(subject)) &&
					logfail("setting content", job.SetContent(content)) &&
					logfail("setting chronos", job.SetChronos(mc)) &&
					logfail("setting next", job.SetNext(next)) {
					outdata.Success = "Changes saved"
				} else {
					outdata.Error = "Could not save everything."
				}
			} else {
				var err error
				if job, err = user.AddJob(subject, content, mc, next); logfail("creating new job", err) {
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
