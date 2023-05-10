// 游标分页组件
package cursor

import (
	"github.com/perpower/goframe/funcs/convert"
	"github.com/perpower/goframe/funcs/ptime"
	"github.com/perpower/goframe/utils/pagination/cursor/base"
)

type RecordsInfo struct {
	Records PageInfo `json:"records" xml:"records"`
}

// 分页数据结果
type PageInfo struct {
	Total      int64                    `json:"total" xml:"total"`           // 数据总数
	List       []map[string]interface{} `json:"list" xml:"list"`             // 分页数据
	NextPage   int                      `json:"nextPage" xml:"nextPage"`     // 是否还有下一页
	NextCursor string                   `json:"nextCursor" xml:"nextCursor"` // 下一页请求游标
}

// 游标信息
type Page struct {
	CursorValue string // 分页游标原始值
	NextTimeAt  int64  // 记录分页发生的时间点, 毫秒时间戳
	PageSize    int    // 分页拉取数量
}

// New 解密游标
// cursor: string 加密后游标值
func New(cursor string, pageSize int) (*Page, error) {
	if cursor == "" {
		return &Page{
			CursorValue: "",
			NextTimeAt:  ptime.TimestampMilli(),
			PageSize:    pageSize,
		}, nil
	}
	decodeBytes, err := base.Base64().Decode(cursor)
	if err != nil {
		return &Page{}, err
	}

	mashRes, err := base.MsgPack().Unmarshal(decodeBytes)
	if err != nil {
		return &Page{}, err
	}

	return &Page{
		CursorValue: convert.String(mashRes["CursorValue"]),
		NextTimeAt:  convert.Int64(mashRes["NextTimeAt"]),
		PageSize:    pageSize,
	}, nil
}

// Generate 输出格式化的分页信息
// total: int64 记录总数
// cursorValue: string 原始游标值
// nextPage: int 下一页状态
// dataList: []map[string]interface{}  查询结果集
func (p *Page) Generate(total int64, cursorValue string, nextPage int, dataList []map[string]interface{}) (RecordsInfo, error) {
	var nextCursor string
	if p.CursorValue == "" && len(dataList) == 0 {
		nextCursor = ""
		dataList = []map[string]interface{}{} // 将结果置为空切片，以达到返回结果为“[]”的目的
	} else {
		if len(dataList) == 0 { // 如果查询结果集是空，则继续使用上一次游标值
			dataList = []map[string]interface{}{} // 将结果置为空切片，以达到返回结果为“[]”的目的
		}
		mashBytes, err := base.MsgPack().Marshal(Page{
			CursorValue: cursorValue,
			NextTimeAt:  ptime.TimestampMilli(),
			PageSize:    p.PageSize,
		})
		if err != nil {
			return RecordsInfo{}, err
		}

		nextCursor, err = base.Base64().Encode(mashBytes)
		if err != nil {
			return RecordsInfo{}, err
		}

	}
	pageinfo := PageInfo{
		Total:      total,
		List:       dataList,
		NextPage:   nextPage,
		NextCursor: nextCursor,
	}
	return RecordsInfo{
		Records: pageinfo,
	}, nil
}
