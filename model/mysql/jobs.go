package mysql

import (
	"database/sql"
	"fmt"
	"github.com/kch42/mailremind/model"
	"github.com/kch42/mailremind/schedule"
	"log"
	"time"
)

type Job struct {
	con *MySQLDBCon

	id      DBID
	user    DBID
	subject string
	content []byte
	next    time.Time
	sched   schedule.MultiSchedule
}

func jobFromSQL(con *MySQLDBCon, s scanner) (*Job, error) {
	var _id, _user uint64
	var subject string
	var content []byte
	var _next int64
	var _msched string

	if err := s.Scan(&_id, &_user, &subject, &content, &_next, &_msched); err != nil {
		return nil, err
	}

	sched, err := schedule.ParseMultiSchedule(_msched)
	if err != nil {
		return nil, err
	}

	return &Job{
		con:     con,
		id:      DBID(_id),
		user:    DBID(_user),
		subject: subject,
		content: content,
		next:    time.Unix(_next, 0),
		sched:   sched,
	}, nil
}

func (u *User) CountJobs() (c int) {
	row := u.con.stmt[qCountJobs].QueryRow(uint64(u.id))
	if err := row.Scan(&c); err != nil {
		log.Printf("Failed counting user's (%d) jobs: %s", u.id, err)
		c = 0
	}
	return
}

func (u *User) Jobs() []model.Job {
	rows, err := u.con.stmt[qJobsOfUser].Query(uint64(u.id))
	if err != nil {
		log.Printf("Failed getting jobs of user %d: %s", u.id, err)
		return nil
	}

	jobs := make([]model.Job, 0)
	for rows.Next() {
		job, err := jobFromSQL(u.con, rows)
		if err != nil {
			log.Printf("Failed getting all jobs of user %d: %s", u.id, err)
			break
		}
		jobs = append(jobs, job)
	}

	return jobs
}

func (u *User) JobByID(_id model.DBID) (model.Job, error) {
	id := _id.(DBID)

	row := u.con.stmt[qJobFromUserAndID].QueryRow(uint64(u.id), uint64(id))
	switch job, err := jobFromSQL(u.con, row); err {
	case nil:
		return job, nil
	case sql.ErrNoRows:
		return nil, model.NotFound
	default:
		return nil, err
	}
}

func (u *User) AddJob(subject string, content []byte, sched schedule.MultiSchedule, next time.Time) (model.Job, error) {
	tx, err := u.con.con.Begin()
	if err != nil {
		return nil, err
	}

	insjob := tx.Stmt(u.con.stmt[qInsertJob])

	res, err := insjob.Exec(uint64(u.id), subject, content, next.Unix(), sched.String())
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_id, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &Job{
		con:     u.con,
		id:      DBID(_id),
		user:    u.id,
		subject: subject,
		content: content,
		next:    next,
		sched:   sched,
	}, nil
}

func (j *Job) ID() model.DBID                   { return j.id }
func (j *Job) Subject() string                  { return j.subject }
func (j *Job) Content() []byte                  { return j.content }
func (j *Job) Schedule() schedule.MultiSchedule { return j.sched }
func (j *Job) Next() time.Time                  { return j.next }

func (j *Job) User() model.User {
	u, err := j.con.UserByID(j.user)
	if err != nil {
		// We panic here, since the user must exist, if the job is there.
		// Since http handlers and the job handler do recover from panics, this should be okay.
		panic(fmt.Errorf("Could not get user (%d) of Job %d: %s", j.user, j.id, err))
	}

	return u
}

func (j *Job) SetSubject(sub string) error {
	if _, err := j.con.stmt[qSetSubject].Exec(sub, uint64(j.id)); err != nil {
		return err
	}

	j.subject = sub
	return nil
}

func (j *Job) SetContent(cont []byte) error {
	if _, err := j.con.stmt[qSetContent].Exec(cont, uint64(j.id)); err != nil {
		return err
	}

	j.content = cont
	return nil
}

func (j *Job) SetSchedule(sched schedule.MultiSchedule) error {
	if _, err := j.con.stmt[qSetSchedule].Exec(sched.String(), uint64(j.id)); err != nil {
		return err
	}

	j.sched = sched
	return nil
}

func (j *Job) SetNext(next time.Time) error {
	if _, err := j.con.stmt[qSetNext].Exec(next.Unix(), uint64(j.id)); err != nil {
		return err
	}

	j.next = next
	return nil
}

func (j *Job) Delete() error {
	_, err := j.con.stmt[qDelJob].Exec(j.id)
	return err
}

func (con *MySQLDBCon) JobsBefore(t time.Time) (jobs []model.Job) {
	rows, err := con.stmt[qJobsBefore].Query(t.Unix())
	if err != nil {
		log.Fatalf("Could not get jobs before %s: %s", t, err)
	}

	for rows.Next() {
		job, err := jobFromSQL(con, rows)
		if err != nil {
			log.Fatalf("Could not get all jobs before %s: %s", t, err)
			break
		}
		jobs = append(jobs, job)
	}

	return
}
