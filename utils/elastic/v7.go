// 针对ElasticSearch V7.X 版本
package elastic

import (
	"context"
	"fmt"

	es7 "github.com/olivere/elastic/v7"
)

type Es7 struct {
	conn *es7.Client
}

// IndexExists 检测指定index 是否存在
// indexName: string 索引名
func (es *Es7) IndexExists(indexName string) (bool, error) {
	exist, err := es.conn.IndexExists(indexName).Do(context.Background())
	fmt.Println("索引检测：", exist, err)
	return exist, err
}

// CreateIndex 创建索引
// indexName: string 索引名
// mappings: interface{} 指定索引映射
func (es *Es7) CreateIndex(indexName string, mappings interface{}) (bool, error) {
	createIndex, err := es.conn.CreateIndex(indexName).BodyJson(mappings).Do(context.Background())
	if err != nil || !createIndex.Acknowledged {
		return false, err
	}
	return true, nil
}

// CreateDoc 创建文档
// indexName: string 索引名
// bodyContent: string
func (es *Es7) CreateDoc(indexName string, bodyContent interface{}) (_id string, err error) {
	res, err := es.conn.Index().Index(indexName).BodyJson(bodyContent).Do(context.Background())
	if err != nil {
		return "", err
	}
	_id = res.Id
	return _id, err
}

// DeleteIndex 删除索引
// indexNames: []string 索引名切片
func (es *Es7) DeleteIndex(indexNames []string) (bool, error) {
	res, err := es.conn.DeleteIndex(indexNames...).Do(context.Background())
	if err != nil || !res.Acknowledged {
		return false, err
	}

	return true, nil
}
