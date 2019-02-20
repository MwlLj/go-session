package redis

type CRedis struct {
}

func (this *CRedis) Dial(rule string) error {
	return nil
}

func (this *CRedis) Create(timeoutS int64) (id *string, e error) {
	return nil, nil
}

func (this *CRedis) Destory(id *string) error {
	return nil
}

func (this *CRedis) IsValid(id *string) (bool, error) {
	return true, nil
}

func (this *CRedis) Reset(id *string, timeoutS *int64) error {
	return nil
}

func New() *CRedis {
	redis := CRedis{}
	return &redis
}
