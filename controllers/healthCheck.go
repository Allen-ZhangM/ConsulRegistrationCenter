package controllers

import (
	"github.com/astaxie/beego"
	"time"
)

type HealthCheckController struct {
	beego.Controller
}

func (this *HealthCheckController) HealthCheck() {
	this.Ctx.WriteString("health check! " + time.Now().String())
}
