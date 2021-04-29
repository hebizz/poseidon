package e

var MsgFlags = map[int]string{
  SUCCESS:        "success",
  ERROR:          "fail",
  INVALID_PARAMS: "请求参数错误",

  ERROR_HOSTNAME:      "读取设备名称失败",
  ERROR_DEVICE:        "读取设备型号失败",
  ERROR_FIRMWARE:      "读取固件版本失败",
  ERROR_SSD:           "挂载ssd失败",
  ERROR_SSD_DEVICE:    "没有可挂载设备",
  ERROR_PASSWORD:      "用户名或密码错误",
  ERROR_IN_SSD:       "软链接失败",

  ERROR_NETWORK_YAML:  "生成yaml文件错误",
  ERROR_COPY_YAML:     "复制yaml文件失败",
  ERROR_APPLY_NETPLAN: "apply netplan错误",
  ERROR_IP:            "ip地址格式错误",
  ERROR_GATEWAY:       "网关地址格式错误",
  ERROR_DNSSERVERS:    "dns服务器格式错误",
  ERROR_IPADDR:        "ip地址已被占用",
  ERROR_APN:           "添加apn失败",
  ERROR_OPERATOR:      "获取运营商失败",
  ERROR_4G:            "未插入4g卡",
  ERROR_ON_4G:         "开启4g网络失败",
  ERROR_OFF_4G:        "关闭4g网络失败",
  ERROR_IP_4G:         "获取4g ip地址失败",

  ERROR_TIMEZONE:     "获取系统时钟失败",
  ERROR_CLOSE_NTP:    "关闭ntp服务失败",
  ERROR_SET_DATETIME: "设置时间失败",
  ERROR_HWCLOCK:      "同步硬件时钟失败",
  ERROR_START_NTP:    "开启ntp服务失败",
  ERROR_JXCORE:       "无jxcore日志",
  ERROR_GET_LOG:      "获取日志失败",

  ERROR_INSERT_LOG:               "插入日志失败",
  ERROR_YAML_CONTAINERNAME_ERROR: "yaml文件中containerName不一致",
  ERROR_NO_CONTAINER:             "No such container",
  ERROR_LOAD_IMAGE:               "导入镜像文件失败",
  ERROR_RUN_IMAGE:                "运行镜像失败",
  ERROR_COMMAND:                  "指令类型错误",
  ERROR_FILE_TYPE:                "上传文件类型错误",
  ERROR_FILE_EXIST:               "上传文件已存在",
  ERROR_USER_OR_PASSWORD:         "用户名或者密码错误",
  ERROR_YAML_FORMAT:              "yaml文件格式错误",
  ERROR_CONNECT_NODE:             "节点连接失败",
}

func GetMsg(code int) string {
  msg, ok := MsgFlags[code]
  if ok {
    return msg
  }
  return MsgFlags[ERROR]
}
