package mysql_session_db

import (
	"bufio"
	"bytes"
	"database/sql"
	"io"
	"os"
	"regexp"
	"strconv"
	"fmt"
	"errors"
)

type CDbHandler struct  {
	m_db *sql.DB
}

func (this *CDbHandler) Connect(host string, port uint, username string, userpwd string, dbname string, dbtype string) (err error) {
	b := bytes.Buffer{}
	b.WriteString(username)
	b.WriteString(":")
	b.WriteString(userpwd)
	b.WriteString("@tcp(")
	b.WriteString(host)
	b.WriteString(":")
	b.WriteString(strconv.FormatUint(uint64(port), 10))
	b.WriteString(")/")
	b.WriteString(dbname)
	var name string
	if dbtype == "mysql" {
		name = b.String()
	} else if dbtype == "sqlite3" {
		name = dbname
	} else {
		return errors.New("dbtype not support")
	}
	this.m_db, err = sql.Open(dbtype, name)
	if err != nil {
		return err
	}
	this.m_db.SetMaxOpenConns(2000)
	this.m_db.SetMaxIdleConns(1000)
	this.m_db.Ping()
	return nil
}

func (this *CDbHandler) ConnectByRule(rule string, dbtype string) (err error) {
	this.m_db, err = sql.Open(dbtype, rule)
	if err != nil {
		return err
	}
	this.m_db.SetMaxOpenConns(2000)
	this.m_db.SetMaxIdleConns(1000)
	this.m_db.Ping()
	return nil
}

func (this *CDbHandler) ConnectByCfg(path string) error {
	fi, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fi.Close()
	br := bufio.NewReader(fi)
	var host string = "localhost"
	var port uint = 3306
	var username string = "root"
	var userpwd string = "123456"
	var dbname string = "test"
	var dbtype string = "mysql"
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		content := string(a)
		r, _ := regexp.Compile("(.*)?=(.*)?")
		ret := r.FindStringSubmatch(content)
		if len(ret) != 3 {
			continue
		}
		k := ret[1]
		v := ret[2]
		switch k {
		case "host":
			host = v
		case "port":
			port_tmp, _ := strconv.ParseUint(v, 10, 32)
			port = uint(port_tmp)
		case "username":
			username = v
		case "userpwd":
			userpwd = v
		case "dbname":
			dbname = v
		case "dbtype":
			dbtype = v
		}
	}
	return this.Connect(host, port, username, userpwd, dbname, dbtype)
}

func (this *CDbHandler) Disconnect() {
	this.m_db.Close()
}

func (this *CDbHandler) Create() (error) {
	var err error = nil
	var _ error = err
	_, err = this.m_db.Exec(`create table if not exists t_session_info (
    sessionuuid varchar(64),
    timeoutS bigint,
    losevalidtime bigint
) charset utf8;`)
	if err != nil {
		return err
	}
	return nil
}

func (this *CDbHandler) AddSession(input0 *CAddSessionInput) (error, uint64) {
	var rowCount uint64 = 0
	tx, _ := this.m_db.Begin()
	var result sql.Result
	var _ = result
	var err error
	var _ error = err
	result, err = this.m_db.Exec(fmt.Sprintf(`insert into t_session_info values(?, ?, ?);`), input0.SessionUuid, input0.TimeoutS, input0.LoseValidTime)
	if err != nil {
		tx.Rollback()
		return err, rowCount
	}
	tx.Commit()
	var _ = result
	return nil, rowCount
}

func (this *CDbHandler) DeleteSession(input0 *[]CDeleteSessionInput) (error, uint64) {
	var rowCount uint64 = 0
	tx, _ := this.m_db.Begin()
	var result sql.Result
	var _ = result
	var err error
	var _ error = err
	for _, v := range *input0 {
		result, err = this.m_db.Exec(fmt.Sprintf(`delete from t_session_info where sessionuuid = ?;`), v.SessionUuid)
		if err != nil {
			tx.Rollback()
			return err, rowCount
		}
		var _ = result
	}
	tx.Commit()
	return nil, rowCount
}

func (this *CDbHandler) UpdateSession(input0 *[]CUpdateSessionInput) (error, uint64) {
	var rowCount uint64 = 0
	tx, _ := this.m_db.Begin()
	var result sql.Result
	var _ = result
	var err error
	var _ error = err
	for _, v := range *input0 {
		result, err = this.m_db.Exec(fmt.Sprintf(`update t_session_info set %s where sessionuuid = ?;`, v.Condition), v.SessionUuid)
		if err != nil {
			tx.Rollback()
			return err, rowCount
		}
		var _ = result
	}
	tx.Commit()
	return nil, rowCount
}

func (this *CDbHandler) GetSession(input0 *CGetSessionInput, output0 *CGetSessionOutput) (error, uint64) {
	var rowCount uint64 = 0
	tx, _ := this.m_db.Begin()
	var result sql.Result
	var _ = result
	var err error
	var _ error = err
	rows0, err := this.m_db.Query(fmt.Sprintf(`select timeoutS, losevalidtime from t_session_info
where sessionuuid = ?;`), input0.SessionUuid)
	if err != nil {
		tx.Rollback()
		return err, rowCount
	}
	tx.Commit()
	defer rows0.Close()
	for rows0.Next() {
		rowCount += 1
		var timeoutS sql.NullInt64
		var loseValidTime sql.NullInt64
		scanErr := rows0.Scan(&timeoutS, &loseValidTime)
		if scanErr != nil {
			continue
		}
		output0.TimeoutS = int64(timeoutS.Int64)
		output0.TimeoutSIsValid = timeoutS.Valid
		output0.LoseValidTime = int64(loseValidTime.Int64)
		output0.LoseValidTimeIsValid = loseValidTime.Valid
	}
	return nil, rowCount
}

func (this *CDbHandler) GetCountBySessionUuid(input0 *CGetCountBySessionUuidInput, output0 *CGetCountBySessionUuidOutput) (error, uint64) {
	var rowCount uint64 = 0
	tx, _ := this.m_db.Begin()
	var result sql.Result
	var _ = result
	var err error
	var _ error = err
	rows0, err := this.m_db.Query(fmt.Sprintf(`select count(0) from t_session_info
where sessionuuid = ?;`), input0.SessionUuid)
	if err != nil {
		tx.Rollback()
		return err, rowCount
	}
	tx.Commit()
	defer rows0.Close()
	for rows0.Next() {
		rowCount += 1
		var count sql.NullInt64
		scanErr := rows0.Scan(&count)
		if scanErr != nil {
			continue
		}
		output0.Count = int(count.Int64)
		output0.CountIsValid = count.Valid
	}
	return nil, rowCount
}

func (this *CDbHandler) DeleteLosetimeRecord() (error, uint64) {
	var rowCount uint64 = 0
	tx, _ := this.m_db.Begin()
	var result sql.Result
	var _ = result
	var err error
	var _ error = err
	result, err = this.m_db.Exec(fmt.Sprintf(`delete from t_session_info where sessionuuid in
(
    select tmp.sessionuuid from
    (
        select sessionuuid from t_session_info where losevalidtime < unix_timestamp(now())
    ) as tmp
);`))
	if err != nil {
		tx.Rollback()
		return err, rowCount
	}
	tx.Commit()
	var _ = result
	return nil, rowCount
}

