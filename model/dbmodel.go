package model

import (
	"errors"
	"fmt"
	"kch42.de/gostuff/mailremind/chronos"
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

	AddJob() Job
	Jobs() []Job
	JobByID(DBID) (Job, error)

	Delete() error
}

type Job interface {
	ID() DBID
	User() User

	Subject() string
	SetSubject(string) error

	Content() []byte
	SetContent([]byte) error

	Chronos() []chronos.Chronos
	SetChronos([]chronos.Chronos) error

	Next() time.Time
	SetNext(time.Time) error

	Delete() error
}

type DBCon interface {
	Close()

	UserByID(DBID) (User, error)
	UserByMail(string) (User, error)

	LastAccess() time.Time
	SetLastAccess(time.Time) error

	JobsBetween(a, b time.Time) ([]Job, error)
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
