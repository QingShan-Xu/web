package test

import "github.com/QingShan-Xu/xjh/bm"

var D database

type database struct {
	Pet Pet
}

type Pet struct {
	bm.Model

	Name   string `json:"name" binding:"required"`
	Status string `json:"status" binding:"required"`
}

var API api

type api struct {
	ReqPet ReqBindGetUser
}

type ReqBindGetUser struct {
	MarketID   int    `uri:"market_id"`
	MarketName string `param:"market_name" binding:"required"`
}
