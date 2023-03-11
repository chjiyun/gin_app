package test

import "testing"

type Item struct {
	Age int
}

// TestRange 测试for循环的坑
func TestRange(t *testing.T) {
	var all []*Item
	items := []Item{
		{Age: 11},
		{Age: 22},
		{Age: 33},
	}
	// item 只在第一次声明并初始化，后面的迭代会改变赋值，也就是地址指向每次迭代的元素
	for _, item := range items {
		all = append(all, &item)
	}
	t.Skipf("items is %v", all)
}
