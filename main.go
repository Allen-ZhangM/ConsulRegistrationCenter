package main

import (
	"RegistrationCenter/conf"
	_ "RegistrationCenter/routers"
	"github.com/astaxie/beego"
)

func main() {
	conf.InitBeegoLogs()
	beego.Run()
}
