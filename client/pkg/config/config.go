package config

import (
	"gitlab.jiangxingai.com/poseidon/client/pkg/utils"
	log "k8s.io/klog"
)

var (
	BuildVersion string
	BuildTime    string
	BuildHash    string
	GoVersion    string
	GoBuildType  string
)

const (
	//config.yaml配置参数, 可变
	serverRunModeKey      = "poseidon.server.RunMode"
	defaultServerRunMode  = "debug"
	serverHttpPortKey     = "poseidon.server.HttpPort"
	defaultServerHttpPort = ":9999"
	NetworkConfigFile     = "networkConfigFile"
	NetworkCard           = "networkCard"
	SsdMountPath          = "ssdMountPath"
	JxcoreLogPath         = "jxcoreLogPath"
	WirelessPyPath        = "wirelessPyPath"
	WssPath               = "wssPath"
	AspUsername           = "aspUsername"

	//系统调用参数
	GolangBirthday            = "2006/01/02 15:04:05"
	TimezoneFilePath          = "/etc/timezone"
	HostNameFilePath          = "/etc/hostname"
	DeviceTypeFilePath        = "/etc/device"
	FirmwareFilePath          = "/etc/firmware"
	AspDataPath               = "/data/local/asp"
	W2sDataPath               = "/data/local/w2s"
	NetworkCardPath           = "/sys/class/net/"
	NetPlanDirPath            = "/etc/netplan"
	SsdConfigFilePath         = "/dev/"
	FstabFilePath             = "/etc/fstab"
	PythonVersion             = "python"
	GetOperatorAtCommand      = "CMD AT+COPS?"
	GetApnAtCommand           = "CMD AT+CGDCONT?"
	GetIdAtCommand            = "CMD AT+CGSN"
	GetSignalAtCommand        = "CMD AT+CSQ"
	UpdateApnAtCommand        = "CMD AT+CGDCONT="
	Default4gInterfaceName    = "usb0"
	Default4gIpAddr           = "192.168.225.1"
	Default4gMetric           = 102
	Switch4gMetric            = 100
	DefaultUserPath           = "/data/webConfig/user/"
	DefaultYamlPath           = "/data/webConfig/yaml/"
	DefaultYamlTmpPath        = "/data/webConfig/yaml/tmp/"
	DefaultContainerTmpLog    = "/data/webConfig/logs/"
	DefaultContainerTmpDocker = "/data/webConfig/docker/"
	DefaultWebConfig          = "/data/webConfig/w2s/"
)

func Version() {
	log.Infof("Version: %s\nBuild Time: %s\nBuild Hash: %s\nGo Version: %s\nGoBuildType: %s\n",
		BuildVersion, BuildTime, BuildHash, GoVersion, GoBuildType)
}

func RunMode() string {
	if ret := utils.ReadString(serverRunModeKey); ret != "" {
		return ret
	} else {
		return defaultServerRunMode
	}
}

func HttpPort() string {
	if ret := utils.ReadString(serverHttpPortKey); ret != "" {
		return ret
	} else {
		return defaultServerHttpPort
	}
}
