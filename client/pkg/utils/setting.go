package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"os/exec"
	"reflect"
	"strconv"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
	log "k8s.io/klog"
)

var (
	viper_ = viper.New()
)

func Setup() {
	//开机netplan, 避免jxcore配置网络影响
	netplanCommand := "netplan apply"
	cmd := exec.Command("bash", "-c", netplanCommand)
	_, err := CommandExecuteLogs(cmd)
	if err != nil {
		log.Error(err)
	}
	//初始化配置config.yaml
	viper_.SetConfigName("config")
	viper_.AddConfigPath("./")
	viper_.SetConfigType("yaml")
	viper_.AutomaticEnv()
	if err = viper_.ReadInConfig(); err != nil {
		panic(err)
	}
}

func ReadString(key string) string {
	if viper_.IsSet(key) {
		if ret := viper_.Get(key); ret != nil {
			return cast.ToString(ret)
		} else {
			return ""
		}
	}
	return ""
}

func ReadMap(key string) map[string]interface{} {
	if viper_.IsSet(key) {
		if ret := viper_.Get(key); ret != nil {
			return cast.ToStringMap(ret)
		} else {
			return map[string]interface{}{}
		}
	}
	return map[string]interface{}{}
}

func CommandExecuteLogs(cmd *exec.Cmd) (string, error) {
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		return stderr.String(), errors.New(stderr.String())
	}
	log.Info("cmd execute result is :\n" + out.String())
	return out.String(), nil
}

func CommandExecuteWithAllLogs(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	log.Info("cmd execute result is :\n" + string(output))
	return string(output), err
}

func GenerateMd5(rawStr string) string {
	data := []byte(rawStr)
	md5Str := fmt.Sprintf("%x", md5.Sum(data))
	return md5Str
}

func GetMD5Str(timestamp int64) string {
	str := strconv.FormatInt(timestamp, 10)
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func IsExistItem(value interface{}, array interface{}) bool {
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(value, s.Index(i).Interface()) {
				return true
			}
		}
	}
	return false
}
