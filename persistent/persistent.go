package persistent

type CPersistent struct {
}

func (this *CPersistent) Dial(rule string) error {
	return nil
}

func (this *CPersistent) Create(timeoutS int64) (id *string, e error) {
	return nil, nil
}

func (this *CPersistent) Destroy(id *string) error {
	return nil
}

func (this *CPersistent) IsValid(id *string) (bool, error) {
	return true, nil
}

func (this *CPersistent) Reset(id *string) error {
	return nil
}

func New(dbType string) *CPersistent {
	db := CPersistent{}
	return &db
}
