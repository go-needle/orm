package clause

import (
	"strings"
)

type Clause struct {
	sql     [9]string
	sqlVars [9][]any
}

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
	UPDATE
	DELETE
	COUNT
)

func (c *Clause) Set(name Type, vars ...any) {
	sql, vars := generators[name](vars...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
}

func (c *Clause) Build(orders ...Type) (string, []any) {
	var sqls []string
	var vars []any
	for _, order := range orders {
		if order < 9 && c.sql[order] != "" {
			sqls = append(sqls, c.sql[order])
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}

func (c *Clause) Clear() {
	c.sql = [9]string{}
	c.sqlVars = [9][]any{}
}
