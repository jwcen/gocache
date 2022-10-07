package main

import (
	"gocache/cache"
	"time"
)

func main() {
	defaultExpiration, _ := time.ParseDuration("0.5h") // 默认
	gcInterval, _ := time.ParseDuration("3s")          // 过期清理周期为3s
	c := cache.NewCache(defaultExpiration, gcInterval)

	k1 := "hello world"
	expiration, _ := time.ParseDuration("5s")

	c.Set("k1", k1, expiration)
	// s, _ := time.ParseDuration("10s")
	// if v, found := c.Get("k1"); found {
	// 	fmt.Println("Found k1:", v)
	// } else {
	// 	fmt.Println("Not found k1")
	// }

	// // 暂停10s，验证k1是否被清理掉
	// time.Sleep(s)
	// if v, found := c.Get("k1"); found {
	// 	fmt.Println("Found k1:", v)
	// } else {
	// 	fmt.Println("Not found k1")
	// }
	err := c.SaveToFile("./cache.txt")
	if err != nil {
		return
	}

	err = c.LoadFile("./cache.txt")
	if err != nil {
		return
	}
	

}
