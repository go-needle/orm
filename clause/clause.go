package clause

import (
	"strings"
)

type Clause struct {
	sql     [6]string
	sqlVars [6][]any
}

type Type int

const (
	INSERT Type = iota
	VALUES
	SELECT
	LIMIT
	WHERE
	ORDERBY
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
		if order < 6 && c.sql[order] != "" {
			sqls = append(sqls, c.sql[order])
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls, " "), vars
}
