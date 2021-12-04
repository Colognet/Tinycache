package lru

import "container/list"

// Cache is a LRU cache. It is not safe for concurrent access.
//数据结构：字典+双向链表
//字典用来查找key对应的val（o1），双向链表用来维护访问的先后关系
type Cache struct {
	maxBytes int64                    //允许使用的最大内存
	nbytes   int64                    //当前使用的内存
	ll       *list.List               //双向链表
	cache    map[string]*list.Element //字典：key字符串，val双向链表对应节点的指针
	// optional and executed when an entry is purged.
	//淘汰队首节点时，需要用 key 从字典中删除对应的映射。
	OnEvicted func(key string, value Value)
}

//双向链表节点的数据类型
type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
//返回值占用内存的大小
type Value interface {
	Len() int
}

// New is the Constructor of Cache
//实例化Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Add adds a value to the cache.
//向Cache中增加值
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		//有key，更新
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		//无key增加
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	//超出限制 删除
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Get look ups a key's value
//查找功能
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// RemoveOldest removes the oldest item
//缓存淘汰
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
