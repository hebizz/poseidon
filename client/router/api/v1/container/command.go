package container

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"

	"gitlab.jiangxingai.com/poseidon/client/router/api/v1/systemManager"
	"gopkg.in/yaml.v3"

	//"github.com/kylelemons/go-gypsy/yaml"
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"gitlab.jiangxingai.com/poseidon/client/pkg/config"
	"gitlab.jiangxingai.com/poseidon/client/pkg/utils"
	"k8s.io/klog"
)

func init() {
	if !Exists(config.DefaultContainerTmpLog) {
		_ = os.MkdirAll(config.DefaultContainerTmpLog, os.ModePerm)
	}

	if !Exists(config.DefaultYamlPath) {
		_ = os.MkdirAll(config.DefaultYamlPath, os.ModePerm)
	}

	if !Exists(config.DefaultYamlTmpPath) {
		_ = os.MkdirAll(config.DefaultYamlTmpPath, os.ModePerm)
	}

	if !Exists(config.DefaultContainerTmpDocker) {
		_ = os.MkdirAll(config.DefaultContainerTmpDocker, os.ModePerm)
	}

	if config.GoBuildType == "amd" {
		if !Exists(config.DefaultWebConfig) {
			_ = os.MkdirAll(config.DefaultWebConfig, os.ModePerm)
		}

		go watchPoseidonYaml()

		go ExecShellToFindIps(nil)
	}
}

const PoseidonConfigYaml = config.DefaultWebConfig + "poseidon.yaml"

/**
  执行脚本扫描ip
*/
func ExecShellToFindIps(fChan chan []string) {
	command := "nmap -sP -PI -PT 10.10.0.0/24 | grep -v '10.10.0.100' | grep -v '10.10.0.200' | grep '10.10.0' | awk '{print $5}' > /tmp/1808.ip"
	_, err := utils.CommandExecuteWithAllLogs(command)
	if err != nil {
		klog.Errorf("find 1808 ips fail: %+v", err)
	}
	data, _ := readFileData("/tmp/1808.ip")
	fChan <- data
}

/**
  获取阵列的ip
*/
func Parse1808IpByShellResult() ([]string, error) {
	filePath := "/tmp/1808.ip"
	res, err := readFileData(filePath)
	if err == nil && len(res) > 0 {
		return res, err
	} else {
		// 重新读取
		// 刷新当前节点，下一次请求有效
		fChan := make(chan []string, 1)
		go ExecShellToFindIps(fChan)
		select {
		case res := <-fChan:
			return res, nil
		}
	}
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
		return res, errors.New("search 1808 ips failed,try again")
	}
	return res, nil
}

/**
  计算model的MD5值
*/
func CalculateHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		klog.Error(err)
		return "", err
	}
	defer file.Close()
	hash := md5.New()
	if _, err = io.Copy(hash, file); err != nil {
		return "", err
	}
	hashInBytes := hash.Sum(nil)[:16]
	md5Str := hex.EncodeToString(hashInBytes)
	return md5Str, nil
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

/**
  是否为目录
*/
func IsDir(path string) bool {
	ss, err := os.Stat(path)
	if err != nil {
		return false
	}
	return ss.IsDir()
}

/**
  创建目录
*/
func Mkdir(filePath string) error {
	if !Exists(filePath) {
		return os.MkdirAll(filePath, os.ModePerm)
	}
	return nil
}

/**
  保存日志到文件
*/
func SaveByteToFile(data []byte, fileName string) error {
	// 2. create file to write
	var f *os.File
	var err error
	exists := Exists(fileName)
	if exists {
		f, err = os.OpenFile(fileName, os.O_APPEND, os.ModePerm)
	} else {
		f, err = os.Create(fileName)
	}
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(data)
	return err
}

/**
 * 通过imageName 截取version
 */
func SplitImageNameToVersion(imageName string) string {
	if imageName != "" {
		splitArray := strings.Split(imageName, ":")
		if len(splitArray) > 1 {
			return splitArray[len(splitArray)-1]
		}
	}
	return ""
}

/**
 * 按照containerName保存yaml文件 (containerName 必须唯一，升级前后必须保持一致)
 */
func ReadContainerNameFromYamlFile(path string) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	var ser YamlServices
	file, err := os.Open(path)
	defer file.Close()
	decode := yaml.NewDecoder(file)
	err = decode.Decode(&ser)
	var container Container
	for _, value := range ser.Services.(map[string]interface{}) {
		klog.Infof("service name is: %+v", value)
		// key is service name
		bytes, _ := json.Marshal(value)
		_ = json.Unmarshal(bytes, &container)
	}
	klog.Infof("yaml file content: %#v", container)
	return container.ContainerName, err
}

/**
  创建w2s调用的poseidon.yaml文件
*/

func initDefaultPoseidonFile() {
	// 获取hostIp
	ethIp, _, _, err := systemManager.GetLocalIp()
	if err != nil {
		klog.Errorf("get local host ip error: %+v", err)
	}
	// 获取containerList
	res, err := queryAiContainerImageInfo()
	if err != nil {
		klog.Errorf("query node ai container name error: %+v", err)
	}
	poseidon := Poseidon{
		HostIp:      ethIp,
		ChipPointer: res,
	}
	data, err := yaml.Marshal(&poseidon)
	if err != nil {
		klog.Info(err)
		return
	}
	poseidonYaml, err := loadPoseidonYaml()
	if err != nil {
		// 插入文件
		err = ioutil.WriteFile(PoseidonConfigYaml, data, os.ModePerm)
		if err != nil {
			klog.Errorf("create poseidon config error: %+v", err)
		}
	} else {
		err := updatePoseidonYamlFile(poseidonYaml, poseidon)
		if err != nil {
			klog.Errorf("update poseidon config error: %+v", err)
		}
	}
}

/**
  加载yaml文件
*/

func loadPoseidonYaml() (Poseidon, error) {
	var poseidon Poseidon
	yamlS, err := ioutil.ReadFile(PoseidonConfigYaml)
	if err != nil {
		return poseidon, err
	}
	if err = yaml.Unmarshal(yamlS, &poseidon); err != nil {
		return poseidon, errors.New("can not parse " + PoseidonConfigYaml + " config")
	}
	return poseidon, nil
}

/*
 *更新yaml文件
 */
func updatePoseidonYamlFile(poseidon Poseidon, pdNew Poseidon) error {
	isChange := false
	if poseidon.HostIp != pdNew.HostIp {
		// 更新
		poseidon.HostIp = pdNew.HostIp
		isChange = true
	}

	if poseidon.ChipPointer != pdNew.ChipPointer {
		poseidon.ChipPointer = pdNew.ChipPointer
		isChange = true
	}

	if isChange {
		d, err := yaml.Marshal(&poseidon)
		if err != nil {
			return err
		}
		return ioutil.WriteFile(PoseidonConfigYaml, d, os.ModePerm)
	}
	return nil
}

/*
   定时刷新poseidon.yaml文件
*/
func watchPoseidonYaml() {
	c := cron.New()
	_ = c.AddFunc("@every 5s", func() {
		initDefaultPoseidonFile()
	})
	c.Start()
	select {}
}

/**
  查询节点上已经部署的ai container name
*/
func queryAiContainerImageInfo() (map[string]AiResponse, error) {

	ipList, err := Parse1808IpByShellResult()
	if err != nil {
		return nil, err
	}
	res := make(map[string]AiResponse)
	aiContainers := make([]AiContainer, 0)

	for _, nodeIp := range ipList {
		// 请求列表
		// todo 用协程方式并发调用api
		info, err, code := queryNodeDockerListInfo(nodeIp)
		if err != nil || code != http.StatusOK || info == nil {
			continue
		}
		resBytes, err := json.Marshal(info)
		if err != nil {
			klog.Errorf("Marshal node response error: %+v", err)
			continue
		}
		var ac []AiContainer
		err = json.Unmarshal(resBytes, &ac)
		if err != nil {
			klog.Errorf("Unmarshal node response error: %+v", err)
			continue
		}
		for _, container := range ac {
			for _, port := range container.Port {
				container.NodeUrl = nodeIp + ":" + port
				// 只添加running状态的ai-server
				if container.State == "running" {
					aiContainers = append(aiContainers, container)
				}
			}
		}
	}
	for _, container := range aiContainers {
		aiRsp := res[container.ImageHash]
		if reflect.DeepEqual(aiRsp, AiResponse{}) {
			rsp := &AiResponse{
				Address:      []string{container.NodeUrl},
				ExtraCommand: make(map[string]interface{}),
			}
			res[container.ImageHash] = *rsp
		} else {
			aiRsp.Address = append(aiRsp.Address, container.NodeUrl)
			res[container.ImageHash] = aiRsp
		}
	}
	return res, nil
}
