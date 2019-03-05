package redis

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"github.com/satori/go.uuid"
	"log"
	"strconv"
)

const (
	fieldTimeStamp string = "time-stamp"
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
	_, err = this.m_conn.Do("hmset", v4Uuid, fieldTimeStamp, strconv.FormatInt(timeoutS, 10))
	if err != nil {
		log.Printf("set session error, id: %s, timeoutS: %d\n", v4Uuid, timeoutS)
		return nil, err
	}
	err = this.Reset(&v4Uuid, &timeoutS)
	if err != nil {
		log.Printf("reset time error, err: %v\n", err)
		return nil, err
	}
	return &v4Uuid, nil
}

func (this *CRedis) CreateWithMap(timeoutS int64, extraInfo *map[string]string) (id *string, e error) {
	uid, err := uuid.NewV4()
	if err != nil {
		log.Println("session create uuid error")
		return nil, err
	}
	v4Uuid := uid.String()
	// map -> array
	arr := []interface{}{}
	arr = append(arr, v4Uuid)
	arr = append(arr, fieldTimeStamp)
	arr = append(arr, strconv.FormatInt(timeoutS, 10))
	if extraInfo != nil {
		for k, v := range *extraInfo {
			arr = append(arr, k)
			arr = append(arr, v)
		}
	}
	// add expire time
	_, err = this.m_conn.Do("hmset", arr...)
	if err != nil {
		log.Printf("set session error, id: %s, timeoutS: %d\n", v4Uuid, timeoutS)
		return nil, err
	}
	err = this.Reset(&v4Uuid, &timeoutS)
	if err != nil {
		log.Printf("reset time error, err: %v\n", err)
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
	result, err := redis.Values(this.m_conn.Do("hgetall", *id))
	if err != nil {
		log.Println("get is exists from redis error, err: %v\n", err)
		return false, err
	}
	if result == nil {
		return false, nil
	}
	length := len(result)
	for i := 0; i < length; i += 2 {
		v := result[i]
		vStr := string(v.([]byte))
		if fieldTimeStamp == vStr {
			valueStr := string(result[i+1].([]byte))
			t, err := strconv.ParseInt(valueStr, 10, 64)
			if err != nil {
				log.Printf("get timeout error, err: %v\n", err)
				return false, err
			}
			err = this.Reset(id, &t)
			if err != nil {
				log.Printf("reset time error, err: %v\n", err)
				return false, err
			}
			continue
		}
	}
	return true, nil
}

func (this *CRedis) IsValidWithMap(id *string) (bool, *map[string]string, error) {
	if id == nil {
		return false, nil, errors.New("isValid id is nil")
	}
	result, err := redis.Values(this.m_conn.Do("hgetall", *id))
	if err != nil {
		log.Println("get is exists from redis error, err: %v\n", err)
		return false, nil, err
	}
	if result == nil {
		return false, nil, nil
	}
	extraValues := make(map[string]string)
	length := len(result)
	for i := 0; i < length; i += 2 {
		v := result[i]
		vStr := string(v.([]byte))
		if fieldTimeStamp == vStr {
			valueStr := string(result[i+1].([]byte))
			t, err := strconv.ParseInt(valueStr, 10, 64)
			if err != nil {
				log.Printf("get timeout error, err: %v\n", err)
				return false, nil, err
			}
			err = this.Reset(id, &t)
			if err != nil {
				log.Printf("reset time error, err: %v\n", err)
				return false, nil, err
			}
			continue
		}
		extraValues[vStr] = string(result[i+1].([]byte))
	}
	return true, &extraValues, nil
}

func (this *CRedis) Reset(id *string, timeoutS *int64) error {
	if id == nil {
		return errors.New("reset id is nil")
	}
	var timeout int64
	if timeoutS == nil {
		t, err := redis.Values(this.m_conn.Do("hmget", *id, fieldTimeStamp))
		if err != nil {
			log.Printf("get id timeout error, err: %v\n", err)
			return err
		}
		if len(t) < 1 {
			log.Println("hmget length is least 0")
			return errors.New("hmget length is least 0, get timestamp error")
		}
		timeoutStr := t[0].([]byte)
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
