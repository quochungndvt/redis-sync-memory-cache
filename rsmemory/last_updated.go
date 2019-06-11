package rsmemory

import "sync"

type LastUpdatedDictionary interface {
	Get(key uint16) int64
	Set(key uint16, data int64)
	init()
}
type lastUpdatedDictionaryImpl struct {
	sync.Mutex
	lastUpdated map[uint16]int64
}

func (this *lastUpdatedDictionaryImpl) init() {
	this.lastUpdated = make(map[uint16]int64)
}
func (this *lastUpdatedDictionaryImpl) Get(key uint16) (data int64) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	if d, ok := this.lastUpdated[key]; ok {
		data = d
	}
	return
}
func (this *lastUpdatedDictionaryImpl) Set(key uint16, data int64) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	this.lastUpdated[key] = data
	return
}
func NewLastUpdatedDictionary() LastUpdatedDictionary {
	lup := new(lastUpdatedDictionaryImpl)
	lup.init()
	return lup
}
