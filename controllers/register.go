package controllers

import (
	"RegistrationCenter/conf"
	"RegistrationCenter/models"
	"RegistrationCenter/services"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	consulapi "github.com/hashicorp/consul/api"
)

type RegisterController struct {
	beego.Controller
}

func (this *RegisterController) RetData(resp map[string]interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

func (c *RegisterController) handleError(msg string, err error, resp map[string]interface{}) {
	if err != nil {
		logs.Error(msg, err)
		resp["msg"] = msg + " err: " + err.Error()
		c.Abort("401")
	}
}

func (c *RegisterController) RegisterConfig() {
	resp := make(map[string]interface{})
	defer c.RetData(resp)

	client, _, err := services.NewRandomCousulClient()
	c.handleError("consul client error  ", err, resp)

	var resultMsg []string
	serviceStruct := conf.ReadServices()
	for _, v := range serviceStruct.Services {
		err = client.Agent().ServiceRegister(v)
		c.handleError("register server error  ", err, resp)

		resultMsg = append(resultMsg, "register service "+v.ID+" success")
	}
	resp["msg"] = resultMsg
}

func (c *RegisterController) DeregisterService() {
	resp := make(map[string]interface{})
	defer c.RetData(resp)
	serviceid := c.Ctx.Input.Param(":id")
	var resultMsg []string
	for _, v := range conf.ServerNodes {
		client, err := services.NewCousulClient(v.Address)
		c.handleError("consul client error  ", err, resp)

		err = client.Agent().ServiceDeregister(serviceid)
		c.handleError("deregister service error   ", err, resp)

		resultMsg = append(resultMsg, "deregister service "+serviceid+" from "+v.Address+" success")
	}
	resp["msg"] = resultMsg
}

func (c *RegisterController) RegisterService() {
	resp := make(map[string]interface{})
	defer c.RetData(resp)

	var cr *models.CatalogRegistration
	//获取传过来的json数据
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &cr)
	c.handleError("json Unmarshal error  ", err, resp)

	// 参数验证
	ok, msg := checkParamRegisterService(cr)
	if !ok {
		resp["msg"] = "param error : " + msg
		return
	}

	var client *consulapi.Client
	var randomAddr string
	if cr.Address == "" {
		client, randomAddr, err = services.NewRandomCousulClient()
	} else {
		client, err = services.NewCousulClient(cr.Address)
	}
	c.handleError("consul client error  ", err, resp)

	asr := models.ToAgentServiceRegistration(cr.Service)
	err = client.Agent().ServiceRegister(asr)

	// 如果请求失败则尝试其他节点
	var resultMsg []string
	isok := false
	if err != nil {
		for _, v := range conf.ServerNodes {
			if randomAddr == v.Address {
				continue
			}
			client, err := services.NewCousulClient(v.Address)

			err = client.Agent().ServiceRegister(asr)

			if err != nil {
				resultMsg = append(resultMsg, "register server error :  "+v.Address)
				continue
			}
			isok = true
			break
		}

		if !isok {
			resp["msg"] = resultMsg
			return
		}
	}

	resp["msg"] = "register service " + cr.Service.ID + " success"
}

func checkParamRegisterService(asr *models.CatalogRegistration) (bool, string) {
	if asr.Service == nil {
		return false, "Service : is empty"
	}
	if asr.Service.Name == "" || asr.Service.ID == "" || asr.Service.Address == "" {
		return false, "Name or ID or Address : is empty"
	}

	return true, ""
}
