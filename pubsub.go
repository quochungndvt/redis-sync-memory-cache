package rsmemory

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/gomodule/redigo/redis"
)

//
func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

//service redis
type Service struct {
	pool *redis.Pool
	conn redis.Conn
}

//new input for constructor
type NewInput struct {
	RedisURL string
}

//new return service
func New(input *NewInput) *Service {
	if input == nil {
		log.Fatal("input is required")
	}
	var redispool *redis.Pool
	redispool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", input.RedisURL)
		},
	}

	// Get a connection
	conn := redispool.Get()
	defer conn.Close()
	// Test the connection
	_, err := conn.Do("PING")
	if err != nil {
		log.Fatalf("can't connect to the redis database, got error:\n%v", err)
	}

	return &Service{
		pool: redispool,
		conn: conn,
	}
}

// Publish publish key value
func (s *Service) Publish(key string, value string) error {
	conn := s.pool.Get()
	defer conn.Close()
	conn.Do("PUBLISH", key, value)
	return nil
}

// Subscribe subscribe
func (s *Service) Subscribe(key string, msg chan []byte) error {
	rc := s.pool.Get()
	psc := redis.PubSubConn{Conn: rc}
	if err := psc.PSubscribe(key); err != nil {
		return err
	}
	go func() {
		for {
			switch v := psc.Receive().(type) {
			case redis.Message:
				msg <- v.Data
			case redis.Subscription:
				fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
			case error:
				fmt.Printf("Error: %v", v)
			}

		}
	}()
	return nil
}

func (s *Service) Do(comman string, args ...interface{}) (reply interface{}, err error) {
	conn := s.pool.Get()
	defer conn.Close()
	reply, err = conn.Do(comman, args...)
	return
}

/*
var luaScript = @"
	local result={}
	result[1] = redis.call('GET', KEYS[1])
	result[2] = redis.call('TTL', KEYS[1])
	return result;
";
*/
func (s *Service) ScriptEvaluate(luaScript string, num_args int, args ...interface{}) (reply interface{}, err error) {
	conn := s.pool.Get()
	defer conn.Close()
	var l = len(args)
	if num_args > 0 {
		l = num_args
	}
	var getScript = redis.NewScript(l, luaScript)

	reply, err = getScript.Do(conn, args...)
	return
}
func (s *Service) Get(key string) (reply string, err error) {
	conn := s.pool.Get()
	defer conn.Close()
	// conn.Do("GET", key)
	reply, err = redis.String(conn.Do("GET", key))
	return
}
func (s *Service) Set(key string, value string) (reply interface{}, err error) {
	conn := s.pool.Get()
	defer conn.Close()
	reply, err = conn.Do("SET", key, value)
	return
}
func (s *Service) SaveCacheToRedis(key string, data Cache) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return
	}
	s.Set(key, string(dataBytes))
}
func (s *Service) GetCacheFromRedis(key string) (data Cache, err error) {
	var result string
	result, err = s.Get(key)
	if err != nil {
		log.Println("[getCacheFromRedis]", err)
		return
	}
	err = json.Unmarshal([]byte(result), &data)
	return
}

//TODO: implement xfetch
// func fetch(name string){
// 	var data,delta,ttl = redis.get(name, delta, ttl)
// 	if (!data or xfetch(delta, time() + ttl))
// 	var data,recompute_time = recompute(name)
// 	redis.set(name, expires, data), redis.set(delta, expires, recompute_time)
// 	return data
// }

//delta – Time to recompute value
//beta – control (default: 1.0, > 1.0 favors earlier recomputation, < 1.0 favors later)
//rand – Random number [ 0.0 … 1.0 ]
func xfetch(delta, expiry int) bool {
	var BETA = 2.0
	/* XFetch is negative; value is being added to time() */
	return float64(time.Now().Unix())-(float64(delta)*float64(BETA)*math.Log(randFloats(0, 1))) >= float64(expiry)
}
func randFloats(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
