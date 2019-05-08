package controllers

import (
	"RegistrationCenter/conf"
	"RegistrationCenter/models"
	"RegistrationCenter/services"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	consulapi "github.com/hashicorp/consul/api"
)

type DiscoverController struct {
	beego.Controller
}

func (this *DiscoverController) RetData(resp map[string]interface{}) {
	this.Data["json"] = resp
	this.ServeJSON()
}

func (c *DiscoverController) handleError(msg string, err error, resp map[string]interface{}) {
	if err != nil {
		logs.Error(msg, err)
		resp["msg"] = msg + " err: " + err.Error()
		c.Abort("401")
	}
}

func (c *DiscoverController) DiscoverServices() {
	resp := make(map[string]interface{})
	defer c.RetData(resp)

	client, randomAddr, err := services.NewRandomCousulClient()
	c.handleError("consul client error  ", err, resp)

	result, _, err := client.Catalog().Services(&consulapi.QueryOptions{
		Datacenter: conf.Datacenters[0],
	})
	// 如果请求失败则尝试其他节点
	var resultMsg []string
	if err != nil {
		for _, v := range conf.ServerNodes {
			if randomAddr == v.Address {
				continue
			}
			client, err := services.NewCousulClient(v.Address)

			result, _, err = client.Catalog().Services(&consulapi.QueryOptions{
				Datacenter: conf.Datacenters[0],
			})
			if err != nil {
				resultMsg = append(resultMsg, "consul Catalog Services error : connection refused "+v.Address)
				continue
			}
			break
		}
	}
	if result == nil {
		resp["msg"] = resultMsg
		return
	}
	resp["services"] = result
}

func (this *DiscoverController) DiscoverById() {
	//设置返回参数
	resp := make(map[string]interface{})
	defer this.RetData(resp)
	//获取传入参数
	serviceid := this.Ctx.Input.Param(":id")
	//设置接口client
	client, randomAddr, err := services.NewRandomCousulClient()
	this.handleError("consul client error  ", err, resp)

	//查询service
	respService, _, err := client.Catalog().Service(serviceid, "", &consulapi.QueryOptions{
		Datacenter: conf.Datacenters[0],
	})
	// 如果请求失败则尝试其他节点
	if err != nil || respService == nil || len(respService) == 0 {
		for _, v := range conf.ServerNodes {
			if randomAddr == v.Address {
				continue
			}
			client, err := services.NewCousulClient(v.Address)

			respService, _, err = client.Catalog().Service(serviceid, "", &consulapi.QueryOptions{
				Datacenter: conf.Datacenters[0],
			})
			if err != nil {
				continue
			}
			if respService == nil || len(respService) == 0 {
				continue
			}
			break
		}
	}
	if respService == nil || len(respService) == 0 {
		resp["msg"] = "service not fond"
		return
	}

	resp["service"] = models.ToSimpleAgentService(*respService[0])

}
