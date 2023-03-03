// 有序集合
package redis

import (
	"github.com/perpower/goframe/funcs/normal"

	"github.com/gomodule/redigo/redis"
)

type Rzset struct {
	conn redis.Conn
}

// ZADD, 将具有指定score的所有指定成员添加到key中，如果key不存在，则创建。此方法不添加"INCR"参数
// 如果指定的成员已经是排序集的成员，则更新score并将元素重新插入正确的位置以确保正确的排序。
// 注意：GT、LT 和 NX 选项相互排斥。
// since: 6.2.0
// key: string 键名
// elementCondition: string , 取值 XX | NX, 不指定传空
//
//	XX：只更新已经存在的元素。 不要添加新元素。
//	NX：只添加新元素。 不要更新已经存在的元素。
//
// scoreCondition: string, 取值 LT | GT, 不指定传空
//
//	LT：如果新分数小于当前分数，则只更新现有元素。 此标志不会阻止添加新元素。
//	GT：如果新分数大于当前分数，则仅更新现有元素。 此标志不会阻止添加新元素。
//
// isCh: bool, 指定该参数时，会将返回值修改，由原先的返回新添加的元素个数，变为返回变更的元素+新添加的元素总数
// isIncr: bool, 指定此选项时，ZADD 的行为类似于 ZINCRBY。 在此模式下只能指定一个score-element对。
// scoreElements: [][2]string{score,element} score-element对数组
// return:
//
//	reply: int 返回成功影响的成员总数
//
// link: https://redis.io/commands/zadd/
func (c *Rzset) Zadd(key string, scoreElements [][2]string, elementCondition, scoreCondition string, isCh bool) (int, error) {
	if len(scoreElements) == 0 {
		return 0, nil
	}

	args := make([]interface{}, 0)
	args = append(args, key)
	if normal.InArray(elementCondition, []string{"XX", "NX"}) {
		args = append(args, elementCondition)
	}
	if normal.InArray(scoreCondition, []string{"LT", "GT"}) {
		args = append(args, scoreCondition)
	}
	if isCh {
		args = append(args, "CH")
	}

	for _, pair := range scoreElements {
		args = append(args, pair[0], pair[1])
	}

	return redis.Int(c.conn.Do("ZADD", args...))
}

// ZADD, 将指定score-element对添加到key中，如果key不存在，则创建。 此方法默认添加 "INCR"参数
// 如果指定的成员已经是排序集的成员，则更新score并将元素重新插入正确的位置以确保正确的排序。
// 注意：GT、LT 和 NX 选项相互排斥。
// since: 6.2.0
// key: string 键名
// score: string score下标值
// element: string 成员值
// elementCondition: string , 取值 XX | NX, 不指定传空
//
//	XX：只更新已经存在的元素。 不要添加新元素。
//	NX：只添加新元素。 不要更新已经存在的元素。
//
// scoreCondition: string, 取值 LT | GT, 不指定传空
//
//	LT：如果新分数小于当前分数，则只更新现有元素。 此标志不会阻止添加新元素。
//	GT：如果新分数大于当前分数，则仅更新现有元素。 此标志不会阻止添加新元素。
//
// isCh: bool, 指定该参数时，会将返回值修改，由原先的返回新添加的元素个数，变为返回变更的元素+新添加的元素总数
// return:
//
//	reply: string 返回新的score值
//
// link: https://redis.io/commands/zadd/
func (c *Rzset) ZaddIncr(key, score, element string, elementCondition, scoreCondition string, isCh bool) (string, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	if normal.InArray(elementCondition, []string{"XX", "NX"}) {
		args = append(args, elementCondition)
	}
	if normal.InArray(scoreCondition, []string{"LT", "GT"}) {
		args = append(args, scoreCondition)
	}
	if isCh {
		args = append(args, "CH")
	}
	args = append(args, "INCR")

	return redis.String(c.conn.Do("ZADD", args...))
}

// ZCARD, 返回有序集合的元素总数, 不存在的key返回0
// key: string 键名
// return:
//
//	reply: int
//
// link: https://redis.io/commands/zcard/
func (c *Rzset) Zcard(key string) (int, error) {
	return redis.Int(c.conn.Do("ZCARD", key))
}

// ZCOUNT, 返回有序集合中指定score区间范围内的元素数量, 不存在的key返回0
// key: string 键名
// min: string score起始值，默认是包含边界值的，可指定为不包括，传值方法例如“(2”，即在score前面加一个“(”符号
//
//	也可以直接传值“-inf”，会自动取集合中最小的score
//
// max: string score截止值，默认是包含边界值的，可指定为不包括，传值方法例如“(2”，即在score前面加一个“(”符号
//
//	也可以直接传值“+inf”，会自动取集合中最大的score
//
// return: reply int
// link: https://redis.io/commands/zcount/
func (c *Rzset) Zcount(key, min, max string) (int, error) {
	return redis.Int(c.conn.Do("ZCOUNT", key, min, max))
}

// ZDIFF, 返回给定的N个有序集合之间的差异，此方法不指定 WITHSCORES 参数
// since: 6.2.0
// key: string 第一个集合键名
// diffKeys: []string 需要比对的键名数组
// return:
//
//	reply []string
//
// link: https://redis.io/commands/zdiff/
func (c *Rzset) Zdiff(key string, diffKeys []string) ([]string, error) {
	if key == "" {
		return []string{}, nil
	}
	args := make([]interface{}, 0)
	args = append(args, len(diffKeys)+1, key)
	for _, v := range diffKeys {
		args = append(args, v)
	}
	return redis.Strings(c.conn.Do("ZDIFF", args...))
}

// ZDIFF, 返回给定的N个有序集合之间的差异，此方法默认指定 WITHSCORES 参数
// since: 6.2.0
// key: string 第一个集合键名
// diffKeys: []string 需要比对的键名数组
// return:
//
//	reply [][2]string{{"score", "element"}, ...}
//
// link: https://redis.io/commands/zdiff/
func (c *Rzset) ZdiffWithScore(key string, diffKeys []string) (arr [][2]string, err error) {
	if key == "" {
		return [][2]string{}, nil
	}
	args := make([]interface{}, 0)
	args = append(args, len(diffKeys)+1, key)
	for _, v := range diffKeys {
		args = append(args, v)
	}
	args = append(args, "WITHSCORES")

	lists, err := redis.Strings(c.conn.Do("ZDIFF", args...))
	if err != nil {
		return [][2]string{}, err
	}
	if len(lists) > 0 {
		// 将数组进行拆分成多个[2]string{score,element}
		var i int
		for i = 0; i < len(lists); i++ {
			arr = append(arr, [2]string{lists[i+1], lists[i]})
			i++
		}
	}

	return arr, err
}

// ZDIFFSTORE, 比对给定的N个有序集合之间的差异，并将结果存储到指定destination中
// 若destination已经存在，则会完全覆盖
// since: 6.2.0
// key: string 第一个集合键名
// diffKeys: []string 需要比对的键名数组
// destination: string 比对结果目标存储键名
// return:
//
//	reply int 返回存储到destination中的元素总数
//
// link: https://redis.io/commands/zdiffstore/
func (c *Rzset) ZdiffStore(key string, diffKeys []string, destination string) (int, error) {
	if key == "" {
		return 0, nil
	}
	args := make([]interface{}, 0)
	args = append(args, destination, len(diffKeys)+1, key)
	for _, v := range diffKeys {
		args = append(args, v)
	}
	return redis.Int(c.conn.Do("ZDIFFSTORE", args...))
}

// ZINCRBY，按增量值递增指定成员的score值，
// 若成员不存在，则把增量值视为score值添加该成员, 若key不存在则创建
// key: string 键名
// increment: string 数值型增量字符串, 可以为负数值，代表递减
// member: string 成员
// return:
//
//	reply: string 返回设置之后该成员的新score值
//
// link: https://redis.io/commands/zincrby/
func (c *Rzset) Zincrby(key, increment, member string) (string, error) {
	return redis.String(c.conn.Do("ZINCRBY", key, increment, member))
}

// ZINTER, 返回多个给定集合的交集，此方法不指定 WITHSCORES 参数
// since: 6.2.0
// keys: []string 集合键名数组
// weights: []string 乘法因子，值必须为数值型字符串
// 使用 WEIGHTS 选项，可以为每个输入排序集指定一个乘法因子。 这意味着在传递给聚合函数之前，
// 每个输入排序集中每个元素的分数都乘以该因子。 当未给出 WEIGHTS 时，倍增因子默认为 1。
// aggregate: string , 取值： SUM | MIN | MAX
// 使用 AGGREGATE 选项, 可以指定结果集元素的score的聚合方式，默认为SUM
// return:
//
//	reply []string
//
// link: https://redis.io/commands/zinter/
func (c *Rzset) Zinter(keys, weights []string, aggregate string) ([]string, error) {
	if len(keys) == 0 {
		return []string{}, nil
	}

	args := make([]interface{}, 0)
	args = append(args, len(keys))
	for _, v := range keys {
		args = append(args, v)
	}

	if len(weights) > 0 {
		args = append(args, "WEIGHTS")
		for _, v := range weights {
			args = append(args, v)
		}
	}

	if !normal.InArray(aggregate, []string{"SUM", "MIN", "MAX"}) {
		aggregate = "SUM"
	}
	args = append(args, "AGGREGATE", aggregate)

	return redis.Strings(c.conn.Do("ZINTER", args...))
}

// ZINTER, 返回多个给定集合的交集，此方法指定 WITHSCORES 参数
// since: 6.2.0
// keys: []string 集合键名数组
// weights: []string 乘法因子，值必须为数值型字符串
// 使用 WEIGHTS 选项，可以为每个输入排序集指定一个乘法因子。 这意味着在传递给聚合函数之前，
// 每个输入排序集中每个元素的分数都乘以该因子。 当未给出 WEIGHTS 时，倍增因子默认为 1。
// aggregate: string , 取值： SUM | MIN | MAX
// 使用 AGGREGATE 选项, 可以指定结果集元素的score的聚合方式，默认为SUM
// return:
//
//	reply [][2]string{{"score", "element"}, ...}
//
// link: https://redis.io/commands/zinter/
func (c *Rzset) ZinterWithScore(keys, weights []string, aggregate string) (arr [][2]string, err error) {
	if len(keys) == 0 {
		return [][2]string{}, nil
	}

	args := make([]interface{}, 0)
	args = append(args, len(keys))
	for _, v := range keys {
		args = append(args, v)
	}

	if len(weights) > 0 {
		args = append(args, "WEIGHTS")
		for _, v := range weights {
			args = append(args, v)
		}
	}

	if !normal.InArray(aggregate, []string{"SUM", "MIN", "MAX"}) {
		aggregate = "SUM"
	}
	args = append(args, "AGGREGATE", aggregate, "WITHSCORES")

	lists, err := redis.Strings(c.conn.Do("ZINTER", args...))

	if err != nil {
		return [][2]string{}, err
	}
	if len(lists) > 0 {
		// 将数组进行拆分成多个[2]string{score,element}
		var i int
		for i = 0; i < len(lists); i++ {
			arr = append(arr, [2]string{lists[i+1], lists[i]})
			i++
		}
	}

	return arr, err
}

// ZINTERCARD, 返回多个给定集合的交集成员数，若其中一个集合为空，则结果必然为0
// since: 7.0.0
// keys: []string 集合键名数组
// limit: int  当提供可选的 LIMIT 参数（默认为 0，表示无限制）时，如果交集基数在计算中途达到极限，
// 算法将退出并产生极限作为基数。 这样的实现确保了limit低于实际交集基数的查询的显著加速。
// return:
//
//	reply: int
//
// link: https://redis.io/commands/zintercard/
func (c *Rzset) ZinterCard(keys []string, limit int) (int, error) {
	if len(keys) == 0 {
		return 0, nil
	}

	args := make([]interface{}, 0)
	args = append(args, len(keys))
	for _, v := range keys {
		args = append(args, v)
	}
	args = append(args, "LIMIT", limit)
	return redis.Int(c.conn.Do("ZINTERCARD", args...))
}

// ZINTERSTORE, 查询多个给定集合的交集，并将结果存储到指定集合中，若destination已经存在，则其值会被完全覆盖
// destination: string 结果存储集合键名
// keys: []string 集合键名数组
// weights: []string 乘法因子，值必须为数值型字符串
// 使用 WEIGHTS 选项，可以为每个输入排序集指定一个乘法因子。 这意味着在传递给聚合函数之前，
// 每个输入排序集中每个元素的分数都乘以该因子。 当未给出 WEIGHTS 时，倍增因子默认为 1。
// aggregate: string , 取值： SUM | MIN | MAX
// 使用 AGGREGATE 选项, 可以指定结果集元素的score的聚合方式，默认为SUM
// return:
//
//	reply int 返回结果集中的成员数
//
// link: https://redis.io/commands/zinterstore/
func (c *Rzset) ZinterStore(destination string, keys, weights []string, aggregate string) (int, error) {
	if len(keys) == 0 {
		return 0, nil
	}

	args := make([]interface{}, 0)
	args = append(args, destination, len(keys))
	for _, v := range keys {
		args = append(args, v)
	}

	if len(weights) > 0 {
		args = append(args, "WEIGHTS")
		for _, v := range weights {
			args = append(args, v)
		}
	}

	if !normal.InArray(aggregate, []string{"SUM", "MIN", "MAX"}) {
		aggregate = "SUM"
	}
	args = append(args, "AGGREGATE", aggregate)

	return redis.Int(c.conn.Do("ZINTERSTORE", args...))
}

// ZLEXCOUNT, 当有序集合中的所有元素都以相同的分数插入时，
// 为了强制按字典顺序排序，返回min和max区间范围内的成员数量
// key: string 键名
// min: string 起始值
// max: string 截止值
// return:
//
//	reply: int
//
// link: https://redis.io/commands/zlexcount/
func (c *Rzset) Zlexcount(key, min, max string) (string, error) {
	return redis.String(c.conn.Do("ZLEXCOUNT", min, max))
}

// ZMPOP, 从提供的键名列表中的第一个非空集合键中移出一个或多个score-member对。
// since: 7.0.0
// keys: []string 键名数组
// condition: string  取值：MIN | MAX
// count: int 返回元素的数量，取非空集合的长度和count两者的较小者
// return:
//
//	keyName: string 移出成员的集合名
//	arr: [][2]string{{score,element}, ...}
//
// link: https://redis.io/commands/zmpop/
func (c *Rzset) Zmpop(keys []string, condition string, count int) (keyName string, arr [][2]string, err error) {
	if condition == "" || !normal.InArray(condition, []string{"MIN", "MAX"}) {
		panic("参数condition传值不正确,值必须为 MIN | MAX")
	}
	if count < 1 {
		count = 1
	}
	args := make([]interface{}, 0)
	args = append(args, len(keys))
	for _, v := range keys {
		args = append(args, v)
	}
	args = append(args, condition, "COUNT", count)
	res, err := redis.Values(c.conn.Do("ZMPOP", args...))
	if len(res) == 0 || err != nil {
		return "", [][2]string{}, err
	}

	keyName, _ = redis.String(res[0], nil)
	lists, _ := redis.Values(res[1], nil)
	if len(lists) > 0 {
		// 将数组重新组装成[2]string{score,element} 格式
		for _, val := range lists {
			valArr, _ := redis.Strings(val, nil)
			arr = append(arr, [2]string{valArr[1], valArr[0]})
		}
	}
	return keyName, arr, err
}

// ZMSCORE, 返回集合中指定成员的score值，不存在的成员score返回空
// since: 6.2.0
// key: string 键名
// members: []string  成员数组
// return:
//
//	reply: map[string]string, 格式为map[element]score
//
// link: https://redis.io/commands/zmscore/
func (c *Rzset) Zmscore(key string, members []string) (map[string]string, error) {
	args := make([]interface{}, 0)
	args = append(args, key)
	for _, v := range members {
		args = append(args, v)
	}
	res, err := redis.Strings(c.conn.Do("ZMSCORE", args...))
	if err != nil {
		return map[string]string{}, err
	}
	// 集合初始化
	maps := make(map[string]string)
	for k, v := range members {
		maps[v] = res[k]
	}

	return maps, err
}

// ZPOPMAX, 移除并返回集合中指定count个数的score最大成员
// key: string 键名
// count: int 移除数量
// return:
//
//	reply: [][2]string{{score,element}, ...}
//
// link: https://redis.io/commands/zpopmax/
func (c *Rzset) ZpopMax(key string, count int) (arr [][2]string, err error) {
	if count < 1 {
		count = 1
	}

	lists, err := redis.Strings(c.conn.Do("ZPOPMAX", key, count))
	if err != nil {
		return [][2]string{}, err
	}
	if len(lists) > 0 {
		// 将数组进行拆分成多个[2]string{score,element}
		var i int
		for i = 0; i < len(lists); i++ {
			arr = append(arr, [2]string{lists[i+1], lists[i]})
			i++
		}
	}

	return arr, err
}

// ZPOPMIN, 移除并返回集合中指定count个数的score最小成员
// key: string 键名
// count: int 移除数量
// return:
//
//	reply: [][2]string{{score,element}, ...}
//
// link: https://redis.io/commands/zpopmin/
func (c *Rzset) ZpopMin(key string, count int) (arr [][2]string, err error) {
	if count < 1 {
		count = 1
	}

	lists, err := redis.Strings(c.conn.Do("ZPOPMIN", key, count))
	if err != nil {
		return [][2]string{}, err
	}
	if len(lists) > 0 {
		// 将数组进行拆分成多个[2]string{score,element}
		var i int
		for i = 0; i < len(lists); i++ {
			arr = append(arr, [2]string{lists[i+1], lists[i]})
			i++
		}
	}

	return arr, err
}

// ZRANK, 返回member在集合中的排名，score从低到高排列，
// 排名从0开始，意味着score最小的成员排名为0, 此方法不传递 WITHSCORE 参数
// key: string 键名
// member: string 指定成员
// return:
//
//	reply: int ,特别说明：member不存在或者key不存在时，也会返回0 ，
//	但是此时error信息不是nil，所以可以通过判断error值来区分是否返回的是正常排名值。
//
// link: https://redis.io/commands/zrank/
func (c *Rzset) Zrank(key, member string) (int, error) {
	return redis.Int(c.conn.Do("ZRANK", key, member))
}

// ZRANK, 返回member在集合中的排名，score从低到高排列，
// 排名从0开始，意味着score最小的成员排名为0, 此方法传递 WITHSCORE 参数
// since: 7.2.0
// key: string 键名
// member: string 指定成员
// return:
//
//	reply: []string
//
// link: https://redis.io/commands/zrank/
func (c *Rzset) ZrankWithScore(key, member string) ([]string, error) {
	return redis.Strings(c.conn.Do("ZRANK", key, member, "WITHSCORE"))
}

// ZREVRANK, 返回member在集合中的排名，score从高到低排列，
// 排名从0开始，意味着score最小的成员排名为0, 此方法不传递 WITHSCORE 参数
// key: string 键名
// member: string 指定成员
// return:
//
//	reply: int ,特别说明：member不存在或者key不存在时，也会返回0 ，
//	但是此时error信息不是nil，所以可以通过判断error值来区分是否返回的是正常排名值。
//
// link: https://redis.io/commands/zrevrank/
func (c *Rzset) Zrevrank(key, member string) (int, error) {
	return redis.Int(c.conn.Do("ZREVRANK", key, member))
}

// ZREVRANK, 返回member在集合中的排名，score从高到低排列，
// 排名从0开始，意味着score最小的成员排名为0, 此方法传递 WITHSCORE 参数
// since: 7.2.0
// key: string 键名
// member: string 指定成员
// return:
//
//	reply: []string
//
// link: https://redis.io/commands/zrevrank/
func (c *Rzset) ZrevrankWithScore(key, member string) ([]string, error) {
	return redis.Strings(c.conn.Do("ZREVRANK", key, member, "WITHSCORE"))
}

// ZREM, 移除集合中指定的N个成员，不存在的成员忽略
// key: string 键名
// members: []string 成员数组
// return:
//
//	reply: int  返回成功移除的成员数，不存在的member不计算在内
//
// link: https://redis.io/commands/zrem/
func (c *Rzset) Zrem(key string, members []string) (int, error) {
	if len(members) == 0 {
		return 0, nil
	}
	args := make([]interface{}, 0)
	args = append(args, key)
	for _, v := range members {
		args = append(args, v)
	}
	return redis.Int(c.conn.Do("ZREM", args...))
}

// ZSCORE, 返回集合中指定member的score值，member或key不存在，返回空
// key: string 键名
// member: string 指定成员
// return:
//
//	reply: string
//
// link: https://redis.io/commands/zscore/
func (c *Rzset) Zscore(key, member string) (string, error) {
	return redis.String(c.conn.Do("ZSCORE", key, member))
}

// ZUNION, 返回多个给定集合的并集，此方法不指定 WITHSCORES 参数
// since: 6.2.0
// keys: []string 集合键名数组
// weights: []string 乘法因子，值必须为数值型字符串
// 使用 WEIGHTS 选项，可以为每个输入排序集指定一个乘法因子。 这意味着在传递给聚合函数之前，
// 每个输入排序集中每个元素的分数都乘以该因子。 当未给出 WEIGHTS 时，倍增因子默认为 1。
// aggregate: string , 取值： SUM | MIN | MAX
// 使用 AGGREGATE 选项, 可以指定结果集元素的score的聚合方式，默认为SUM
// return:
//
//	reply []string
//
// link: https://redis.io/commands/zunion/
func (c *Rzset) Zunion(keys, weights []string, aggregate string) ([]string, error) {
	if len(keys) == 0 {
		return []string{}, nil
	}

	args := make([]interface{}, 0)
	args = append(args, len(keys))
	for _, v := range keys {
		args = append(args, v)
	}

	if len(weights) > 0 {
		args = append(args, "WEIGHTS")
		for _, v := range weights {
			args = append(args, v)
		}
	}

	if !normal.InArray(aggregate, []string{"SUM", "MIN", "MAX"}) {
		aggregate = "SUM"
	}
	args = append(args, "AGGREGATE", aggregate)

	return redis.Strings(c.conn.Do("ZUNION", args...))
}

// ZUNION, 返回多个给定集合的并集，此方法指定 WITHSCORES 参数
// since: 6.2.0
// keys: []string 集合键名数组
// weights: []string 乘法因子，值必须为数值型字符串
// 使用 WEIGHTS 选项，可以为每个输入排序集指定一个乘法因子。 这意味着在传递给聚合函数之前，
// 每个输入排序集中每个元素的分数都乘以该因子。 当未给出 WEIGHTS 时，倍增因子默认为 1。
// aggregate: string , 取值： SUM | MIN | MAX
// 使用 AGGREGATE 选项, 可以指定结果集元素的score的聚合方式，默认为SUM
// return:
//
//	reply [][2]string{{"score", "element"}, ...}
//
// link: https://redis.io/commands/zunion/
func (c *Rzset) ZunionWithScore(keys, weights []string, aggregate string) (arr [][2]string, err error) {
	if len(keys) == 0 {
		return [][2]string{}, nil
	}

	args := make([]interface{}, 0)
	args = append(args, len(keys))
	for _, v := range keys {
		args = append(args, v)
	}

	if len(weights) > 0 {
		args = append(args, "WEIGHTS")
		for _, v := range weights {
			args = append(args, v)
		}
	}

	if !normal.InArray(aggregate, []string{"SUM", "MIN", "MAX"}) {
		aggregate = "SUM"
	}
	args = append(args, "AGGREGATE", aggregate, "WITHSCORES")

	lists, err := redis.Strings(c.conn.Do("ZUNION", args...))

	if err != nil {
		return [][2]string{}, err
	}
	if len(lists) > 0 {
		// 将数组进行拆分成多个[2]string{score,element}
		var i int
		for i = 0; i < len(lists); i++ {
			arr = append(arr, [2]string{lists[i+1], lists[i]})
			i++
		}
	}

	return arr, err
}

// ZUNIONSTORE, 查询多个给定集合的并集，并将结果存储到指定集合中，若destination已经存在，则其值会被完全覆盖
// destination: string 结果存储集合键名
// keys: []string 集合键名数组
// weights: []string 乘法因子，值必须为数值型字符串
// 使用 WEIGHTS 选项，可以为每个输入排序集指定一个乘法因子。 这意味着在传递给聚合函数之前，
// 每个输入排序集中每个元素的分数都乘以该因子。 当未给出 WEIGHTS 时，倍增因子默认为 1。
// aggregate: string , 取值： SUM | MIN | MAX
// 使用 AGGREGATE 选项, 可以指定结果集元素的score的聚合方式，默认为SUM
// return:
//
//	reply int 返回结果集中的成员数
//
// link: https://redis.io/commands/zunionstore/
func (c *Rzset) ZunionStore(destination string, keys, weights []string, aggregate string) (int, error) {
	if len(keys) == 0 {
		return 0, nil
	}

	args := make([]interface{}, 0)
	args = append(args, destination, len(keys))
	for _, v := range keys {
		args = append(args, v)
	}

	if len(weights) > 0 {
		args = append(args, "WEIGHTS")
		for _, v := range weights {
			args = append(args, v)
		}
	}

	if !normal.InArray(aggregate, []string{"SUM", "MIN", "MAX"}) {
		aggregate = "SUM"
	}
	args = append(args, "AGGREGATE", aggregate)

	return redis.Int(c.conn.Do("ZUNIONSTORE", args...))
}

// ZRANDMEMBER, 从集合中随机返回N个成员，此方法默认不指定 WITHSCORES 参数
// since: 6.2.0
// key: string 键名
// count: int 指定返回的成员数
//
//	正数：返回结果数量是count与集合成员总数两者的较小者
//	负数：返回的结果数量是count的绝对值，此时的结果可能会存在相同成员
//
// return:
//
//	reply: []string
//
// link: https://redis.io/commands/zrandmember/
func (c *Rzset) ZrandMember(key string, count int) ([]string, error) {
	if count == 0 {
		count = 1
	}
	return redis.Strings(c.conn.Do("ZRANDMEMBER", key, count))
}

// ZRANDMEMBER, 从集合中随机返回N个成员，此方法默认指定 WITHSCORES 参数
// since: 6.2.0
// key: string 键名
// count: int 指定返回的成员数
//
//	正数：返回结果数量是count与集合成员总数两者的较小者
//	负数：返回的结果数量是count的绝对值，此时的结果可能会存在相同成员
//
// return:
//
//	reply: [][2]string{{score,element}, ...}
//
// link: https://redis.io/commands/zrandmember/
func (c *Rzset) ZrandMemberWithScore(key string, count int) (arr [][2]string, err error) {
	if count == 0 {
		count = 1
	}
	res, err := redis.Strings(c.conn.Do("ZRANDMEMBER", key, count, "WITHSCORES"))

	if err != nil {
		return [][2]string{}, err
	}
	if len(res) > 0 {
		// 将数组进行拆分成多个[2]string{score,element}
		var i int
		for i = 0; i < len(res); i++ {
			arr = append(arr, [2]string{res[i+1], res[i]})
			i++
		}
	}

	return arr, err
}
