package lru

// lru 缓存淘汰策略

import (
	"container/list"
	"geecache/utils"
)

type Cache struct {
	maxBytes int64                    //允许使用的最大内存
	nbytes   int64                    //当前已使用的内存
	ll       *list.List               //双向链表
	Cache    map[string]*list.Element //键是字符串，值是双向链表中对应节点的指针
	// optional and executed when an entry is purged.
	onEvicted func(key string, value Value) //某条记录被移除时的回调函数
}

// entry 双向链表节点的数据类型
// 淘汰队首节点时，需要用 key 从字典中删除对应的映射。
type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
// 为了通用性，我们允许值是实现了 Value 接口的任意类型，该接口只包含了一个方法 Len() int，用于返回值所占用的内存大小
type Value interface {
	Len() int
}

// NewCache is the Constructor of Cache
func NewCache(maxBytes int64, onEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		Cache:     make(map[string]*list.Element),
		onEvicted: onEvicted,
	}
}

//查找主要有 2 个步骤，第一步是从字典中找到对应的双向链表的节点，第二步，将该节点移动到队首

// Get look ups a key's value
// 如果键对应的链表节点存在，则将对应节点移动到队首，并返回查找到的值
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.Cache[key]; ok {
		c.ll.MoveToFront(ele)
		en := ele.Value.(*entry)
		utils.DPrintf("[Get]:%+v", en)
		return en.value, true
	}
	return
}

// RemoveOldest removes the oldest item
// 这里的删除，实际上是缓存淘汰。即移除最近最少访问的节点（队尾）
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		en := ele.Value.(*entry)
		delete(c.Cache, en.key)
		c.nbytes -= int64(len(en.key)) + int64(en.value.Len())
		utils.DPrintf("[RemoveOldest] delete element:%+v,nbytes:%d,maxBytes:%d", ele, c.nbytes, c.maxBytes)
		if c.onEvicted != nil {
			c.onEvicted(en.key, en.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.Cache[key]; ok {
		c.ll.MoveToFront(ele)
		en := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(en.value.Len())
		utils.DPrintf("[Add] update element:%+v,nbytes:%d,maxBytes:%d", ele, c.nbytes, c.maxBytes)
		en.value = value
	} else {
		ele := c.ll.PushFront(&entry{
			key:   key,
			value: value,
		})
		c.Cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
		utils.DPrintf("[Add] add element:%+v,nbytes:%d,maxBytes:%d", ele, c.nbytes, c.maxBytes)
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
