package session

import (
	"errors"
	"github.com/MwlLj/go-session/memory"
	"github.com/MwlLj/go-session/persistent"
	"github.com/MwlLj/go-session/redis"
)

var Memory_type_sqlite string = "sqlite"
var Memory_type_mysql string = "mysql"
var Memory_type_memory string = "memory"
var Memory_type_redis string = "redis"

type ISession interface {
	Dial(rule string) error
	Create(timeoutS int64) (id *string, e error)
	Destroy(id *string) error
	IsValid(id *string) (bool, error)
	Reset(id *string) error
}

func New(memoryType *string) (ISession, error) {
	if memoryType == nil {
		return nil, errors.New("memoryType is nil")
	}
	if *memoryType == Memory_type_memory {
		return memory.New(), nil
	} else if *memoryType == Memory_type_mysql {
		return persistent.New("mysql"), nil
	} else if *memoryType == Memory_type_sqlite {
		return persistent.New("sqlite"), nil
	} else if *memoryType == Memory_type_redis {
		return redis.New(), nil
	}
	return nil, errors.New("memoryType is not support")
}
