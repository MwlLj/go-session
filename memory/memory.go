package memory

import (
	"github.com/satori/go.uuid"
	"sync"
	"time"
)

type CData struct {
	m_time     time.Time
	m_timeoutS int64
}

type CMemory struct {
	m_sessions map[string]CData
	m_mutex    sync.Mutex
}

func (this *CMemory) init() {
	this.m_sessions = make(map[string]CData)
}

func (this *CMemory) Dial(rule string) error {
	return nil
}

func (this *CMemory) Create(timeoutS int64) (id *string, e error) {
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}
	v4Uuid := uid.String()
	this.m_mutex.Lock()
	data := CData{
		m_time:     time.Now(),
		m_timeoutS: timeoutS,
	}
	this.m_sessions[v4Uuid] = data
	this.m_mutex.Unlock()
	return &v4Uuid, nil
}

func (this *CMemory) Destroy(id *string) error {
	return nil
}

func (this *CMemory) IsValid(id *string) (bool, error) {
	return true, nil
}

func (this *CMemory) Reset(id *string) error {
	return nil
}

func New() *CMemory {
	memory := CMemory{}
	memory.init()
	return &memory
}
