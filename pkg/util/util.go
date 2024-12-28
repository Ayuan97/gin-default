package util

import (
	"crypto/md5"
	"encoding/hex"
	"justus/pkg/setting"
	"sort"
	"strings"
)

// Setup Initialize the util
func Setup() {
	jwtSecret = []byte(setting.AppSetting.JwtSecret)
}

//数字转换
func NumTransform(number int) int {
	return number
	//if number > 1000000{
	//	number2 := float64(number)/1000000
	//	numStr := fmt.Sprintf("%.2f",number2)
	//	return numStr+"M"
	//}else if number > 1000{
	//	number2 := float64(number)/1000
	//	numStr := fmt.Sprintf("%.2f",number2)
	//	return numStr+"K"
	//}
	//return fmt.Sprintf("%d",number)
}

//获取图片链接
func GetImageUrl(imagePath string) string {
	if strings.Contains(imagePath, "http") {
		return imagePath
	} else {
		return setting.AppSetting.ImageUrl + "/" + imagePath
	}
}

// int 是否存在数组中
func In(i int, intArray []int) int {

	sort.Ints(intArray)

	index := sort.SearchInts(intArray, i)

	if index < len(intArray) && intArray[index] == i {

		return 1

	}

	return 0
}

//通过map键的唯一性去重
func RemoveRepeatedElement(s []int) []int {
	result := make([]int, 0)
	m := make(map[int]bool) //map的值不重要
	for _, v := range s {
		if _, ok := m[v]; !ok {
			result = append(result, v)
			m[v] = true
		}
	}
	return result
}

// 返回一个32位md5加密后的字符串
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// 根据map的键名排序
func SortMapByKey(m map[string]interface{}) (map[string]interface{}, []string) {
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}
	// 字符串数组排序 倒序
	//sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	sort.Strings(keys)
	var sortedMap []string
	var newMap = make(map[string]interface{})
	for _, k := range keys {
		sortedMap = append(sortedMap, k)
		newMap[k] = m[k]
		//fmt.Println(k, m[k])
	}
	return newMap, sortedMap
}
