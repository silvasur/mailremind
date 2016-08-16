package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/silvasur/mailremind/model"
	"strconv"
)

type scanner interface {
	Scan(dest ...interface{}) error
}

type DBID uint64

func (id DBID) String() string {
	return strconv.FormatUint(uint64(id), 16)
}

func parseDBID(s string) (model.DBID, error) {
	_id, err := strconv.ParseUint(s, 16, 64)
	return DBID(_id), err
}

type MySQLDBCon struct {
	con  *sql.DB
	stmt []*sql.Stmt
}

func connect(dbconf string) (model.DBCon, error) {
	con, err := sql.Open("mysql", dbconf)
	if err != nil {
		return nil, err
	}

	dbc := &MySQLDBCon{
		con:  con,
		stmt: make([]*sql.Stmt, qEnd),
	}

	for i := 0; i < qEnd; i++ {
		stmt, err := con.Prepare(queries[i])
		if err != nil {
			con.Close()
			return nil, fmt.Errorf("Failed to prepare statement %d : <%s>: %s", i, queries[i], err)
		}
		dbc.stmt[i] = stmt
	}

	return dbc, nil
}

func init() {
	model.Register("mysql", model.DBInfo{
		Connect:   connect,
		ParseDBID: parseDBID,
	})
}

func (con *MySQLDBCon) Close() {
	con.con.Close()
}

func rollbackAfterFail(err error, tx *sql.Tx) error {
	if rberr := tx.Rollback(); rberr != nil {
		return fmt.Errorf("Rollback error: <%s>, Original error: %s", rberr, err)
	}
	return err
}

func i2b(i int) bool {
	if i == 0 {
		return false
	}
	return true
}

func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}
