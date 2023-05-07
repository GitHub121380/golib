# extend 模块参照github项目https://github.com/thinkeridea/go-extend
主要是字符串操作避免频繁复制，节省性能等功能函数

# exbytes模块
主要包含函数模块
```
//有些场景源数据生命周期非常短，且可以原地替换，如果这么实现可以减少极大的内存分配。
func Replace(s, old, new []byte, n int) []byte {
}

// ToString 把 []byte 转换为 string 没有多余的内存开销。
func ToString(s []byte) string {
}
```

# exatomic：float类型的原子操作
主要实现了float类型的原子操作，使用较少，具体参见文件

# exstring模块

## convert模块
主要包含函数模块
```
//UnsafeToBytes 把 string 转换为 []byte 没有多余的内存开销。
func UnsafeToBytes(s string) []byte {
}
```
## bytes模块：输入string返回[]byte
主要包含函数模块
```
// UnsafeReplaceToBytes 替换字符串，并返回 []byte，减少类型转换
func UnsafeReplaceToBytes(s, old, new string, n int) []byte {
}

// ReplaceToBytes 替换字符串，并返回 []byte，减少类型转换
func ReplaceToBytes(s, old, new string, n int) []byte {
}

//RepeatToBytes 返回由字符串s的计数副本组成的 []byte。该方法是对标准库 strings.Repeat 修改，对于创建大字符串能有效减少内存分配。
//如果计数为负或 len(s) * count 溢出将触发panic。
func RepeatToBytes(s string, count int) []byte {
}

// JoinToBytes 使用 sep 连接 a 的字符串并返回 []byte
// 该方法是对标准库 strings.Join 修改，配合 unsafe 包能有效减少内存分配。
func JoinToBytes(a []string, sep string) []byte {
}
```
## string模块
主要包含函数模块
```
// 字符串buffer拼接，避免循环中使用次方法
func MultiJoinString(str ...string) string {
}

// Reverse 反转字符串
func Reverse(s string) string {
}

// 该方法是对标准库 strings.Replace 修改，配合 unsafe 包能有效减少内存分配。
func Replace(s, old, new string, n int) string {
}

/*Repeat 返回由字符串s的计数副本组成的新字符串。
该方法是对标准库 strings.Repeat 修改，对于创建大字符串能有效减少内存分配。
如果计数为负或 len(s) * count 溢出将触发panic。
*/
func Repeat(s string, count int) string {
}

// 该方法是对标准库 strings.Join 修改，配合 unsafe 包能有效减少内存分配。
func Join(a []string, sep string) string {
}
```

## unsafe模块
主要包含函数模块
```
/*
UnsafeRepeat 返回由字符串s的计数副本组成的新字符串。
该方法是对标准库 strings.Repeat 修改，对于创建大字符串能有效减少内存分配。
如果计数为负或 len(s) * count 溢出将触发panic。
与标准库的性能差异（接近标准库性能的两倍）：
	BenchmarkUnsafeRepeat-8            	   50000	     28003 ns/op	  303104 B/op	       1 allocs/op
	BenchmarkStandardLibraryRepeat-8   	   30000	     50619 ns/op	  606208 B/op	       2 allocs/op
Deprecated: 不在使用 Unsafe 前缀，保持与标准库相同的命名
*/
func UnsafeRepeat(s string, count int) string {
}

// UnsafeJoin 使用 sep 连接 a 的字符串。
// 该方法是对标准库 strings.Join 修改，配合 unsafe 包能有效减少内存分配。
// Deprecated: 不再使用 Unsafe 前缀，保持与标准库相同的命名
func UnsafeJoin(a []string, sep string) string {
}

// UnsafeReplace 替换字符串
// 该方法是对标准库 strings.Replace 修改，配合 unsafe 包能有效减少内存分配。
//
// Deprecated: 不再使用 Unsafe 前缀，保持与标准库相同的命名
func UnsafeReplace(s, old, new string, n int) string {
}
```

## pad模块，字符串填充
主要包含函数模块
```
// repeat 重复字符串到指定长度, []byte必须有充足的容量。
func repeat(b []byte, pad string, padLen int) {
}

/*
Pad 使用另一个字符串填充字符串为指定长度。
该函数返回 s 被从左端、右端或者同时两端被填充到制定长度后的结果。
填充方向由 falg 控制，可选值：PadLeft、PadBoth、PadRight。
在两边填充字符串为指定长度，如果补充长度是奇数，右边的字符会更多一些。
*/
func Pad(s, pad string, c, falg int) string {
}

// LeftPad 使用另一个字符串从左端填充字符串为指定长度。
func LeftPad(s, pad string, c int) string {
}

// RightPad 使用另一个字符串从右端填充字符串为指定长度。
func RightPad(s, pad string, c int) string {
}

// BothPad 使用另一个字符串从两端填充字符串为指定长度，
// 如果补充长度是奇数，右边的字符会更多一些。
func BothPad(s, pad string, c int) string {
}

/*
UnsafePad 使用另一个字符串填充字符串为指定长度。
该函数使用 unsafe 包转换数据类型，降低内存分配。
该函数返回 s 被从左端、右端或者同时两端被填充到制定长度后的结果。
填充方向由 falg 控制，可选值：PadLeft、PadBoth、PadRight。
在两边填充字符串为指定长度，如果补充长度是奇数，右边的字符会更多一些。
Deprecated: 不在使用 Unsafe 前缀
*/
func UnsafePad(s, pad string, c, falg int) string {
}

// UnsafeLeftPad 使用另一个字符串从左端填充字符串为指定长度。
// 该函数使用 unsafe 包转换数据类型，降低内存分配。
// Deprecated: 不再使用 Unsafe 前缀
func UnsafeLeftPad(s, pad string, c int) string {
}

// UnsafeBothPad 使用另一个字符串从两端填充字符串为指定长度，
// 如果补充长度是奇数，右边的字符会更多一些。
// 该函数使用 unsafe 包转换数据类型，降低内存分配。
// Deprecated: 不再使用 Unsafe 前缀
func UnsafeBothPad(s, pad string, c int) string {
```

## join_int模块，整数拼接，使用了bufferpool进行拼接
主要包含函数模块：
```
// JoinInts 使用 sep 连接 []int 并返回连接的字符串
func JoinInts(i []int, sep string) string {
}

// JoinInt8s 使用 sep 连接 []int8 并返回连接的字符串
func JoinInt8s(i []int8, sep string) string {
}

// JoinInt16s 使用 sep 连接 []int16 并返回连接的字符串
func JoinInt16s(i []int16, sep string) string {
}

// JoinInt32s 使用 sep 连接 []int32 并返回连接的字符串
func JoinInt32s(i []int32, sep string) string {
}

// JoinInt64s 使用 sep 连接 []int64 并返回连接的字符串
func JoinInt64s(i []int64, sep string) string {
}

// JoinUints 使用 sep 连接 []uint 并返回连接的字符串
func JoinUints(i []uint, sep string) string {
}

// JoinUint8s 使用 sep 连接 []uint8 并返回连接的字符串
func JoinUint8s(i []uint8, sep string) string {
}

// JoinUint16s 使用 sep 连接 []uint16 并返回连接的字符串
func JoinUint16s(i []uint16, sep string) string {
}

// JoinUint32s 使用 sep 连接 []uint32 并返回连接的字符串
func JoinUint32s(i []uint32, sep string) string {
}

// JoinUint64s 使用 sep 连接 []uint64 并返回连接的字符串
func JoinUint64s(i []uint64, sep string) string {
}

```