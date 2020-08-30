package test

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// DataType 测试数据类型
func For(c *gin.Context) {
	x := "hello"
	// for _, x := range x {
	// 	n := x + 'A' - 'a'
	// 	fmt.Printf("%q\n", n) // "HELLO" (one letter per iteration)
	// }

	// for i, r := range "Hello, 世界" {
	//   fmt.Printf("%d\t%q\t%d\n", i, r, r)
	// }

	// i := 0;
	// for ; i < 100; i++ { // 相当于while
	// 	fmt.Println(i)
	// }

	// for y:= 1; y <= 9; y++ {
	// 	for x := 1; x <= y; x++ {
	// 		fmt.Printf("%d*%d=%d\t", x, y, x*y)
	// 	}
	// 	fmt.Println()
	// }

OuterLoop:
	for i := 0; i < 2; i++ {
		for j := 0; j < 5; j++ {
			switch j {
			case 2, 3:
				fmt.Println(i, j)
				break OuterLoop
			}
		}
	}

	c.JSON(http.StatusOK, x)
}
