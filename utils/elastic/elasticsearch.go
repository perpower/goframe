package elastic

import (
	es7 "github.com/olivere/elastic/v7"
	"github.com/perpower/goframe/funcs/normal"
)

type ElastiConfig struct {
	Nodes string   // 服务节点，多个节点用“|”链接
	Auth  struct { // 身份鉴权
		Username string // 用户名
		Password string // 密码
	}
	HealthCheck     bool   // 节点健康检查
	DefaultProtocol string // 请求协议
	Version         int    // 主版本号，需与服务端用的主版本号保持一致, 例如 7=Elasticsearch7.xx, 8=Elasticsearch8.xx
}

type Client struct {
	V7 *Es7
}

// Instance 客户端链接实例化
func Instance(conf ElastiConfig) (c *Client, err error) {
	// 将节点字符串切割为数组
	nodes := normal.SplitAndTrim(conf.Nodes, "|")
	var client *es7.Client
	switch conf.Version {
	case 7:
		client, err = es7.NewClient(
			es7.SetURL(nodes...),
			es7.SetBasicAuth(conf.Auth.Username, conf.Auth.Password),
			es7.SetScheme(conf.DefaultProtocol),
			es7.SetHealthcheck(conf.HealthCheck),
			es7.SetSniff(false), // 这里设置成不需要地址自动转换，否则会报错context deadline exceeded
		)
	default:
		panic("ElasticSearch版本配置有误")
	}

	if err != nil {
		return nil, err
	}

	c = &Client{
		V7: &Es7{
			conn: client,
		},
	}

	return c, err
}
