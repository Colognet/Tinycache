package geecache

//存储节点的真实值
type ByteView struct {
	b []byte
}

//返回占用内存的大小
func (v ByteView) Len() int {
	return len(v.b)
}

//返回一个拷贝，防止缓存的值被外部程序修改
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
