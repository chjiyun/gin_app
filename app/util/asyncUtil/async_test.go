package asyncUtil

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestGroup_Limit(t *testing.T) {
	g, _ := NewGroup(context.Background(), 3)
	cost := []int{5, 2, 3, 2, 4, 5, 2, 5, 1}
	for i, v := range cost {
		v := v
		i := i
		g.Go(func() error {
			fmt.Printf("i=%v, v=%v, time=%v\n", i, v, time.Now().Format(time.TimeOnly))
			time.Sleep(time.Second * time.Duration(v))
			return nil
		})
	}
	_ = g.Wait()
	t.Skip("测试并发任务分组执行")
}
