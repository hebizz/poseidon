package container

import (
	"gitlab.jiangxingai.com/poseidon/client/pkg/config"
	"gitlab.jiangxingai.com/poseidon/client/pkg/docker"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	e "gitlab.jiangxingai.com/poseidon/client/pkg/response"
	"gitlab.jiangxingai.com/poseidon/client/pkg/utils"
	"k8s.io/klog"
)

/**
  上传docker文件
  注： image 文件必须是通过 docker save imageName:imageVersion > target.tar/target.tar.gz 方式导出
       否则load image 后 imageName  和 imageVersion 都为None
*/
func UploadModelHandler(c *gin.Context) {
	response := e.Gin{C: c}
	header, err := c.FormFile("file")
	headerYaml, err := c.FormFile("yaml")

	if err != nil || header == nil || headerYaml == nil {
		klog.Error("upload file error:", err)
		response.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil, "")
		return
	}
	originName := header.Filename
	if !strings.HasSuffix(originName, ".tar.gz") && !strings.HasSuffix(originName, ".tar") && !strings.HasSuffix(originName, ".zip") {
		klog.Error("image file type error")
		response.Response(http.StatusBadRequest, e.ERROR_FILE_TYPE, nil, "")
		return
	}
	yamlFileName := headerYaml.Filename
	if !strings.HasSuffix(yamlFileName, ".yaml") {
		klog.Error("yaml file type error")
		response.Response(http.StatusBadRequest, e.ERROR_FILE_TYPE, nil, "")
		return
	}
	// 保存yaml文件到指定路径并重命名为docker-compose.yaml
	unixNano := utils.GetMD5Str(time.Now().UnixNano() / 1e6)
	tmpYamlFilePath := config.DefaultYamlTmpPath + unixNano + "/docker-compose.yaml"
	_ = Mkdir(config.DefaultYamlTmpPath + unixNano)

	if err = c.SaveUploadedFile(headerYaml, tmpYamlFilePath); err != nil {
		response.Response(http.StatusInternalServerError, e.ERROR, err, "")
		err := os.RemoveAll(tmpYamlFilePath)
		klog.Error("remove dir error:", err)
		return
	}

	// 从yaml 文件中读取containerName
	containerName, err := ReadContainerNameFromYamlFile(tmpYamlFilePath)
	if err != nil {
		klog.Errorf("read yaml content error: %+v", err)
		response.Response(http.StatusBadRequest, e.ERROR_YAML_FORMAT, err, "")
		return
	}
	// 创建目录
	if err = Mkdir(config.DefaultYamlPath + containerName); err != nil {
		klog.Error("mk yaml dir error:", err)
		response.Response(http.StatusInternalServerError, e.ERROR, err, "")
		return
	}

	realYamlFilePath := config.DefaultYamlPath + containerName + "/docker-compose.yaml"
	// 将文件移动到当前目录下
	if err := os.Rename(tmpYamlFilePath, realYamlFilePath); err != nil {
		response.Response(http.StatusInternalServerError, e.ERROR, err, "")
		klog.Error("mv yaml dir error:", err)
		return
	}
	// 删除临时yaml目录
	_ = os.RemoveAll(config.DefaultYamlTmpPath + unixNano)

	dir := config.DefaultContainerTmpDocker + containerName + "/"
	if !Exists(dir) {
		if err = Mkdir(dir); err != nil {
			klog.Error("create file dir error:", err)
			response.Response(http.StatusInternalServerError, e.ERROR, err, "")
			return
		}
	}

	// 将imageFile 保存为临时文件
	tmpDir := config.DefaultContainerTmpDocker + containerName + "/" + unixNano + "/"
	imageFileTmpPath := tmpDir + originName
	_ = Mkdir(tmpDir)

	if err = c.SaveUploadedFile(header, imageFileTmpPath); err != nil {
		response.Response(http.StatusInternalServerError, e.ERROR, err, "")
		err := os.RemoveAll(tmpDir)
		klog.Error("remove dir error:", err)
		return
	}

	originPath := dir + originName
	if Exists(originPath) {
		originModelMd5, _ := CalculateHash(originPath)
		newModelMd5, _ := CalculateHash(imageFileTmpPath)
		if originModelMd5 == newModelMd5 {
			// imageFile 已存在，只需判断对应的容器是否已经处于running的状态，running: 直接返回， exited: up
			containerJSON, err := docker.InspectByContainerName(containerName)
			if err != nil {
				klog.Errorf("inspect container error:%+v", err)
			} else {
				state := containerJSON.State.Status
				if state == "running" {
					_ = os.RemoveAll(tmpDir)
				} else {
					//up
					//docker-compose start docker container
					result, err := DockerComposeCommand(realYamlFilePath, "up", "-d")
					klog.Info(string(result))
					if err != nil {
						klog.Error("docker compose up error: +%+v",err)
						response.Response(http.StatusBadRequest, e.ERROR_RUN_IMAGE, err, "")
						return
					}
				}
				response.Response(http.StatusOK, e.SUCCESS, nil, "")
				return
			}
		}
	}

	// 将文件移动到当前目录下
	if err := os.Rename(imageFileTmpPath, originPath); err != nil {
		klog.Errorf("mv image file error:%+v", err)
		response.Response(http.StatusInternalServerError, e.ERROR, err, "")
		return
	}
	// 删除临时文件文目录
	_ = os.RemoveAll(tmpDir)

	file, err := os.Open(originPath)
	if err != nil {
		klog.Errorf("open image file error: %+v", err)
		response.Response(http.StatusInternalServerError, e.ERROR, err, "")
		return
	}
	defer file.Close()
	klog.Info("start load image ---")

	// docker load image
	if err = docker.LoadImage(file); err != nil {
		klog.Errorf("load image error: %+v", err)
		response.Response(http.StatusInternalServerError, e.ERROR_RUN_IMAGE, err, "")
		return
	}
	klog.Info("load image finish ---")
	//docker-compose start docker container
	result, err := DockerComposeCommand(realYamlFilePath, "up", "-d")
	klog.Info(string(result))
	if err != nil {
		klog.Errorf("docker-compose up container error: %+v", err)
		response.Response(http.StatusInternalServerError, e.ERROR, err, "")
		return
	}
	response.Response(http.StatusOK, e.SUCCESS, nil, "")
}



/**
  下载nginx服务器下的docker image
*/
func DownloadImageByRequestNginx(c *gin.Context) {






}
