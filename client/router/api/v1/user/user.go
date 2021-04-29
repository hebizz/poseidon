package user

import (
	"github.com/gin-gonic/gin"
	"gitlab.jiangxingai.com/poseidon/client/driver"
	"gitlab.jiangxingai.com/poseidon/client/interfaces"
	e "gitlab.jiangxingai.com/poseidon/client/pkg/response"
	"k8s.io/klog"
	"net/http"
	"time"
)

type requestBody struct {
	UserName string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateBody struct {
	requestBody
	NewPassword string `json:"newPassword" binding:"required"`
}

const (
	userNormal = 1
	userAdmin  = 2
)

/**
  登录
*/
func LoginHandler(c *gin.Context) {
	response := &e.Gin{C: c}
	var body requestBody
	if err := c.BindJSON(&body); err != nil {
		response.Response(http.StatusBadRequest, e.INVALID_PARAMS, err, "")
		return
	}
	user, err := loadFromYaml()
	if err != nil {
		response.Response(http.StatusInternalServerError, e.ERROR, err, "")
		return
	}
	accounts := user.Acc
	for _, account := range accounts {
		if body.UserName == account.UserName && body.Password == account.Password {
			response.Response(http.StatusOK, e.SUCCESS, nil, gin.H{"accountType": account.Role})
			return
		}
	}
	// 记录登录日志
	if err = driver.InsertLog(interfaces.LogMsg{Timestamp: int(time.Now().Unix()), EventType: "操作", Description: "用户上线"}); err != nil {
		klog.Error(err)
	}
	response.Response(http.StatusBadRequest, e.ERROR_USER_OR_PASSWORD, nil, "")
}

/**
  修改密码
*/
func UpdatePasswordHandler(c *gin.Context) {
	response := &e.Gin{C: c}
	var body UpdateBody
	if err := c.BindJSON(&body); err != nil {
		response.Response(http.StatusBadRequest, e.INVALID_PARAMS, err, "")
		return
	}
	user, err := loadFromYaml()
	if err != nil {
		response.Response(http.StatusInternalServerError, e.ERROR, err, "")
		return
	}
	accounts := user.Acc
	for _, account := range accounts {
		if body.UserName == account.UserName && body.Password == account.Password {
			var err error
			if account.Role == "normal" {
				err = setNewValueToYamlFile(userNormal, body.NewPassword)
			} else {
				err = setNewValueToYamlFile(userAdmin, body.NewPassword)
			}
			if err != nil {
				response.Response(http.StatusInternalServerError, e.ERROR, err, "")
				return
			}
			response.Response(http.StatusOK, e.SUCCESS, nil, "")
			return
		}
	}
	response.Response(http.StatusBadRequest, e.ERROR_USER_OR_PASSWORD, nil, "")
}

/**
  获取用户列表
*/
func QueryUserListHandler(c *gin.Context) {
	response := &e.Gin{C: c}
	userList, err := loadFromYaml()
	if err != nil {
		response.Response(http.StatusInternalServerError, e.ERROR, err, "")
		return
	}
	response.Response(http.StatusOK, e.SUCCESS, nil, userList)
}
