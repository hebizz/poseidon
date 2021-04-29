package main

import (
	"github.com/gin-gonic/gin"
	"gitlab.jiangxingai.com/poseidon/client/driver"
	"gitlab.jiangxingai.com/poseidon/client/pkg/config"
	"gitlab.jiangxingai.com/poseidon/client/pkg/ssdp"
	"gitlab.jiangxingai.com/poseidon/client/pkg/utils"
	"gitlab.jiangxingai.com/poseidon/client/router/api/router"
	"k8s.io/klog"
)

func init() {
	config.Version()
	if config.GoBuildType == "amd" {
		utils.Setup()
		driver.Setup()
		go ssdp.SsdpClient()
	}
}

func main() {
	gin.SetMode(config.RunMode())
	engine := gin.Default()
	router.InitRouter(engine)
	if err := engine.Run(config.HttpPort()); err != nil {
		klog.Info(err)
	}
}
