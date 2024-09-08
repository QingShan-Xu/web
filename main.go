package main

import (
	"github.com/QingShan-Xu/xjh/api"
	"github.com/QingShan-Xu/xjh/cf"
	"github.com/QingShan-Xu/xjh/rt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var router = rt.Router{
	Name:   "人员",
	Path:   "user/:uri_pet",
	Method: "POST",
	Bind:   api.API.ReqBindGetUser,
	WHERE: map[string]interface{}{
		"id = ?": api.API.ReqBindGetUser.UriPet,
	},
}

func main() {
	gin := gin.Default()
	rootGroup := gin.Group("/")
	cf.Init(rootGroup, &gorm.DB{}, &cf.CfgRegist{})
	rt.Init(&router)

	gin.Run(":8080")
}
