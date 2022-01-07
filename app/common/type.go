// 自定义类型
package common

import (
	"fmt"
	"time"
)

type Date time.Time

func (d Date) MarshalJSON() ([]byte, error) {
	dateTime := fmt.Sprintf("%q", time.Time(d).Format("2006-01-02"))
	return []byte(dateTime), nil
}
