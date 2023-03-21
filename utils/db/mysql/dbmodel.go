// 数据库操作方法统一封装
// Author: sywen
package mysql

import (
	"log"
	"reflect"

	"github.com/perpower/goframe/funcs/ptime"

	"gorm.io/gorm"
)

type Db struct {
	Prefix string // 数据库表前缀
	Conn   *gorm.DB
	Error  error
}

// func(query, args...) 类方法传参结构体
type QueryArgs struct {
	Query interface{}   // 参数化query语句
	Args  []interface{} // 参数切片
}

// 定义查询条件入参结构体
type FilterParams struct {
	Table    string        // 主表名
	Fields   []string      // 指定字段
	Where    QueryArgs     // 指定where条件
	Or       QueryArgs     // or条件
	Not      QueryArgs     // not条件
	Join     []QueryArgs   // 链表查询, 传值格式[]QueryArgs{}
	Order    []interface{} // 字段排序
	Group    string        // 指定分组查询
	Having   QueryArgs     // 指定having()条件
	Limit    [2]int        // 限制[limit,offset]
	Distinct []interface{} // 指定选择 distinct values
}

var (
	defaultBatchSize = 100 // CreateInBatch 默认每批创建的数据量
)

// Create 创建单条数据
// datas: interface{} 待插入的数据，传值需加地址符&
//
//	集合：map[string]interface{} , 根据 map 创建记录时，association 不会被调用，且主键也不会自动填充
//	结构体：struct{} , 此方式会触发grom的自动补充值机制
//
// return:
//
//	num: int 操作影响的行数
func (db *Db) Create(datas interface{}) (num int64, err error) {
	result := db.Conn.Create(datas)

	num = result.RowsAffected
	err = result.Error

	if err != nil {
		log.Println(result.Error)
	}
	return num, err
}

// CreateBatch 批量插入数据
// datas: interface{} 待插入的数据切片，传值需加地址符&
//
//	集合：map[string]interface{} , 根据 map 创建记录时，association 不会被调用，且主键也不会自动填充
//	结构体：struct{} , 此方式会触发grom的自动补充值机制
//
// return: int64  成功插入的数据条数
func (db *Db) CreateBatch(datas interface{}) (num int64, err error) {
	result := db.Conn.Create(datas)

	num = result.RowsAffected
	err = result.Error

	if err != nil {
		log.Println(result.Error)
	}
	return num, err
}

// CreateInBatch 批量分批插入数据
// datas: interface{} 待插入的数据切片，此处传值不需加地址符&
//
//	集合：map[string]interface{} , 根据 map 创建记录时，association 不会被调用，且主键也不会自动填充
//	结构体：struct{} , 此方式会触发grom的自动补充值机制
//
// batchSize: int 指定每批的数量
// return: int64  成功插入的数据条数
func (db *Db) CreateInBatch(datas interface{}, batchSize int) (num int64, err error) {
	if batchSize <= 0 {
		batchSize = defaultBatchSize
	}
	result := db.Conn.CreateInBatches(datas, batchSize)

	num = result.RowsAffected
	err = result.Error

	if err != nil {
		log.Println(result.Error)
	}
	return num, err
}

// 查询单条数据
// params: FilterParams 查询条件
// res: map[string]interface{} 查询结果
func (db *Db) GetOne(params FilterParams) (res map[string]interface{}, err error) {
	conn, _ := db.FilterWhere(params)
	conn.Limit(1)
	result := conn.Find(&res)

	err = result.Error
	return res, err
}

// GetList 查询数据列表
// params: FilterParams 查询条件
// return:
//
//	res: []map[string]interface{} 查询结果
//	count: int64 匹配查询条件的总记录数
func (db *Db) GetList(params FilterParams) (res []map[string]interface{}, count int64, err error) {
	conn, filters := db.FilterWhere(params)

	//统计满足匹配条件的数据总数
	conn.Count(&count)

	if _, ok := filters.FieldByName("Limit"); ok {
		conn.Limit(params.Limit[0])
		conn.Offset(params.Limit[1])
	}

	result := conn.Find(&res)

	err = result.Error

	return res, count, err
}

// GetListCursor 游标法查询数据列表
// params: FilterParams 查询条件
// limitSize: int 查询数据条数
// cursorWhere: []QueryArgs  游标过滤条件
// return:
//
//	res: []map[string]interface{} 查询结果
//	count: int64 匹配查询条件的总记录数
func (db *Db) GetListCursor(params FilterParams, limitSize int, cursorWhere ...QueryArgs) (res []map[string]interface{}, count int64, err error) {
	conn, _ := db.FilterWhere(params)

	//统计满足匹配条件的数据总数
	conn.Count(&count)

	// 加入游标过滤条件
	if len(cursorWhere) > 0 {
		for _, val := range cursorWhere {
			conn.Where(val.Query, val.Args...)
		}
	}
	conn.Limit(limitSize)

	result := conn.Find(&res)

	err = result.Error

	return res, count, err
}

// FilterWhere 统一处理查询条件
// params: FilterParams  查询条件
func (db *Db) FilterWhere(params FilterParams) (*gorm.DB, reflect.Type) {
	conn := db.Conn.Table(params.Table)
	filters := reflect.TypeOf(params)
	if _, ok := filters.FieldByName("Fields"); ok {
		if len(params.Fields) == 0 {
			// 未指定查询字段，默认查询所有
			params.Fields = []string{"*"}
		}
		conn.Select(params.Fields)
	}

	if _, ok := filters.FieldByName("Join"); ok {
		if len(params.Join) > 0 {
			for _, value := range params.Join {
				conn.Joins(value.Query.(string), value.Args...)
			}
		}
	}

	//查询条件如果未传值或传空值，默认以delStatus=1作为过滤条件
	if _, ok := filters.FieldByName("Where"); ok {
		if !reflect.DeepEqual(params.Where, QueryArgs{}) { // 判断是否为空结构体
			conn.Where(params.Where.Query, params.Where.Args...)
		} else {
			conn.Where(map[string]interface{}{
				"delStatus": 1,
			})
		}
	} else {
		conn.Where(map[string]interface{}{
			"delStatus": 1,
		})
	}

	if _, ok := filters.FieldByName("Or"); ok {
		if !reflect.DeepEqual(params.Or, QueryArgs{}) { // 判断是否为空结构体
			conn.Or(params.Or)
		}
	}

	if _, ok := filters.FieldByName("Not"); ok {
		if !reflect.DeepEqual(params.Not, QueryArgs{}) { // 判断是否为空结构体
			conn.Not(params.Not)
		}
	}

	//排序如果未传值或者传的是一个空值，默认按createAt降序
	if _, ok := filters.FieldByName("Order"); ok {
		if len(params.Order) > 0 {
			for _, value := range params.Order {
				conn.Order(value)
			}
		} else {
			conn.Order("createdAt desc")
		}
	} else {
		conn.Order("createdAt desc")
	}

	if _, ok := filters.FieldByName("Group"); ok {
		if params.Group != "" {
			conn.Group(params.Group)
		}
	}

	if _, ok := filters.FieldByName("Having"); ok {
		if !reflect.DeepEqual(params.Having, QueryArgs{}) { // 判断是否为空结构体
			conn.Having(params.Having)
		}
	}

	if _, ok := filters.FieldByName("Distinct"); ok {
		if len(params.Distinct) > 0 { // 判断是否为空数组
			conn.Distinct(params.Distinct...)
		}
	}

	return conn, filters
}

// Update 更新数据--gorm支持Update更新单个字段，Updates更新多个字段，统一使用Updates()方法
// params: FilterParams{Table: "", Where: map[string]interface{}}
// datas: 有两种传值方式：指定更新字段，传值需加地址符&
//
//	集合：map[string]interface{} , 根据 map 创建记录时，association 不会被调用，且主键也不会自动填充
//	结构体：struct{} , 此方式会触发grom的自动补充值机制
//
// return: int64  返回操作影响的行数
func (db *Db) Update(params FilterParams, datas interface{}) (num int64, err error) {
	conn := db.Conn.Table(params.Table)
	filters := reflect.TypeOf(params)
	if _, ok := filters.FieldByName("Where"); ok {
		if !reflect.DeepEqual(params.Where, QueryArgs{}) { //不允许传空值条件
			result := conn.Where(params.Where.Query, params.Where.Args...).Updates(datas)
			num = result.RowsAffected
			err = result.Error

			if err != nil {
				log.Println(result.Error)
			}

			return num, err
		}
	}

	return 0, nil
}

// Delete 支持真删除和软删除
// params: FilterParams  匹配条件
// model: interface{} 需要操作的数据表模型，传值需加地址符&
// scoped: bool 真删除=true 软删除=false
func (db *Db) Delete(model interface{}, params FilterParams, scoped bool) (num int64, err error) {
	filters := reflect.TypeOf(params)
	if _, ok := filters.FieldByName("Where"); ok {
		if !reflect.DeepEqual(params.Where, QueryArgs{}) { //不允许传空值条件
			var result *gorm.DB
			if scoped { //真删除
				db.Conn = db.Conn.Unscoped()
				result = db.Conn.Where(params.Where.Query, params.Where.Args...).Delete(model)
			} else {
				result = db.Conn.Table(params.Table).Where(params.Where.Query, params.Where.Args...).Updates(map[string]interface{}{
					"deletedAt": ptime.TimestampMilli(),
					"delStatus": 2,
				})
			}

			num = result.RowsAffected
			err = result.Error

			if err != nil {
				log.Println(result.Error)
			}

			return num, err
		}
	}

	return 0, nil
}

// Raw 执行原生sql查询
// res: interface{} 存储查询结果，传值需加地址符&
// sql: string 原生sql语句
// values: ...interface{}  多参数值映射
func (db *Db) Raw(res interface{}, sql string, values ...interface{}) (interface{}, error) {
	if sql != "" {
		result := db.Conn.Raw(sql, values...).Scan(res)
		err := result.Error
		return res, err
	}
	return res, nil
}

// Exec 执行原生sql(非查询操作)
// sql: string 原生sql语句
// values: ...interface{}  多参数值映射
// num: 影响的行数
// err: 错误信息
func (db *Db) Exec(sql string, values ...interface{}) (num int64, err error) {
	if sql != "" {
		result := db.Conn.Exec(sql, values...)

		num = result.RowsAffected
		err = result.Error

		if err != nil {
			log.Println(result.Error)
		}

		return num, err
	}

	return 0, nil
}
