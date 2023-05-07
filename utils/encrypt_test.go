package utils

import "testing"

func TestMd5(t *testing.T) {
	plainText := "this is a Test"
	m := Md5(plainText)
	t.Logf("input data: %s ,md5: %s \n", plainText, m)
}

func TestDecryptDesEcb(t *testing.T) {
	input, key := "12345678", "iVPed<7K"
	encryptStr, err := EncryptDesEcb(input, key, PaddingTypePKCS7)
	if err != nil {
		t.Errorf("EncryptDesEcb error: %s\n", err.Error())
		return
	}

	t.Logf("input data: %s ,key: %s \n", input, key)
	t.Logf("DecryptDesEcb: %s \n", encryptStr)
	decryptStr, err := DecryptDesEcb(encryptStr, key, PaddingTypePKCS7)
	if err != nil {
		t.Errorf("DecryptDesEcb error: %s\n", err.Error())
		return
	}
	t.Logf("phpBin2hex: %s \n", decryptStr)
}
