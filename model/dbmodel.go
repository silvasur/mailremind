package model

import (
	"errors"
	"fmt"
	"github.com/kch42/mailremind/schedule"
	"sync"
	"time"
)

var (
	NotFound = errors.New("Not Found")
)

type DBID interface {
	fmt.Stringer
}

type User interface {
	ID() DBID
	Email() string

	PWHash() []byte
	SetPWHash([]byte) error

	AddJob(subject string, content []byte, sched schedule.MultiSchedule, next time.Time) (Job, error)
	Jobs() []Job
	JobByID(DBID) (Job, error)
	CountJobs() int

	Active() bool
	SetActive(bool) error

	ActivationCode() string
	SetActivationCode(string) error

	Location() *time.Location
	SetLocation(*time.Location) error

	Delete() error
}

type Job interface {
	ID() DBID
	User() User

	Subject() string
	SetSubject(string) error

	Content() []byte
	SetContent([]byte) error

	Schedule() schedule.MultiSchedule
	SetSchedule(schedule.MultiSchedule) error

	Next() time.Time
	SetNext(time.Time) error

	Delete() error
}

type DBCon interface {
	Close()

	UserByID(DBID) (User, error)
	UserByMail(string) (User, error)

	AddUser(email string, pwhash []byte, location *time.Location, active bool, acCode string) (User, error)

	InactiveUsers(olderthan time.Time) []DBID

	JobsBefore(t time.Time) []Job // Get Jobs with next <= t
}

type DBInfo struct {
	Connect   func(dbconf string) (DBCon, error)
	ParseDBID func(string) (DBID, error)
}

var dbinfos map[string]DBInfo
var dbinfoInit sync.Once

func Register(name string, dbinfo DBInfo) {
	dbinfoInit.Do(func() {
		dbinfos = make(map[string]DBInfo)
	})

	dbinfos[name] = dbinfo
}

func GetDBInfo(name string) (DBInfo, bool) {
	dbinfo, ok := dbinfos[name]
	return dbinfo, ok
}

func AllDatabases() []string {
	names := make([]string, 0, len(dbinfos))
	for name := range dbinfos {
		names = append(names, name)
	}
	return names
}
