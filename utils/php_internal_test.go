package utils

import "testing"

func TestPack(t *testing.T) {
	qid := 42949672978
	//qid := 123456
	packedData, err := pack("NnCCVvC", qid, qid%MagicNum, 0, 0, qid, qid%MagicNum2, 0)
	if err != nil {
		t.Errorf("pack error: %s\n", err.Error())
		return
	}

	unpackedData, err := unpack("NnCCVvC", string(packedData))
	if err != nil {
		t.Errorf("unpack error: %s\n", err.Error())
		return
	}

	t.Logf("packData: %d %d %d %d %d %d %d\n", qid, qid%MagicNum, 0, 0, qid, qid%MagicNum2, 0)
	t.Logf("unpackedData: %v\n", unpackedData)
}

func TestHex2bin(t *testing.T) {
	hexStr := "123456"
	binStr, err := Hex2bin(hexStr)
	t.Logf("input data: %s \n", hexStr)
	if err != nil {
		t.Errorf("phpHex2bin error: %s\n", err.Error())
		return
	}
	t.Logf("phpHex2bin: %s \n", binStr)

	retStr := Bin2hex(binStr)
	t.Logf("phpBin2hex: %s \n", retStr)
}

func TestBase64Encode(t *testing.T) {
	qid := "123456"
	e := Base64Encode(qid)
	t.Logf("encode: %s \n", qid)
	t.Logf("encodeRet: %s \n", e)

	d, err := Base64Decode(e)
	if err != nil {
		t.Errorf("Base64Decode error: %s\n", err.Error())
		return
	}
	t.Logf("decodeRet: %s \n", d)
}
