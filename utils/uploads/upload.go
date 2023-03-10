// 文件上传工具
package uploads

type Upload struct {
	Local      *localFile
	TencentCos *tencentCos
}

// Instance
// conf: interface{} 上传方式配置
func Instance(conf interface{}) (c Upload, typ interface{}) {
	switch confType := conf.(type) { // 考虑到switch类型断言的问题，将结果分配给一个变量，否则可能会触发panic
	case LocalConfig:
		c.Local = &localFile{}
		typ = confType
	case CosConfig:
		cosConf := conf.(CosConfig)
		c.TencentCos = &tencentCos{
			config: &cosConf,
		}
		typ = confType
	}
	return c, typ
}
