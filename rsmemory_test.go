package rsmemory

import (
	"fmt"
	"testing"
)

func TestMCache(t *testing.T) {
	key := fmt.Sprintf("App:Employment:%d", 1)
	multiLevelCache := NewRedisMultilevelCache(&Configrsmemory{})
	gameData := `NewRedisMultilevelCache`
	multiLevelCache.Set(key, gameData)
	data, _type := multiLevelCache.Get(key)

	fmt.Println("multiLevelCache get ", data, _type)

}
