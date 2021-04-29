package container

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	e "gitlab.jiangxingai.com/poseidon/client/pkg/response"
	"gitlab.jiangxingai.com/poseidon/client/pkg/utils"
	"io/ioutil"
	"k8s.io/klog"
	"net/http"
	"os"
	"runtime"
	"strings"
)

type nodeIpBody struct {
	NodeIp string `form:"nodeIp" binding:"required"`
}

type requestGetBody struct {
	RequestGetBody
	nodeIpBody
}
type requestDeleteBody struct {
	RequestDeleteBody
	nodeIpBody
}

type requestPostBody struct {
	CommonRequestPostBody
	nodeIpBody
}

/**
  获取阵列1808AI ip+port
*/

func QueryNodesModelHostListHandler(c *gin.Context) {
	response := e.Gin{C: c}
	res, err := queryAiContainerImageInfo()
	if err != nil {
		response.Response(http.StatusBadRequest, e.ERROR, err, "")
		return
	}
	response.Response(http.StatusOK, e.SUCCESS, err, res)
}

/**
  获取阵列1808的ip列表
*/

func QueryNodeIpsHandler(c *gin.Context) {
	response := e.Gin{C: c}
	res, err := Parse1808IpByShellResult()
	if err != nil {
		response.Response(http.StatusBadRequest, e.ERROR, err, "")
		return
	}
	response.Response(http.StatusOK, e.SUCCESS, err, res)
}

/**
  通过IP 获取阵列1808上的docker list
*/

func QueryNodeContainerListHandler(c *gin.Context) {
	response := e.Gin{C: c}
	var body nodeIpBody
	if err := c.ShouldBindWith(&body, binding.Query); err != nil {
		response.Response(http.StatusBadRequest, e.INVALID_PARAMS, err, "")
		return
	}
	data, err, code := queryNodeDockerList(body.NodeIp)
	if err != nil || code != http.StatusOK {
		response.Response(http.StatusBadGateway, e.ERROR_CONNECT_NODE, err, "")
		return
	}
	response.Response(http.StatusOK, e.SUCCESS, err, data)
}

/**
  开启 / 关闭 节点container
*/

func SwitchNodeContainerHandler(c *gin.Context) {
	response := e.Gin{C: c}
	rawData, err := c.GetRawData()
	// 将读过的字节流重新放到body
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawData))
	var body requestPostBody
	if err := c.BindJSON(&body); err != nil {
		response.Response(http.StatusBadRequest, e.INVALID_PARAMS, err, "")
		return
	}
	klog.Info(len(rawData))
	data, err, code := switchNodeDockerStatus(rawData, body.NodeIp)
	if err != nil || code != http.StatusOK {
		response.Response(http.StatusBadGateway, e.ERROR_CONNECT_NODE, err, "")
		return
	}
	response.Response(http.StatusOK, e.SUCCESS, err, data)
}

/**
  卸载节点container
*/

func UninstallNodeContainerHandler(c *gin.Context) {
	response := e.Gin{C: c}
	var body requestDeleteBody
	if err := c.ShouldBindWith(&body, binding.Query); err != nil {
		response.Response(http.StatusBadRequest, e.INVALID_PARAMS, err, "")
		return
	}
	data, err, code := uninstallNodeDockerStatus(body.ContainerName, body.NodeIp)
	if err != nil || code != http.StatusOK {
		response.Response(http.StatusBadGateway, e.INVALID_PARAMS, err, "")
		return
	}
	response.Response(http.StatusOK, e.SUCCESS, err, data)
}

/**
  查看节点日志
*/

func QueryNodeContainerLogsHandler(c *gin.Context) {
	response := e.Gin{C: c}
	var body requestGetBody
	if err := c.ShouldBindWith(&body, binding.Query); err != nil {
		response.Response(http.StatusBadRequest, e.INVALID_PARAMS, err, "")
		return
	}
	data, err, code := downLoadNodeDockerLogs(body.ContainerId, body.NodeIp)
	if code == fileTypeLogsCode {
		// 返回文件类型的日志
		res := data.(map[string]string)
		c.FileAttachment(res["filePath"], res["fileName"])
		return
	}
	if err != nil || code != http.StatusOK {
		response.Response(http.StatusBadGateway, e.INVALID_PARAMS, err, "")
		return
	}
	response.Response(http.StatusOK, e.SUCCESS, err, data)
}

/**
  升级节点container
*/

func UpgradeNodeContainerHandler(c *gin.Context) {
	response := e.Gin{C: c}
	nodeIp := c.PostForm("nodeIp")

	if nodeIp == "" {
		response.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, "")
		return
	}
	header, err := c.FormFile("file")
	headerYaml, err := c.FormFile("yaml")
	if err != nil || header == nil || headerYaml == nil {
		klog.Error("upload file error:", err)
		response.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, "")
		return
	}
	// 保存到临时文件
	imagePath := "/tmp/" + header.Filename
	if err = c.SaveUploadedFile(header, imagePath); err != nil {
		response.Response(http.StatusBadRequest, e.ERROR, nil, "")
		return
	}
	yamlPath := "/tmp/" + headerYaml.Filename
	if err = c.SaveUploadedFile(headerYaml, yamlPath); err != nil {
		response.Response(http.StatusBadRequest, e.ERROR, nil, "")
		return
	}
	data, err, code := upgradeNodeDockerContainer(nodeIp, imagePath, yamlPath, )
	if err != nil || code != http.StatusOK {
		response.Response(http.StatusBadGateway, e.INVALID_PARAMS, err, "")
		return
	}
	// 删除临时文件
	_ = os.Remove(imagePath)
	_ = os.Remove(yamlPath)
	response.Response(http.StatusOK, e.SUCCESS, err, data)
}

/**
  批量升级节点container
*/
func UpgradeNodesContainerHandler(c *gin.Context) {
	response := e.Gin{C: c}
	nodeForm := c.PostForm("nodeIpList")
	if nodeForm == "" {
		response.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, "")
		return
	}
	nodeIps := strings.Split(nodeForm, ",")
	header, err := c.FormFile("file")
	headerYaml, err := c.FormFile("yaml")
	if err != nil || header == nil || headerYaml == nil {
		klog.Error("upload file error:", err)
		response.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, "")
		return
	}
	klog.Info("start save file to /tmp")
	// 保存到临时文件
	imagePath := "/tmp/" + header.Filename
	if !Exists(imagePath) {
		// file exist
		if err = c.SaveUploadedFile(header, imagePath); err != nil {
			response.Response(http.StatusBadRequest, e.ERROR, nil, "")
			return
		}

	}
	yamlPath := "/tmp/" + headerYaml.Filename
	if !Exists(yamlPath) {
		if err = c.SaveUploadedFile(headerYaml, yamlPath); err != nil {
			response.Response(http.StatusBadRequest, e.ERROR, nil, "")
			return
		}
	}

	klog.Info("------------------------start to upload-----------------------------------")

	channel := make(chan UpgradeResult, len(nodeIps))
	result := make([]UpgradeResult, 0)

	// 协程池大小，若文件过大会导致系统kill当前应用
	pool, err := utils.NewPool(1)
	if err != nil {
		klog.Infof("create go routing pool error: %+v", err)
	}

	for _, nodeIp := range nodeIps {
		_ = pool.Put(&utils.Task{
			Handler: func(v ...interface{}) {
				for _, ip := range v {
					klog.Infof("upload ------------- : %+v", ip)
					data, err, code := upgradeNodeDockerContainer(ip.(string), imagePath, yamlPath)

					klog.Infof("upload file code : %d", code)
					klog.Infof("upload file result : %+v", data)

					res := UpgradeResult{
						Data:   data,
						Code:   code,
						NodeIp: ip.(string),
						Err:    err,
					}
					go runtime.GC()
					channel <- res
				}
			},
			Params: []interface{}{nodeIp},
		})
	}

	for range nodeIps {
		res := <-channel
		result = append(result, res)
	}

	klog.Infof("--------------------upload result----------------------%+v", result)
	go runtime.GC()
	// 删除临时文件
	_ = os.Remove(imagePath)
	_ = os.Remove(yamlPath)
	response.Response(http.StatusOK, e.SUCCESS, err, result)
}

/**
  通过校验文件md5判断文件是否存在
*/

func CheckImageExistByMD5(c *gin.Context) {

}

/**
  批量升级节点container (利用Nginx做静态服务器实现节点对文件的下载)
*/
func UpgradeNodesContainerHandlerV2(c *gin.Context) {
	response := e.Gin{C: c}
	nodeForm := c.PostForm("nodeIpList")
	if nodeForm == "" {
		response.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, "")
		return
	}
	nodeIps := strings.Split(nodeForm, ",")
	header, err := c.FormFile("file")
	headerYaml, err := c.FormFile("yaml")
	if err != nil || header == nil || headerYaml == nil {
		klog.Error("upload file error:", err)
		response.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, "")
		return
	}
	klog.Info("start save file to /tmp")
	// 保存到临时文件
	imagePath := "/tmp/" + header.Filename
	if !Exists(imagePath) {
		// file exist
		if err = c.SaveUploadedFile(header, imagePath); err != nil {
			response.Response(http.StatusBadRequest, e.ERROR, nil, "")
			return
		}

	}
	yamlPath := "/tmp/" + headerYaml.Filename
	if !Exists(yamlPath) {
		if err = c.SaveUploadedFile(headerYaml, yamlPath); err != nil {
			response.Response(http.StatusBadRequest, e.ERROR, nil, "")
			return
		}
	}

	klog.Info("------------------------start to upload-----------------------------------")

	channel := make(chan UpgradeResult, len(nodeIps))
	result := make([]UpgradeResult, 0)

	// 协程池大小，若文件过大会导致系统kill当前应用
	pool, err := utils.NewPool(1)
	if err != nil {
		klog.Infof("create go routing pool error: %+v", err)
	}

	for _, nodeIp := range nodeIps {
		_ = pool.Put(&utils.Task{
			Handler: func(v ...interface{}) {
				for _, ip := range v {
					klog.Infof("upload ------------- : %+v", ip)
					data, err, code := upgradeNodeDockerContainer(ip.(string), imagePath, yamlPath)

					klog.Infof("upload file code : %d", code)
					klog.Infof("upload file result : %+v", data)

					res := UpgradeResult{
						Data:   data,
						Code:   code,
						NodeIp: ip.(string),
						Err:    err,
					}
					go runtime.GC()
					channel <- res
				}
			},
			Params: []interface{}{nodeIp},
		})
	}

	for range nodeIps {
		res := <-channel
		result = append(result, res)
	}

	klog.Infof("--------------------upload result----------------------%+v", result)
	go runtime.GC()
	// 删除临时文件
	_ = os.Remove(imagePath)
	_ = os.Remove(yamlPath)
	response.Response(http.StatusOK, e.SUCCESS, err, result)
}
