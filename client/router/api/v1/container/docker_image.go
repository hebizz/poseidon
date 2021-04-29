package container

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gitlab.jiangxingai.com/poseidon/client/pkg/config"
	"gitlab.jiangxingai.com/poseidon/client/pkg/docker"
	e "gitlab.jiangxingai.com/poseidon/client/pkg/response"
	"gitlab.jiangxingai.com/poseidon/client/pkg/utils"
	"io/ioutil"
	"k8s.io/klog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

/**
  docker run
*/

type RequestGetBody struct {
	ContainerId string `form:"containerId" binding:"required"`
}

type RequestDeleteBody struct {
	ContainerName string `form:"containerName" binding:"required"`
}

type CommonRequestPostBody struct {
	ContainerName string `json:"containerName" binding:"required"`
	Switch        string `json:"switch" binding:"required"`
	ContainerId   string `json:"containerId"`
}

type UpgradeRequestBody struct {
	ContainerId   string `json:"containerId"`
	ContainerName string `json:"containerName"`
	ImageName     string `json:"imageName"`
	Password      string `json:"password"`
	UserName      string `json:"userName"`
	Url           string `json:"url"`
}

/**
  查询所有container列表,并按照ai model 分类
*/
func QueryDockerListToClassifyHandler(c *gin.Context) {
	response := e.Gin{C: c}
	containers, err := docker.ListContainer(true)
	if err != nil {
		response.Response(http.StatusBadRequest, e.ERROR, err, "")
		return
	}
	res := make([]*AiContainer, 0)
	for _, container := range containers {
		act := new(AiContainer)
		containerJSON, err := docker.InspectByContainerName(container.ID)
		if err != nil {
			continue
		}
		split := strings.Split(container.ImageID, ":")
		if len(split) == 2 && len(split[1]) >= 10 {
			rs := []rune(split[1])
			act.ImageHash = string(rs[0:10])
		}
		portMap := containerJSON.HostConfig.PortBindings
		act.Port = make([]string, 0)
		for _, hostPorts := range portMap {
			act.Port = append(act.Port, hostPorts[0].HostPort)
		}
		act.State = container.State
		res = append(res, act)
	}
	response.Response(http.StatusOK, e.SUCCESS, nil, res)
}

/**
  查询所有container列表
*/
func QueryDockerListHandler(c *gin.Context) {
	response := e.Gin{C: c}
	containers, err := docker.ListContainer(true)
	if err != nil {
		response.Response(http.StatusBadRequest, e.ERROR, err, "")
		return
	}
	res := make([]*Docker, 0)
	for _, container := range containers {
		dck := new(Docker)
		dck.ContainerId = container.ID
		dck.CreateTime = container.Created
		dck.ImageID = container.ImageID
		dck.Name = container.Image
		dck.State = container.State
		dck.Version = SplitImageNameToVersion(container.Image)
		dck.ContainerName = strings.Replace(container.Names[0], "/", "", 1)
		// 判断目录是否存在
		path := config.DefaultYamlPath + dck.ContainerName
		if Exists(path) {
			dck.IsEdit = true
		} else {
			dck.IsEdit = false
		}
		res = append(res, dck)
	}
	response.Response(http.StatusOK, e.SUCCESS, nil, res)
}

/**
  查看运行日志
*/
func QueryDockerRunningLogsHandler(c *gin.Context) {
	response := e.Gin{C: c}
	var body RequestGetBody
	queryBinding := binding.Query
	klog.Info(queryBinding.Name())
	if err := c.ShouldBindWith(&body, queryBinding); err != nil {
		response.Response(http.StatusBadRequest, e.INVALID_PARAMS, err, "")
		return
	}
	t := time.Now()
	tm1 := time.Date(t.Year(), t.Month(), t.Day()-1, 0, 0, 0, 0, t.Location())
	tm2 := tm1.AddDate(0, 0, 1)

	logReader, err := docker.QueryContainerLogsByDate(body.ContainerId, strconv.FormatInt(tm2.Unix(), 10), strconv.FormatInt(t.Unix(), 10))
	if err != nil {
		response.Response(http.StatusBadRequest, e.ERROR_NO_CONTAINER, err, "")
		return
	}
	defer logReader.Close()
	content, err := ioutil.ReadAll(logReader)
	length := len(content)

	// header -> Content-Type = text/plain
	if length > 15000 {
		filePath := config.DefaultContainerTmpLog + body.ContainerId + ".txt"
		if err := SaveByteToFile(content, filePath); err != nil {
			klog.Errorf("save docker logs error: +%V", err)
			// 删除日志文件
			_ = os.Remove(filePath)
		} else {
			c.FileAttachment(filePath, body.ContainerId+".txt")
			klog.Info(" file type")
			// 删除日志文件
			_ = os.Remove(filePath)
			return
		}
	}
	all := strings.ReplaceAll(string(content), "\ufffd", "")
	all = strings.ReplaceAll(all, "\u0000", "")
	all = strings.ReplaceAll(all, "\u0001", "")
	all = strings.ReplaceAll(all, "\u007f", "")

	// header -> Content-Type = application/json
	response.Response(http.StatusOK, e.SUCCESS, nil, all)
}

/**
  镜像升级
*/

func UpgradeDockerImageHandler(c *gin.Context) {
	response := e.Gin{C: c}
	headerYaml, err := c.FormFile("yaml")
	containerName := c.PostForm("containerName")

	if err != nil || headerYaml == nil {
		response.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, "")
		return
	}
	// 保存文件到临时目录下
	unixNano := utils.GetMD5Str(time.Now().UnixNano() / 1e6)
	tmpDir := config.DefaultYamlTmpPath + unixNano
	tmpPath := tmpDir + "/docker-compose.yaml"
	_ = Mkdir(tmpDir)
	if err = c.SaveUploadedFile(headerYaml, tmpPath); err != nil {
		response.Response(http.StatusBadRequest, e.ERROR, err, "")
		return
	}
	// 上传yaml文件
	containerNameFromYaml, err := ReadContainerNameFromYamlFile(tmpPath)
	if err != nil || containerNameFromYaml == "" {
		response.Response(http.StatusBadRequest, e.ERROR_YAML_FORMAT, err, "")
		return
	}
	// containerName 升级前后必须一致

	klog.Info("=----------------", containerName)
	klog.Info("=----------------", containerNameFromYaml)

	if containerName != containerNameFromYaml {
		response.Response(http.StatusBadRequest, e.ERROR_YAML_CONTAINERNAME_ERROR, err, "")
		return
	}
	dir := config.DefaultYamlPath + containerNameFromYaml
	path := dir + "/docker-compose.yaml"
	_ = Mkdir(dir)

	// 使用固定的fileName,保存yaml时已经重命名为docker-compose.yaml
	_, err = DockerComposeCommand(path, "down")
	if err != nil {
		// 如果 down 失败，继续执行up指令
		klog.Errorf("docker-compose down fail:  %+v", err)
	}
	// 将文件移动到当前目录下
	if err := os.Rename(tmpPath, path); err != nil {
		response.Response(http.StatusInternalServerError, e.ERROR, err, "")
		klog.Errorf("mv yaml dir error: %+v", err)
		return
	}
	// 删除临时文件
	_ = os.Remove(tmpPath)

	// docker-compose 运行docker
	_, err = DockerComposeCommand(path, "up", "-d")
	if err != nil {
		klog.Errorf("docker-compose up error: %+v", err)
		response.Response(http.StatusBadRequest, e.ERROR, err, "")
		return
	}
	response.Response(http.StatusOK, e.SUCCESS, nil, "")
}

/**
  开启 / 停止container
*/
func SwitchContainerHandler(c *gin.Context) {
	response := e.Gin{C: c}
	var body CommonRequestPostBody
	if err := c.BindJSON(&body); err != nil {
		response.Response(http.StatusBadRequest, e.INVALID_PARAMS, err, "")
		return
	}
	switchStatus := body.Switch
	command := "start"
	if switchStatus == "off" {
		command = "stop"
	}
	// 查询当前container status (预留逻辑)

	// 使用固定的fileName,保存yaml时已经重命名为docker-compose.yaml
	path := config.DefaultYamlPath + body.ContainerName + "/docker-compose.yaml"

	if !Exists(path) {
		response.Response(http.StatusBadRequest, e.ERROR_NO_CONTAINER, nil, "")
		return
	}
	dockerComposeCommand, err := DockerComposeCommand(path, command)
	klog.Info(string(dockerComposeCommand))
	if err != nil {
		response.Response(http.StatusBadRequest, e.ERROR, err, "")
		return
	}
	response.Response(http.StatusOK, e.SUCCESS, nil, "")
}

/**
  卸载container
*/
func UninstallContainerHandler(c *gin.Context) {
	response := e.Gin{C: c}
	var body RequestDeleteBody
	if err := c.ShouldBindWith(&body, binding.Query); err != nil {
		response.Response(http.StatusBadRequest, e.INVALID_PARAMS, err, "")
		return
	}
	// 使用固定的fileName,保存yaml时已经重命名为docker-compose.yaml
	path := config.DefaultYamlPath + body.ContainerName + "/docker-compose.yaml"

	if !Exists(path) {
		response.Response(http.StatusBadRequest, e.ERROR_NO_CONTAINER, nil, "")
		return
	}
	dockerComposeCommand, err := DockerComposeCommand(path, "down")
	if err != nil {
		response.Response(http.StatusBadRequest, e.ERROR, err, "")
		return
	}
	klog.Info(string(dockerComposeCommand))
	// 删除tar包
	imageFilePath := config.DefaultContainerTmpDocker + body.ContainerName
	_ = os.RemoveAll(imageFilePath)

	yamlFilePath := config.DefaultYamlPath + body.ContainerName
	_ = os.RemoveAll(yamlFilePath)
	response.Response(http.StatusOK, e.SUCCESS, nil, "")
}

/**
  通过yaml安装image (master only)
*/
func InstallContainerByYaml(c *gin.Context) {
	response := e.Gin{C: c}
	headerYaml, err := c.FormFile("yaml")
	if err != nil || headerYaml == nil {
		response.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, "")
		return
	}
	// check yaml
	// 保存文件到临时目录下
	unixNano := utils.GetMD5Str(time.Now().UnixNano() / 1e6)
	tmpDir := config.DefaultYamlTmpPath + unixNano
	tmpPath := tmpDir + "/docker-compose.yaml"
	_ = Mkdir(tmpDir)
	if err = c.SaveUploadedFile(headerYaml, tmpPath); err != nil {
		response.Response(http.StatusBadRequest, e.ERROR, err, "")
		return
	}
	// 上传yaml文件
	containerNameFromYaml, err := ReadContainerNameFromYamlFile(tmpPath)
	if err != nil || containerNameFromYaml == "" {
		response.Response(http.StatusBadRequest, e.ERROR_YAML_FORMAT, err, "")
		return
	}
	dir := config.DefaultYamlPath + containerNameFromYaml
	path := dir + "/docker-compose.yaml"
	_ = Mkdir(dir)
	// 将文件移动到当前目录下
	if err := os.Rename(tmpPath, path); err != nil {
		response.Response(http.StatusInternalServerError, e.ERROR, err, "")
		klog.Error("mv yaml dir error:", err)
		return
	}
	// 删除临时文件
	_ = os.Remove(tmpPath)
	// docker-compose 运行docker
	output, err := DockerComposeCommand(path, "up", "-d")
	klog.Errorf("docker-compose up error: %+v", string(output))
	if err != nil {
		response.Response(http.StatusBadRequest, e.ERROR, err, "")
		return
	}
	response.Response(http.StatusOK, e.SUCCESS, nil, "")
}
