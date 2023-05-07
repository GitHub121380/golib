package utils

import (
	"errors"
)

const NapiQuestionCryptKey = "iVPed<7K"
const NapiArticleCryptKey = "^.vAy$TT"
const QbQuestionCryptKey = "lVPed<8K"
const ZybChargeCryptKey = "^.vAy$TG"
const ZybUidCryptKey = "^.vB>$TS"
const ZybLectureCryptKey = "*)<~0YZS"
const ZybOpenidCryptKey = "^.iC$eSC"
const MagicNum = 65521
const MagicNum2 = 65519
const INT32MAX = 4294967296

func EncodeQid(qid, qType int) (e string, err error) {
	var data []byte
	if qid >= INT32MAX {
		data, err = pack("NNnCCVVvC", qid>>32, qid%INT32MAX, qid%MagicNum, 0, 0, qid>>32, qid%INT32MAX, qid%MagicNum2, 0)
	} else {
		data, err = pack("NnCCVvC", qid, qid%MagicNum, 0, 0, qid, qid%MagicNum2, 0)
	}
	if err != nil {
		return e, err
	}

	key := QbQuestionCryptKey
	if qType == 0 {
		key = NapiQuestionCryptKey
	}
	return EncryptDesEcb(string(data), key, PaddingTypePKCS7)
}

func DecodeQid(qstr string, qType int) (qid int, err error) {
	key := QbQuestionCryptKey
	if qType == 0 {
		key = NapiQuestionCryptKey
	}

	decryptStr, err := DecryptDesEcb(qstr, key, PaddingTypePKCS7)
	if err != nil {
		return qid, err
	}
	decryptLen := len(decryptStr)
	if decryptLen == 23 {
		unPackedArr, err := unpack("NNnCCVVvC", decryptStr)
		if err != nil {
			return 0, err
		}
		if len(unPackedArr) < 9 {
			return qid, errors.New("id decode failed")
		}

		qid := unPackedArr[0]*INT32MAX + unPackedArr[1]
		qid2 := unPackedArr[5]*INT32MAX + unPackedArr[6]
		if unPackedArr[3] != 0 || unPackedArr[4] != 0 || unPackedArr[8] != 0 {
			return qid, errors.New("id decode failed")
		}
		if qid%MagicNum == unPackedArr[2] && qid2%MagicNum2 == unPackedArr[7] && qid == qid2 {
			return qid, nil
		}
	} else if decryptLen == 15 {
		unPackedArr, err := unpack("NnCCVvC", decryptStr)
		if err != nil {
			return 0, err
		}
		if len(unPackedArr) < 7 {
			return qid, errors.New("id decode failed")
		}

		qid := unPackedArr[0]
		qid2 := unPackedArr[4]
		if unPackedArr[2] != 0 || unPackedArr[3] != 0 || unPackedArr[6] != 0 {
			return qid, errors.New("id decode failed")
		}
		if qid%MagicNum == unPackedArr[1] && qid2%MagicNum2 == unPackedArr[5] && qid == qid2 {
			return qid, nil
		}
	}
	return qid, errors.New("id decode failed")
}

func EncodeAQid(qid int) (e string, err error) {
	var data []byte
	data, err = pack("NnCCVvC", qid, qid%MagicNum, 0, 0, qid, qid%MagicNum2, 0)
	if err != nil {
		return e, err
	}
	return EncryptDesEcb(string(data), NapiArticleCryptKey, PaddingTypeZero)
}

func DecodeAQid(qStr string) (qid int, err error) {
	decryptStr, err := DecryptDesEcb(qStr, NapiArticleCryptKey, PaddingTypeNoPadding)
	if err != nil {
		return qid, err
	}
	unPackedArr, err := unpack("NnCCVvC", decryptStr)
	if err != nil {
		return qid, err
	}
	if len(unPackedArr) < 7 {
		return qid, errors.New("id decode failed")
	}

	if unPackedArr[2] != 0 || unPackedArr[3] != 0 || unPackedArr[6] != 0 {
		return qid, errors.New("id decode failed")
	}

	qid = unPackedArr[0]
	qid2 := unPackedArr[4]
	if qid%MagicNum == unPackedArr[1] && qid2%MagicNum2 == unPackedArr[5] && qid == qid2 {
		return qid, nil
	}

	return qid, errors.New("id decode failed")
}

func EncodeUid(qid int) (e string, err error) {
	var data []byte
	data, err = pack("NnCCVvC", qid, qid%MagicNum, 0, 0, qid, qid%MagicNum2, 0)
	if err != nil {
		return e, err
	}
	return EncryptDesEcb(string(data), ZybUidCryptKey, PaddingTypeZero)
}

func DecodeUid(qStr string) (qid int, err error) {
	decryptStr, err := DecryptDesEcb(qStr, ZybUidCryptKey, PaddingTypeZero)
	if err != nil {
		return qid, err
	}
	unPackedArr, err := unpack("NnCCVvC", decryptStr)
	if err != nil {
		println("err: ", err.Error())
		return qid, err
	}
	if len(unPackedArr) < 7 {
		return qid, errors.New("id decode failed")
	}

	if unPackedArr[2] != 0 || unPackedArr[3] != 0 || unPackedArr[6] != 0 {
		return qid, errors.New("id decode failed")
	}

	qid = unPackedArr[0]
	qid2 := unPackedArr[4]
	if qid%MagicNum == unPackedArr[1] && qid2%MagicNum2 == unPackedArr[5] && qid == qid2 {
		return qid, nil
	}

	return qid, errors.New("id decode failed")
}

func EncodeCid(qid int) (e string, err error) {
	var data []byte
	data, err = pack("NnCCVvC", qid, qid%MagicNum, 0, 0, qid, qid%MagicNum2, 0)
	if err != nil {
		return e, err
	}
	return EncryptDesEcb(string(data), ZybChargeCryptKey, PaddingTypeZero)
}

func DecodeCid(qStr string) (qid int, err error) {
	decryptStr, err := DecryptDesEcb(qStr, ZybChargeCryptKey, PaddingTypeZero)
	if err != nil {
		return qid, err
	}
	unPackedArr, err := unpack("NnCCVvC", decryptStr)
	if err != nil {
		println("err: ", err.Error())
		return qid, err
	}
	if len(unPackedArr) < 7 {
		return qid, errors.New("id decode failed")
	}

	if unPackedArr[2] != 0 || unPackedArr[3] != 0 || unPackedArr[6] != 0 {
		return qid, errors.New("id decode failed")
	}

	qid = unPackedArr[0]
	qid2 := unPackedArr[4]
	if qid%MagicNum == unPackedArr[1] && qid2%MagicNum2 == unPackedArr[5] && qid == qid2 {
		return qid, nil
	}

	return qid, errors.New("id decode failed")
}

func EncodeLid(qid int) (e string, err error) {
	var data []byte
	data, err = pack("NnCCVvC", qid, qid%MagicNum, 0, 0, qid, qid%MagicNum2, 0)
	if err != nil {
		return e, err
	}
	return EncryptDesEcb(string(data), ZybLectureCryptKey, PaddingTypeZero)
}

func DecodeLid(qStr string) (qid int, err error) {
	decryptStr, err := DecryptDesEcb(qStr, ZybLectureCryptKey, PaddingTypeZero)
	if err != nil {
		return qid, err
	}
	unPackedArr, err := unpack("NnCCVvC", decryptStr)
	if err != nil {
		println("err: ", err.Error())
		return qid, err
	}
	if len(unPackedArr) < 7 {
		return qid, errors.New("id decode failed")
	}

	if unPackedArr[2] != 0 || unPackedArr[3] != 0 || unPackedArr[6] != 0 {
		return qid, errors.New("id decode failed")
	}

	qid = unPackedArr[0]
	qid2 := unPackedArr[4]
	if qid%MagicNum == unPackedArr[1] && qid2%MagicNum2 == unPackedArr[5] && qid == qid2 {
		return qid, nil
	}

	return qid, errors.New("id decode failed")
}

func EncodeOuid(qid int) (e string, err error) {
	var data []byte
	data, err = pack("NnCCVvC", qid, qid%MagicNum, 0, 0, qid, qid%MagicNum2, 0)
	if err != nil {
		return e, err
	}
	return EncryptDesEcb(string(data), ZybOpenidCryptKey, PaddingTypeZero)
}

func DecodeOuid(qStr string) (qid int, err error) {
	decryptStr, err := DecryptDesEcb(qStr, ZybOpenidCryptKey, PaddingTypeZero)
	if err != nil {
		return qid, err
	}
	unPackedArr, err := unpack("NnCCVvC", decryptStr)
	if err != nil {
		println("err: ", err.Error())
		return qid, err
	}
	if len(unPackedArr) < 7 {
		return qid, errors.New("id decode failed")
	}

	if unPackedArr[2] != 0 || unPackedArr[3] != 0 || unPackedArr[6] != 0 {
		return qid, errors.New("id decode failed")
	}

	qid = unPackedArr[0]
	qid2 := unPackedArr[4]
	if qid%MagicNum == unPackedArr[1] && qid2%MagicNum2 == unPackedArr[5] && qid == qid2 {
		return qid, nil
	}

	return qid, errors.New("id decode failed")
}
