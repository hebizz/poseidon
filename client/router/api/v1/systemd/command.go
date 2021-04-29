package systemd

import (
	"bufio"
	"errors"
	"gitlab.jiangxingai.com/poseidon/client/pkg/utils"
	"io"
	"k8s.io/klog"
	"os"
	"os/exec"
	"strings"
)

var (
	serviceCmdFormat string // 服务管理命令和使用方式，如：systemCtl/service
)

func init() {
	serviceCmdFormat = "systemctl {action} {name}"
	if _, err := ExecCmdWithShell("which systemctl"); err != nil {
		serviceCmdFormat = "service {name} {action}"
	}
}

/**
  查询所有systemd管理的service
*/
func ExecuteFindServiceCommand() ([]string, error) {
	systemDCommand := "systemctl list-units --type=service --all | grep -E 'running|dead' | grep -v 'not-found' |awk '{print $1\":\"$2\":\"$3\":\"$4}' > /tmp/systemd1.txt"
	_, err := utils.CommandExecuteWithAllLogs(systemDCommand)
	if err != nil {
		klog.Errorf("find systemd1 service fail: %+v", err)
	}
	return readFileData("/tmp/systemd1.txt")
}

/**
  查询开机自启动 or not
*/
func ExecuteFindEnableServiceCommand() ([]string, error) {
	systemDCommand := "systemctl list-unit-files --type=service --all | grep -E 'disable|enable' |awk '{print $1\":\"$2}' > /tmp/systemd2.txt"
	_, err := utils.CommandExecuteWithAllLogs(systemDCommand)
	if err != nil {
		klog.Errorf("find systemd2 service fail: %+v", err)
	}
	return readFileData("/tmp/systemd2.txt")
}

/**
  按行读取文件
*/
func readFileData(filePath string) ([]string, error) {
	res := make([]string, 0)
	// 按行读取
	if Exists(filePath) {
		file, err := os.Open(filePath)
		if err != nil {
			return res, err
		}
		reader := bufio.NewReader(file)
		for {
			line, _, err := reader.ReadLine()
			if err == io.EOF {
				break
			}
			res = append(res, string(line))
		}
	} else {
		return res, errors.New("query systemd service failed,try again")
	}
	return res, nil
}

/**
  文件 / 目录是否存在
*/
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func EnableService(name string) ([]byte, error) {
	cmd := genServiceCmd("enable", name)
	return ExecCmdWithShell(cmd)
}

func DisableService(name string) ([]byte, error) {
	cmd := genServiceCmd("disable", name)
	return ExecCmdWithShell(cmd)
}

func StartService(name string) ([]byte, error) {
	cmd := genServiceCmd("start", name)
	return ExecCmdWithShell(cmd)
}

func StopService(name string) ([]byte, error) {
	cmd := genServiceCmd("stop", name)
	return ExecCmdWithShell(cmd)
}

func RestartService(name string) ([]byte, error) {
	cmd := genServiceCmd("restart", name)
	return ExecCmdWithShell(cmd)
}

func genServiceCmd(action, name string) string {
	return strings.NewReplacer("{action}", action, "{name}", name).Replace(serviceCmdFormat)
}

func ExecCmdWithShell(name string, args ...string) ([]byte, error) {
	tempArgs := make([]string, 0)
	tempArgs = append(tempArgs, "-c") // bash 的选项参数
	tempArgs = append(tempArgs, name)
	for _, arg := range args {
		tempArgs = append(tempArgs, arg)
	}
	return exec.Command("/bin/bash", tempArgs...).CombinedOutput()
}
