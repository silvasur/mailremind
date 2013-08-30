package mysql

import (
	"database/sql"
	"kch42.de/gostuff/mailremind/model"
	"log"
	"time"
)

type User struct {
	con *MySQLDBCon

	id                    DBID
	email, passwd, acCode string
	location              *time.Location
	added                 time.Time
	active                bool
}

func userFromSQL(con *MySQLDBCon, s scanner) (*User, error) {
	var id uint64
	var added int64
	var email, passwd, _loc, acCode string
	var active int

	switch err := s.Scan(&id, &email, &passwd, &_loc, &active, &acCode, &added); err {
	case nil:
	case sql.ErrNoRows:
		return nil, model.NotFound
	default:
		return nil, err
	}

	user := &User{
		con:    con,
		id:     DBID(id),
		email:  email,
		passwd: passwd,
		acCode: acCode,
		added:  time.Unix(added, 0),
		active: i2b(active),
	}

	loc, err := time.LoadLocation(_loc)
	if err != nil {
		loc = time.UTC
	}
	user.location = loc

	return user, nil
}

func (con *MySQLDBCon) UserByID(_id model.DBID) (model.User, error) {
	id := _id.(DBID)

	row := con.stmt[qUserByID].QueryRow(uint64(id))
	return userFromSQL(con, row)
}

func (con *MySQLDBCon) UserByMail(email string) (model.User, error) {
	row := con.stmt[qUserByEmail].QueryRow(email)
	return userFromSQL(con, row)
}

func (u *User) ID() model.DBID         { return u.id }
func (u *User) Email() string          { return u.email }
func (u *User) PWHash() []byte         { return []byte(u.passwd) }
func (u *User) Active() bool           { return u.active }
func (u *User) ActivationCode() string { return u.acCode }

func (u *User) SetPWHash(_pwhash []byte) error {
	pwhash := string(_pwhash)

	if _, err := u.con.stmt[qSetPWHash].Query(pwhash, uint64(u.id)); err != nil {
		return err
	}

	u.passwd = string(_pwhash)
	return nil
}

func (u *User) SetActive(b bool) error {
	if _, err := u.con.stmt[qSetActive].Query(b2i(b), uint64(u.id)); err != nil {
		return err
	}

	u.active = b
	return nil
}

func (u *User) SetActivationCode(c string) error {
	if _, err := u.con.stmt[qSetAcCode].Query(c, uint64(u.id)); err != nil {
		return err
	}

	u.acCode = c
	return nil
}

func (u *User) Delete() error {
	tx, err := u.con.con.Begin()
	if err != nil {
		return err
	}

	id := uint64(u.id)

	deljobs := tx.Stmt(u.con.stmt[qDelUsersJobs])
	deluser := tx.Stmt(u.con.stmt[qDelUser])

	if _, err := deljobs.Query(id); err != nil {
		return rollbackAfterFail(err, tx)
	}

	if _, err := deluser.Query(id); err != nil {
		return rollbackAfterFail(err, tx)
	}

	return tx.Commit()
}

func (con *MySQLDBCon) InactiveUsers(olderthan time.Time) []model.DBID {
	ids := make([]model.DBID, 0)

	rows, err := con.stmt[qGetOldInactiveUsers].Query(olderthan.Unix())
	if err != nil {
		log.Printf("Failed to get old, inactive users: %s", err)
		return ids
	}

	for rows.Next() {
		var _id uint64

		if err := rows.Scan(&_id); err != nil {
			log.Printf("Failed to get old, inactive users: %s", err)
			return ids
		}

		ids = append(ids, DBID(_id))
	}

	return ids
}

func (con *MySQLDBCon) AddUser(email string, pwhash []byte, location *time.Location, active bool, acCode string) (model.User, error) {
	now := time.Now()

	tx, err := con.con.Begin()
	if err != nil {
		return nil, err
	}

	insjob := tx.Stmt(con.stmt[qInsertUser])

	res, err := insjob.Exec(email, string(pwhash), location.String(), b2i(active), acCode, now.Unix())
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

	return &User{
		con:      con,
		id:       DBID(_id),
		email:    email,
		passwd:   string(pwhash),
		acCode:   acCode,
		location: location,
		added:    now,
		active:   active,
	}, nil
}
