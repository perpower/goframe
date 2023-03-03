// 地理位置(本质是有序集合)
package redis

import (
	"github.com/perpower/goframe/funcs/convert"
	"github.com/perpower/goframe/funcs/normal"

	"github.com/gomodule/redigo/redis"
)

type Rgeo struct {
	conn redis.Conn
}

// GEOADD, 将指定的地理空间项（经度、纬度、名称）添加到指定的键
// since: 6.2.0
// key: string 键名
// elementCondition: string , 取值 XX | NX, 不指定传空
//
//	XX：只更新已经存在的元素。 不要添加新元素。
//	NX：只添加新元素。 不要更新已经存在的元素。
//
// isCh: bool, 指定该参数时，会将返回值修改，由原先的返回新添加的元素个数，变为返回变更的元素+新添加的元素总数
// points: [][3]string{longitude, latitude, name}数组
// return: reply int
// link: https://redis.io/commands/geoadd/
func (c *Rgeo) GeoAdd(key string, points [][3]string, elementCondition string, isCh bool) (int, error) {
	if len(points) == 0 {
		return 0, nil
	}

	args := make([]interface{}, 0)
	args = append(args, key)
	if normal.InArray(elementCondition, []string{"XX", "NX"}) {
		args = append(args, elementCondition)
	}
	if isCh {
		args = append(args, "CH")
	}

	for _, pair := range points {
		args = append(args, pair[0], pair[1], pair[2])
	}

	return redis.Int(c.conn.Do("GEOADD", args...))
}

// GEODIST, 返回的地理空间索引中两个成员之间的距离。
// key: string 键名
// members: [2]string{member1,member2}
// unit: string 单位，默认为M
//
//	M for 米.
//	KM for 公里.
//	MI for 英里.
//	FT for 英尺.
//
// return:
//
//	reply: string
//
// link: https://redis.io/commands/geodist/
func (c *Rgeo) GeoDist(key string, members [2]string, unit string) (string, error) {
	if len(members) != 2 {
		panic("参数members必须包含两个值")
	}
	if !normal.InArray(unit, []string{"M", "KM", "MI", "FT"}) {
		unit = "M"
	}

	return redis.String(c.conn.Do("GEODIST", key, members[0], members[1], unit))
}

// GEOHASH, 返回给定的N个元素在集合中的有效的 Geohash 字符串
// key: string 键名
// members: []string 指定元素数组
// return:
//
//	reply: map[string]string, 格式 map[member]hash
//
// link: https://redis.io/commands/geohash/
func (c *Rgeo) GeoHash(key string, members []string) (map[string]string, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	for _, v := range members {
		args = append(args, v)
	}
	res, err := redis.Strings(c.conn.Do("GEOHASH", args...))
	maps := make(map[string]string)
	if err != nil {
		return maps, err
	}
	for k, v := range members {
		maps[v] = res[k]
	}

	return maps, err
}

// GEOPOS, 返回给定的N个成员的经纬度位置信息
// key: string 键名
// members: []string 指定成员数组
// return:
//
//	reply: []map[string]string, 格式 []map[string]string{"member":"", "longitude":"","latitude"}, ...}
//
// link: https://redis.io/commands/geopos/
func (c *Rgeo) GeoPos(key string, members []string) (arr []map[string]string, err error) {
	if len(members) == 0 {
		return []map[string]string{}, nil
	}
	args := make([]interface{}, 0)
	args = append(args, key)
	for _, v := range members {
		args = append(args, v)
	}
	res, err := redis.Values(c.conn.Do("GEOPOS", args...))
	for k, v := range members {
		poi, _ := redis.Strings(res[k], nil)
		if len(poi) > 0 {
			arr = append(arr, map[string]string{
				"member":    v,
				"longitude": poi[0],
				"latitude":  poi[1],
			})
		} else { // 当member不存在时
			arr = append(arr, nil)
		}
	}

	return arr, err
}

// GEOSEARCH, 此方法通过FROMMEMBER方式搜索,同时不指定任何WITH选项，即使用给定的集合内的某个成员
// since: 7.0.0
// key: string 键名
// member: string 指定集合内的成员作为中心点
// byType: []string 指定搜索范围,两种方式取其一
//
//	BYRADIUS: []string, 根据给定的<radius>在圆形区域内搜索, 传值格式：[]string{"BYRADIUS", radius, unit}
//	BYBOX: []string, 在轴对齐的矩形内搜索, 传值格式：[]string{"BYBOX", width, height, unit}
//
// sort: string, 排序方式，取值：ASC | DESC
//
//	ASC：相对于中心点，从最近到最远进行排序。
//	DESC：相对于中心点，从最远到最近进行排序。
//
// count: [2]int, 默认返回所有匹配项。 使用COUNT参数将结果限制为前 N 个匹配项
//
//	不使用ANY选项时：传值格式[2]int{count,0}
//	使用ANY选项时：传值格式[2]int{count,1}
//
// 特别说明：使用 ANY 选项时，一旦找到足够的匹配项，命令就会返回。 这意味着返回的结果可能不是最接近指定点的结果，但服务器为生成它们所投入的努力要少得多。
// return:
//
//	reply: []string
//
// link: https://redis.io/commands/geosearch/
func (c *Rgeo) GeoSearchFromMemberNoWith(key, member string, byType []string, sort string, count [2]int) ([]string, error) {
	args := make([]interface{}, 0)
	args = append(args, key, "FROMMEMBER", member)
	for _, v := range byType {
		args = append(args, v)
	}

	if !normal.InArray(sort, []string{"ASC", "DESC"}) {
		sort = "ASC"
	}
	args = append(args, sort)

	if len(count) == 2 {
		if count[1] == 1 { // 使用ANY选项
			args = append(args, "COUNT", count[0], "ANY")
		} else {
			args = append(args, "COUNT", count[0])
		}
	}

	return redis.Strings(c.conn.Do("GEOSEARCH", args...))
}

// GEOSEARCH, 此方法通过FROMMEMBER方式搜索,同时指定所有WITH选项，即使用给定的集合内的某个成员
// since: 7.0.0
// key: string 键名
// member: string 指定集合内的成员作为中心点
// byType: []string 指定搜索范围,两种方式取其一
//
//	BYRADIUS: []string, 根据给定的<radius>在圆形区域内搜索, 传值格式：[]string{"BYRADIUS", radius, unit}
//	BYBOX: []string, 在轴对齐的矩形内搜索, 传值格式：[]string{"BYBOX", width, height, unit}
//
// sort: string, 排序方式，取值：ASC | DESC
//
//	ASC：相对于中心点，从最近到最远进行排序。
//	DESC：相对于中心点，从最远到最近进行排序。
//
// count: [2]int, 默认返回所有匹配项。 使用COUNT参数将结果限制为前 N 个匹配项
//
//	不使用ANY选项时：传值格式[2]int{count,0}
//	使用ANY选项时：传值格式[2]int{count,1}
//
// 特别说明：使用 ANY 选项时，一旦找到足够的匹配项，命令就会返回。 这意味着返回的结果可能不是最接近指定点的结果，但服务器为生成它们所投入的努力要少得多。
// WITH选项说明：
// WITHCOORD: 同时返回匹配项的经度和纬度。
// WITHDIST: 同样返回成员距指定中心点的距离。 返回的距离与为半径或高度和宽度参数指定的单位相同
// WITHHASH: 以 52 位无符号整数的形式返回成员的原始 geohash 编码排序集score
// return:
//
//	reply: []map[string]string, map格式{"member":"", "dist":"","hash":"","longitude":"","latitude":""}
//
// link: https://redis.io/commands/geosearch/
func (c *Rgeo) GeoSearchFromMemberHasWith(key, member string, byType []string, sort string, count [2]int) ([]map[string]string, error) {
	args := make([]interface{}, 0)
	args = append(args, key, "FROMMEMBER", member)
	for _, v := range byType {
		args = append(args, v)
	}

	if !normal.InArray(sort, []string{"ASC", "DESC"}) {
		sort = "ASC"
	}
	args = append(args, sort)

	if len(count) == 2 {
		if count[1] == 1 { // 使用ANY选项
			args = append(args, "COUNT", count[0], "ANY")
		} else {
			args = append(args, "COUNT", count[0])
		}
	}

	args = append(args, "WITHCOORD", "WITHDIST", "WITHHASH")

	res, err := redis.Values(c.conn.Do("GEOSEARCH", args...))

	arr := make([]map[string]string, 0)
	if (err != nil) || (len(res) == 0) {
		return arr, err
	}

	// 将结果转换成指定格式的map
	for _, v := range res {
		lists, _ := redis.Values(v, nil)
		member, _ := redis.String(lists[0], nil)
		dist, _ := redis.String(lists[1], nil)
		hash, _ := redis.Int64(lists[2], nil)
		poi, _ := redis.Strings(lists[3], nil)

		arr = append(arr, map[string]string{
			"member":    member,
			"dist":      dist,
			"hash":      convert.String(hash),
			"longitude": poi[0],
			"latitude":  poi[1],
		})
	}

	return arr, err
}

// GEOSEARCH, 此方法通过FROMLONLAT方式搜索,同时不指定任何WITH选项，即使用给定的经纬度
// since: 7.0.0
// key: string 键名
// poi: [2]string{longitude, latitude} 指定经纬度作为中心点
// byType: []string 指定搜索范围,两种方式取其一
//
//	BYRADIUS: []string, 根据给定的<radius>在圆形区域内搜索, 传值格式：[]string{"BYRADIUS", radius, unit}
//	BYBOX: []string, 在轴对齐的矩形内搜索, 传值格式：[]string{"BYBOX", width, height, unit}
//
// sort: string, 排序方式，取值：ASC | DESC
//
//	ASC：相对于中心点，从最近到最远进行排序。
//	DESC：相对于中心点，从最远到最近进行排序。
//
// count: [2]int, 默认返回所有匹配项。 使用COUNT参数将结果限制为前 N 个匹配项
//
//	不使用ANY选项时：传值格式[2]int{count,0}
//	使用ANY选项时：传值格式[2]int{count,1}
//
// 特别说明：使用 ANY 选项时，一旦找到足够的匹配项，命令就会返回。 这意味着返回的结果可能不是最接近指定点的结果，但服务器为生成它们所投入的努力要少得多。
// return:
//
//	reply: []string
//
// link: https://redis.io/commands/geosearch/
func (c *Rgeo) GeoSearchFromLonlatNoWith(key string, poi [2]string, byType []string, sort string, count [2]int) ([]string, error) {
	if len(poi) != 2 {
		return []string{}, nil
	}
	args := make([]interface{}, 0)
	args = append(args, key, "FROMLONLAT", poi[0], poi[1])
	for _, v := range byType {
		args = append(args, v)
	}

	if !normal.InArray(sort, []string{"ASC", "DESC"}) {
		sort = "ASC"
	}
	args = append(args, sort)

	if len(count) == 2 {
		if count[1] == 1 { // 使用ANY选项
			args = append(args, "COUNT", count[0], "ANY")
		} else {
			args = append(args, "COUNT", count[0])
		}
	}

	return redis.Strings(c.conn.Do("GEOSEARCH", args...))
}

// GEOSEARCH, 此方法通过FROMLONLAT方式搜索,同时指定所有WITH选项，即使用给定的经纬度
// since: 7.0.0
// key: string 键名
// poi: [2]string{longitude, latitude} 指定经纬度作为中心点
// byType: []string 指定搜索范围,两种方式取其一
//
//	BYRADIUS: []string, 根据给定的<radius>在圆形区域内搜索, 传值格式：[]string{"BYRADIUS", radius, unit}
//	BYBOX: []string, 在轴对齐的矩形内搜索, 传值格式：[]string{"BYBOX", width, height, unit}
//
// sort: string, 排序方式，取值：ASC | DESC
//
//	ASC：相对于中心点，从最近到最远进行排序。
//	DESC：相对于中心点，从最远到最近进行排序。
//
// count: [2]int, 默认返回所有匹配项。 使用COUNT参数将结果限制为前 N 个匹配项
//
//	不使用ANY选项时：传值格式[2]int{count,0}
//	使用ANY选项时：传值格式[2]int{count,1}
//
// 特别说明：使用 ANY 选项时，一旦找到足够的匹配项，命令就会返回。 这意味着返回的结果可能不是最接近指定点的结果，但服务器为生成它们所投入的努力要少得多。
// WITH选项说明：
// WITHCOORD: 同时返回匹配项的经度和纬度。
// WITHDIST: 同样返回成员距指定中心点的距离。 返回的距离与为半径或高度和宽度参数指定的单位相同
// WITHHASH: 以 52 位无符号整数的形式返回成员的原始 geohash 编码排序集score
// return:
//
//	reply: []map[string]string, map格式{"member":"", "dist":"","hash":"","longitude":"","latitude":""}
//
// link: https://redis.io/commands/geosearch/
func (c *Rgeo) GeoSearchFromLonlatHasWith(key string, poi [2]string, byType []string, sort string, count [2]int) ([]map[string]string, error) {
	if len(poi) != 2 {
		return []map[string]string{}, nil
	}
	args := make([]interface{}, 0)
	args = append(args, key, "FROMLONLAT", poi[0], poi[1])
	for _, v := range byType {
		args = append(args, v)
	}

	if !normal.InArray(sort, []string{"ASC", "DESC"}) {
		sort = "ASC"
	}
	args = append(args, sort)

	if len(count) == 2 {
		if count[1] == 1 { // 使用ANY选项
			args = append(args, "COUNT", count[0], "ANY")
		} else {
			args = append(args, "COUNT", count[0])
		}
	}

	args = append(args, "WITHCOORD", "WITHDIST", "WITHHASH")

	res, err := redis.Values(c.conn.Do("GEOSEARCH", args...))

	arr := make([]map[string]string, 0)
	if (err != nil) || (len(res) == 0) {
		return arr, err
	}

	// 将结果转换成指定格式的map
	for _, v := range res {
		lists, _ := redis.Values(v, nil)
		member, _ := redis.String(lists[0], nil)
		dist, _ := redis.String(lists[1], nil)
		hash, _ := redis.Int64(lists[2], nil)
		poi, _ := redis.Strings(lists[3], nil)

		arr = append(arr, map[string]string{
			"member":    member,
			"dist":      dist,
			"hash":      convert.String(hash),
			"longitude": poi[0],
			"latitude":  poi[1],
		})
	}

	return arr, err
}

// GEOSEARCHSTORE, 此方法通过FROMMEMBER方式搜索,即使用给定的集合内的某个成员, 将结果存储到指定key中
// since: 7.0.0
// destination: string 目标键名
// source: string 源键名
// member: string 指定集合内的成员作为中心点
// byType: []string 指定搜索范围,两种方式取其一
//
//	BYRADIUS: []string, 根据给定的<radius>在圆形区域内搜索, 传值格式：[]string{"BYRADIUS", radius, unit}
//	BYBOX: []string, 在轴对齐的矩形内搜索, 传值格式：[]string{"BYBOX", width, height, unit}
//
// sort: string, 排序方式，取值：ASC | DESC
//
//	ASC：相对于中心点，从最近到最远进行排序。
//	DESC：相对于中心点，从最远到最近进行排序。
//
// count: [2]int, 默认返回所有匹配项。 使用COUNT参数将结果限制为前 N 个匹配项
//
//	不使用ANY选项时：传值格式[2]int{count,0}
//	使用ANY选项时：传值格式[2]int{count,1}
//
// 特别说明：使用 ANY 选项时，一旦找到足够的匹配项，命令就会返回。 这意味着返回的结果可能不是最接近指定点的结果，但服务器为生成它们所投入的努力要少得多。
// storeDist: bool 是否同时存储与指定中心点的距离，默认false不存储
// return:
//
//	reply: int 存储到结果集中的成员数量
//
// link: https://redis.io/commands/geosearchstore/
func (c *Rgeo) GeoSearchFromMemberStore(destination, source, member string, byType []string, sort string, count [2]int, storeDist bool) (int, error) {
	args := make([]interface{}, 0)
	args = append(args, destination, source, "FROMMEMBER", member)
	for _, v := range byType {
		args = append(args, v)
	}

	if !normal.InArray(sort, []string{"ASC", "DESC"}) {
		sort = "ASC"
	}
	args = append(args, sort)

	if len(count) == 2 {
		if count[1] == 1 { // 使用ANY选项
			args = append(args, "COUNT", count[0], "ANY")
		} else {
			args = append(args, "COUNT", count[0])
		}
	}

	if storeDist {
		args = append(args, "STOREDIST")
	}
	return redis.Int(c.conn.Do("GEOSEARCHSTORE", args...))
}

// GEOSEARCHSTORE, 此方法通过FROMLONLAT方式搜索,即使用给定的经纬度, 将结果存储到指定key中
// since: 7.0.0
// destination: string 目标键名
// source: string 源键名
// poi: [2]string{longitude, latitude} 指定经纬度作为中心点
// byType: []string 指定搜索范围,两种方式取其一
//
//	BYRADIUS: []string, 根据给定的<radius>在圆形区域内搜索, 传值格式：[]string{"BYRADIUS", radius, unit}
//	BYBOX: []string, 在轴对齐的矩形内搜索, 传值格式：[]string{"BYBOX", width, height, unit}
//
// sort: string, 排序方式，取值：ASC | DESC
//
//	ASC：相对于中心点，从最近到最远进行排序。
//	DESC：相对于中心点，从最远到最近进行排序。
//
// count: [2]int, 默认返回所有匹配项。 使用COUNT参数将结果限制为前 N 个匹配项
//
//	不使用ANY选项时：传值格式[2]int{count,0}
//	使用ANY选项时：传值格式[2]int{count,1}
//
// 特别说明：使用 ANY 选项时，一旦找到足够的匹配项，命令就会返回。 这意味着返回的结果可能不是最接近指定点的结果，但服务器为生成它们所投入的努力要少得多。
// storeDist: bool 是否同时存储与指定中心点的距离，默认false不存储
// return:
//
//	reply: int 存储到结果集中的成员数量
//
// link: https://redis.io/commands/geosearchstore/
func (c *Rgeo) GeoSearchFromLonlatStore(destination, source string, poi [2]string, byType []string, sort string, count [2]int, storeDist bool) (int, error) {
	if len(poi) != 2 {
		return 0, nil
	}
	args := make([]interface{}, 0)
	args = append(args, destination, source, "FROMLONLAT", poi[0], poi[1])
	for _, v := range byType {
		args = append(args, v)
	}

	if !normal.InArray(sort, []string{"ASC", "DESC"}) {
		sort = "ASC"
	}
	args = append(args, sort)

	if len(count) == 2 {
		if count[1] == 1 { // 使用ANY选项
			args = append(args, "COUNT", count[0], "ANY")
		} else {
			args = append(args, "COUNT", count[0])
		}
	}

	if storeDist {
		args = append(args, "STOREDIST")
	}
	return redis.Int(c.conn.Do("GEOSEARCHSTORE", args...))
}
