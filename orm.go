package orm

import (
	"database/sql"
	"fmt"
	"github.com/go-needle/orm/dialect"
	"github.com/go-needle/orm/log"
	"github.com/go-needle/orm/session"
	"strings"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	if err != nil {
		log.Error(err)
		return
	}
	// Send a ping to make sure the database connection is alive.
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}
	// make sure the specific dialect exists
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		log.Errorf("dialect %s Not Found", driver)
		return
	}
	e = &Engine{db: db, dialect: dial}
	log.Info("Connect database success")
	return
}

func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db, engine.dialect)
}

type TxFunc func(*session.Session) (any, error)

func (engine *Engine) Transaction(f TxFunc, isDebug bool) (result any, err error) {
	s := engine.NewSession()
	if isDebug {
		s = s.Debug()
	}
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = s.Rollback() // err is non-nil; don't change it
		} else {
			err = s.Commit() // err is nil; if Commit returns error update err
		}
	}()

	return f(s)
}

// difference returns a - b
func difference(a []string, b []string) (diff []string) {
	setB := make(map[string]struct{})
	for _, v := range b {
		setB[v] = struct{}{}
	}
	for _, v := range a {
		if _, ok := setB[v]; !ok {
			diff = append(diff, v)
		}
	}
	return
}

// Migrate table
func (engine *Engine) Migrate(value any, isDebug bool) error {
	_, err := engine.Transaction(func(s *session.Session) (result any, err error) {
		if !s.Model(value).HasTable() {
			log.Infof("table %s doesn't exist", s.RefTable().Name)
			return nil, s.CreateTable()
		}
		table := s.RefTable()
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).QueryRows()
		columns, _ := rows.Columns()
		addCols := difference(table.MappingFieldNames, columns)
		delCols := difference(columns, table.MappingFieldNames)
		log.Infof("added cols %v, deleted cols %v", addCols, delCols)
		for _, col := range addCols {
			f := table.GetField(col)
			sqlStr := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", table.Name, f.MappingName, f.Type)
			if _, err = s.Raw(sqlStr).Exec(); err != nil {
				return
			}
		}

		if len(delCols) == 0 {
			return
		}
		tmp := "tmp_" + table.Name
		fieldStr := strings.Join(table.MappingFieldNames, ", ")
		s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s FROM %s;", tmp, fieldStr, table.Name))
		s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name))
		s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmp, table.Name))
		_, err = s.Exec()
		return
	}, isDebug)
	return err
}
