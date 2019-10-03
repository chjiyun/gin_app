package test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DataType(c *gin.Context) {
	x := "hello"
	for _, x := range x {
		x := x + 'A' - 'a'
		fmt.Println(x) // "HELLO" (one letter per iteration)
	}
	// c.JSON(http.StatusOK, {hello: x})
}
