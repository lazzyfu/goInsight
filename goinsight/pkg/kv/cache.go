/*
@Time    :   2022/07/06 10:12:48
@Author  :   xff
@Desc    :   None
*/

package kv

import (
	"sync"
)

// Cache 存储单次请求的上下文信息
type Cache struct {
	value interface{}
}

// KVCache 存储多个Cache实例的集合
type KVCache struct {
	sync.RWMutex //继承读写锁，用于并发控制
	Id           string
	Items        map[string]*Cache // K-V存储
}

// Put 写入
func (c *KVCache) Put(key string, value interface{}) {
	c.Lock()
	defer c.Unlock()
	c.Items[key] = &Cache{value: value}
}

// Get 查询
func (c *KVCache) Get(key string) interface{} {
	c.RLock()
	defer c.RUnlock()
	if item, ok := c.Items[key]; ok {
		return item.value
	}
	return nil
}

// Delete 删除
func (c *KVCache) Delete(key string) {
	c.Lock()
	defer c.Unlock()
	delete(c.Items, key)
}

// NewKVCache 新建缓存
func NewKVCache(Id string) *KVCache {
	return &KVCache{Id: Id, Items: make(map[string]*Cache)}
}
