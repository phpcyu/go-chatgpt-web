package strBase

type BaseParams struct {
	Page  int `form:"page"`
	Limit int `form:"limit"`
}

type Result struct {
	Code int
	Msg  string
}

type HandleResult struct {
	Status bool
}

type ApiResult struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type IdParams struct {
	Id int `json:"id" form:"id" binding:"required,gt=0"`
}
