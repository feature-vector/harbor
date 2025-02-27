package keylock

import (
	"fmt"
	"sync"
)

var (
	memoryLocker = &memoryLockerImpl{
		locksMap: map[string]*lockVal{},
	}
)

func Memory() KeyLocker {
	return memoryLocker
}

type memoryLockerImpl struct {
	locksMap map[string]*lockVal
	gLock    sync.Mutex
}

type lockVal struct {
	m   *sync.Mutex
	cnt int
}

func (m *memoryLockerImpl) Lock(key string) {
	m.gLock.Lock()
	lock, ok := m.locksMap[key]
	if !ok {
		lock = &lockVal{
			m:   &sync.Mutex{},
			cnt: 1,
		}
		m.locksMap[key] = lock
	} else {
		lock.cnt++
	}
	m.gLock.Unlock()
	lock.m.Lock()
}

func (m *memoryLockerImpl) Unlock(key string) {
	m.gLock.Lock()
	lock, ok := m.locksMap[key]
	if !ok {
		panic(fmt.Sprintf("no lock found with key [%s]", key))
	}
	lock.cnt--
	if lock.cnt == 0 {
		delete(m.locksMap, key)
	}
	m.gLock.Unlock()
	lock.m.Unlock()
}

func (m *memoryLockerImpl) TryLock(key string) bool {
	m.gLock.Lock()
	lock, ok := m.locksMap[key]
	if ok {
		m.gLock.Unlock()
		return false
	}
	lock = &lockVal{
		m:   &sync.Mutex{},
		cnt: 1,
	}
	m.locksMap[key] = lock
	m.gLock.Unlock()
	lock.m.Lock()
	return true
}
