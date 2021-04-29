package activate

import (
	"fmt"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"gitlab.jiangxingai.com/poseidon/client/pkg/config"
	"gitlab.jiangxingai.com/poseidon/client/pkg/utils"
	log "k8s.io/klog"
)

/**
  注册到iot
*/

func registerToIot(host string, key string) (string, error) {
	//registerToIotCommand2 := "/edge/jxcore/bin/jxcore bootstrap -t jiangxing123 -m openvpn"
	registerToIotCommand := fmt.Sprintf("/edge/jxcore/bin/jxcore bootstrap -t %s -m openvpn --host=%s", key, host)
	log.Info("register to iot command: ", registerToIotCommand)
	res, err := utils.CommandExecuteWithAllLogs(registerToIotCommand)
	if err != nil {
		log.Infof("register to iot error: %+v", err)
		return "", err
	}
	// 注册后重启jxCore
	restartJxCoreCommand := "service jxcore restart"
	_, err = utils.CommandExecuteWithAllLogs(restartJxCoreCommand)
	if err != nil {
		log.Infof("restart jxCore error: %+v", err)
		return "", err
	}
	return res, nil
}

/**
  注册到综合平台
*/
func registerToCommonSystem(data RequestBody) (string, error) {
  err = UpdateAspConfigYaml(data.Key, data.Host)
  if err != nil {
    return "", err
  }
  // 避免手动重启wss的影响
  Asp.Msg = ""
  restartWssCommand := fmt.Sprintf("cd %s && docker-compose restart", utils.ReadString(config.WssPath))
  cmd := exec.Command("bash", "-c", restartWssCommand)
  _, err = utils.CommandExecuteLogs(cmd)
  if err != nil {
   return "", err
  }
  for {
    if Asp.Msg != "" {
      break
    }
  }
  if Asp.Status != 200 {
    return Asp.Msg, errors.New("register fail")
  }
  return Asp.Msg, nil
}

func UpdateAspConfigYaml(key string, addr string) error {
	viper_asp := viper.New()
	viper_asp.AddConfigPath("/etc/asp")
	viper_asp.SetConfigName("config")
	viper_asp.Set("pass", key)
	viper_asp.Set("address", addr)
	viper_asp.Set("access", utils.ReadString(config.AspUsername))
	err := viper_asp.WriteConfigAs("/etc/asp/config.yaml")
	if err != nil {
		return err
	}
	return nil
}
