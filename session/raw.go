package session

import (
	"database/sql"
	"github.com/go-needle/orm/clause"
	"github.com/go-needle/orm/dialect"
	"github.com/go-needle/orm/log"
	"github.com/go-needle/orm/schema"
	"strings"
)

type Session struct {
	db       *sql.DB
	tx       *sql.Tx
	sql      strings.Builder
	dialect  dialect.Dialect
	clause   clause.Clause
	refTable *schema.Schema
	sqlVars  []any
	isDebug  bool
}

// CommonDB is a minimal function set of db
type CommonDB interface {
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Exec(query string, args ...any) (sql.Result, error)
}

var _ CommonDB = (*sql.DB)(nil)
var _ CommonDB = (*sql.Tx)(nil)

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{
		db:      db,
		dialect: dialect,
	}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.clause.Clear()
	s.sqlVars = nil
}

// DB returns tx if a tx begins. otherwise return *sql.DB
func (s *Session) DB() CommonDB {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

func (s *Session) Raw(sql string, values ...any) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlVars = append(s.sqlVars, values...)
	return s
}

// Exec raw sql with sqlVars
func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	if s.isDebug {
		s.debugSql(s.sql.String(), s.sqlVars...)
	}
	if result, err = s.DB().Exec(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

// QueryRow gets a record from db
func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	if s.isDebug {
		s.debugSql(s.sql.String(), s.sqlVars...)
	}
	return s.DB().QueryRow(s.sql.String(), s.sqlVars...)
}

// QueryRows gets a list of records from db
func (s *Session) QueryRows() (rows *sql.Rows, err error) {
	defer s.Clear()
	if s.isDebug {
		s.debugSql(s.sql.String(), s.sqlVars...)
	}
	if rows, err = s.DB().Query(s.sql.String(), s.sqlVars...); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) Debug() *Session {
	s.isDebug = true
	return s
}

func (s *Session) debugSql(query string, args ...any) {
	if len(args) == 0 {
		log.Debug(query)
	} else {
		log.Debugf(strings.Replace(query, "?", "%v", len(args)), args...)
	}
}
