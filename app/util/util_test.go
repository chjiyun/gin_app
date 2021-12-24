package util

import (
	"errors"
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

func TestToString(t *testing.T) {
	v := map[int]interface{}{
		1: "1234",
		2: true,
		3: 2.34,
		4: nil,
	}
	v1 := []interface{}{0, true, "hello", nil, 1.23}
	v2 := struct {
		Name   string  `json:"name"`
		Age    int     `json:"age"`
		Height float32 `json:"height"`
		Male   bool    `json:"male"`
	}{
		Name:   "LIHUA",
		Age:    23,
		Height: 1.68,
		Male:   false,
	}
	expect := `{"1":"1234","2":true,"3":2.34,"4":null}`
	expect1 := `[0,true,"hello",null,1.23]`
	expect2 := `{"name":"LIHUA","age":23,"height":1.68,"male":false}`
	result := ToString(v)
	result1 := ToString(v1)
	result2 := ToString(v2)

	if result != expect {
		t.Errorf("result = %s, v= %v", result, v)
	}
	if result1 != expect1 {
		t.Errorf("result1 = %s, v= %v", result1, v1)
	}
	if result2 != expect2 {
		t.Errorf("result2 = %s, v= %v", result2, v2)
	}
}

func TestHandleData(t *testing.T) {
	expect := "this is a error"
	err := errors.New(expect)
	HandleData(&err)
	result := err.Error()

	if result != expect {
		t.Errorf("result = %v, expect = %v", result, expect)
	}
}
