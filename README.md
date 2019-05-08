# Registration Center (RC)

## Usage

### Console 配置与启动

1. 下载

https://releases.hashicorp.com/consul/
2. 安装
```
    unzip consul_${CONSUL_VERSION}_linux_amd64.zip
    sudo mv consul /usr/local/bin/
    consul --version
```
3. 启动
 
- -bind 指定部署的主机地址
- -node 指定名称
- --client=0.0.0.0 指定外网访问
- -join 10.1.100.159 加入节点
```
    nohup consul agent -server -bootstrap-expect 3  -data-dir /tmp/consul -node=n1 -config-dir=/etc/consul.d/ -bind=10.1.100.159 --client=0.0.0.0 -ui &

    nohup consul agent -server -bootstrap-expect 3 -data-dir /tmp/consul -node=n2 -config-dir=/etc/consul.d/ -bind=10.1.100.92 --client=0.0.0.0 -join 10.1.100.159 &

    nohup consul agent -server -bootstrap-expect 3 -data-dir /tmp/consul -node=n3 -config-dir=/etc/consul.d/ -bind=10.1.100.91 --client=0.0.0.0 -join 10.1.100.159 &
```
4. 查看
```
//查看任务
jobs
//查看进程
ps -aux | grep consul

//查看consul状态-UI  加了-ui参数的服务器查看UI监控
http://10.1.100.159:8500

//查看选举状态
consul operator raft list-peers
//查看成员
consul members
```

### RegistrationCenter 配置与启动

1. 依赖
```
//安装govendor
go get github.com/kardianos/govendor
//同步依赖
govendor sync
```
2. 配置 

主配置文件 `conf/app.conf`   **配置时要删除注释**
```
appname = registrationcenter 
httpport = 8080 //端口
runmode = dev 
autorender = false
copyrequestbody = true
EnableDocs = true

consul_server = 10.1.100.159:8500 //consul集群中的一个服务器地址
consul_server_port = 8500 //consul端口默认8500
```

项目启动时会把自己作为一个服务注册在`conf/app.conf--consul_server`上，需要配置服务注册信息，配置文件`conf/services.json`。也可以在其中添加其他默认注册的服务，配置文件内容参考下面--**服务注册 - 配置文件方式**

3. 启动

```
bee run
```


## 接口说明

### 服务注册

#### 服务注册 - 配置文件方式

在某一节点服务器 `/etc/consul.d` 目录下 `.json` 文件中追加服务配置，或者新建 `.json` 文件

配置如下, 其中 Servces 段为单个 service 的配置:

```json
{
  "services": [
    {
      "id": "backupConstfun",
      "name": "backupConstfun",
      "tags": [
        "primary"
      ],
      "address": "10.1.100.115",
      "port": 10101,
      "checks": [
        {
          "http": "http://10.1.100.115:10101/check",
          "tls_skip_verify": false,
          "method": "Get",
          "interval": "60s"
        }
      ]
    }
  ]
}
```

#### 参数说明

- **name (string: required)** - 服务的逻辑名，可重复
- **address (string: required)** - 指定服务器的地址
- **id (string: required)** - 服务的id，在每个agent唯一，重复则会覆盖
- tags (array(string): nil) - 指定要分配给服务的tag列表
- meta (map(string|string): nil) - 可指定任意KV数据
- checks (Check: nil) - 指定健康检查，地址、周期和超时

### 服务注册 - API方式 - 动态注册

请求 RC 的注册接口: `http://10.1.100.91:8080/api/register/service`

#### Service配置参数说明

- **service (Service: required)** - 参数内容和上面`services`数组中的一个`service`值一致
- datacenter (string: "") - 指定一个数据中心，一个集群时默认为 dc1
- address (string: "") - 指定一个节点注册，缺省则随机

```json
{
    "datacenter": "dc1",
    "address": "10.1.100.159:8500",
    "service": {
        "id": "RegistrationCenter",
        "name": "RegistrationCenter",
        "tags": [
            "primary"
        ],
        "address": "10.1.100.240",
        "port": 8080,
        "checks": [
            {
                "http": "http://10.1.100.240:8080/api/health/check",
                "tls_skip_verify": false,
                "method": "GET",
                "interval": "10s",
                "header": {
                    "content-type": [
                        "application/x-www-form-urlencoded",
                        "baz"
                    ]
                }
            }
        ]
    }
}
```

### 查询全部服务列表

请求 RC 接口: `http://10.1.100.91:8080/api/discover/services`

返回数据

```json
{
    "services": {
        "RegistrationCenter": [
            "primary"
        ],
        "backupConstfun": [
            "primary"
        ],
        "backup_constfun": [
            "primary"
        ],
        "consul": [],
        "cryptoCheck": [
            "primary",
            "excel-250"
        ]
    }
}
```

- key : 服务名称
- value : 服务的tags内容

### 查询服务信息

请求 RC 接口 GET : `http://10.1.100.91:8080/api/discover/service/:serviceid`

返回数据

```json
{
  "service": {
    "id": "backup_constfun",
    "name": "backup_constfun",
    "tags": [
      "primary"
    ],
    "meta": {},
    "port": 10101,
    "address": "10.1.100.115"
  }
}
```

### 删除服务

清除所有agent节点上的`serviceid`服务

请求 RC 接口 POST : `http://10.1.100.91:8080/api/deregister/:serviceid`

返回数据

```
{
  "msg": [
    "deregister service backupConstfun2 from 10.1.100.159:8500 success",
    "deregister service backupConstfun2 from 10.1.100.92:8500 success",
    "deregister service backupConstfun2 from 10.1.100.91:8500 success"
  ]
}
```