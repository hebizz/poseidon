package container

import (
	"fmt"
	"os/exec"
	"strings"
)

var (
	dockerComposeBin string // dockerCompose 命令
)

func init() {
	dockerComposeBin = "docker-compose"
}

/**
   执行 docker-compose 相关指令
 */

func DockerComposeCommand(filename, command string, args ...string) ([]byte, error) {
	cmd := genDockerComposeCmd(filename, command, args...)
	return ExecCmdWithShell(cmd)
}

func genDockerComposeCmd(filename, command string, args ...string) string {
	cmd := fmt.Sprintf("%s -f %s %s", dockerComposeBin, filename, command)
	if args != nil {
		cmd += " " + strings.Join(args, " ")
	}
	return cmd
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
