package util

import "testing"

func TestEncrypt(t *testing.T) {
	// 原文
	data := "i am a chinese!"
	key := []byte("e310b2af80c4128c67a2356902421c8f01d89152")
	result := Encrypt(data, key)

	// 每次加密结果会不一样
	t.Skipf("encrypt txt: %s", result)
}

func TestDecrypt(t *testing.T) {
	data := "456F29842E6B9BD8E67A98EC87CE3062"
	key := []byte("e310b2af80c4128c67a2356902421c8f01d89152")
	result := Decrypt(data, key)
	expect := "i am a chinese!"

	if result != expect {
		t.Errorf("result = %s, expect = %s", result, expect)
	}
}
