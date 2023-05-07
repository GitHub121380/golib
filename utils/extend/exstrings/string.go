package exstrings

import (
	"bytes"

	"github.com/GitHub121380/golib/utils/extend/exbytes"
)

// var joinStringBufferPoll = bufferpool.NewBufferPool(512)

// // 字符串buffer拼接，使用缓存池方式，避免循环中使用次方法
// func MultiJoinStringByBufferPool(str ...string) string {
// 	buf := joinBufferPoll.Get()
// 	length := len(str)
// 	for i := 0; i < length; i++ {
// 		buf.WriteString(str[i])
// 	}
// 	defer joinStringBufferPoll.Put(buf)
// 	return buf.String()
// }

// 字符串buffer拼接，避免循环中使用次方法
func MultiJoinString(str ...string) string {
	var b bytes.Buffer
	for i := 0; i < len(str); i++ {
		b.WriteString(str[i])
	}
	return b.String()
}

// Reverse 反转字符串，通过 https://golang.org/doc/code.html#Library 收集
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

// Replace 替换字符串
// 该方法是对标准库 strings.Replace 修改，配合 unsafe 包能有效减少内存分配。
func Replace(s, old, new string, n int) string {
	return exbytes.ToString(UnsafeReplaceToBytes(s, old, new, n))
}

/*
Repeat 返回由字符串s的计数副本组成的新字符串。
该方法是对标准库 strings.Repeat 修改，对于创建大字符串能有效减少内存分配。
如果计数为负或 len(s) * count 溢出将触发panic。
*/
func Repeat(s string, count int) string {
	return exbytes.ToString(RepeatToBytes(s, count))
}

// Join 使用 sep 连接 a 的字符串。
// 该方法是对标准库 strings.Join 修改，配合 unsafe 包能有效减少内存分配。
func Join(a []string, sep string) string {
	switch len(a) {
	case 0:
		return ""
	case 1:
		return a[0]
	case 2:
		// Special case for common small values.
		// Remove if golang.org/issue/6714 is fixed
		return a[0] + sep + a[1]
	case 3:
		// Special case for common small values.
		// Remove if golang.org/issue/6714 is fixed
		return a[0] + sep + a[1] + sep + a[2]
	}

	return exbytes.ToString(JoinToBytes(a, sep))
}
