package session

import (
	"errors"
	"github.com/go-needle/orm/clause"
	"reflect"
	"strings"
)

func (s *Session) Insert(values ...any) (int64, error) {
	recordValues := make([]any, 0)
	for _, value := range values {
		table := s.Model(value).RefTable()
		s.clause.Set(clause.INSERT, table.Name, table.MappingFieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}
	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (s *Session) Find(values any) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	s.clause.Set(clause.SELECT, table.Name, table.MappingFieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}

	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []any
		for _, field := range table.Fields {
			values = append(values, dest.FieldByName(field.Name).Addr().Interface())
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}

func (s *Session) First(value any) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	dest.Set(destSlice.Index(0))
	return nil
}

// support map[string]any
// also support kv list: "Name", "Tom", "Age", 18, ....
func (s *Session) Update(kv ...any) (int64, error) {
	m, ok := kv[0].(map[string]any)
	if !ok {
		m = make(map[string]any)
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}
	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Save(value any) (int64, error) {
	m := make(map[string]any)
	modelValue := reflect.Indirect(reflect.ValueOf(value))
	modelType := modelValue.Type()
	for i := 0; i < modelValue.NumField(); i++ {
		if !modelValue.Field(i).IsZero() {
			m[s.refTable.GetField(modelType.Field(i).Name).MappingName] = modelValue.Field(i).Interface()
		}
	}
	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete records with where clause
func (s *Session) Delete() (int64, error) {
	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Count records with where clause
func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}

// Limit adds limit condition to clause
func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

// Where adds limit condition to clause
func (s *Session) Where(desc any, args ...any) *Session {
	var vars []any
	if reflect.TypeOf(desc).Kind() == reflect.String {
		vars = append(vars, desc)
		vars = append(vars, args...)
		s.clause.Set(clause.WHERE, vars...)
	} else {
		modelValue := reflect.Indirect(reflect.ValueOf(desc))
		modelType := modelValue.Type()
		var description strings.Builder
		var values []any
		for i := 0; i < modelValue.NumField(); i++ {
			if !modelValue.Field(i).IsZero() {
				mappingName := s.refTable.GetField(modelType.Field(i).Name).MappingName
				description.WriteString(mappingName)
				description.WriteString(" = ?")
				values = append(values, modelValue.Field(i).Interface())
			}
		}
		vars = append(vars, description.String())
		vars = append(vars, values...)
		s.clause.Set(clause.WHERE, vars...)
	}
	return s
}

// OrderBy adds order by condition to clause
func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}
