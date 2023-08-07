package src

type ByteView struct {
	b []byte
}

func (bv ByteView) String() string {
	return string(bv.b)
}

func (bv ByteView) Len() int {
	return len(bv.b)
}

func (bv ByteView) ByteSlice() []byte {
	return cloneBytes(bv.b)
}

func cloneBytes(b []byte) []byte {
	res := make([]byte, len(b))
	copy(res, b)
	return res
}
