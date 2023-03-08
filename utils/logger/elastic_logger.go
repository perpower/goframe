// elasticsearch 日志存储组件
package logger

import (
	"strconv"
	"time"

	"github.com/perpower/goframe/utils/elastic"
)

var (
	esclient *elastic.Client
)

// InitElastic 链接ElasticSearch服务
func InitElastic(conf elastic.ElastiConfig) {
	es, _ := elastic.Instance(conf)
	esclient = es
}

// CreateElasticLog 创建日志文档
// level: string 错误等级
// IndexName: string 索引名称
// msg: string 消息文本
// filedSlice: []ExtendFields  额外参数
func CreateElasticLog(level, IndexName, msg string, filedSlice ...ExtendFields) (string, error) {
	if IndexName == "" {
		year, month, day := time.Now().Date()
		IndexName = strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day)
	}
	if createIndex(IndexName) {
		params := requestInfo()
		docContent := map[string]interface{}{
			"logLevel":      level,
			"requestTime":   params.RequestTime,
			"requestMethod": params.RequestMethod,
			"requestHost":   params.RequestHost,
			"requestUri":    params.RequestUri,
			"userAgent":     params.UserAgent,
			"clientIp":      params.ClientIp,
			"requestBody":   params.RequestBody,
			"extraDatas":    filedSlice,
			"message":       msg,
		}
		return esclient.V7.CreateDoc(IndexName, docContent)
	}
	return "", nil
}

// createIndex 创建索引
// indexName: string 索引名称
func createIndex(indexName string) bool {
	// 先判断索引是否存在
	status, err := esclient.V7.IndexExists(indexName)
	if err != nil {
		return false
	}
	if !status { // 如果索引不存在，则创建
		mappings := `{
			"mappings": {
				"_doc": {
					"properties": {
						"logLevel": {
							"type": "keyword"
						}
						"requestTime": {
							"type": "date",
							"format": "yyyy-MM-dd HH:mm:ss||yyyy-MM-dd||epoch_millis"
						},
						"requestMethod": {
							"type": "text"
						},
						"requestHost": {
							"type": "keyword"
						},
						"requestUri": {
							"type": "text"
						},
						"userAgent": {
							"type": "text"
						},
						"clientIp": {
							"type": "IP"
						},
						"requestBody": {
							"type": "object"
						},
						"extraDatas": {
							"type": "arrays"
						}
						"message": {
							"type": "text"
						}
					}
				}
			}
		}`
		status, err := esclient.V7.CreateIndex(indexName, mappings)
		if err != nil {
			return false
		}

		return status
	}
	return true
}
