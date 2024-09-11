package main

import (
	"github.com/QingShan-Xu/xjh/bm"
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
		// {
		// 	Name:   "新建",
		// 	Path:   ":user_id",
		// 	Bind:   test.API.ReqCreatePet,
		// 	MODEL:  test.D.Pet,
		// 	Method: "POST",
		// 	BeforeFinisher: func(bind interface{}) interface{} {
		// 		newBind := bind.(*test.ReqCreatePet)
		// 		var pet test.Pet
		// 		pet.Name = newBind.Name
		// 		pet.Status = newBind.Status
		// 		return pet
		// 	},
		// 	Finisher: rt.Finisher.Create,
		// },
		// {
		// 	Name: "查询",
		// 	Path: ":id",
		// 	Bind: struct {
		// 		ID int `uri:"id" binding:"required"`
		// 	}{},
		// 	MODEL:    test.D.Pet,
		// 	Method:   "GET",
		// 	WHERE:    map[string]string{"id": "ID"},
		// 	Finisher: rt.Finisher.First,
		// },
		// {
		// 	Name:   "更新",
		// 	Path:   ":id",
		// 	Bind:   test.API.ReqUpdatePet,
		// 	MODEL:  test.D.Pet,
		// 	Method: "PUT",
		// 	WHERE:  map[string]string{"id": "ID"},
		// 	BeforeFinisher: func(bind interface{}) interface{} {
		// 		newBind := bind.(*test.ReqUpdatePet)
		// 		var pet test.Pet
		// 		pet.Name = newBind.Name
		// 		return pet
		// 	},
		// 	Finisher: rt.Finisher.Update,
		// },
		// {
		// 	Name: "删除",
		// 	Path: ":id",
		// 	Bind: struct {
		// 		ID int `uri:"id" binding:"required"`
		// 	}{},
		// 	MODEL:    test.D.Pet,
		// 	Method:   "DELETE",
		// 	WHERE:    map[string]string{"id": "ID"},
		// 	Finisher: rt.Finisher.Delete,
		// },
		{
			Name: "查列表",
			Path: "list",
			Bind: struct {
				bm.Pagination
				ID int `form:"id" binding:"required"`
			}{},
			MODEL:  test.D.Pet,
			Method: "GET",
			Type:   rt.TYPE.GET_LIST,
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
