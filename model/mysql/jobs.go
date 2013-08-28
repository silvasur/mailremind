package mysql

import (
	"database/sql"
	"fmt"
	"kch42.de/gostuff/mailremind/chronos"
	"kch42.de/gostuff/mailremind/model"
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
	chron   []chronos.Chronos
}

func jobFromSQL(con *MySQLDBCon, s scanner) (*Job, error) {
	var _id, _user uint64
	var subject string
	var content []byte
	var _next int64
	var _mchron string

	if err := s.Scan(&_id, &_user, &subject, &content, &_next, &_mchron); err != nil {
		return nil, err
	}

	chron, err := chronos.ParseMultiChronos(_mchron)
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
		chron:   chron,
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

func (u *User) AddJob(subject string, content []byte, chron chronos.MultiChronos, next time.Time) (model.Job, error) {
	tx, err := u.con.con.Begin()
	if err != nil {
		return nil, err
	}

	insjob := tx.Stmt(u.con.stmt[qInsertJob])

	res, err := insjob.Exec(uint64(u.id), subject, content, next.Unix(), chron.String())
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
		chron:   chron,
	}, nil
}

func (j *Job) ID() model.DBID                { return j.id }
func (j *Job) Subject() string               { return j.subject }
func (j *Job) Content() []byte               { return j.content }
func (j *Job) Chronos() chronos.MultiChronos { return j.chron }
func (j *Job) Next() time.Time               { return j.next }

func (j *Job) User() model.User {
	u, err := j.con.UserByID(j.user)
	if err != nil {
		// TODO: Should we really panic here? If yes, we need to recover panics!
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

func (j *Job) SetChronos(chron chronos.MultiChronos) error {
	if _, err := j.con.stmt[qSetChronos].Exec(chron.String(), uint64(j.id)); err != nil {
		return err
	}

	j.chron = chron
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

func (con *MySQLDBCon) JobsBefore(t time.Time) []model.DBID {
	rows, err := con.stmt[qJobsBefore].Query(t.Unix())
	if err != nil {
		log.Fatalf("Could not get jobs before %s: %s", t, err) // TODO: Really fatal?
	}

	ids := make([]model.DBID, 0)
	for rows.Next() {
		var _id uint64
		if err := rows.Scan(&_id); err != nil {
			log.Printf("Could not get all jobs before %s: %s", t, err)
			break
		}
		ids = append(ids, DBID(_id))
	}

	return ids
}
