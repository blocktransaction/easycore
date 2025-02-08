package oss

type Option struct {
	//oss访问地址
	Endpoint string
	//oss访问id
	AccessKeyId string
	//oss访问秘钥
	AccessKeySecret string
	//oss存储桶
	BucketName string
	//基础路径
	BasePath string
	//最大文件大小
	MaxFileSize int
	//前缀
	prefix string
}

// 初始化
func (o *Option) init() {
	if o.MaxFileSize == 0 {
		o.MaxFileSize = 16 * 1024 * 1024 //默认16M
	}
	o.prefix = "https://"
}
