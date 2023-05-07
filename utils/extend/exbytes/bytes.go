package exbytes

import (
	"bytes"
)

/*
Replace 思路来源于 bytes.Replace，bytes.Replace 总是返回 s 的副本，
有些场景源数据生命周期非常短，且可以原地替换，如果这么实现可以减少极大的内存分配。
len(old) >= len(new) 会执行原地替换，这回浪费一部分空间，但是会减少内存分配，
建议输入生命周期较短的数据， len(old) < len(new) 会调用 bytes.Replace 并返回一个替换后的副本。
最佳实践是使用 Replace 的结果覆盖源变量，避免再次对源数据引用，导致访问过时的数据，并且数据内容错乱，如下：

	var s []byte
	s = exbytes.Replace(s, []byte(" "), []byte(""), -1)

关于字符串可以结合 exstrings.UnsafeToBytes 来实现，要避免常量字符串和字面量字符串，否者会产生运行时错误。
*/
func Replace(s, old, new []byte, n int) []byte {
	if n == 0 {
		return s
	}

	if len(old) < len(new) {
		return bytes.Replace(s, old, new, n)
	}

	if n < 0 {
		n = len(s)
	}

	var wid, i, j, w int
	for i, j = 0, 0; i < len(s) && j < n; j++ {
		wid = bytes.Index(s[i:], old)
		if wid < 0 {
			break
		}

		w += copy(s[w:], s[i:i+wid])
		w += copy(s[w:], new)
		i += wid + len(old)
	}

	w += copy(s[w:], s[i:])
	return s[0:w]
}
