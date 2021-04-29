package logManager

import (
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gitlab.jiangxingai.com/poseidon/client/driver"
	"gitlab.jiangxingai.com/poseidon/client/pkg/config"
	e "gitlab.jiangxingai.com/poseidon/client/pkg/response"
	"gitlab.jiangxingai.com/poseidon/client/pkg/utils"
	log "k8s.io/klog"
)

//下载jxcore日志
func GetJxcoreLog(c *gin.Context) {
	app := e.Gin{C: c}
	jxcoreLogPath := utils.ReadString(config.JxcoreLogPath)
	_, err := os.Stat(jxcoreLogPath)
	if err != nil {
		app.Response(http.StatusBadRequest, e.ERROR_JXCORE, nil, nil)
		return
	}
	c.Header("Content-Type", "text/plain")
	c.Header("Content-Disposition", "attachment; filename=jxcore_event.txt")
	c.File(jxcoreLogPath)
}

//获取事件日志
func GetEventLog(c *gin.Context) {
	app := e.Gin{C: c}
	startTime, errS := strconv.ParseInt(c.Query("startTime"), 10, 64)
	endTime, errE := strconv.ParseInt(c.Query("endTime"), 10, 64)
	eventType := c.Query("type")
	limit, errL := strconv.Atoi(c.Query("limit"))
	offset, errO := strconv.Atoi(c.Query("offset"))
	if errS != nil || errE != nil || errL != nil || errO != nil {
		log.Error(errS, errE, errL, errO)
		app.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, nil)
		return
	}
	res, count, err := driver.QueryLog(startTime, endTime, eventType, limit, offset)
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR_GET_LOG, err, nil)
		return
	}
	app.Response(http.StatusOK, e.SUCCESS, nil, gin.H{"data":res, "count": count})
}
