[toc]

## WebConfig

### 查询容器列表

**url**: `/api/v1/master/container/list`

**method**: `GET`

**url params**: None


**success response**:

```
{
    "c": "200",
    "msg": "success",
    "data": [
        {
            "name": "redis",
            "version": "",
            "status": "running",
            "createTime": 1605756935,
            "containerId": "d1bb56e0de0286af3f90b5710cf0840478a37f4b4feda55f9acecbaa35ddf260",
            "imageId": "sha256:74d107221092875724ddb06821416295773bee553bbaf8d888ababe9be7b947f",
            "containerName": "reverent_clarke"
        }
    ]
}
```

**error response**

```
{
     "c": "500",
     "msg": "fail",
     "data": ""
}
```

---

#### 查看运行日志

**url**: `/api/v1/master/container/logs`

**method**: `GET`

**url params**:

```
- containerId=xxxx                                 # 从list列表中获取的containerId           [必填]
```

**request body**: None

**success response**

```
# 如果response header 中 Content-Type = application/json
{
   "c": "200",
    "msg": "success",
    "data": "
    M1:C 19 Nov 2020 03:35:36.088 # oO0OoO0OoO0Oo Redis is starting \n
    n1:C 19 Nov 2020 03:35:36.088 # Redis version=6.0.9, bits=64, commit=00000000, modified=0\n"
}

#  如果response header 中 Content-Type = text/plain
   返回的是file 格式，前端需要转成文件并提示用户保存到本地
    
```

**error response**

```
{
    "c": "30007",
    "msg": "No such container",
    "data": ""
}
```

---

#### 开启 / 关闭容器

**url**: `/api/v1/master/container/switch`

**method**: `POST`

**url params**: None

**request body**: 

```
{
     "containerName": xx                # 容器name      [必填]
     "switch": "on" / "off"             # 开启 / 关闭   [必填]
}
```

**success response**

```
{
    "c": "200",
    "msg": "success",
    "data": "
}
```

**error response**

```
{
    "c": "500",
    "msg": "fail",
    "data": ""
}
```

---

#### 卸载容器

**url**: `/api/v1/master/container/uninstall`

**method**: `DELETE`

**url params**: 

```
{
     "containerName": xx                # 容器name      [必填]
}
```

**request body**: None

**success response**

```
{
    "c": "200",
    "msg": "success",
    "data": "
}
```

**error response**

```
{
    "c": "30007",
    "msg": "No such container",
    "data": ""
}
```

---

#### 镜像升级

**url**: `/api/v1/master/container/upgrade`

**method**: `POST`

**url params**: None

**request body**: 

```
{
     "containerName": xx                # 容器name      [必填]
     "yaml"  File 类型                   # yaml文件，yaml文件中的containerName 与 传参containerName 必须保持一致      [必填]
}
```

**success response**

```
{
    "c": "200",
    "msg": "success",
    "data": "
}
```

**error response**

```
{
    "c": "30007",
    "msg": "No such container",
    "data": ""
}
```

---

#### 镜像上传/升级(上传zip or .tar.gz 格式的文件)

**url**: `/api/v1/master/container/upload`

**method**: `POST`

**url params**: None

**request body**: 

```
{
     "file": File 类型                   # docker image file （只支持.tar.gz 和 .zip 格式）                         [必填]
     "yaml"  File 类型                   # yaml文件，yaml文件中的containerName 与 传参containerName 必须保持一致      [必填]
}
```

**success response**

```
{
    "c": "200",
    "msg": "success",
    "data": "
}
```

**error response**

```
{
    "c": "500",
    "msg": "fail",
    "data": ""
}
```

---

#### 安装镜像(上传yaml)

**url**: `/api/v1/master/container/install`

**method**: `POST`

**url params**: None

**request body**: 

```
{
     "yaml"  File 类型                   # yaml文件                  [必填]
}
```

**success response**

```
{
    "c": "200",
    "msg": "success",
    "data": "
}
```

**error response**

```
{
    "c": "500",
    "msg": "fail",
    "data": ""
}
```

---

#### 获取节点ip列表

**url**: `/api/v1/node/list`

**method**: `GET`

**url params**: None

**request body**: None

**success response**

```
{
    "c": "200",
    "msg": "success",
    "data": [
        "10.10.0.70",
        "10.10.0.71",
        "10.10.0.72",
        "10.10.0.73"
    ]
}
```

**error response**

```
{
    "c": "500",
    "msg": "fail",
    "data": ""
}
```

---

#### 获取节点容器列表

**url**: `/api/v1/node/container/list`


**method**: `GET`

**url params**: 

``` 
- nodeIp=xxxx                                                       [必填]
```

**request body**: None

**success response**

```
{
    "c": "200",
    "msg": "success",
    "data": [
        {
            "name": "redis",
            "version": "",
            "status": "running",
            "createTime": 1605756935,
            "containerId": "d1bb56e0de0286af3f90b5710cf0840478a37f4b4feda55f9acecbaa35ddf260",
            "imageId": "sha256:74d107221092875724ddb06821416295773bee553bbaf8d888ababe9be7b947f",
            "containerName": "reverent_clarke"
        }
    ]
}
```

**error response**

```
{
    "c": "500",
    "msg": "fail",
    "data": ""
}
```

---

#### 查看节点运行日志

**url**: `/api/v1/node/container/logs`

**method**: `GET`

**url params**:

```
-  containerId=xxxx                                 # 从list列表中获取的containerId           [必填]
-  nodeIp : xxxx                                    # 节点ip        [必填]
```

**request body**: None

**success response**

```
# 如果response header 中 Content-Type = application/json
{
   "c": "200",
    "msg": "success",
    "data": "
    M1:C 19 Nov 2020 03:35:36.088 # oO0OoO0OoO0Oo Redis is starting \n
    n1:C 19 Nov 2020 03:35:36.088 # Redis version=6.0.9, bits=64, commit=00000000, modified=0\n"
}

#  如果response header 中 Content-Type = text/plain
   返回的是file 格式，前端需要转成文件并提示用户保存到本地
    
```

**error response**

```
{
    "c": "30007",
    "msg": "No such container",
    "data": ""
}
```

---

#### 开启 / 关闭容器

**url**: `/api/v1/node/container/switch`

**method**: `POST`

**url params**: None

**request body**: 

```
{
     "containerName": xx                # 容器name      [必填]
     "switch": "on" / "off"             # 开启 / 关闭   [必填]
     "nodeIp" : xxxx                    # 节点ip        [必填]
}
```

**success response**

```
{
    "c": "200",
    "msg": "success",
    "data": "
}
```

**error response**

```
{
    "c": "500",
    "msg": "fail",
    "data": ""
}
```

---

#### 卸载节点容器

**url**: `/api/v1/node/container/uninstall`

**method**: `DELETE`

**url params**: 

```
- containerName: xx                # 容器name      [必填]
- nodeIp : xxxx                    # 节点ip        [必填]

```

**request body**: None

**success response**

```
{
    "c": "200",
    "msg": "success",
    "data": "
}
```

**error response**

```
{
    "c": "30007",
    "msg": "No such container",
    "data": ""
}
```

---

#### 节点镜像升级 / 安装

**url**: `/api/v1/node/container/upgrade`

**method**: `POST`

**url params**: None

**request body**: 

```
{
     "containerName": xx                # 容器name      [必填]
     "yaml"  File 类型                   # yaml文件，yaml文件中的containerName 与 传参containerName 必须保持一致      [必填]
     "nodeIp" : xxxx                    # 节点ip        [必填]
}
```

**success response**

```
{
    "c": "200",
    "msg": "success",
    "data": "
}
```

**error response**

```
{
    "c": "30007",
    "msg": "No such container",
    "data": ""
}
```

---

#### 节点批量 升级 / 安装

**url**: `/api/v1/node/container/upgradeMany`

**method**: `POST`

**url params**: None

**request body**: 

```
{
     "file": File 类型                                  # 镜像文件                  [必填]
     "yaml"  File 类型                                  # yaml文件，yaml文件中的containerName 与 传参containerName 必须保持一致      [必填]
     "nodeIpList" : 10.10.0.70,10.10.0.71               # 节点ip列表,文本类型, ip之间用逗号隔开      [必填]
}
```

**success response**

```
{
    "c": "200",
    "msg": "success",
    "data": "
}
```

**error response**

```
{
    "c": "30007",
    "msg": "No such container",
    "data": ""
}
```

---

#### 设备激活到iot

**url**: `/api/v1/device/activate`

**method**: `POST`

**url params**: None

**request body**: 

```
{
    "platform":  注册平台  1 iot 2 综合平台           [必填]
    "key":       注册key                            [必填]
    "host":      注册平台host    ip:port             [必填]
}

```

**success response**

```
{
    "c": "200",
    "msg": "success",
    "data": "xxxx"
}
```

**error response**

```
{
    "c": "500",
    "msg": "fail",
    "data": ""
}
```

---

#### 软硬件检测

**url**: `/api/v1/device/detection`

**method**: `POST`

**url params**: None

**request body**: 

```
{
     "type": 1                           # 检测类型  1 硬件检测,  2 软件检测（未实现）      [必填]
}
```

**success response**

```
{
    "c": "200",
    "msg": "success",
    "data": "xxxx"                      // 输出检测结果
}
```

**error response**

```
{
    "c": "500",
    "msg": "fail",
    "data": ""
}
```

---

#### 获取systemd 管理的service 列表

**url**: `/api/v1/systemd/service/list`

**method**: `GET`

**url params**: None

**request body**: None

**success response**

```
{
    "c": "200",
    "msg": "success",
    "data": [
        {
            "unit": "NetworkManager-wait-online.service",
            "state": "enabled",                //enable 开机自启动 disable !开机自启动
            "status": "exited"                 // running 运行 dead 关闭 exited 退出
        },
        {
            "unit": "ondemand.service",
            "state": "enabled",
            "status": "dead"
        }
    ]
}
```

**error response**

```
{
    "c": "500",
    "msg": "fail",
    "data": ""
}
```

---

#### 通过systemd操作service

**url**: `/api/v1/systemd/service/manager`

**method**: `POST`

**url params**: None

**request body**: 

```
{
     "unit":        "irqbalance.service"                                          [必填]
     "command"      "stop"                //三个选项 stop, start, restart  （针对status 状态）        [必填]
}                                         //两个选项 enable, disable      (针对state 状态)
```

**success response**

```
{
    "c": "200",
    "msg": "success",
    "data": ""
}
```

**error response**

```
{
    "c": "500",
    "msg": "fail",
    "data": ""
}
```

---

#### 登录

**url**: `/api/v1/user/login`


**method**: `POST`

**url params**: None

**request body**: 

```
## 系统内置两个账号
管理员账号: admin/admin
普通账号:   jiangxing/123456

{
     "userName": "admin"                           [必填]
     "password"  "admin"                           [必填]
}
```

**success response**

```
{
    "c": "200",
    "msg": "success",
    "data": {
        "accountType":"root"                  # accountType: root(管理员) normal(普通用户)
    }
}
```

**error response**

```
{
    "c": "20006",
    "msg": "用户名或者密码错误",
    "data": ""
}
```

---

#### 修改密码

**url**: `/api/v1/user/password/reset`


**method**: `PUT`

**url params**: None

**request body**: 

```
## 系统内置两个账号
管理员账号: admin/admin
普通账号:   jiangxing/123456

{
     "userName":    "admin"                           [必填]
     "password"     "admin"                           [必填]
     "newPassword"  "admin123"                        [必填]
}
```

**success response**

```
{
    "c": "200",
    "msg": "success",
    "data": ""
}
```

**error response**

```
{
    "c": "20006",
    "msg": "用户名或者密码错误",
    "data": ""
}
```

---

#### 获取账号列表

**url**: `/api/v1/user/account/list`


**method**: `GET`

**url params**: None

**request body**:  None

**success response**

```
{
    "c": "200",
    "msg": "success",
    "data": {
        "accounts": [
            {
                "userName": "jiangxing",
                "password": "123456",
                "role": "normal"
            },
            {
                "userName": "admin",
                "password": "admin123",
                "role": "root"
            }
        ]
        }
}
```

**error response**

```
{
    "c": "20006",
    "msg": "用户名或者密码错误",
    "data": ""
}
```


---

#### 获取节点ai类型和访问地址 (后端调用)

**url**: `/api/v1/node/container/hostList`

**method**: `GET`

**url params**: None

**request body**: None

**success response**

```
{
    "c": "200",
    "msg": "success",
    "data": {
        "ada9e2da90": {                    // image sha256
            "address": [
                "10.10.0.70:5000"
            ],
            "extra_command": {}
        }
    }
}
```

**error response**

```
{
    "c": "500",
    "msg": "fail",
    "data": ""
}
```

---