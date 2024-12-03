package session

import (
	"database/sql"
	"github.com/go-needle/log"
	"github.com/go-needle/orm/dialect"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

type User struct {
	Name string `orm:"name:user_name;constraint:PRIMARY KEY"`
	Age  int
}

func TestSession_CreateTable(t *testing.T) {
	db, err := sql.Open("sqlite3", "gee.db")
	if err != nil {
		log.Error(err)
		return
	}
	d, _ := dialect.GetDialect("sqlite3")
	s := New(db, d).Model(&User{}).Table("sys_user")
	_ = s.DropTable()
	_ = s.CreateTable()
	if !s.HasTable() {
		t.Fatal("Failed to create table User")
	}
}
