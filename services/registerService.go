package services

import (
	"RegistrationCenter/conf"
	"github.com/astaxie/beego/logs"
	consulapi "github.com/hashicorp/consul/api"
)

func init() {
	RegisterConfig()
}

func NewRandomCousulClient() (*consulapi.Client, string, error) {
	config := consulapi.DefaultConfig()
	config.Address = conf.RandomServer().Address
	client, err := consulapi.NewClient(config)
	return client, config.Address, err
}

func NewCousulClient(address string) (*consulapi.Client, error) {
	config := consulapi.DefaultConfig()
	config.Address = address
	return consulapi.NewClient(config)
}

func RegisterConfig() {
	client, err := NewCousulClient(conf.G_consul_server)
	if err != nil {
		logs.Error("consul client error : ", err)
	}

	serverStruct := conf.ReadServices()
	for _, v := range serverStruct.Services {
		err = client.Agent().ServiceRegister(v)
		if err != nil {
			logs.Error("register server error : ", err)
		}
		logs.Info("register service " + v.ID + " success")
	}
}
