package main

import (
	"gocache/caches"
	"gocache/servers"
)

func main() {
	cache := caches.NewCache()
	err := servers.NewHTTPServer(cache).Run(":8888")
	if err != nil {
		panic(err)
	}
}
