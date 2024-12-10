package schema

import (
	"github.com/go-needle/orm/dialect"
	"go/ast"
	"reflect"
	"strings"
)

// Field represents a column of database
type Field struct {
	Name        string
	MappingName string
	Type        string
	Constraint  string
}

// Schema represents a table of database
type Schema struct {
	Model             any
	Name              string
	Fields            []*Field
	MappingFieldNames []string
	fieldMap          map[string]*Field
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

func Parse(dest any, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name:        p.Name,
				MappingName: p.Name,
				Type:        d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("orm"); ok {
				tags := strings.Split(v, ";")
				for _, tag := range tags {
					temp := strings.Split(tag, ":")
					if len(temp) != 2 {
						continue
					}
					key, value := strings.TrimSpace(temp[0]), strings.TrimSpace(temp[1])
					switch key {
					case "name":
						field.MappingName = value
					case "constraint":
						field.Constraint = value
					}
				}
			}
			schema.Fields = append(schema.Fields, field)
			schema.MappingFieldNames = append(schema.MappingFieldNames, field.MappingName)
			schema.fieldMap[p.Name] = field
			schema.fieldMap[field.MappingName] = field
		}
	}
	return schema
}

func (schema *Schema) RecordValues(dest any) []any {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []any
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}
