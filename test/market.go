package test

import "github.com/QingShan-Xu/xjh/bm"

var D database

type database struct {
	Pet Market
}

type Market struct {
	bm.Model

	Title string
	Value string
}

var API api

type api struct {
	ReqPet ReqBindGetUser
}

type ReqBindGetUser struct {
	MarketID   int    `uri:"market_id"`
	MarketName string `param:"market_name" binding:"required"`
}
