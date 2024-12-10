package session

import (
	"github.com/go-needle/orm/log"
)

// Hooks constants
const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

type IBeforeQuery interface {
	BeforeQuery(s *Session) error
}

type IAfterQuery interface {
	AfterQuery(s *Session) error
}
type IBeforeUpdate interface {
	BeforeUpdate(s *Session) error
}
type IAfterUpdate interface {
	AfterUpdate(s *Session) error
}
type IBeforeDelete interface {
	BeforeDelete(s *Session) error
}
type IAfterDelete interface {
	AfterDelete(s *Session) error
}
type IBeforeInsert interface {
	BeforeInsert(s *Session) error
}
type IAfterInsert interface {
	AfterInsert(s *Session) error
}

// CallMethod calls the registered hooks
func (s *Session) CallMethod(method string, value any) {
	dest := value
	if dest == nil {
		dest = s.RefTable().Model
	}
	var err error
	switch method {
	case BeforeQuery:
		if v, ok := dest.(IBeforeQuery); ok {
			err = v.BeforeQuery(s)
		}
	case AfterQuery:
		if v, ok := dest.(IAfterQuery); ok {
			err = v.AfterQuery(s)
		}
	case BeforeUpdate:
		if v, ok := dest.(IBeforeUpdate); ok {
			err = v.BeforeUpdate(s)
		}
	case AfterUpdate:
		if v, ok := dest.(IAfterUpdate); ok {
			err = v.AfterUpdate(s)
		}
	case BeforeDelete:
		if v, ok := dest.(IBeforeDelete); ok {
			err = v.BeforeDelete(s)
		}
	case AfterDelete:
		if v, ok := dest.(IAfterDelete); ok {
			err = v.AfterDelete(s)
		}
	case BeforeInsert:
		if v, ok := dest.(IBeforeInsert); ok {
			err = v.BeforeInsert(s)
		}
	case AfterInsert:
		if v, ok := dest.(IAfterInsert); ok {
			err = v.AfterInsert(s)
		}
	default:
		log.Fatal("Unsupported hook method")
	}
	if err != nil {
		log.Error(err)
	}
	return
}
