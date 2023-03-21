// offset位移分页组件
package page

type RecordsInfo struct {
	Records PageInfo `json:"records" xml:"records"`
}

// 分页数据结果
type PageInfo struct {
	Total int64                    `json:"total" xml:"total"` // 数据总数
	List  []map[string]interface{} `json:"list" xml:"list"`   // 分页数据
}

// Generate 输出格式化的分页信息
// total: int64 记录总数
// dataList: []map[string]interface{}  查询结果集
func Generate(total int64, dataList []map[string]interface{}) (RecordsInfo, error) {
	pageinfo := PageInfo{
		Total: total,
		List:  dataList,
	}
	return RecordsInfo{
		Records: pageinfo,
	}, nil
}
