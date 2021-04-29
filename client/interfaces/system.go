package interfaces

type SystemMsg struct {
	DeviceName   string `json:"deviceName"`   //设备名称
	DeviceSerial string `json:"deviceSerial"` //设备序列号
	DeviceType   string `json:"deviceType"`   //设备型号
	Firmware     string `json:"firmware"`     //固件版本
	WebVersion   string `json:"webVersion"`   //web版本
}

type NetworkMsg struct {
	Name       string `json:"name"`      //网卡名称
	Ip         string `json:"ip"`        //ip地址
	Gateway    string `json:"gateway"`   //网关地址
	DnsServers string `json:"dnsServer"` //dns服务器
	Mode       string `json:"mode"`      //网络模式(static, dhcp)
}

type WirelessMsg struct {
	Switch    bool   `json:"switch"`    //4g开关 (true, false)
	Operator  string `json:"operator"`  //运营商
	Apn       string `json:"apn"`       //APN
	Id        string `json:"id"`        //物联网序列号
	SimStatus string `json:"simStatus"` //SIM卡状态
	Status    string `json:"status"`    //拨号状态
	Ip        string `json:"ip"`        //IP地址
	Signal    int    `json:"signal"`    // 信号强度(0, 1, 2, 3)
}

type PasswordMsg struct {
	Account     string `json:"account"   binding:"required"`   //帐号
	RawPassword string `json:"rawPassword" binding:"required"` //原密码
	NewPassword string `json:"newPassword" binding:"required"` //新密码
}

type DateAndTimezone struct {
	Date     string `json:"date"`     //时间
	Timezone string `json:"timezone"` //时钟
}
