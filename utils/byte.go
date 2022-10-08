package utils

// Copy 复制 src 到新的 []byte 中并返回
func Copy(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}
