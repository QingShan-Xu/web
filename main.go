package main

import (
	"net/http"

	"github.com/QingShan-Xu/web/gm"
	"github.com/QingShan-Xu/web/rt"
)

// import (
// 	"github.com/QingShan-Xu/web/bm"
// 	"github.com/QingShan-Xu/web/cf"
// 	"github.com/QingShan-Xu/web/rt"
// 	"github.com/QingShan-Xu/web/test"
// 	"github.com/gin-gonic/gin"
// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// )

// var router = rt.Router{
// 	Name: "宠物",
// 	Path: "pet",
// 	Children: []rt.Router{
// 		{
// 			Name:   "新建",
// 			Path:   "",
// 			Bind:   test.D.Pet,
// 			MODEL:  test.D.Pet,
// 			Method: "POST",
// 			Type:   rt.TYPE.CREATE_ONE,
// 		},
// 		{
// 			Name: "查询",
// 			Path: ":id",
// 			Bind: struct {
// 				ID int `uri:"id" binding:"required"`
// 			}{},
// 			MODEL:  test.D.Pet,
// 			Method: "GET",
// 			WHERE:  map[string]string{"id = ?": "ID"},
// 			Type:   rt.TYPE.GET_ONE,
// 		},
// 		{
// 			Name:   "更新",
// 			Path:   ":id",
// 			Bind:   test.API.ReqUpdatePet,
// 			MODEL:  test.D.Pet,
// 			Method: "PUT",
// 			WHERE:  map[string]string{"id": "ID"},
// 			SELECT: map[string]string{
// 				"name": "Name",
// 			},
// 			Type: rt.TYPE.UPDATE_ONE,
// 		},
// 		{
// 			Name: "删除",
// 			Path: ":id",
// 			Bind: struct {
// 				ID int `uri:"id" binding:"required"`
// 			}{},
// 			MODEL:  test.D.Pet,
// 			Method: "DELETE",
// 			WHERE:  map[string]string{"id": "ID"},
// 			Type:   rt.TYPE.DELETE_ONE,
// 		},
// 		{
// 			Name: "查列表",
// 			Path: "list",
// 			Bind: struct {
// 				bm.Pagination
// 			}{},
// 			MODEL:  test.D.Pet,
// 			Method: "GET",
// 			Type:   rt.TYPE.GET_LIST,
// 		},
// 	},
// }

var router = rt.Router{
	Path: "/",
	Children: []rt.Router{
		{
			Name: "宠物",
			Path: "pet",
			Children: []rt.Router{
				{
					Name:   "新建",
					Path:   "",
					Method: "POST",
					Handler: func(w http.ResponseWriter, r *http.Request) error {
						return nil
					},
				},
				{
					Name:   "拿列表",
					Path:   "",
					Method: "GET",
				},
			},
		},
	},
}

func main() {
	gm.Start(
		gm.Cfg{
			FileName: "app",
			FilePath: []string{"."},
		},
		router,
	)
}
