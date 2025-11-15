package conf

type Jwt struct {
	Expire int    `yaml:"expire"` // 过期时间（单位：小时或秒，看你的配置）
	Secret string `yaml:"secret"` // 签名密钥
	Issuer string `yaml:"issuer"` // 签发者
}
