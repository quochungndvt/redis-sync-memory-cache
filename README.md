# redis-sync-memory-cache
Using redis to synchronizing memory caches

Classic Redis use case with lots of advantages and incredibly fast but as we know local RAM access is still many times faster than network I/O
BTW memory caches in decentralized system have some challenges
1. Data consistency
2. Data lag between the memory caches and Redis
2. Don’t blow up the network!

This module implement Techniques for Synchronizing In-Memory Caches with Redis

### Documentation
-------------

- [API Reference](https://godoc.org/github.com/quochungndvt/redis-sync-memory-cache/rsmemory)

### Installation
-------------

Install using the "go get" command:

  go get github.com/quochungndvt/redis-sync-memory-cache/rsmemory
  
  
### Examples
-------------

try this [example](https://github.com/quochungndvt/redis-sync-memory-cache/tree/master/examples/server)
- server1: go run main.go -port 8081
- server2: go run main.go -port 8082

Open in browser 

- A: http://localhost:8081/read/1/1
- B: http://localhost:8082/read/1/1
- Update cache http://localhost:8081/write/1/update
refresh A, B and see


### Feature
-------------
- Support more Redis type
- Support cache expire
