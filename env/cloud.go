package env

import (
	"os"
	"strings"
)

// 机房相关信息
const (
	CLOUD_TEST    string = "test"
	CLOUD_BAIDU   string = "baidu"
	CLOUD_TENCENT string = "tencent"
)

const (
	IDC_TEST    string = "test"
	IDC_BJCQ    string = "bjcq"
	IDC_BJDD    string = "bjdd"
	IDC_QCPMBJ4 string = "qcpmbj4"
	IDC_QCVMBJ4 string = "qcvmbj4"
	IDC_QCVMBJ5 string = "qcvmbj5"
)

var (
	//idc与CLOUD对应关系
	IDC_CLOUD_MAPS = map[string]string{
		IDC_TEST:    CLOUD_TEST,
		IDC_BJCQ:    CLOUD_BAIDU,
		IDC_BJDD:    CLOUD_BAIDU,
		IDC_QCPMBJ4: CLOUD_TENCENT,
		IDC_QCVMBJ4: CLOUD_TENCENT,
		IDC_QCVMBJ5: CLOUD_TENCENT,
	}
)

// 获取本机所属IDC
func queryIDC() (string, string) {
	if r := os.Getenv(AppIDC); r == "test" {
		return IDC_TEST, CLOUD_TEST
	}

	divide := strings.Split(Hostname, ".")
	if len(divide) > 2 {
		idc := divide[1]
		if c, ok := IDC_CLOUD_MAPS[idc]; ok {
			return idc, c
		}
	}
	return IDC_TEST, CLOUD_TEST
}
