package exstrings

import (
	"github.com/GitHub121380/golib/utils"
	"strconv"
)

var joinBufferPoll = utils.NewBufferPool(64)

// JoinInts 使用 sep 连接 []int 并返回连接的字符串
func JoinInts(i []int, sep string) string {
	buf := joinBufferPoll.Get()
	for _, v := range i {
		if buf.Len() > 0 {
			buf.WriteString(sep)
		}
		buf.WriteString(strconv.FormatInt(int64(v), 10))
	}

	defer joinBufferPoll.Put(buf)
	return buf.String()
}

// JoinInt8s 使用 sep 连接 []int8 并返回连接的字符串
func JoinInt8s(i []int8, sep string) string {
	buf := joinBufferPoll.Get()
	for _, v := range i {
		if buf.Len() > 0 {
			buf.WriteString(sep)
		}
		buf.WriteString(strconv.FormatInt(int64(v), 10))
	}

	defer joinBufferPoll.Put(buf)
	return buf.String()
}

// JoinInt16s 使用 sep 连接 []int16 并返回连接的字符串
func JoinInt16s(i []int16, sep string) string {
	buf := joinBufferPoll.Get()
	for _, v := range i {
		if buf.Len() > 0 {
			buf.WriteString(sep)
		}
		buf.WriteString(strconv.FormatInt(int64(v), 10))
	}

	defer joinBufferPoll.Put(buf)
	return buf.String()
}

// JoinInt32s 使用 sep 连接 []int32 并返回连接的字符串
func JoinInt32s(i []int32, sep string) string {
	buf := joinBufferPoll.Get()
	for _, v := range i {
		if buf.Len() > 0 {
			buf.WriteString(sep)
		}
		buf.WriteString(strconv.FormatInt(int64(v), 10))
	}

	defer joinBufferPoll.Put(buf)
	return buf.String()
}

// JoinInt64s 使用 sep 连接 []int64 并返回连接的字符串
func JoinInt64s(i []int64, sep string) string {
	buf := joinBufferPoll.Get()
	for _, v := range i {
		if buf.Len() > 0 {
			buf.WriteString(sep)
		}
		buf.WriteString(strconv.FormatInt(v, 10))
	}

	defer joinBufferPoll.Put(buf)
	return buf.String()
}

// JoinUints 使用 sep 连接 []uint 并返回连接的字符串
func JoinUints(i []uint, sep string) string {
	buf := joinBufferPoll.Get()
	for _, v := range i {
		if buf.Len() > 0 {
			buf.WriteString(sep)
		}
		buf.WriteString(strconv.FormatUint(uint64(v), 10))
	}

	defer joinBufferPoll.Put(buf)
	return buf.String()
}

// JoinUint8s 使用 sep 连接 []uint8 并返回连接的字符串
func JoinUint8s(i []uint8, sep string) string {
	buf := joinBufferPoll.Get()
	for _, v := range i {
		if buf.Len() > 0 {
			buf.WriteString(sep)
		}
		buf.WriteString(strconv.FormatUint(uint64(v), 10))
	}

	defer joinBufferPoll.Put(buf)
	return buf.String()
}

// JoinUint16s 使用 sep 连接 []uint16 并返回连接的字符串
func JoinUint16s(i []uint16, sep string) string {
	buf := joinBufferPoll.Get()
	for _, v := range i {
		if buf.Len() > 0 {
			buf.WriteString(sep)
		}
		buf.WriteString(strconv.FormatUint(uint64(v), 10))
	}

	defer joinBufferPoll.Put(buf)
	return buf.String()
}

// JoinUint32s 使用 sep 连接 []uint32 并返回连接的字符串
func JoinUint32s(i []uint32, sep string) string {
	buf := joinBufferPoll.Get()
	for _, v := range i {
		if buf.Len() > 0 {
			buf.WriteString(sep)
		}
		buf.WriteString(strconv.FormatUint(uint64(v), 10))
	}

	defer joinBufferPoll.Put(buf)
	return buf.String()
}

// JoinUint64s 使用 sep 连接 []uint64 并返回连接的字符串
func JoinUint64s(i []uint64, sep string) string {
	buf := joinBufferPoll.Get()
	for _, v := range i {
		if buf.Len() > 0 {
			buf.WriteString(sep)
		}
		buf.WriteString(strconv.FormatUint(v, 10))
	}

	defer joinBufferPoll.Put(buf)
	return buf.String()
}
