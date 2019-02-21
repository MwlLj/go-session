package memory

import (
	"errors"
	"github.com/satori/go.uuid"
	"log"
	"sync"
	"time"
)

type CData struct {
	m_time     time.Time
	m_timeoutS int64
}

type CMemory struct {
	m_sessions map[string]CData
	m_mutex    sync.RWMutex
}

func (this *CMemory) init() {
	this.m_sessions = make(map[string]CData)
	this.check()
}

func (this *CMemory) check() {
	go func() {
		for {
			this.m_mutex.RLock()
			for id, data := range this.m_sessions {
				if data.m_time.Add(time.Duration(data.m_timeoutS) * time.Second).Before(time.Now()) {
					// timeout
					this.m_mutex.RUnlock()
					this.m_mutex.Lock()
					delete(this.m_sessions, id)
					this.m_mutex.Unlock()
					this.m_mutex.RLock()
					log.Printf("timeout, sessionid = %s\n", id)
				}
			}
			this.m_mutex.RUnlock()
			time.Sleep(1 * time.Second)
		}
	}()
}

func (this *CMemory) Dial(rule string) error {
	return nil
}

func (this *CMemory) Create(timeoutS int64) (id *string, e error) {
	uid, err := uuid.NewV4()
	if err != nil {
		log.Println("create uuid error")
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

func (this *CMemory) Destory(id *string) error {
	if id == nil {
		return errors.New("id is nil")
	}
	this.m_mutex.RLock()
	if _, ok := this.m_sessions[*id]; ok {
		this.m_mutex.RUnlock()
		this.m_mutex.Lock()
		delete(this.m_sessions, *id)
		this.m_mutex.Unlock()
		this.m_mutex.RLock()
		log.Printf("destroy, sessionid = %s\n", *id)
	}
	this.m_mutex.RUnlock()
	return nil
}

func (this *CMemory) IsValid(id *string) (bool, error) {
	if id == nil {
		return false, errors.New("id is nil")
	}
	this.m_mutex.RLock()
	defer this.m_mutex.RUnlock()
	if _, ok := this.m_sessions[*id]; ok {
		return true, nil
	}
	return false, nil
}

func (this *CMemory) Reset(id *string, timeoutS *int64) error {
	if id == nil {
		return errors.New("id is nil")
	}
	this.m_mutex.RLock()
	if d, ok := this.m_sessions[*id]; ok {
		this.m_mutex.RUnlock()
		this.m_mutex.Lock()
		timeout := timeoutS
		if timeout == nil {
			timeout = &d.m_timeoutS
		}
		data := CData{
			m_time:     time.Now(),
			m_timeoutS: *timeout,
		}
		this.m_sessions[*id] = data
		this.m_mutex.Unlock()
		this.m_mutex.RLock()
		log.Printf("reset, sessionid = %s\n", *id)
	}
	this.m_mutex.RUnlock()
	return nil
}

func New() *CMemory {
	memory := CMemory{}
	memory.init()
	return &memory
}
