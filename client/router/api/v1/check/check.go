package check

import (
	"github.com/gin-gonic/gin"
	e "gitlab.jiangxingai.com/poseidon/client/pkg/response"
	"net/http"
)

const (
	hardware = 1
	soft     = 2
)

var (
	stdout string
	err    error
	res    interface{}
)

type RequestBody struct {
	Type int `json:"type"binding:"required"`
}

/**
  硬件检测(节点侧master上执行)
*/

func HardwareDetectHandlerForNode(c *gin.Context) {
	response := e.Gin{C: c}
	stdout, err = hardwareDetect()
	// 忽略err
	response.Response(http.StatusOK, e.SUCCESS, err, stdout)
}

/**
  硬件检测
*/

func HardwareDetectHandlerForMaster(c *gin.Context) {
	response := e.Gin{C: c}
	var body RequestBody
	err = c.BindJSON(&body)
	if err != nil || (body.Type != 1 && body.Type != 2) {
		response.Response(http.StatusBadRequest, e.INVALID_PARAMS, err, "")
		return
	}
	if body.Type == hardware {
		nodeIps, err, code := queryNodeIps()
		if err != nil {
			response.Response(code, e.ERROR, err, "")
			return
		}
		ips := nodeIps.([]interface{})
		if len(ips) == 0 {
			response.Response(http.StatusInternalServerError, e.ERROR_CONNECT_NODE, err, "")
			return
		}
		res, err, code = hardwareCheck(ips[0].(string))
		if err != nil {
			response.Response(code, e.ERROR, err, "")
			return
		}
	} else if body.Type == soft {
		//todo  no implement
		res, err = softDetect()
	}
	response.Response(http.StatusOK, e.SUCCESS, nil, res)
}
