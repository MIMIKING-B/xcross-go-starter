package convert

// UniqueSlice 切片去重
// 泛型函数：支持所有可比较类型（int/string/float/bool等）的切片去重
// 参数：
//
//	languages - 需要去重的原始切片
//
// 返回值：
//
//	[]K - 去重后的新切片，保持原始元素出现顺序
//
// 实现原理：使用map记录已存在元素，遍历切片只添加未出现过的元素
func UniqueSlice[K comparable](languages []K) []K {
	// 初始化结果切片，预分配容量提升性能
	result := make([]K, 0, len(languages))
	// 空结构体不占用内存，用于标记元素是否已存在
	temp := map[K]struct{}{}

	for _, item := range languages {
		// 判断元素是否已存在于map中
		if _, ok := temp[item]; !ok {
			// 不存在则标记为已存在
			temp[item] = struct{}{}
			// 追加到结果切片
			result = append(result, item)
		}
	}
	return result
}

// Remove 根据自定义条件删除切片中第一个匹配的元素
// 非泛型：适用于[]interface{}类型切片
// 参数：
//
//	sl - 原始切片
//	f - 匹配函数，返回true表示需要删除该元素
//
// 返回值：
//
//	[]interface{} - 删除元素后的切片（修改原切片底层数组）
//
// 注意：仅删除**第一个**匹配成功的元素，删除效率高（替换为最后一个元素+截断）
func Remove(sl []interface{}, f func(v1 interface{}) bool) []interface{} {
	for k, v := range sl {
		// 调用匹配函数，判断是否需要删除当前元素
		if f(v) {
			// 将最后一个元素覆盖到当前位置
			sl[k] = sl[len(sl)-1]
			// 截断切片，删除最后一个元素（已复制到前面）
			sl = sl[:len(sl)-1]
			// 只删除第一个匹配项，直接返回
			return sl
		}
	}
	// 无匹配元素，直接返回原切片
	return sl
}

// RemoveSlice 删除切片中第一个等于指定值的元素
// 泛型函数：支持所有可比较类型切片
// 参数：
//
//	src - 原始切片
//	sub - 需要删除的目标元素
//
// 返回值：
//
//	[]K - 删除元素后的切片（修改原切片底层数组）
//
// 实现原理：找到元素后，将后面元素向前复制一位，再截断切片
func RemoveSlice[K comparable](src []K, sub K) []K {
	for k, v := range src {
		// 找到第一个匹配的元素
		if v == sub {
			// 将k+1及之后的元素向前复制一位
			copy(src[k:], src[k+1:])
			// 截断最后一位元素，完成删除
			return src[:len(src)-1]
		}
	}
	// 未找到目标元素，返回原切片
	return src
}

// DifferenceSlice 计算两个切片的差集（只在s2中存在、不在s1中存在的元素）
// 泛型函数：支持所有可比较类型
// 示例：
//
//	slice1 := []int{1,2,3,4,5}
//	slice2 := []int{4,5,6,7,8}
//	输出：[6 7 8]
//
// 参数：
//
//	s1 - 基准切片
//	s2 - 对比切片
//
// 返回值：
//
//	[]T - 差集切片（元素在s2中，不在s1中）
func DifferenceSlice[T comparable](s1, s2 []T) []T {
	// 用map存储s1的所有元素，用于O(1)时间复杂度查找
	m := make(map[T]bool)
	for _, item := range s1 {
		m[item] = true
	}

	// 存储差集结果
	var diff []T
	for _, item := range s2 {
		// 元素不在s1中，则加入差集
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return diff
}
