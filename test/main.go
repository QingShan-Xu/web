package test

// package main

// import (
// 	"net/http"

// 	"github.com/QingShan-Xu/web/bm"
// 	"github.com/QingShan-Xu/web/gm"
// 	"github.com/QingShan-Xu/web/rt"
// )

// type Pet struct {
// 	bm.Model
// 	Type string
// 	Name string
// }

// var router = rt.Router{
// 	Path: "/",
// 	Children: []rt.Router{
// 		{
// 			Path: "/pet",
// 			Children: []rt.Router{
// 				{
// 					Name:   "新建",
// 					Path:   "/",
// 					Method: "POST",
// 					MODEL:  Pet{},
// 					Bind: struct {
// 						Name string `bind:"name" validate:"required"`
// 					}{},
// 					CREATE_ONE: map[string]string{
// 						"Name": "Name",
// 					},
// 				},
// 				{
// 					Name:   "修改",
// 					Path:   "/{id}",
// 					Method: http.MethodPut,
// 					MODEL:  Pet{},
// 					Bind: struct {
// 						ID   string `bind:"id" validate:"required"`
// 						Name string `bind:"name"`
// 						Type *int   `bind:"type"`
// 					}{},
// 					UPDATE_ONE: map[string]string{
// 						"Name": "Name",
// 						"Type": "Type",
// 					},
// 					WHERE: [][]string{
// 						{"id", "ID"},
// 					},
// 				},
// 				{
// 					Name:   "详情",
// 					Path:   "/{id}",
// 					Method: http.MethodGet,
// 					MODEL:  Pet{},
// 					Bind: struct {
// 						ID string `bind:"id" validate:"required"`
// 					}{},
// 					GET_ONE: true,
// 					WHERE: [][]string{
// 						{"id", "ID"},
// 					},
// 				},
// 				{
// 					Name:   "删除",
// 					Path:   "/{id}",
// 					Method: http.MethodDelete,
// 					MODEL:  Pet{},
// 					Bind: struct {
// 						ID string `bind:"id" validate:"required"`
// 					}{},
// 					DELETE_ONE: true,
// 					WHERE: [][]string{
// 						{"id", "ID"},
// 					},
// 				},
// 				{
// 					Name:   "拿列表",
// 					Path:   "/",
// 					Method: "GET",
// 					MODEL:  Pet{},
// 					Bind: struct {
// 						bm.Pagination
// 					}{},
// 					GET_LIST: true,
// 				},
// 			},
// 		},
// 		{
// 			Name: "宠物2",
// 			Path: "/pet2",
// 			Children: []rt.Router{
// 				{
// 					Name:   "新建",
// 					Path:   "/",
// 					Method: "POST",
// 				},
// 				{
// 					Name:   "拿列表",
// 					Path:   "/",
// 					Method: "GET",
// 				},
// 			},
// 		},
// 	},
// }

// func main() {
// 	gm.Start(
// 		gm.Cfg{
// 			FileName: "app",
// 			FilePath: []string{"."},
// 		},
// 		&router,
// 	)
// }
