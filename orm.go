package orm

import (
	"database/sql"
	"github.com/go-needle/log"
	"github.com/go-needle/orm/dialect"
	"github.com/go-needle/orm/session"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
	Log     *log.Logger
}

func NewEngine(driver, source string) (e *Engine, err error) {
	db, err := sql.Open(driver, source)
	logger := log.New()
	if err != nil {
		logger.Error(err)
		return
	}
	// Send a ping to make sure the database connection is alive.
	if err = db.Ping(); err != nil {
		logger.Error(err)
		return
	}
	// make sure the specific dialect exists
	dial, ok := dialect.GetDialect(driver)
	if !ok {
		logger.Errorf("dialect %s Not Found", driver)
		return
	}
	e = &Engine{db: db, dialect: dial, Log: logger}
	logger.Info("Connect database success")
	return
}

func (engine *Engine) Close() {
	if err := engine.db.Close(); err != nil {
		engine.Log.Error("Failed to close database")
	}
	engine.Log.Info("Close database success")
}

func (engine *Engine) NewSession() *session.Session {
	return session.New(engine.db, engine.dialect, engine.Log)
}
