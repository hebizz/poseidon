package systemManager

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gitlab.jiangxingai.com/poseidon/client/interfaces"
	"gitlab.jiangxingai.com/poseidon/client/pkg/config"
	e "gitlab.jiangxingai.com/poseidon/client/pkg/response"
)

//获取当前时间和时钟
func GetDateAndTimezone(c *gin.Context) {
	app := e.Gin{C: c}
	var data interfaces.DateAndTimezone
	data.Date = time.Now().Format(config.GolangBirthday)
	timezone, err := ioutil.ReadFile(config.TimezoneFilePath)
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR_TIMEZONE, nil, nil)
		return
	}
	data.Timezone = strings.Replace(string(timezone), "\n", "", 1)
	app.Response(http.StatusOK, e.SUCCESS, nil, data)
}

//手动校时
func ManualDate(c *gin.Context) {
	app := e.Gin{C: c}
	datetime := c.PostForm("date")
	if datetime == "" {
		app.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, nil)
		return
	}
	err := CloseNtp()
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR_CLOSE_NTP, err, nil)
		return
	}
	err = SetDateTime(datetime)
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR_SET_DATETIME, err, nil)
		return
	}
	err = HwClock()
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR_HWCLOCK, err, nil)
		return
	}
	app.Response(http.StatusOK, e.SUCCESS, nil, datetime)
}

//ntp校时
func NtpDate(c *gin.Context) {
	app := e.Gin{C: c}
	err := StartNtp()
	if err != nil {
		app.Response(http.StatusInternalServerError, e.ERROR_START_NTP, err, nil)
		return
	}
	//等待ntp服务器校准时间
	time.Sleep(time.Duration(500)*time.Millisecond)
	datetime := time.Now().Format(config.GolangBirthday)
	app.Response(http.StatusOK, e.SUCCESS, nil, datetime)
}