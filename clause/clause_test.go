package clause

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func testSelect(t *testing.T) {
	var clause Clause
	go clause.Set(LIMIT, 3)
	go clause.Set(SELECT, "User", []string{"*"})
	go clause.Set(WHERE, "Name = ?", "Tom")
	go clause.Set(ORDERBY, "Age ASC")
	sql, vars := clause.Build(SELECT, WHERE, ORDERBY, LIMIT)
	t.Log(sql, vars)
	if sql != "SELECT * FROM User WHERE Name = ? ORDER BY Age ASC LIMIT ?" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []any{"Tom", 3}) {
		t.Fatal("failed to build SQLVars")
	}
}

func TestClause_Build(t *testing.T) {
	t.Run("select", func(t *testing.T) {
		testSelect(t)
	})
}

func test() {
	var clause Clause
	clause.Set(LIMIT, 3)
	clause.Set(SELECT, "User", []string{"*"})
	clause.Set(WHERE, "Name = ?", "Tom")
	clause.Set(ORDERBY, "Age ASC")
	_, _ = clause.Build(SELECT, WHERE, ORDERBY, LIMIT)
}
func TestTime(t *testing.T) {
	start := time.Now()
	for i := 0; i < 100000; i++ {
		test()
	}
	fmt.Println(time.Since(start).String())
}
