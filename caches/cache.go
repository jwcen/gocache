// 用于存放有关 cache 底层结构的代码
package caches

import (
	"gocache/utils"
	"sync"
)

// Cache 是一个结构体，用于封装缓存底层结构
type Cache struct {
	// data 是一个map，存储了所有的数据
	// value 类型使用[]byte，以便网络传输
	data map[string][]byte

	// count 记录data中键值对的个数
	// 这是一个冗余设计，直接使用len(data)就行
	// 使用count记录是为了更快得到结果
	count int64

	// lock 用于保证并发安全
	lock *sync.RWMutex
}

// NewCache 返回一个缓存对象
func NewCache() *Cache {
	return &Cache{
		// 预先分配256个槽位，避免后续因容量不足导致map扩容
		// 扩容会分配内存，影响性能；而且槽位少了，哈希冲突几率就大，map查找性能下降
		// 256 并非最佳值，需根据实际情况而定
		data:  make(map[string][]byte, 256),
		count: 0,
		lock:  &sync.RWMutex{},
	}
}

// Set 保存 key 和 value 到缓存中
func (c *Cache) Set(key string, value []byte) {
	// Set 操作会改变数据的状态，需要保证串行执行，故使用写锁
	c.lock.Lock()
	defer c.lock.Unlock()
	// 查询是否已经存在该元素, 不存在则计数++
	if _, ok := c.data[key]; !ok {
		c.count++
	}
	// 该 Copy 方法会将 value 拷贝一份
	// 这样即使传进来的 value 被修改或者清空了也不会影响缓存里面的数据
	c.data[key] = utils.Copy(value)
}

// Get 返回指定的 key 的 value， 如果找不到则返回 false
func (c *Cache) Get(key string) ([]byte, bool) {
	// 查询数据不会改变数据的状态，故可并发执行。
	// 使用读锁，加快读取速度
	c.lock.RLock()
	defer c.lock.RUnlock()
	value, ok := c.data[key]
	return value, ok
}

// Delete 删除指定 key 的键值对数据
func (c *Cache) Delete(key string) {
	// Delete 操作会改变数据状态，需要保证串行执行，使用写锁
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, ok := c.data[key]; ok {
		c.count--
		delete(c.data, key)
	}
}

// Count 返回键值对数据的个数
func (c *Cache) Count() int64 {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.count
}
