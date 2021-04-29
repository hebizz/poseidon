package systemd

import (
	"errors"
	"github.com/gin-gonic/gin"
	e "gitlab.jiangxingai.com/poseidon/client/pkg/response"
	"net/http"
	"reflect"
	"strings"
)

type Service struct {
	Unit  string `json:"unit"`
	State string `json:"state"`
	Sub   string `json:"status"`
}

type RequestBody struct {
	Command string `json:"command" binding:"required"`
	Unit    string `json:"unit" binding:"required"`
}

const (
	CommandStart   = "start"
	CommandStop    = "stop"
	CommandRestart = "restart"
	CommandEnable  = "enable"
	CommandDisable = "disable"
)

/**
  获取systemD 管理的所有service
*/
func QueryAllServiceBySystemDHandler(c *gin.Context) {
	response := e.Gin{C: c}
	// 读取systemD service 列表
	serviceList, err := ExecuteFindServiceCommand()
	if err != nil {
		response.Response(http.StatusInternalServerError, e.ERROR, err, "")
		return
	}
	// 读取systemD 开机启动项
	enableServiceList, err := ExecuteFindEnableServiceCommand()
	if err != nil {
		response.Response(http.StatusInternalServerError, e.ERROR, err, "")
		return
	}
	res := make(map[string]Service)
	for _, line := range enableServiceList {
		split := strings.Split(line, ":")
		if len(split) != 2 {
			continue
		}
		service := Service{Unit: split[0], State: split[1]}
		res[split[0]] = service
	}
	for _, line := range serviceList {
		split := strings.Split(line, ":")
		if len(split) != 4 {
			continue
		}
		service := res[split[0]]
		if reflect.DeepEqual(service, Service{}) {
			service = Service{Unit: split[0], Sub: split[3]}
		} else {
			service.Sub = split[3]
		}
		res[split[0]] = service
	}
	//取两个查询列表的交集
	result := make([]Service, 0)
	for _, re := range res {
		if re.Sub == "" || re.State == "" {
			continue
		}
		result = append(result, re)
	}
	response.Response(http.StatusOK, e.SUCCESS, nil, result)
}

/**
  通过systemD操作service
*/
func OperateSystemDHandler(c *gin.Context) {
	response := e.Gin{C: c}
	var body RequestBody
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Response(http.StatusBadRequest, e.INVALID_PARAMS, err, "")
		return
	}
	var err error
	switch body.Command {
	case CommandStart:
		_, err = StartService(body.Unit)
		break
	case CommandStop:
		_, err = StopService(body.Unit)
		break
	case CommandRestart:
		_, err = RestartService(body.Unit)
		break
	case CommandEnable:
		_, err = EnableService(body.Unit)
		break
	case CommandDisable:
		_, err = DisableService(body.Unit)
		break
	default:
		response.Response(http.StatusBadRequest, e.ERROR_COMMAND, errors.New("invalid command type"), "")
		return
	}
	if err != nil {
		response.Response(http.StatusInternalServerError, e.ERROR, err, "")
		return
	}
	response.Response(http.StatusOK, e.SUCCESS, nil, "")
}
