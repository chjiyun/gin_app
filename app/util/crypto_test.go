package util

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var hashkey = []byte("e310b2af80c4128c67a2356902421c8f01d89152")

func TestEncrypt(t *testing.T) {
	// 原文
	var data string
	src := os.Args[len(os.Args)-1]
	prefix := "src="
	if strings.HasPrefix(src, prefix) {
		data = strings.TrimPrefix(src, prefix)
		basedir, _ := filepath.Abs(".")
		filename := filepath.Join(basedir, "../../hashkey.txt")
		hashkey, _ = ioutil.ReadFile(filename)
	} else {
		data = "i am a chinese!"
	}
	result := Encrypt(data, hashkey)

	t.Skipf("encrypt txt: %s", result)
}

func TestDecrypt(t *testing.T) {
	var data string
	src := os.Args[len(os.Args)-1]
	prefix := "src="
	if strings.HasPrefix(src, prefix) {
		data = strings.TrimPrefix(src, prefix)
		basedir, _ := filepath.Abs(".")
		filename := filepath.Join(basedir, "../../hashkey.txt")
		hashkey, _ = ioutil.ReadFile(filename)
	} else {
		data = "456F29842E6B9BD8E67A98EC87CE3062"
	}
	result := Decrypt(data, hashkey)

	t.Skipf("original txt: %s", result)
}
