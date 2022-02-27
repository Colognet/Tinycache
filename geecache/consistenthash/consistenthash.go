package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

//定义函数类型，采取依赖注入的方式，允许用于替换成自定义的 Hash 函数
type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int                               //虚拟节点倍数
	keys     []int                             //哈希环
	hashMap  map[int]string                    //虚拟节点与真实节点的映射表
}

//自定义虚拟节点的倍数、哈希函数
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE             //默认哈希函数
	}
	return m
}

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		//每个节点创建replicas个虚拟节点
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))      //虚拟节点名称：strconv.Itoa(i) + key，通过添加编号来区分不同的虚拟节点
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key                                   //虚拟节点与真实节点映射
		}
	}
	sort.Ints(m.keys)                                               //环上的哈希值排序
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	//计算key的哈希值
	hash := int(m.hash([]byte(key)))
	//顺时针找到第一个匹配的虚拟节点的下标 idx，从 m.keys 中获取到对应的哈希值。
	//如果 idx == len(m.keys)，说明应选择 m.keys[0]，因为 m.keys 是一个环状结构，所以用取余数的方式来处理这种情况。
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	//通过hashMap 映射得到真实的节点。
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
