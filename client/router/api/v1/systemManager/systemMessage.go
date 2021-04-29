package systemManager

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gitlab.jiangxingai.com/poseidon/client/interfaces"
	"gitlab.jiangxingai.com/poseidon/client/pkg/config"
	e "gitlab.jiangxingai.com/poseidon/client/pkg/response"
	"gitlab.jiangxingai.com/poseidon/client/pkg/utils"
)

//获取系统基本信息
func GetSystemMessage(c *gin.Context) {
	app := e.Gin{C: c}
	var data interfaces.SystemMsg
	deviceName, err := ioutil.ReadFile(config.HostNameFilePath)
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR_HOSTNAME, err, nil)
		return
	}
	deviceSerial := deviceName
	deviceType, err := ioutil.ReadFile(config.DeviceTypeFilePath)
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR_DEVICE, err, nil)
		return
	}
	firmware, err := ioutil.ReadFile(config.FirmwareFilePath)
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR_FIRMWARE, err, nil)
		return
	}
	data.DeviceName = strings.Replace(string(deviceName), "\n", "", 1)
	data.DeviceSerial = strings.Replace(string(deviceSerial), "\n", "", 1)
	data.DeviceType = strings.Replace(string(deviceType), "\n", "", 1)
	data.Firmware = strings.Replace(string(firmware), "\n", "", 1)
	data.WebVersion = config.BuildVersion

	app.Response(http.StatusOK, e.SUCCESS, nil, data)
}

//挂载ssd, 挂载目录待確定
func MountSsd(c *gin.Context) {
	app := e.Gin{C: c}
	ssdMountPath := utils.ReadString(config.SsdMountPath)
	ssdPath, err := GetMountDevice()
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR_SSD, err, nil)
		return
	}
	err = MountDeviceHandler(ssdPath, ssdMountPath)
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR_SSD_DEVICE, err, nil)
		return
	}
	err = InAspW2s(ssdMountPath)
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR_IN_SSD, err, nil)
		return
	}
	app.Response(http.StatusOK, e.SUCCESS, nil, nil)
}


//修改密碼 todo: 校驗/etc/shadow
func UpdatePasswd(c *gin.Context) {
	app := e.Gin{C: c}
	var data interfaces.PasswordMsg
	err := c.ShouldBindJSON(&data)
	if err != nil {
		app.Response(http.StatusBadRequest, e.INVALID_PARAMS, err, nil)
		return
	}
	err = UpdatePasswordHandler(data)
	if err != nil {
		app.Response(http.StatusBadRequest, e.ERROR_PASSWORD, err, nil)
		return
	}
	app.Response(http.StatusOK, e.SUCCESS, nil, nil)

}
