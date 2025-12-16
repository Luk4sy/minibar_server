package conf

import "fmt"

type Config struct {
	System System `yaml:"system"`
	Log    Log    `yaml:"log"`
	DB     []DB   `yaml:"db"` // 数据库连接列表
	Jwt    Jwt    `yaml:"jwt"`
	Redis  Redis  `yaml:"redis"`
	Site   Site   `yaml:"site"`
	Email  Email  `yaml:"email"`
	QQ     QQ     `yaml:"qq"`
	QiNiu  QiNiu  `yaml:"qiNiu"`
	Ai     Ai     `yaml:"ai"`
	Upload Upload `yaml:"upload"`
	ES     ES     `yaml:"es"`
	River  River  `yaml:"river"`
}

func (e ES) EsUrl() string {
	if e.IsHttps {
		return fmt.Sprintf("https://%s", e.Addr)
	}
	return fmt.Sprintf("http://%s", e.Addr)
}
