package router

import (
	"github.com/gin-gonic/gin"
	"gitlab.jiangxingai.com/poseidon/client/router/api/v1/activate"
	"gitlab.jiangxingai.com/poseidon/client/router/api/v1/check"
	"gitlab.jiangxingai.com/poseidon/client/router/api/v1/container"
	"gitlab.jiangxingai.com/poseidon/client/router/api/v1/logManager"
	"gitlab.jiangxingai.com/poseidon/client/router/api/v1/systemManager"
	"gitlab.jiangxingai.com/poseidon/client/router/api/v1/systemd"
	"gitlab.jiangxingai.com/poseidon/client/router/api/v1/user"
)

func InitRouter(r *gin.Engine) *gin.Engine {
	//AI 外部调用接口
	r.GET("/ip", systemManager.GetIp)
	system := r.Group("/api/v1/system")
	system.GET("/message", systemManager.GetSystemMessage)
	system.GET("/date", systemManager.GetDateAndTimezone)
	system.POST("/date", systemManager.ManualDate)
	system.PUT("/date", systemManager.NtpDate)
	system.GET("/network", systemManager.GetNetworkMsg)
	system.PUT("/network/dhcp", systemManager.DhcpNetwork)
	system.POST("/network/static", systemManager.UpdateStaticIp)
	system.GET("/wireless", systemManager.GetWirelessMsg)
	system.PUT("/wireless", systemManager.SwitchWireless)
	system.POST("/wireless", systemManager.UpdateApn)
	system.GET("/ssd", systemManager.MountSsd)
	system.POST("/passwd", systemManager.UpdatePasswd)

	log := r.Group("/api/v1/log")
	log.GET("/jxcore", logManager.GetJxcoreLog)
	log.GET("/event", logManager.GetEventLog)

	v1 := r.Group("/api/v1")
	{
		// 用户相关
		v1.POST("/user/login", user.LoginHandler)
		v1.PUT("/user/password/reset", user.UpdatePasswordHandler)
		v1.GET("/user/account/list", user.QueryUserListHandler)

		// 设备激活
		v1.POST("/device/activate", activate.RegisterToIotHandler)
		v1.POST("/device/register", activate.InstallAspRegisterHandler)

		// master 容器相关操作
		v1.GET("/master/container/list", container.QueryDockerListHandler)
		v1.GET("/master/container/logs", container.QueryDockerRunningLogsHandler)
		v1.POST("/master/container/switch", container.SwitchContainerHandler)
		v1.POST("/master/container/install", container.InstallContainerByYaml)
		v1.DELETE("/master/container/uninstall", container.UninstallContainerHandler)
		v1.POST("/master/container/upgrade", container.UpgradeDockerImageHandler)
		v1.POST("/master/container/upload", container.UploadModelHandler)
		v1.GET("/master/container/info", container.QueryDockerListToClassifyHandler)

		// 节点容器相关操作
		v1.GET("/node/container/hostList", container.QueryNodesModelHostListHandler)
		v1.GET("/node/list", container.QueryNodeIpsHandler)
		v1.GET("/node/container/list", container.QueryNodeContainerListHandler)
		v1.GET("/node/container/logs", container.QueryNodeContainerLogsHandler)
		v1.POST("/node/container/switch", container.SwitchNodeContainerHandler)
		v1.DELETE("/node/container/uninstall", container.UninstallNodeContainerHandler)
		v1.POST("/node/container/upgrade", container.UpgradeNodeContainerHandler)
		v1.POST("/node/container/upgradeMany", container.UpgradeNodesContainerHandler)

		// 软硬件检测
		v1.POST("/device/detection", check.HardwareDetectHandlerForMaster)
		// 内部接口
		v1.GET("/device/detection/node", check.HardwareDetectHandlerForNode)

		// 启动项管理
		v1.GET("/systemd/service/list", systemd.QueryAllServiceBySystemDHandler)
		v1.POST("/systemd/service/manager", systemd.OperateSystemDHandler)
	}
	return r
}
