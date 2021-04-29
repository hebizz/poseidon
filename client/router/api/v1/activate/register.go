package activate

import (
	"net/http"

	"github.com/gin-gonic/gin"
	e "gitlab.jiangxingai.com/poseidon/client/pkg/response"
	log "k8s.io/klog"
)

const (
	platformIot          = 1
	platformCommonSystem = 2
)

var (
	Asp    AspBody
	output string
	err    error
)

type RequestBody struct {
	Platform int    `json:"platform" binding:"required"`
	Key      string `json:"key" binding:"required"`
	Host     string `json:"host" binding:"required"`
}

type AspBody struct {
	Status int    `json:"status" binding:"required"`
	Msg    string `json:"msg" binding:"required"`
}

/**
  设备激活到iotEdge or 综合平台
*/
func RegisterToIotHandler(c *gin.Context) {
	response := e.Gin{C: c}
	var body RequestBody
	err = c.BindJSON(&body)
	if err != nil {
		response.Response(http.StatusBadRequest, e.INVALID_PARAMS, err, nil)
		return
	}
	if body.Platform == platformIot {
		output, err = registerToIot(body.Host, body.Key)
	} else if body.Platform == platformCommonSystem {
		output, err = registerToCommonSystem(body)
	}
	if err != nil {
		response.Response(http.StatusInternalServerError, e.ERROR, err, output)
		return
	}
	log.Info(output)
	response.Response(http.StatusOK, e.SUCCESS, nil, output)
}

/**
  上传wss注册信息
*/
func InstallAspRegisterHandler(c *gin.Context) {
	response := e.Gin{C: c}
	err = c.BindJSON(&Asp)
	log.Info(Asp)
	if err != nil {
		response.Response(http.StatusBadRequest, e.INVALID_PARAMS, err, nil)
		return
	}
	response.Response(http.StatusOK, e.SUCCESS, nil, nil)
}
