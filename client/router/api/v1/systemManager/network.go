package systemManager

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.jiangxingai.com/poseidon/client/interfaces"
	"gitlab.jiangxingai.com/poseidon/client/pkg/config"
	e "gitlab.jiangxingai.com/poseidon/client/pkg/response"
)

//获取所有网卡信息
func GetNetworkMsg(c *gin.Context) {
	app := e.Gin{C: c}
	var data []interfaces.NetworkMsg
	//获取物理网卡名
	netDeviceList := GetLocalNetDeviceNames()
	for _, netDeviceName := range netDeviceList {
		res, err := GetNetDeviceMsg(interfaces.NetworkMsg{Name: netDeviceName})
		if err != nil {
			app.Response(http.StatusInternalServerError, e.ERROR, err, nil)
			return
		}
		data = append(data, res)
	}
	app.Response(http.StatusOK, e.SUCCESS, nil, data)
}

//修改指定网卡静态ip
func UpdateStaticIp(c *gin.Context) {
	var data interfaces.NetworkMsg
	app := e.Gin{C: c}
	data.Name = c.PostForm("name")
	data.Ip = c.PostForm("ip")
	data.Gateway = c.PostForm("gateway")
	data.DnsServers = c.PostForm("dnsServer")
	if data.Ip == "" || data.Gateway == "" || data.DnsServers == "" || data.Name == "" {
		app.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, nil)
		return
	}
	flag := ParseNetworkMsg(data)
	if flag != e.SUCCESS {
		app.Response(http.StatusBadRequest, flag, nil, nil)
		return
	}
	err := GenerateStaticNetworkYaml(data)
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR_NETWORK_YAML, err, nil)
		return
	}
	err = ExecuteNetworkYaml()
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR_APPLY_NETPLAN, err, nil)
		return
	}
	app.Response(http.StatusOK, e.SUCCESS, nil, "")
}

//dhcp 指定网卡
func DhcpNetwork(c *gin.Context) {
	app := e.Gin{C: c}
	data := c.PostForm("name")
	netDeviceList := GetLocalNetDeviceNames()
	flag := false
	for _, deviceName := range netDeviceList {
		if deviceName == data {
			flag = true
			break
		}
	}
	if !flag {
		app.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, nil)
		return
	}
	err := GenerateDhcpNetworkYaml(data)
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR_NETWORK_YAML, err, nil)
		return
	}
	err = ExecuteNetworkYaml()
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR_APPLY_NETPLAN, err, nil)
		return
	}
	app.Response(http.StatusOK, e.SUCCESS, nil, nil)
}

//获取4g信息
func GetWirelessMsg(c *gin.Context) {
	app := e.Gin{C: c}
	var data interfaces.WirelessMsg
	var err error
	data.Switch = GetWirelessStatus()
	data.Operator = GetWirelessOperator()
	//未插入4g卡
	if data.Operator == "" {
		app.Response(http.StatusOK, e.SUCCESS, nil, data)
		return
	}
	data.Apn = GetWirelessApn()
	data.Id = GetWirelessId()
	//通过apn状态来判断4g状态
	if data.Apn != "" {
		data.SimStatus = "正常"
		data.Status = "拨号成功"
	} else {
		data.SimStatus = "失败"
		data.Status = "拨号失败"
	}
	data.Ip, err = GetWirelessIp()
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR_IP_4G, err, nil)
		return
	}
	data.Signal = GetWirelessSignal()
	app.Response(http.StatusOK, e.SUCCESS, nil, data)
}

//切换4g和有线网络状态
func SwitchWireless(c *gin.Context) {
	app := e.Gin{C: c}
	flag, err := strconv.ParseBool(c.PostForm("switch"))
	if err != nil {
		app.Response(http.StatusBadRequest, e.INVALID_PARAMS, err, nil)
		return
	}
	if GetWirelessOperator() == "" {
		app.Response(http.StatusBadRequest, e.ERROR_4G, nil, nil)
		return
	}
	if flag {
		err := Switch4g(config.Switch4gMetric)
		if err != nil {
			app.Response(http.StatusInternalServerError, e.ERROR_ON_4G, err, nil)
			return
		}
	} else {
		err := Switch4g(config.Default4gMetric)
		if err != nil {
			app.Response(http.StatusInternalServerError, e.ERROR_OFF_4G, err, nil)
			return
		}
	}
	app.Response(http.StatusOK, e.SUCCESS, nil, nil)
}

//添加apn
func UpdateApn(c *gin.Context) {
	app := e.Gin{C: c}
	apn := c.PostForm("apn")
	err := UpdateWirelessApn(apn)
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR_APN, err, nil)
		return
	}
	app.Response(http.StatusOK, e.SUCCESS, nil, nil)
}

// 获取本机ip, tun0, mac地址, 需指定特定网卡
func GetIp(c *gin.Context) {
	app := e.Gin{C: c}
	ethIp, tunIp, mac, err := GetLocalIp()
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR, err, nil)
		return
	}
	app.Response(http.StatusOK, e.SUCCESS, nil, gin.H{"eth0": ethIp, "tun0": tunIp, "mac": mac})
}
