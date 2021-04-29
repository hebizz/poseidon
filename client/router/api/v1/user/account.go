package user

type Acc struct {
	UserName string `yaml:"userName" json:"userName"`
	Password string `yaml:"password" json:"password"`
	Role     string `yaml:"role" json:"role"`
}

type Account struct {
	Acc []Acc `yaml:"account" json:"accounts"`
}

const (
	NormalName     = "jiangxing"
	NormalRole     = "normal"
	NormalPassword = "123456"
	AdminName      = "admin"
	AdminRole      = "root"
	AdminPassword  = "admin123"
)
