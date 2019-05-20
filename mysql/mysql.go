package mysql

import (
	"bytes"
	"errors"
	"github.com/MwlLj/go-session/common"
	db "github.com/MwlLj/go-session/mysql/session_db"
	"github.com/satori/go.uuid"
	"log"
	"strconv"
	"time"
)

type CMysql struct {
	m_dbHandler db.CDbHandler
	m_isConnect bool
}

func (this *CMysql) init() {
	this.m_isConnect = false
	this.check()
}

func (this *CMysql) check() {
	go func() {
		for {
			if this.m_isConnect == false {
				time.Sleep(10 * time.Millisecond)
				continue
			}
			err, _ := this.m_dbHandler.DeleteLosetimeRecord()
			if err != nil {
				log.Println("delete losrtime record from db error")
			}
			time.Sleep(1 * time.Second)
		}
	}()
}

func (this *CMysql) Dial(rule string) error {
	err := this.m_dbHandler.ConnectByRule(rule, "mysql")
	if err != nil {
		log.Fatalf("connect db error, rule: %s, err: %v\n", rule, err)
	}
	err = this.m_dbHandler.Create()
	if err != nil {
		log.Fatalln("create table error")
	}
	this.m_isConnect = true
	return err
}

func (this *CMysql) Create(timeoutS int64) (id *string, e error) {
	uid, err := uuid.NewV4()
	if err != nil {
		log.Println("session create uuid error")
		return nil, err
	}
	v4Uuid := uid.String()
	input := db.CAddSessionInput{}
	input.SessionUuid = v4Uuid
	input.TimeoutS = timeoutS
	input.LoseValidTime = common.GetNowTimeStampS() + timeoutS
	err, _ = this.m_dbHandler.AddSession(&input)
	if err != nil {
		log.Println("add session to db error")
		return nil, err
	}
	return &v4Uuid, nil
}

func (this *CMysql) CreateWithMap(timeoutS int64, extraInfo *map[string]string) (id *string, e error) {
	return nil, errors.New("not support")
}

func (this *CMysql) Destory(id *string) error {
	if id == nil {
		return errors.New("delete session id is nil")
	}
	input := []db.CDeleteSessionInput{}
	in := db.CDeleteSessionInput{}
	in.SessionUuid = *id
	input = append(input, in)
	err, _ := this.m_dbHandler.DeleteSession(&input)
	if err != nil {
		log.Println("delete session from db error")
		return err
	}
	return nil
}

func (this *CMysql) IsValid(id *string) (bool, error) {
	if id == nil {
		return false, errors.New("isvalid id is nil")
	}
	input := db.CGetCountBySessionUuidInput{}
	output := db.CGetCountBySessionUuidOutput{}
	input.SessionUuid = *id
	err, _ := this.m_dbHandler.GetCountBySessionUuid(&input, &output)
	if err != nil {
		log.Println("getCountBySessionUuid from db error")
		return false, err
	}
	if output.Count == 0 {
		return false, nil
	}
	return true, nil
}

func (this *CMysql) IsValidWithMap(id *string) (bool, *map[string]string, error) {
	return false, nil, errors.New("not support")
}

func (this *CMysql) Reset(id *string, timeoutS *int64) error {
	if id == nil {
		return errors.New("reset session id is nil")
	}
	input := []db.CUpdateSessionInput{}
	in := db.CUpdateSessionInput{}
	in.SessionUuid = *id
	buffer := bytes.Buffer{}
	if timeoutS != nil {
		buffer.WriteString("timeoutS = ")
		buffer.WriteString(strconv.FormatInt(*timeoutS, 10))
		buffer.WriteString(", ")
		buffer.WriteString("losevalidtime = unix_timestamp(now()) + ")
		buffer.WriteString(strconv.FormatInt(*timeoutS, 10))
	} else {
		buffer.WriteString("losevalidtime = unix_timestamp(now()) + timeoutS")
	}
	in.Condition = buffer.String()
	input = append(input, in)
	err, _ := this.m_dbHandler.UpdateSession(&input)
	if err != nil {
		log.Println("update session from db error")
		return err
	}
	return nil
}

func (this *CMysql) KeyTimeoutNtf() <-chan *string {
	return nil
}

func (this *CMysql) StartExpiredEventListen() {
}

func New() *CMysql {
	db := CMysql{}
	db.init()
	return &db
}
