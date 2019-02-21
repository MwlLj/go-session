package sqlite3

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

func (this *CSqlite3) Destory(id *string) error {
	return nil
}

func (this *CSqlite3) IsValid(id *string) (bool, error) {
	return true, nil
}

func (this *CSqlite3) Reset(id *string, timeoutS *int64) error {
	return nil
}

func New() *CSqlite3 {
	db := CSqlite3{}
	db.init()
	return &db
}
