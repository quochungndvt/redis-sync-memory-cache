package rsmemory

const (
	SYNC_CHANNEL_NAME       string = "RedisMultilevelCache_Sync"
	READ_FROM_MEMORY        string = "READ_FROM_MEMORY"
	READ_FROM_REDIS         string = "READ_FROM_REDIS"
	OUTDATE_READ_FROM_REDIS string = "OUTDATE_READ_FROM_REDIS"
	READ_DEFAULT            string = "READ_DEFAULT"

	DEFAULT_REDIS_URI string = "localhost:6379"
)
