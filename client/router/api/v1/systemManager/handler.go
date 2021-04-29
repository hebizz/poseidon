package systemManager

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/go-ping/ping"
	"github.com/spf13/viper"
	"gitlab.jiangxingai.com/poseidon/client/interfaces"
	"gitlab.jiangxingai.com/poseidon/client/pkg/config"
	e "gitlab.jiangxingai.com/poseidon/client/pkg/response"
	"gitlab.jiangxingai.com/poseidon/client/pkg/utils"
	log "k8s.io/klog"
)

func CheckIp(target string) bool {
	pinger, err := ping.NewPinger(target)
	if err != nil {
		return true
	}
	pinger.Count = -1
	pinger.Timeout = time.Second
	pinger.SetPrivileged(true)
	pinger.Run() // blocks until finished
	stats := pinger.Statistics()
	fmt.Println(stats)
	if stats.PacketsRecv >= 1 {
		return true
	}
	return false
}

func ParseNetworkMsg(data interfaces.NetworkMsg) int {
	ip, ipNet, err := net.ParseCIDR(data.Ip)
	if err != nil {
		log.Error(err)
		return e.ERROR_IP
	}
	log.Infof("ip is %s, ipnet is %s ", ip, ipNet)
	res := CheckIp(ip.String())
	if res {
		return e.ERROR_IPADDR
	}
	gateway := net.ParseIP(data.Gateway)
	if gateway == nil {
		return e.ERROR_GATEWAY
	}
	log.Infof("gateway is %s", gateway.String())
	dnsServers := net.ParseIP(data.DnsServers)
	if dnsServers == nil {
		return e.ERROR_DNSSERVERS
	}
	log.Infof("dnsServers is %s", dnsServers.String())
	return e.SUCCESS
}

func GenerateStaticNetworkYaml(data interfaces.NetworkMsg) error {
	var insertIp, insertDns []string
	insertIp = append(insertIp, data.Ip)
	insertDns = append(insertDns, data.DnsServers)
	viper_static := viper.New()
	viper_static.AddConfigPath("./")
	addressKey := fmt.Sprintf("network.ethernets.%s.addresses", data.Name)
	gatewayKey := fmt.Sprintf("network.ethernets.%s.gateway4", data.Name)
	dnsKey := fmt.Sprintf("network.ethernets.%s.nameservers.addresses", data.Name)
	log.Info(addressKey, insertIp, gatewayKey, data.Gateway, dnsKey, insertDns)
	viper_static.Set(addressKey, insertIp)
	viper_static.Set(gatewayKey, data.Gateway)
	viper_static.Set(dnsKey, insertDns)
	err := viper_static.WriteConfigAs(utils.ReadString(config.NetworkConfigFile))
	if err != nil {
		return err
	}
	return nil
}

func GenerateDhcpNetworkYaml(name string) error {
	viper_dhcp := viper.New()
	viper_dhcp.AddConfigPath("./")
	dhcpKey := fmt.Sprintf("network.ethernets.%s.dhcp4", name)
	viper_dhcp.Set(dhcpKey, true)
	err := viper_dhcp.WriteConfigAs(utils.ReadString(config.NetworkConfigFile))
	if err != nil {
		return err
	}
	return nil
}

func ExecuteNetworkYaml() error {
	copyCommand := "cp " + utils.ReadString(config.NetworkConfigFile) + " " + config.NetPlanDirPath
	cmd := exec.Command("bash", "-c", copyCommand)
	_, err := utils.CommandExecuteLogs(cmd)
	if err != nil {
		return err
	}
	netplanCommand := "netplan apply"
	cmd2 := exec.Command("bash", "-c", netplanCommand)
	_, err = utils.CommandExecuteLogs(cmd2)
	if err != nil {
		return err
	}
	return nil
}

func GetLocalNetDeviceNames() []string {
	cmd := exec.Command("ls", config.NetworkCardPath)
	buf, _ := cmd.Output()
	output := string(buf)
	var res []string
	for _, device := range strings.Split(output, "\n") {
		if len(device) > 1 {
			if device != "lo" && strings.Index(device, "e") == 0 {
				res = append(res, device)
			}
		}
	}
	return res
}

func GetNetDeviceMsg(data interfaces.NetworkMsg) (interfaces.NetworkMsg, error) {
	inter, err := net.InterfaceByName(data.Name)
	if err != nil {
		log.Error(err)
		return data, err
	}
	addrs, err := inter.Addrs()
	if err != nil {
		log.Error(err)
		return data, err
	}
	if len(addrs) > 0 {
		data.Ip = addrs[0].String()
	} else {
		data.Ip = ""
	}

	gatewayCommand := "netstat -rn |grep " + data.Name + " |grep  0.0.0.0  |awk '{print $1}'" + " |tail -n 1"
	cmd3 := exec.Command("bash", "-c", gatewayCommand)
	gateway, err := utils.CommandExecuteLogs(cmd3)
	if err != nil {
		gateway = ""
	}
	res := strings.Split(gateway, "\n")
	if len(res) > 0 {
		gateway = res[0]
	}
	data.Gateway = gateway

	dnsCommand := "tail -n 1 /etc/dnsmasq.resolv.conf"
	cmd4 := exec.Command("bash", "-c", dnsCommand)
	dns, err := utils.CommandExecuteLogs(cmd4)
	if err != nil {
		dns = ""
	} else {
		dnsList := strings.Split(dns, " ")
		if len(dnsList) > 0 {
			dns = dnsList[1]
		} else {
			dns = ""
		}
	}
	data.DnsServers = dns

	data.Mode = GetNetDeviceMode(data.Name)
	return data, nil
}

func GetNetDeviceMode(name string) string {
	var mode string
	viper_read := viper.New()
	viper_read.AddConfigPath("/etc/netplan")
	viper_read.SetConfigName("01-network-manager-all")
	viper_read.SetConfigType("yaml")
	if err := viper_read.ReadInConfig(); err != nil {
		log.Info("network file not found")
		mode = "dhcp"
	}
	ipKey := fmt.Sprintf("network.ethernets.%s.addresses", name)
	flag := viper_read.IsSet(ipKey)
	if flag {
		mode = "static"
	} else {
		mode = "dhcp"
	}
	return mode
}

func GetWirelessStatus() bool {
	routeCommand := "ip route |sed -n '1p'"
	cmd := exec.Command("bash", "-c", routeCommand)
	res, err := utils.CommandExecuteLogs(cmd)
	if err != nil {
		return false
	}
	return strings.Contains(res, "usb0")
}

func GetWirelessApn() string {
	var apn string
	getApnCommand := fmt.Sprintf("%s %s %s",
		config.PythonVersion, utils.ReadString(config.WirelessPyPath), config.GetApnAtCommand)
	cmd := exec.Command("bash", "-c", getApnCommand)
	res, err := utils.CommandExecuteLogs(cmd)
	if err != nil {
		log.Error(err)
		return apn
	}
	resList := strings.Split(res, ":")
	if len(resList) > 1 {
		apnList := strings.Split(resList[1], "\n")
		apn = apnList[0]
		if strings.Contains(apn, "\r") {
			apn = strings.Replace(apn, "\r", "", 1)
		}
		if strings.Contains(apn, " ") {
			apn = strings.Replace(apn, " ", "", 1)
		}
	} else {
		apn = ""
	}
	return apn
}

func GetWirelessId() string {
	var id string
	getIotIdCommand := fmt.Sprintf("%s %s %s",
		config.PythonVersion, utils.ReadString(config.WirelessPyPath), config.GetIdAtCommand)
	cmd := exec.Command("bash", "-c", getIotIdCommand)
	res, err := utils.CommandExecuteLogs(cmd)
	if err != nil {
		log.Error(err)
		return id
	}
	resList := strings.Split(res, "\n")
	if len(resList) > 1 {
		id = resList[1]
		if strings.Contains(id, "\r") {
			id = strings.Replace(id, "\r", "", 1)
		}
	} else {
		id = ""
	}
	return id
}

func GetWirelessOperator() string {
	var operator string
	getOperatorCommand := fmt.Sprintf("%s %s %s",
		config.PythonVersion, utils.ReadString(config.WirelessPyPath), config.GetOperatorAtCommand)
	cmd := exec.Command("bash", "-c", getOperatorCommand)
	res, err := utils.CommandExecuteLogs(cmd)
	if err != nil {
		log.Error(err)
		return operator
	}
	resList := strings.Split(res, `"`)
	if len(resList) > 1 {
		operator = resList[1]
	} else {
		operator = ""
	}
	return operator
}

func GetWirelessIp() (string, error) {
	var ip string
	inter, err := net.InterfaceByName(config.Default4gInterfaceName)
	if err != nil {
		return ip, err
	}
	addrs, err := inter.Addrs()
	if err != nil {
		log.Error(err)
		return ip, err
	}
	return addrs[0].String(), nil
}

func GetWirelessSignal() int {
	var signal int
	getSignalCommand := fmt.Sprintf("%s %s %s",
		config.PythonVersion, utils.ReadString(config.WirelessPyPath), config.GetSignalAtCommand)
	cmd := exec.Command("bash", "-c", getSignalCommand)
	res, err := utils.CommandExecuteLogs(cmd)
	if err != nil {
		log.Error(err)
		return signal
	}
	resList := strings.Split(res, ",")
	if len(resList) > 1 {
		resList2 := strings.Split(resList[0], " ")
		if len(resList2) > 1 {
			res = resList2[1]
		}
	}
	re, _ := strconv.Atoi(res)
	switch {
	case re == 0:
		signal = 0
	case re == 1:
		signal = 1
	case re >= 2 && re <= 31:
		signal = 2
	case re > 31:
		signal = 3
	}
	return signal
}

func UpdateWirelessApn(apn string) error {
	updateApnCommand := fmt.Sprintf("%s %s %s'%s'",
		config.PythonVersion, utils.ReadString(config.WirelessPyPath), config.UpdateApnAtCommand, apn)
	log.Info(updateApnCommand)
	cmd := exec.Command("bash", "-c", updateApnCommand)
	res, err := utils.CommandExecuteLogs(cmd)
	if err != nil {
		log.Error(err)
	}
	if strings.Contains(res, "ERROR") {
		return errors.New("update apn fail")
	}
	return nil
}

func GetMountDevice() (string, error) {
	lsblkCommand := "lsblk -S -o Name,TYPE |grep disk"
	cmd := exec.Command("bash", "-c", lsblkCommand)
	res, err := utils.CommandExecuteLogs(cmd)
	if err != nil {
		return "", err
	}
	res = strings.ReplaceAll(res, "disk", "")
	res = strings.ReplaceAll(res, "\n", "")
	res = strings.TrimSpace(res)
	resList := strings.Split(res, "  ")
	log.Info(resList)
	log.Info(len(resList))
	device, flag := IsMounted(resList)
	if flag {
		return "", errors.New("没有可挂载设备")
	}

	ssdPath := fmt.Sprintf("/dev/%s", device)
	log.Info(ssdPath)
	mkfsCommand := fmt.Sprintf("mkfs -t ext4 %s", ssdPath)
	cmd = exec.Command("bash", "-c", mkfsCommand)
	_, err = utils.CommandExecuteLogs(cmd)
	if err != nil {
		//todo: 解除格式化分区, 根据mkfs返回日志判断
		log.Error(err)
	}
	return ssdPath, nil
}

func MountDeviceHandler(ssdPath string, ssdMountPath string) error {
	_ = os.Mkdir(ssdMountPath, 0777)
	mountCommand := fmt.Sprintf("mount %s %s", ssdPath, ssdMountPath)
	cmd := exec.Command("bash", "-c", mountCommand)
	_, err := utils.CommandExecuteLogs(cmd)
	if err != nil {
		return err
	}

	fstabCommand := fmt.Sprintf("\n%s %s auto defaults,nofail,comment=cloudconfig  0  2\n", ssdPath, ssdMountPath)
	f, err := os.OpenFile(config.FstabFilePath, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(fstabCommand)
	if err != nil {
		return err
	}
	return nil
}

func IsMounted(devices []string) (string, bool) {
	for _, device := range devices {
		//todo: 优化判断挂载点
		devicePath := config.SsdConfigFilePath + device + "1"
		if !Exist(devicePath) {
			log.Info(devicePath)
			return device, false
		}
	}
	return "", true
}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func InAspW2s(ssdMountPath string) error {
	aspPath := fmt.Sprintf("%s/asp", ssdMountPath)
	w2sPath := fmt.Sprintf("%s/w2s", ssdMountPath)
	_ = os.Mkdir(aspPath, 0777)
	_ = os.Mkdir(w2sPath, 0777)
	lnAspCommand := fmt.Sprintf("ln -s %s %s", aspPath, config.AspDataPath)
	cmd := exec.Command("bash", "-c", lnAspCommand)
	_, err := utils.CommandExecuteLogs(cmd)
	if err != nil {
		return err
	}
	lnW2sCommand := fmt.Sprintf("ln -s %s %s", w2sPath, config.W2sDataPath)
	cmd = exec.Command("bash", "-c", lnW2sCommand)
	_, err = utils.CommandExecuteLogs(cmd)
	if err != nil {
		return err
	}
	return nil
}

func Switch4g(metric int) error {
	//有可能不存在4g默认路由，则会报错
	deleteIpRouteCommand := fmt.Sprintf("ip route del  default via %s dev usb0", config.Default4gIpAddr)
	cmd := exec.Command("bash", "-c", deleteIpRouteCommand)
	_, err := utils.CommandExecuteLogs(cmd)
	if err != nil {
		log.Error(err)
	}

	addIpRouteCommand := fmt.Sprintf("ip route add  default via %s dev usb0 metric %d",
		config.Default4gIpAddr, metric)
	cmd = exec.Command("bash", "-c", addIpRouteCommand)
	_, err = utils.CommandExecuteLogs(cmd)
	if err != nil {
		return err
	}
	return nil
}

func StartNtp() error {
	startNtpCommand := "timedatectl set-ntp yes"
	cmd := exec.Command("bash", "-c", startNtpCommand)
	_, err := utils.CommandExecuteLogs(cmd)
	if err != nil {
		return err
	}
	return nil
}

func CloseNtp() error {
	closeNtpCommand := "timedatectl set-ntp no"
	cmd := exec.Command("bash", "-c", closeNtpCommand)
	_, err := utils.CommandExecuteLogs(cmd)
	if err != nil {
		return err
	}
	return nil
}

func SetDateTime(datetime string) error {
	dateCommand := "date --set \"" + datetime + "\""
	log.Info(dateCommand)
	cmd := exec.Command("bash", "-c", dateCommand)
	_, err := utils.CommandExecuteLogs(cmd)
	if err != nil {
		return err
	}
	return nil
}

func HwClock() error {
	hwClockCommand := "hwclock -w"
	cmd := exec.Command("bash", "-c", hwClockCommand)
	_, err := utils.CommandExecuteLogs(cmd)
	if err != nil {
		return err
	}
	return nil
}

func UpdatePasswordHandler(data interfaces.PasswordMsg) error {
	updatePasswordCommand := fmt.Sprintf("echo \"%s:%s\" |chpasswd",
		data.Account, data.NewPassword)
	cmd := exec.Command("bash", "-c", updatePasswordCommand)
	_, err := utils.CommandExecuteLogs(cmd)
	if err != nil {
		return err
	}
	return nil
}

func GetLocalIp() (string, string, string, error) {

	defer func() {
		if err := recover(); err != nil {
			//log.Errorf("%+v", err)
			return
		}
	}()

	var ethIp, tunIp, mac string
	inter, err := net.InterfaceByName(utils.ReadString(config.NetworkCard))
	if err != nil {
		return "", "", "", err
	}
	addrs, err := inter.Addrs()

	if err != nil {
		return "", "", "", err
	}
	addr := addrs[0]
	if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() {
		if ip.IP.To4() != nil {
			log.Info(ip.IP)
			ethIp = ip.IP.String()
		}
	}

	inter2, err := net.InterfaceByName("tun0")

	if err != nil {
		//log.Errorf("%+v", err)
		tunIp = ""
	} else {
		addrs2, err := inter2.Addrs()
		if err != nil {
			log.Errorf("%+v", err)
		}
		addr2 := addrs2[0]
		if ip2, ok := addr2.(*net.IPNet); ok && !ip2.IP.IsLoopback() {
			if ip2.IP.To4() != nil {
				log.Info(ip2.IP)
				tunIp = ip2.IP.String()
			}
		}
	}
	ifs, err := net.Interfaces()

	if err != nil {
		log.Errorf("%+v", err)
	}
	for _, inter := range ifs {
		if inter.Name == utils.ReadString(config.NetworkCard) {
			log.Info(inter.Name, inter.HardwareAddr)
			mac = inter.HardwareAddr.String()
		}
	}
	return ethIp, tunIp, mac, nil
}
