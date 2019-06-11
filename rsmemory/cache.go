package rsmemory

import (
	"encoding/json"
	"fmt"
	"sync"
	"unsafe"
)

type Cache struct {
	HashSlot  uint16
	TimeStamp int64
	Value     interface{}
}
type InProcessCache interface {
	Get(key string) *Cache
	Set(key string, data *Cache)
	Size() int
	init()
}
type inProcessCacheImpl struct {
	sync.Mutex
	cache map[string]*Cache
}

func (this *inProcessCacheImpl) init() {
	this.cache = make(map[string]*Cache)
}

//Just test
func (this *inProcessCacheImpl) Size() int {
	size := int(unsafe.Sizeof(*this))
	//FIXME:
	for _, cache := range this.cache {
		// size_cache := reflect.TypeOf(cache.Value).Elem().Size()
		json_str, _ := json.Marshal(cache.Value)
		size_cache := len(string(json_str))
		size += int(size_cache)
		fmt.Println("+++++++++++++", size_cache)
	}
	return size
}

//get data from memory cache
//
//need lock & unlock cache
func (this *inProcessCacheImpl) Get(key string) (data *Cache) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	if d, ok := this.cache[key]; ok {
		data = d
	}
	return
}

//set data to memory cache
//
//need lock & unlock cache
func (this *inProcessCacheImpl) Set(key string, data *Cache) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	this.cache[key] = data
	return
}

//
//init new memory cache
//
func NewInProcessCache() InProcessCache {
	cache := new(inProcessCacheImpl)
	cache.init()
	return cache
}
