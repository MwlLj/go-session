package sqlite3

import (
	"errors"
)

type CSqlite3 struct {
}

func (this *CSqlite3) init() {
}

func (this *CSqlite3) Dial(rule string) error {
	return nil
}

func (this *CSqlite3) Create(timeoutS int64) (id *string, e error) {
	return nil, nil
}

func (this *CSqlite3) CreateWithMap(timeoutS int64, extraInfo *map[string]string) (id *string, e error) {
	return nil, errors.New("not support")
}

func (this *CSqlite3) Destory(id *string) error {
	return nil
}

func (this *CSqlite3) IsValid(id *string) (bool, error) {
	return true, nil
}

func (this *CSqlite3) IsValidWithMap(id *string) (bool, *map[string]string, error) {
	return false, nil, errors.New("not support")
}

func (this *CSqlite3) Reset(id *string, timeoutS *int64) error {
	return nil
}

func (this *CSqlite3) KeyTimeoutNtf() <-chan *string {
	return nil
}

func (this *CSqlite3) StartExpiredEventListen() {
}

func New() *CSqlite3 {
	db := CSqlite3{}
	db.init()
	return &db
}
