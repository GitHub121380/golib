package utils

import (
	"testing"
)

func TestEncodeQid(t *testing.T) {
	qid := 42949672978
	//qid := 123456
	output, err := EncodeQid(qid, 1)
	if err != nil {
		t.Errorf("encodeQid error: %s\n", err.Error())
		return
	}
	t.Logf("EncodeQid : %s", output)

	ret, err := DecodeQid(output, 1)
	if err != nil {
		t.Errorf("decodeQid error: %s\n", err.Error())
		return
	}
	t.Logf("decodeQid : %d", ret)
}

func TestEncodeAQid(t *testing.T) {
	qid := 123456
	output, err := EncodeAQid(qid)
	if err != nil {
		t.Errorf("EncodeAQid error: %s\n", err.Error())
		return
	}
	t.Logf("EncodeAQid : %s", output)

	ret, err := DecodeAQid(output)
	if err != nil {
		t.Errorf("decodeQid error: %s\n", err.Error())
		return
	}
	t.Logf("decodeQid : %d", ret)
}

func TestEncodeUid(t *testing.T) {
	qid := 123456
	output, err := EncodeUid(qid)
	if err != nil {
		t.Errorf("EncodeUid error: %s\n", err.Error())
		return
	}
	t.Logf("EncodeUid : %s", output)

	ret, err := DecodeUid(output)
	if err != nil {
		t.Errorf("DecodeUid error: %s\n", err.Error())
		return
	}
	t.Logf("DecodeUid : %d", ret)
}

func TestEncodeCid(t *testing.T) {
	qid := 123456
	output, err := EncodeCid(qid)
	if err != nil {
		t.Errorf("EncodeCid error: %s\n", err.Error())
		return
	}
	t.Logf("EncodeCid : %s", output)

	ret, err := DecodeCid(output)
	if err != nil {
		t.Errorf("DecodeCid error: %s\n", err.Error())
		return
	}
	t.Logf("DecodeCid : %d", ret)
}

func TestEncodeLid(t *testing.T) {
	qid := 123456
	output, err := EncodeLid(qid)
	if err != nil {
		t.Errorf("EncodeLid error: %s\n", err.Error())
		return
	}
	t.Logf("EncodeLid : %s", output)

	ret, err := DecodeLid(output)
	if err != nil {
		t.Errorf("DecodeLid error: %s\n", err.Error())
		return
	}
	t.Logf("DecodeLid : %d", ret)
}

func TestEncodeOuid(t *testing.T) {
	qid := 123456
	output, err := EncodeOuid(qid)
	if err != nil {
		t.Errorf("EncodeOuid error: %s\n", err.Error())
		return
	}
	t.Logf("EncodeOuid : %s", output)

	ret, err := DecodeOuid(output)
	if err != nil {
		t.Errorf("EncodeOuid error: %s\n", err.Error())
		return
	}
	t.Logf("EncodeOuid : %d", ret)
}
