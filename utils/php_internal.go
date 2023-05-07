package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
)

/*
	php 内部函数实现，为了兼容phplib中部分函数而封装，不推荐使用
*/

func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}
func Base64Decode(str string) (string, error) {
	switch len(str) % 4 {
	case 2:
		str += "=="
	case 3:
		str += "="
	}

	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func Hex2bin(hexStr string) (string, error) {
	oldLen := len(hexStr)
	targetLength := oldLen >> 1

	str := make([]byte, targetLength)
	for i, j := 0, 0; i < targetLength; i++ {
		c := hexStr[j]
		j += 1
		if c >= '0' && c <= '9' {
			str[i] = (c - '0') << 4
		} else if c >= 'a' && c <= 'f' {
			str[i] = (c - 'a' + 10) << 4
		} else if c >= 'A' && c <= 'F' {
			str[i] = (c - 'A' + 10) << 4
		} else {
			return "", errors.New("invalid hex string")
		}

		c = hexStr[j]
		j += 1
		if c >= '0' && c <= '9' {
			str[i] |= c - '0'
		} else if c >= 'a' && c <= 'f' {
			str[i] |= c - 'a' + 10
		} else if c >= 'A' && c <= 'F' {
			str[i] |= c - 'A' + 10
		} else {
			return "", errors.New("invalid hex string!")
		}
	}

	return string(str), nil
}

func Bin2hex(binStr string) string {
	hexConvTab := []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'}
	oldLen := len(binStr)

	ret := make([]byte, 2*oldLen)
	for i, j := 0, 0; i < oldLen; i++ {
		ret[j] = hexConvTab[binStr[i]>>4]
		j += 1
		ret[j] = hexConvTab[binStr[i]&15]
		j += 1
	}

	return string(ret)
}

// 不完全实现php pack 方法，慎用
func unpack(format string, binStr string) (outputs []int, err error) {
	buf := bytes.NewBufferString(binStr)
	byteOrder := binary.BigEndian
	for _, val := range []rune(format) {
		fmtKind := string(val)
		var err error
		switch fmtKind {
		case "N":
			var tmp uint32
			err = binary.Read(buf, byteOrder, &tmp)
			outputs = append(outputs, int(tmp))
		case "n":
			var tmp uint16
			err = binary.Read(buf, byteOrder, &tmp)
			outputs = append(outputs, int(tmp))
		case "c":
			var tmp int8
			err = binary.Read(buf, binary.LittleEndian, &tmp)
			outputs = append(outputs, int(tmp))
		case "C":
			var tmp uint8
			err = binary.Read(buf, binary.LittleEndian, &tmp)
			outputs = append(outputs, int(tmp))
		case "V":
			var tmp uint32
			err = binary.Read(buf, binary.LittleEndian, &tmp)
			outputs = append(outputs, int(tmp))
		case "v":
			var tmp uint16
			err = binary.Read(buf, binary.LittleEndian, &tmp)
			outputs = append(outputs, int(tmp))
		}

		if err != nil {
			return outputs, err
		}
	}

	return outputs, nil
}

func pack(format string, args ...int) (outputs []byte, err error) {
	buf := new(bytes.Buffer)
	byteOrder := binary.BigEndian
	for idx, val := range []rune(format) {
		fmtKind := string(val)
		switch fmtKind {
		case "N":
			err = binary.Write(buf, byteOrder, uint32(args[idx]))
		case "n":
			err = binary.Write(buf, byteOrder, uint16(args[idx]))
		case "c":
			err = binary.Write(buf, binary.LittleEndian, int8(args[idx]))
		case "C":
			err = binary.Write(buf, binary.LittleEndian, uint8(args[idx]))
		case "V":
			err = binary.Write(buf, binary.LittleEndian, uint32(args[idx]))
		case "v":
			err = binary.Write(buf, binary.LittleEndian, uint16(args[idx]))
		}
		if err != nil {
			return outputs, err
		}
	}

	outputs = buf.Bytes()
	return outputs, nil
}
