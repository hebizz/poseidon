package container

import (
	"bytes"
	"encoding/json"
	"github.com/pkg/errors"
	"gitlab.jiangxingai.com/poseidon/client/pkg/config"
	e "gitlab.jiangxingai.com/poseidon/client/pkg/response"
	"io"
	"io/ioutil"
	"k8s.io/klog"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	_ "runtime/pprof"
	"strings"
	"time"
)

const (
	fileTypeLogsCode             = -1
	queryNodeDockerListUrl       = "/api/v1/master/container/list"
	switchNodeContainerStatusUrl = "/api/v1/master/container/switch"
	uninstallNodeContainerUrl    = "/api/v1/master/container/uninstall"
	downloadNodeContainerLogsUrl = "/api/v1/master/container/logs"
	uploadNodeContainerUrl       = "/api/v1/master/container/upload"
	queryNodeContainerInfosUrl   = "/api/v1/master/container/info"
)

var masterPort = strings.Replace(config.HttpPort(), ":", "", 1)

/**
  查询节点 ip 列表
*/
func queryNodeDockerList(nodeIp string) (interface{}, error, int) {
	hostPort := net.JoinHostPort(nodeIp, masterPort)
	url := "http://" + hostPort + queryNodeDockerListUrl
	connect := nodeConnectStatusCheck(hostPort)
	if !connect {
		return nil, errors.New("节点访问失败"), e.ERROR_CONNECT_NODE
	}
	klog.Info("query node docker list  url:", url)
	return doGet(url, "GET")
}

/**
  查询节点container详情
*/
func queryNodeDockerListInfo(nodeIp string) (interface{}, error, int) {
	hostPort := net.JoinHostPort(nodeIp, masterPort)
	url := "http://" + hostPort + queryNodeContainerInfosUrl
	connect := nodeConnectStatusCheck(hostPort)
	if !connect {
		return nil, errors.New("节点访问失败"), e.ERROR_CONNECT_NODE
	}
	klog.Info("query node docker info   url:", url)
	return doGet(url, "GET")
}

/**
  开启/关闭 节点container
*/
func switchNodeDockerStatus(rawData []byte, nodeIp string) (interface{}, error, int) {
	// 转发节点前先确认节点是否连通
	hostPort := net.JoinHostPort(nodeIp, masterPort)
	url := "http://" + hostPort + switchNodeContainerStatusUrl
	connect := nodeConnectStatusCheck(hostPort)
	if !connect {
		return nil, errors.New("节点访问失败"), e.ERROR_CONNECT_NODE
	}
	klog.Info("switch node container status  url:", url)
	return doPost(url, rawData, "application/json")
}

/**
  卸载节点container
*/
func uninstallNodeDockerStatus(containerName string, nodeIp string) (interface{}, error, int) {
	hostPort := net.JoinHostPort(nodeIp, masterPort)
	url := "http://" + hostPort + uninstallNodeContainerUrl + "?containerName=" + containerName
	connect := nodeConnectStatusCheck(hostPort)
	if !connect {
		return nil, errors.New("节点访问失败"), e.ERROR_CONNECT_NODE
	}
	klog.Info("uninstall node container  url:", url)
	return doGet(url, "DELETE")
}

/**
  查看节点container logs
*/
func downLoadNodeDockerLogs(containerId string, nodeIp string) (interface{}, error, int) {
	hostPort := net.JoinHostPort(nodeIp, masterPort)
	url := "http://" + hostPort + downloadNodeContainerLogsUrl + "?containerId=" + containerId
	connect := nodeConnectStatusCheck(hostPort)
	if !connect {
		return nil, errors.New("节点访问失败"), e.ERROR_CONNECT_NODE
	}
	klog.Info("download node container logs url:", url)
	return doGetLogs(url, "GET", containerId)
}

/**
  升级/安装 节点container
*/
func upgradeNodeDockerContainer(nodeIp string, imagePath string, yamlPath string) (interface{}, error, int) {
	hostPort := net.JoinHostPort(nodeIp, masterPort)
	url := "http://" + hostPort + uploadNodeContainerUrl
	connect := nodeConnectStatusCheck(hostPort)
	if !connect {
		return nil, errors.New("节点访问失败"), e.ERROR_CONNECT_NODE
	}
	klog.Info("upload node container url:", url)
	return NewFileUploadRequestV2(url, "file", imagePath, "yaml", yamlPath)
}

func NewFileUploadRequestV2(url string, paramName1, path1 string, paramName2, path2 string) (interface{}, error, int) {

	bodyBuffer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuffer)

	// file1
	file1, _ := os.Open(path1)
	fileWriter1, _ := bodyWriter.CreateFormFile(paramName1, strings.Split(file1.Name(), "/")[2])
	defer file1.Close()
	_, _ = io.Copy(fileWriter1, file1)

	// file2
	file2, _ := os.Open(path2)
	fileWriter2, _ := bodyWriter.CreateFormFile(paramName2, strings.Split(file2.Name(), "/")[2])
	defer file2.Close()
	_, _ = io.Copy(fileWriter2, file2)

	_ = bodyWriter.Close()
	request, err := http.NewRequest("POST", url, bodyBuffer)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	request.Header.Add("Content-Type", bodyWriter.FormDataContentType())
	request.Header.Add("Connection", "close")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	//err = pprof.WriteHeapProfile(fileWriter1)
	return parseHttpResult(resp)
}

func doGetLogs(url string, method string, containerId string) (interface{}, error, int) {
	req, _ := http.NewRequest(method, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	defer resp.Body.Close()
	contentType := resp.Header.Get("content-Type")
	// 返回文件类型
	if strings.Contains(contentType, "text/plain") {
		if resp.StatusCode == http.StatusOK {
			result, _ := ioutil.ReadAll(resp.Body)
			filePath := config.DefaultContainerTmpLog + containerId + ".txt"
			// if file exist, remove first
			if !IsDir(filePath) && Exists(filePath) {
				_ = os.Remove(filePath)
			}
			if err := SaveByteToFile(result, filePath); err != nil {
				return nil, err, http.StatusInternalServerError
			}
			res := map[string]string{
				"filePath": filePath,
				"fileName": containerId + ".txt",
			}
			return res, nil, fileTypeLogsCode
		}
	}
	return parseHttpResult(resp)
}

func doGet(url string, method string) (interface{}, error, int) {
	req, _ := http.NewRequest(method, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	defer resp.Body.Close()
	return parseHttpResult(resp)
}

func doPost(url string, rawData []byte, contentType string) (interface{}, error, int) {
	resp, err := http.Post(url, contentType, bytes.NewBuffer(rawData))
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	defer resp.Body.Close()
	return parseHttpResult(resp)
}

func parseHttpResult(resp *http.Response) (interface{}, error, int) {
	defer resp.Body.Close()
	statusCode := resp.StatusCode
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err, statusCode
	}
	var res e.Response
	err = json.Unmarshal(result, &res)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	klog.Infof("response data: %+v", res.Data)
	return res.Data, err, statusCode
}

/**
  节点网络状态探测
*/
func nodeConnectStatusCheck(host string) bool {
	_, err := net.DialTimeout("tcp", host, time.Duration(1*time.Second))
	if err != nil {
		klog.Info("Site unreachable, error: ", err)
		return false
	}
	return true
}
