package main

import (
	"github.com/QingShan-Xu/xjh/cf"
	"github.com/QingShan-Xu/xjh/rt"
	"github.com/QingShan-Xu/xjh/test"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var router = rt.Router{
	Name: "宠物",
	Path: "pet",
	Children: []rt.Router{
		{
			Name:     "新建",
			Path:     "",
			Bind:     test.D.Pet,
			MODEL:    test.D.Pet,
			Method:   "POST",
			Finisher: rt.Finisher.Create,
		},
	},
}

func main() {
	gin := gin.Default()
	rootGroup := gin.Group("/")
	DB, _ := gorm.Open(mysql.Open("root:xjh123321@tcp(127.0.0.1:3306)/go-learn?charset=utf8mb4&parseTime=true&loc=Local"), &gorm.Config{CreateBatchSize: 1000})
	cf.Init(rootGroup, DB.Debug(), &cf.CfgRegist{})
	rt.Init(&router)

	gin.Run(":8080")
}
