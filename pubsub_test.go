package rsmemory

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"

	"time"
)

const (
	KEY_TEST  string = "test_key"
	REDIS_URL string = "localhost:6379"
)

func TestPublish(t *testing.T) {
	t.Parallel()
	svc := New(&NewInput{
		RedisURL: REDIS_URL,
	})

	err := svc.Publish("test/foo", "bar")
	if err != nil {
		log.Fatal(err)
	}
}

func TestSubscribe(t *testing.T) {
	t.Parallel()
	svc := New(&NewInput{
		RedisURL: REDIS_URL,
	})
	channel := fmt.Sprintf("test/%s", time.Now().Add(10*time.Second).String())
	val := time.Now().String()

	reply := make(chan []byte)
	err := svc.Subscribe(channel, reply)
	if err != nil {
		log.Fatal(err)
	}

	err = svc.Publish(channel, val)
	if err != nil {
		log.Fatal(err)
	}

	t.Run("message", func(t *testing.T) {
		msg := <-reply
		if string(msg) != val {
			t.Fatal("expected correct reply message")
		}
		log.Printf("recieved %q", string(msg))
	})
}

func TestLua(t *testing.T) {
	// var luaScript string = `
	//     local result={}
	//     result[1] = redis.call('GET', KEYS[1])
	// 	result[2] = redis.call('TTL', KEYS[1])
	//     return result;
	//   `

	var luaScript string = `
        return redis.call('GET', KEYS[1])
	  `
	svc := New(&NewInput{
		RedisURL: REDIS_URL,
	})
	r, err := svc.ScriptEvaluate(luaScript, 1, "key")

	json_str, err := json.Marshal(r)

	fmt.Println("+++++", r, err, string(json_str))
	b := []byte{182}
	data := []interface{}{}
	json.Unmarshal(json_str, &data)
	fmt.Println("++++++", string(b), data)

}
func TestSetLua(t *testing.T) {
	var luaScript string = `
		redis.call('SET', KEYS[1], ARGV[1], 'EX',  ARGV[2])
		redis.call('PUBLISH', ARGV[3], ARGV[4])
	  `
	svc := New(&NewInput{
		RedisURL: REDIS_URL,
	})
	r, err := svc.ScriptEvaluate(luaScript, 4, "key1", "key2", "key3", "key4", "argv1", 300, "SYNC_CHANNEL_NAME", "SYNC_CHANNEL_NAME_VALUE")
	// data := []interface{}{}
	json_str, err := json.Marshal(r)
	// json.Unmarshal(r.([]byte), &data)
	fmt.Println("+++++", r, err, string(json_str))
}
func TestRand(t *testing.T) {
	r := xfetch(300, 200)
	t.Log(r)
}
