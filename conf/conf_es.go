package conf

type ES struct {
	Url      string `yaml:"url"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Addr     string `yaml:"addr"`
	IsHttps  bool   `yaml:"is_https"`
}
