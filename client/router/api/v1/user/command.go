package user

import (
	"errors"
	"gitlab.jiangxingai.com/poseidon/client/pkg/config"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"k8s.io/klog"
	"os"
)

var userYamlFilePath = config.DefaultUserPath + "user.yaml"

func init() {
	_ = InitDefaultConfig()
}

/**
  初始化yaml文件配置
*/

func InitDefaultConfig() error {
	var err error
	// 写入默认账号到yaml 文件
	if !Exists(config.DefaultUserPath) {
		err = os.MkdirAll(config.DefaultUserPath, os.ModePerm)
	}
	if !Exists(userYamlFilePath) {
		err = createDefaultAccount()
	}
	return err
}

/**
  创建yaml文件，并写入默认账号
*/

func createDefaultAccount() error {
	acc := Account{}
	subAcc1 := Acc{UserName: NormalName, Role: NormalRole, Password: NormalPassword}
	subAcc2 := Acc{UserName: AdminName, Role: AdminRole, Password: AdminPassword}
	acc.Acc = append(acc.Acc, subAcc1, subAcc2)
	data, err := yaml.Marshal(&acc)
	if err != nil {
		klog.Info(err)
		return err
	}
	return ioutil.WriteFile(userYamlFilePath, data, os.ModePerm)
}

/**
  加载yaml文件
*/

func loadFromYaml() (*Account, error) {
	var user Account
	yamlS, err := ioutil.ReadFile(userYamlFilePath)
	if err != nil {
		if err = InitDefaultConfig(); err != nil {
			return &user, err
		} else {
			yamlS, _ = ioutil.ReadFile(userYamlFilePath)
		}
	}
	if err = yaml.Unmarshal(yamlS, &user); err != nil {
		return &user, errors.New("can not parse " + userYamlFilePath + " config")
	}
	return &user, nil
}

/**
  更新yaml文件
*/

func setNewValueToYamlFile(userType int, password string) error {
	acc := Account{}
	yamlS, err := ioutil.ReadFile(userYamlFilePath)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(yamlS, &acc); err != nil {
		return err
	}
	if userType == userNormal {
		acc.Acc[0].Password = password
	} else {
		acc.Acc[1].Password = password
	}
	d, err := yaml.Marshal(&acc)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(userYamlFilePath, d, os.ModePerm)
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
