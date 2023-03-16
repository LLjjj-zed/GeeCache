package geecache

// 缓存值的抽象与封装

// A ByteView holds an immutable view of bytes.
// ByteView 只有一个数据成员，b []byte，b 将会存储真实的缓存值。
// 选择 byte 类型是为了能够支持任意的数据类型的存储，例如字符串、图片等。
type ByteView struct {
	b []byte
}

// Len returns the view's length
func (v ByteView) Len() int {
	return len(v.b)
}

// ByteSlice returns a copy of the data as a byte slice.
// b 是只读的，使用 ByteSlice() 方法返回一个拷贝，防止缓存值被外部程序修改
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

// String returns the data as a string, making a copy if necessary.
func (v ByteView) String() string {
	return string(v.b)
}
