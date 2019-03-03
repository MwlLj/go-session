package session

import (
	"errors"
	"github.com/MwlLj/go-session/memory"
	"github.com/MwlLj/go-session/mysql"
	"github.com/MwlLj/go-session/redis"
	"github.com/MwlLj/go-session/sqlite3"
)

var Memory_type_sqlite string = "sqlite"
var Memory_type_mysql string = "mysql"
var Memory_type_memory string = "memory"
var Memory_type_redis string = "redis"

type ISession interface {
	/*
		@name Dial
		@params
			rule: if mysql, username:userpwd@tcp(host:port)/dbname
				  if sqlite, dbpath
	*/
	Dial(rule string) error
	Create(timeoutS int64) (id *string, e error)
	CreateWithMap(timeoutS int64, extraInfo *map[string]string) (id *string, e error)
	Destory(id *string) error
	IsValid(id *string) (bool, error)
	IsValidWithMap(id *string) (bool, *map[string]string, error)
	/*
		@name Reset
		@params
			id: Create return id
			timeoutS: if this field is nil, keep last timeoutS
	*/
	Reset(id *string, timeoutS *int64) error
}

func New(memoryType *string) (ISession, error) {
	if memoryType == nil {
		return nil, errors.New("memoryType is nil")
	}
	if *memoryType == Memory_type_memory {
		return memory.New(), nil
	} else if *memoryType == Memory_type_mysql {
		return mysql.New(), nil
	} else if *memoryType == Memory_type_sqlite {
		return sqlite3.New(), nil
	} else if *memoryType == Memory_type_redis {
		return redis.New(), nil
	}
	return nil, errors.New("memoryType is not support")
}
