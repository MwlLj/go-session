package redis

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"github.com/satori/go.uuid"
	"log"
	"strconv"
)

type CRedis struct {
	m_conn redis.Conn
}

func (this *CRedis) Dial(rule string) error {
	var err error = nil
	this.m_conn, err = redis.Dial("tcp", rule)
	if err != nil {
		log.Fatalf("connect redis server error, rule: %s, err: %v\n", rule, err)
		return err
	}
	return nil
}

func (this *CRedis) Create(timeoutS int64) (id *string, e error) {
	uid, err := uuid.NewV4()
	if err != nil {
		log.Println("session create uuid error")
		return nil, err
	}
	v4Uuid := uid.String()
	_, err = this.m_conn.Do("set", v4Uuid, strconv.FormatInt(timeoutS, 10), "ex", timeoutS)
	if err != nil {
		log.Println("set session error, id: %s, timeoutS: %d\n", v4Uuid, timeoutS)
		return nil, err
	}
	return &v4Uuid, nil
}

func (this *CRedis) Destory(id *string) error {
	if id == nil {
		return errors.New("destory id is nil")
	}
	_, err := this.m_conn.Do("del", *id)
	if err != nil {
		log.Println("delete id from redis error, err: %v\n", err)
		return err
	}
	return nil
}

func (this *CRedis) IsValid(id *string) (bool, error) {
	if id == nil {
		return false, errors.New("isValid id is nil")
	}
	result, err := this.m_conn.Do("exists", *id)
	if err != nil {
		log.Println("get is exists from redis error, err: %v\n", err)
		return false, err
	}
	isExist := result.(int64)
	if isExist == 0 {
		return false, nil
	}
	return true, nil
}

func (this *CRedis) Reset(id *string, timeoutS *int64) error {
	if id == nil {
		return errors.New("reset id is nil")
	}
	var timeout int64
	if timeoutS == nil {
		t, err := this.m_conn.Do("get", *id)
		if err != nil {
			log.Printf("get id timeout error, err: %v\n", err)
			return err
		}
		timeoutStr := t.([]uint8)
		timeout, err = strconv.ParseInt(string(timeoutStr), 10, 64)
		if err != nil {
			log.Println("parse timeout from string to int error")
			return err
		}
	} else {
		timeout = *timeoutS
	}
	_, err := this.m_conn.Do("expire", *id, timeout)
	if err != nil {
		log.Println("update timeout to redis error, err: %v\n", err)
		return err
	}
	return nil
}

func New() *CRedis {
	redis := CRedis{}
	return &redis
}
