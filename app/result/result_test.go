package result

import (
	"errors"
	"testing"
)

func TestHandleData(t *testing.T) {
	data := Result{
		Code: 200,
		Msg:  "test",
		Data: errors.New("i am error"),
	}
	var data1 interface{}
	handleData(&data.Data, &data1)

	if data1 == nil {
		t.Errorf("result = %v", data)
	}
}
