package redis

import (
	"errors"
	// "github.com/garyburd/redigo/redis"
	"github.com/go-redis/redis"
	"github.com/satori/go.uuid"
	"log"
	"strconv"
	"time"
)

const (
	fieldTimeStamp string = "time-stamp"
)

type CRedis struct {
	m_conn              *redis.Client
	m_redisExpiredEvent <-chan *redis.Message
	m_expiredEvent      chan *string
}

func (this *CRedis) Dial(rule string) error {
	var err error = nil
	this.m_conn = redis.NewClient(&redis.Options{
		Addr:     rule,
		Password: "",
		DB:       0,
	})
	_, err = this.m_conn.Ping().Result()
	if err != nil {
		log.Fatalf("connect redis server error, rule: %s, err: %v\n", rule, err)
		return err
	}
	// init event notify
	pb := this.m_conn.PSubscribe("__keyevent@*__expired")
	_, err = pb.Receive()
	if err != nil {
		log.Fatalln("receive subscribe error")
	}
	this.m_redisExpiredEvent = pb.Channel()
	go func() {
		for {
			message := <-this.m_redisExpiredEvent
			s := message.String()
			this.m_expiredEvent <- &s
		}
	}()
	return nil
}

func (this *CRedis) Create(timeoutS int64) (id *string, e error) {
	uid, err := uuid.NewV4()
	if err != nil {
		log.Println("session create uuid error")
		return nil, err
	}
	v4Uuid := uid.String()
	err = this.m_conn.HMSet(v4Uuid, map[string]interface{}{fieldTimeStamp: strconv.FormatInt(timeoutS, 10)}).Err()
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
	arr := map[string]interface{}{}
	arr[fieldTimeStamp] = strconv.FormatInt(timeoutS, 10)
	if extraInfo != nil {
		for k, v := range *extraInfo {
			arr[k] = v
		}
	}
	// add expire time
	err = this.m_conn.HMSet(v4Uuid, arr).Err()
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
	err := this.m_conn.HDel(*id).Err()
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
	result, err := this.m_conn.HGetAll(*id).Result()
	if err != nil {
		log.Println("get is exists from redis error, err: %v\n", err)
		return false, err
	}
	length := len(result)
	if length == 0 {
		return false, nil
	}
	if valueStr, ok := result[fieldTimeStamp]; ok {
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
	}
	return true, nil
}

func (this *CRedis) IsValidWithMap(id *string) (bool, *map[string]string, error) {
	if id == nil {
		return false, nil, errors.New("isValid id is nil")
	}
	result, err := this.m_conn.HGetAll(*id).Result()
	if err != nil {
		log.Println("get is exists from redis error, err: %v\n", err)
		return false, nil, err
	}
	length := len(result)
	if length == 0 {
		return false, nil, errors.New("sessionid not found")
	}
	if valueStr, ok := result[fieldTimeStamp]; ok {
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
	}
	return true, &result, nil
}

func (this *CRedis) Reset(id *string, timeoutS *int64) error {
	if id == nil {
		return errors.New("reset id is nil")
	}
	var timeout int64
	if timeoutS == nil {
		t, err := this.m_conn.HMGet(*id, fieldTimeStamp).Result()
		if err != nil {
			log.Printf("get id timeout error, err: %v\n", err)
			return err
		}
		if t == nil || len(t) < 1 {
			log.Println("hmget length is least 0")
			return errors.New("hmget length is least 0, get timestamp error")
		}
		timeoutStr := t[0].(string)
		timeout, err = strconv.ParseInt(string(timeoutStr), 10, 64)
		if err != nil {
			log.Println("parse timeout from string to int error")
			return err
		}
	} else {
		timeout = *timeoutS
	}
	err := this.m_conn.Expire(*id, time.Duration(timeout)*time.Second).Err()
	if err != nil {
		log.Println("update timeout to redis error, err: %v\n", err)
		return err
	}
	return nil
}

func (this *CRedis) KeyTimeoutNtf() <-chan *string {
	return this.m_expiredEvent
	// _, err := pubsub.Receive()
	// if err != nil {
	// 	panic(err)
	// }

	// // Go channel which receives messages.
	// ch := pubsub.Channel()

	// // Publish a message.
	// err = redisdb.Publish("mychannel1", "hello").Err()
	// if err != nil {
	// 	panic(err)
	// }

	// time.AfterFunc(time.Second, func() {
	// 	// When pubsub is closed channel is closed too.
	// 	_ = pubsub.Close()
	// })

	// // Consume messages.
	// for msg := range ch {
	// 	fmt.Println(msg.Channel, msg.Payload)
	// }
	// return nil
}

func New() *CRedis {
	redis := CRedis{}
	redis.m_expiredEvent = make(chan *string)
	return &redis
}
