package util

import (
	"path/filepath"
	"testing"
)

func TestWriteString(t *testing.T) {
	str := []string{"i", " am", " a", " teacher"}
	expect := "i am a teacher"
	result := WriteString(str...)

	if result != expect {
		t.Errorf("result = %s, expect = %s", result, expect)
	}
}

func TestGetFileBasename(t *testing.T) {
	dirname := "../router"
	expect := 1
	result := GetFileBasename(dirname, []string{"go"})

	if len(result) < expect {
		t.Skipf("result = %v, no files", result)
	}
	t.Logf("result = %v", result)
}

func TestBasename(t *testing.T) {
	str := "D://go/workspace/ttt.jpg"
	expect := "ttt"
	name := filepath.Base(str)
	result := Basename(name)

	if result != expect {
		t.Errorf("result = %s, expect = %s", result, expect)
	}
}
