package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/quochungndvt/redis-sync-memory-cache/rsmemory"
)

type Person struct {
	ID   int64  `uri:"id" binding:"required"`
	Mess string `uri:"mess" binding:"required"`
}

var (
	port            = flag.String("port", "8080", "server port")
	multiLevelCache rsmemory.RedisMultilevelCache
)

func init() {
	config := &rsmemory.Configrsmemory{}
	multiLevelCache = rsmemory.NewRedisMultilevelCache(config)
}

func main() {
	flag.Parse()

	r := gin.Default()
	//in process memory cache
	r.GET("/read/:id/:mess", func(c *gin.Context) {
		var person Person
		if err := c.ShouldBindUri(&person); err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}

		key := fmt.Sprintf("test-sync-memory-cache:%d", person.ID)
		data, t := multiLevelCache.Get(key)

		c.JSON(200, gin.H{"mess": person.Mess, "id": person.ID, "cache": data, "action": t})
		return

	})
	//in redis cache
	r.GET("/read-redis/:id/:mess", func(c *gin.Context) {
		var person Person
		if err := c.ShouldBindUri(&person); err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}

		key := fmt.Sprintf("test-sync-memory-cache:%d", person.ID)
		//get data from redis
		data, err := multiLevelCache.GetRedisService().GetCacheFromRedis(key)

		c.JSON(200, gin.H{"mess": person.Mess, "id": person.ID, "cache": data, "action": "READ_FROM_REDIS", "error": err})
	})
	r.GET("/write/:id/:mess", func(c *gin.Context) {
		var person Person
		if err := c.ShouldBindUri(&person); err != nil {
			c.JSON(400, gin.H{"msg": err})
			return
		}
		key := fmt.Sprintf("test-sync-memory-cache:%d", person.ID)
		multiLevelCache.Set(key, person.Mess)
		// data, t := multiLevelCache.Get(key)
		c.JSON(200, gin.H{"mess": person.Mess, "id": person.ID, "cache": person.Mess})
	})
	r.GET("/test", func(c *gin.Context) {
		t := do_some_thing()
		c.JSON(200, gin.H{"t": t})
	})
	r.Run(fmt.Sprintf("0.0.0.0:%s", *port)) // listen and serve on 0.0.0.0:8080
}
func do_some_thing() string {
	messages := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		messages <- "ping"
	}()

	//sleep 1
	go goroutine(1)
	//sleep 2
	go goroutine(2)
	//sleep 3
	go goroutine(3)
	//sleep 4
	go goroutine(4)
	msg := <-messages
	fmt.Println(msg)
	return msg
}
func goroutine(t int) {
	time.Sleep(time.Duration(t) * time.Second)
	fmt.Println("pong", t, time.Now().Nanosecond())
}
