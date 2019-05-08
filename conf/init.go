package conf

import (
	"github.com/astaxie/beego/logs"
)

func InitBeegoLogs() {
	logs.SetLogger(logs.AdapterFile, `{"filename":"logs/logs.log","level":3,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)
	//异步 缓冲
	logs.Async(10000)
	logs.EnableFuncCallDepth(true)
}
