package mysql

type CMysql struct {
	m_dbType *string
}

func (this *CMysql) init() {
}

func (this *CMysql) Dial(rule string) error {
	return nil
}

func (this *CMysql) Create(timeoutS int64) (id *string, e error) {
	return nil, nil
}

func (this *CMysql) Destory(id *string) error {
	return nil
}

func (this *CMysql) IsValid(id *string) (bool, error) {
	return true, nil
}

func (this *CMysql) Reset(id *string, timeoutS *int64) error {
	return nil
}

func New() *CMysql {
	db := CMysql{}
	db.init()
	return &db
}
