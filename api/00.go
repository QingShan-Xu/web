package api

var API api

type api struct {
	ReqBindGetUser ReqBindGetUser
}

type ReqBindGetUser struct {
	FormAge   int    `form:"form_age" binding:"required"`
	UriPet    string `uri:"uri_pet" binding:"required"`
	ParamFood bool   `param:"param_food" binding:"required"`
}
