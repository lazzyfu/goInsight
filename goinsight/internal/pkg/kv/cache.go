/*
@Time    :   2022/07/06 10:12:48
@Author  :   xff
@Desc    :   None
*/

package kv

import (
	"sync"
)

type Cache struct {
	value interface{}
}

type KVCache struct {
	sync.RWMutex //继承读写锁，用于并发控制
	Id           string
	Items        map[string]*Cache // K-V存储
}

// 写入
func (c *KVCache) Put(key string, value interface{}) {
	c.Lock()
	defer c.Unlock()
	c.Items[key] = &Cache{value: value}
}

// 查询
func (c *KVCache) Get(key string) interface{} {
	c.RLock()
	defer c.RUnlock()
	if item, ok := c.Items[key]; ok {
		return item.value
	}
	return nil
}

// 新建缓存
func NewKVCache(Id string) *KVCache {
	return &KVCache{Id: Id, Items: map[string]*Cache{}}
}
