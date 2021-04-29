package container

type Docker struct {
	Name          string `json:"name" binding:"required"`
	Version       string `json:"version" binding:"required"`
	State         string `json:"status"`
	CreateTime    int64  `json:"createTime"`
	ContainerId   string `json:"containerId"`
	ImageID       string `json:"imageId"`
	ContainerName string `json:"containerName"`
	IsEdit        bool   `json:"isEdit"`
}

type AiContainer struct {
	ImageHash string   `json:"imageHash"`
	State     string   `json:"status"`
	Port      []string `json:"port"`
	NodeUrl   string   `json:"nodeUrl"`
}

type Port struct {
	HostIP   string
	HostPort string
}

type AiResponse struct {
	Address      [] string   `json:"address" yaml:"address"`
	ExtraCommand interface{} `json:"extra_command" yaml:"extra_command"`
}

type UpgradeResult struct {
	Data   interface{}
	Code   int
	NodeIp string
	Err    error
}

type YamlServices struct {
	Services interface{} `yaml:"services"`
}

type Container struct {
	ContainerName string `yaml:"container_name" json:"container_name"`
	Image         string `yaml:"image"`
	HostName      string `yaml:"hostname"`
}

type Poseidon struct {
	HostIp      string      `yaml:"host_ip" json:"host_ip"`
	ChipPointer interface{} `yaml:"chip_pointer" json:"chip_pointer"`
}
