package check

import (
	"encoding/json"
	"errors"
	"gitlab.jiangxingai.com/poseidon/client/pkg/config"
	e "gitlab.jiangxingai.com/poseidon/client/pkg/response"
	"gitlab.jiangxingai.com/poseidon/client/pkg/utils"
	"io/ioutil"
	"k8s.io/klog"
	"net"
	"net/http"
	"strings"
	"time"
)

const (
	hardwareCheckUrl = "/api/v1/device/detection/node"
	queryNodeIpsUrl  = "/api/v1/node/list"
)

var masterPort = strings.Replace(config.HttpPort(), ":", "", 1)

/**
  硬件检测
*/
func hardwareDetect() (string, error) {
	netmaskCommand := "bash /data/webConfig/board_check.sh"
	stdout, err := utils.CommandExecuteWithAllLogs(netmaskCommand)
	all := strings.ReplaceAll(string(stdout), "\u001b", "")
	return all, err
}

/**
  软件检测
*/
func softDetect() (string, error) {

	// todo  逻辑预留
	return "", nil
}

/**
  阵列节点检测
*/
func hardwareCheck(nodeIp string) (interface{}, error, int) {
	hostPort := net.JoinHostPort(nodeIp, masterPort)
	url := "http://" + hostPort + hardwareCheckUrl
	connect := nodeConnectStatusCheck(hostPort)
	if !connect {
		return nil, errors.New("节点访问失败"), e.ERROR_CONNECT_NODE
	}
	klog.Info("hardware check url:", url)
	return doGet(url, "GET")
}

func queryNodeIps() (interface{}, error, int) {
	url := "http://" + net.JoinHostPort("localhost", masterPort) + queryNodeIpsUrl
	return doGet(url, "GET")
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

func parseHttpResult(resp *http.Response) (interface{}, error, int) {
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
	return res.Data, err, statusCode
}

/**
  节点网络状态探测
*/
func nodeConnectStatusCheck(host string) bool {
	_, err := net.DialTimeout("tcp", host, time.Duration(time.Second))
	if err != nil {
		klog.Info("Site unreachable, error: ", err)
		return false
	}
	return true
}
