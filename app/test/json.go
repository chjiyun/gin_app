package test

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Stu struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Height  int    `json:"height"`
	Class   *Class `json:"class"` //指针变量，更快且能节省内存空间
	Records []int  `json:"records"`
}
type Class struct {
	Name  string
	Grade int
}
type StuRead struct {
	Name  interface{}
	Type  interface{}
	Class json.RawMessage `json:"class"` //注意这里
}

func Json(c *gin.Context) {
	// input := `["a", "b", ["c", "d"], "e"]`

	// output := make([]interface{}, 0)

	// if err := json.Unmarshal([]byte(input), &output); err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Printf("Unmarshalled slice %v\n", output)

	// // marshal it back
	// back, err := json.Marshal(output)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("Back to JSON: %v\n", string(back))

	// 创建对象数组
	students := []Stu{
		{
			Name:   "Sasha",
			Age:    25,
			Height: 170,
			Class: &Class{ //结构体实例化
				Name:  "一班",
				Grade: 3,
			},
			Records: []int{88, 66, 74, 95, 92},
		},
		{
			Name:   "Coin",
			Age:    28,
			Height: 171,
		},
	}
	students = append(students, Stu{Name: "晨曦", Age: 26, Height: 168})
	fmt.Println(students)

	// 解析json字符串
	data := `
		{
			"Type": 1,
			"Class":{
				"Name":"一班",
				"Grade": 3
			}
		}`
	stu := StuRead{}
	// json.Unmarshal([]byte(data), &stu)
	if err := json.Unmarshal([]byte(data), &stu); err != nil {
		panic(err)
	}
	fmt.Println("stu:", stu)
	fmt.Println("stu.Class:", string(stu.Class))

	cla := &Class{}
	json.Unmarshal(stu.Class, cla)
	fmt.Println("class:", cla)

	c.JSON(http.StatusOK, students)
}
