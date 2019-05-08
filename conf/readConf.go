package conf

import (
	"RegistrationCenter/utils"
	"encoding/json"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	consulapi "github.com/hashicorp/consul/api"
	"io/ioutil"
	"math/rand"

	"time"
)

func init() {
	InitConfig()
	InitServerNodes()
	InitDatacenters()
}

//定义config变量
var (
	servicesPath         = "conf/services.json" //配置service的文件
	G_consul_server      string                 //consul服务器ip地址
	G_consul_server_port string                 //consul服务器端口 默认8500
	Datacenters          []string               //启动时获取全部dc
	ServerNodes          []ServerNode           //启动时获取全部节点，ip:port
)

type ServicesStruct struct {
	Services []*consulapi.AgentServiceRegistration
}

type ServerNode struct {
	Address string
	Name    string
}

func ReadServices() ServicesStruct {
	v := ServicesStruct{}
	load(servicesPath, &v)
	return v
}

func load(filename string, v interface{}) {
	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, v)
	if err != nil {
		return
	}
}

func RandomServer() ServerNode {
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(len(ServerNodes))
	return ServerNodes[random]
}

func InitConfig() {
	//从配置文件读取配置信息
	appconf, err := config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		logs.Debug(err)
		return
	}
	G_consul_server = appconf.String("consul_server")
	G_consul_server_port = appconf.String("consul_server_port")
}

func InitServerNodes() {
	nodes := []consulapi.Node{}
	bytes := utils.HttpGet("http://" + G_consul_server + "/v1/catalog/nodes")
	json.Unmarshal(bytes, &nodes)
	for i, _ := range nodes {
		serverNode := ServerNode{
			Name:    nodes[i].Node,
			Address: nodes[i].Address + ":" + G_consul_server_port,
		}
		ServerNodes = append(ServerNodes, serverNode)
	}
}

func InitDatacenters() {
	dcs := []string{}
	bytes := utils.HttpGet("http://" + G_consul_server + "/v1/catalog/datacenters")
	json.Unmarshal(bytes, &dcs)
	for i, _ := range dcs {
		Datacenters = append(Datacenters, dcs[i])
	}
}
