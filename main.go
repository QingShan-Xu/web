package main

import (
	"github.com/QingShan-Xu/xjh/cf"
	"github.com/QingShan-Xu/xjh/rt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var router = rt.Router{
	Name: "宠物",
	Path: "user/:market_id",
	Children: []rt.Router{
		{
			Name:   "新建",
			Path:   "",
			Method: "POST",
		},
	},
}

func main() {
	gin := gin.Default()
	rootGroup := gin.Group("/")
	DB, _ := gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:3306)/go-learn?charset=utf8mb4&parseTime=true&loc=Local"), &gorm.Config{CreateBatchSize: 1000})
	cf.Init(rootGroup, DB.Debug(), &cf.CfgRegist{})
	rt.Init(&router)

	gin.Run(":8080")
}
