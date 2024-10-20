package main

import (
	"net/http"

	"github.com/QingShan-Xu/web/bm"
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

//	var router = rt.Router{
//		Name: "宠物",
//		Path: "pet",
//		Children: []rt.Router{
//			{
//				Name:   "新建",
//				Path:   "",
//				Bind:   test.D.Pet,
//				MODEL:  test.D.Pet,
//				Method: "POST",
//				Type:   rt.TYPE.CREATE_ONE,
//			},
//			{
//				Name: "查询",
//				Path: ":id",
//				Bind: struct {
//					ID int `uri:"id" binding:"required"`
//				}{},
//				MODEL:  test.D.Pet,
//				Method: "GET",
//				WHERE:  map[string]string{"id = ?": "ID"},
//				Type:   rt.TYPE.GET_ONE,
//			},
//			{
//				Name:   "更新",
//				Path:   ":id",
//				Bind:   test.API.ReqUpdatePet,
//				MODEL:  test.D.Pet,
//				Method: "PUT",
//				WHERE:  map[string]string{"id": "ID"},
//				SELECT: map[string]string{
//					"name": "Name",
//				},
//				Type: rt.TYPE.UPDATE_ONE,
//			},
//			{
//				Name: "删除",
//				Path: ":id",
//				Bind: struct {
//					ID int `uri:"id" binding:"required"`
//				}{},
//				MODEL:  test.D.Pet,
//				Method: "DELETE",
//				WHERE:  map[string]string{"id": "ID"},
//				Type:   rt.TYPE.DELETE_ONE,
//			},
//			{
//				Name: "查列表",
//				Path: "list",
//				Bind: struct {
//					bm.Pagination
//				}{},
//				MODEL:  test.D.Pet,
//				Method: "GET",
//				Type:   rt.TYPE.GET_LIST,
//			},
//		},
//	}

type Pet struct {
	bm.Model
	Type string
	Name string
}

var router = rt.Router{
	Path: "/",
	Children: []rt.Router{
		{
			Path: "/pet",
			Children: []rt.Router{
				{
					Name:   "新建",
					Path:   "/",
					Method: "POST",
					MODEL:  Pet{},
					Bind: struct {
						Name string `bind:"name" validate:"required"`
					}{},
					CREATE_ONE: map[string]string{
						"Name": "Name",
					},
				},
				{
					Name:   "修改",
					Path:   "/{id}",
					Method: http.MethodPut,
					MODEL:  Pet{},
					Bind: struct {
						ID   string `bind:"id" validate:"required"`
						Name string `bind:"name"`
						Type *int   `bind:"type"`
					}{},
					UPDATE_ONE: map[string]string{
						"Name": "Name",
						"Type": "Type",
					},
					WHERE: [][]string{
						{"id", "ID"},
					},
				},
				{
					Name:   "详情",
					Path:   "/{id}",
					Method: http.MethodGet,
					MODEL:  Pet{},
					Bind: struct {
						ID string `bind:"id" validate:"required"`
					}{},
					GET_ONE: true,
					WHERE: [][]string{
						{"id", "ID"},
					},
				},
				{
					Name:   "删除",
					Path:   "/{id}",
					Method: http.MethodDelete,
					MODEL:  Pet{},
					Bind: struct {
						ID string `bind:"id" validate:"required"`
					}{},
					DELETE_ONE: true,
					WHERE: [][]string{
						{"id", "ID"},
					},
				},
				{
					Name:   "拿列表",
					Path:   "/",
					Method: "GET",
					MODEL:  Pet{},
					Bind: struct {
						bm.Pagination
					}{},
					GET_LIST: true,
				},
			},
		},
		{
			Name: "宠物2",
			Path: "/pet2",
			Children: []rt.Router{
				{
					Name:   "新建",
					Path:   "/",
					Method: "POST",
				},
				{
					Name:   "拿列表",
					Path:   "/",
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
		&router,
	)
}
