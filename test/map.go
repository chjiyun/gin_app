package test

import (
	// "encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	// "log"
	"net/http"
)

type Response struct {
	Msg       string `json:"msg"`
	Code      int16  `json:"code"`
	Countries map[string]string
	Movies    []Movie
}

type Movie struct {
	Title  string
	Year   int  `json:"released"`
	Color  bool `json:"color,omitempty"`
	Actors []string
}

// Index 测试语法的 API
func Map(c *gin.Context) {
	// 如果不初始化 map，那么就会创建一个 nil map。nil map 不能用来存放键值对
	countryCapitalMap := map[string]string{"France": "巴黎", "Italy": "罗马"}
	// countryCapitalMap = make(map[string]string{})

	// for country, val := range countryCapitalMap {
	// 	fmt.Println(country, "首都是", val)
	// }

	// 查看元素在集合中是否存在
	capital, ok := countryCapitalMap["American"]
	// capital 此时是零值: ""
	fmt.Println("American首都是", capital, "存在与否：", ok)

	var movies = []Movie{
		{
			Title: "Casablanca", Year: 1943, Color: false,
			Actors: []string{"Humphrey Bogart", "Ingrid Bergman"},
		},
		{
			Title: "Cool Hand Luke", Year: 1967, Color: true,
			Actors: []string{"Paul Newman"},
		},
		{
			Title: "Bullitt", Year: 1968, Color: true,
			Actors: []string{"Steve McQueen", "Jacqueline Bisset"},
		},
	}

	// 要初始化
	res := Response{}

	res.Msg = "success"
	res.Code = http.StatusOK
	res.Countries = countryCapitalMap
	res.Movies = movies

	// data, err := json.Marshal(res)
	// if err != nil {
	// 	log.Fatalf("JSON marshaling failed: %s", err)
	// }
	// fmt.Println(data)

	c.JSON(http.StatusOK, res)
}
