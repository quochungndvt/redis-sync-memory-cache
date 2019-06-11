package rsmemory

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

type Configrsmemory struct {
	RedisURL          string
	HashSlot          uint16
	SYNC_CHANNEL_NAME string
}
type RedisMultilevelCache interface {
	Get(key string) (interface{}, string)
	Set(key string, data interface{})
	GetRedisService() *Service
	Size() int
	init(config *Configrsmemory)
}
type redisMultilevelCacheImpl struct {
	Mutex       sync.RWMutex
	cacheMemory InProcessCache
	lastUpdated LastUpdatedDictionary
	hslt        HashSlotCalculator
	svc         *Service
	config      *Configrsmemory
	close       chan bool
}

func mergeConfig(config *Configrsmemory) {
	if config.RedisURL == "" {
		config.RedisURL = DEFAULT_REDIS_URI
	}
	if config.SYNC_CHANNEL_NAME == "" {
		config.SYNC_CHANNEL_NAME = SYNC_CHANNEL_NAME
	}
}
func (t *redisMultilevelCacheImpl) init(config *Configrsmemory) {
	t.config = config
	//FIXME:
	mergeConfig(t.config)
	t.cacheMemory = NewInProcessCache()
	t.lastUpdated = NewLastUpdatedDictionary()
	t.svc = New(&NewInput{
		RedisURL: t.config.RedisURL,
	})
	t.hslt = NewHashSlotCalculator(t.config.HashSlot)
	t.close = make(chan bool)
	go t.Sub()

}
func (t *redisMultilevelCacheImpl) Size() int {
	return t.cacheMemory.Size()
}
func (t *redisMultilevelCacheImpl) Sub() {
	reply := make(chan []byte)
	err := t.svc.Subscribe(t.config.SYNC_CHANNEL_NAME, reply)
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case msg := <-reply:
			log.Printf("[Sub] recieved %q", string(msg))
			fmt.Printf("#goroutines: %d\n", runtime.NumGoroutine())
			//TODO: detect pub in current sv and don't update
			hashSlot := t.CalculateHashSlot(string(msg))
			t.SetLastUpdatedTime(hashSlot, t.GetCurrentTime())
		case <-t.close:
			break
		default:
		}
	}
}
func (t *redisMultilevelCacheImpl) Pub(data string) {
	t.svc.Publish(t.config.SYNC_CHANNEL_NAME, data)
}
func (t *redisMultilevelCacheImpl) CalculateHashSlot(key string) uint16 {
	return t.hslt.CalculateHashSlot(key)
}
func (t *redisMultilevelCacheImpl) GetCurrentTime() int64 {
	return time.Now().Unix() - 1
}
func (t *redisMultilevelCacheImpl) GetMemoryCache(key string) *Cache {
	return t.cacheMemory.Get(key)
}
func (t *redisMultilevelCacheImpl) SetMemoryCache(key string, data *Cache) {
	t.cacheMemory.Set(key, data)
}
func (t *redisMultilevelCacheImpl) GetCacheFromRedis(key string) (Cache, error) {
	redisData, err := t.svc.GetCacheFromRedis(key)
	return redisData, err
}
func (t *redisMultilevelCacheImpl) SaveCacheToRedis(key string, data Cache) {
	t.svc.SaveCacheToRedis(key, data)

}
func (t *redisMultilevelCacheImpl) GetLastUpdatedTime(hashSlot uint16) int64 {
	return t.lastUpdated.Get(hashSlot)
}
func (t *redisMultilevelCacheImpl) SetLastUpdatedTime(hashSlot uint16, currentTime int64) {
	t.lastUpdated.Set(hashSlot, currentTime)
}
func (t *redisMultilevelCacheImpl) Get(key string) (result interface{}, _type string) {
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()
	hashSlot := t.CalculateHashSlot(key)
	currentTime := t.GetCurrentTime()

	//get data from memory cache
	inProcessCacheEntry := t.GetMemoryCache(key)
	//if has data in memory cache
	if inProcessCacheEntry != nil {
		//lastUpdated dictionary timestamp greater than the cache entry timestamp?
		// lastUpdatedTime, ok := lastUpdated[data.HashSlot]
		lastUpdatedTime := t.GetLastUpdatedTime(inProcessCacheEntry.HashSlot)
		hashSlot = inProcessCacheEntry.HashSlot
		if lastUpdatedTime > inProcessCacheEntry.TimeStamp {
			//get data from redis
			redisData, err := t.GetCacheFromRedis(key)
			if err == nil {
				//update time in lastUpdated
				// lastUpdated[data.HashSlot] = currentTime
				t.SetLastUpdatedTime(hashSlot, currentTime)
				//save data to cacheMemory
				data := &Cache{HashSlot: hashSlot, TimeStamp: currentTime, Value: redisData.Value}
				t.SetMemoryCache(key, data)
				return data.Value, OUTDATE_READ_FROM_REDIS
			}
			return nil, READ_DEFAULT

		}
		//return data to client

		return inProcessCacheEntry.Value, READ_FROM_MEMORY

	}
	//else has no data in memory cache
	//get data from redis
	redisData, err := t.GetCacheFromRedis(key)
	if err == nil {
		//update time in lastUpdated
		t.SetLastUpdatedTime(hashSlot, currentTime)
		//save data to cacheMemory
		data := &Cache{HashSlot: hashSlot, TimeStamp: currentTime, Value: redisData.Value}
		t.SetMemoryCache(key, data)
		return data.Value, READ_FROM_REDIS
	}
	return nil, READ_DEFAULT

}
func (t *redisMultilevelCacheImpl) Set(key string, data interface{}) {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	hashSlot := t.CalculateHashSlot(key)
	currentTime := t.GetCurrentTime()

	lastUpdatedTime := t.GetLastUpdatedTime(hashSlot)
	dataCache := &Cache{HashSlot: hashSlot, TimeStamp: currentTime, Value: data}
	if lastUpdatedTime > 0 {
		//save data to cacheMemory
		t.SetMemoryCache(key, dataCache)
		//save data to redis
		//TODO: using lua script to batch TCP
		t.SaveCacheToRedis(key, *dataCache)
		//Publish update message to all clients
		t.Pub(key)
		//TODO: detect pub in current sv and don't update
	} else {
		//update time in lastUpdated
		t.SetLastUpdatedTime(hashSlot, currentTime)
		t.SetMemoryCache(key, dataCache)
		//save data to redis
		//TODO: using lua script to batch TCP
		t.SaveCacheToRedis(key, *dataCache)
		//Publish update message to all clients
		t.Pub(key)
		//TODO: detect pub in current sv and don't update
	}
}
func (t *redisMultilevelCacheImpl) GetRedisService() *Service {
	return t.svc
}
func NewRedisMultilevelCache(config *Configrsmemory) RedisMultilevelCache {
	cache := new(redisMultilevelCacheImpl)
	cache.init(config)
	return cache
}
