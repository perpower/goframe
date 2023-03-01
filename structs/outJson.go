package structs

//定义接口数据返回结构体
type Outjson struct {
	Code int    `json:"code" xml:"code" yml:"code"`
	Msg  string `json:"msg" xml:"msg" yml:"msg"`
	Data any    `json:"data" xml:"data" yml:"data"`
}
