package recommend_service

import (
	"fmt"
	"justus/dao"
)

// DeleteRecommendInfo 删除过期的推荐信息
func DeleteRecommendInfo() {
	//查询每个语言的推荐数量
	languageData, err := dao.GetLanguageRecommendCount()
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	if len(languageData) > 0 {
		for _, v := range languageData {
			//fmt.Println("language:", v.Lange)
			//fmt.Println("count:", v.Count)
			if v.Count > 2000 {
				//查找2001条推荐信息 并删除
				dao.GetRecommendDeleteInfo(v.Lange, 2000)
			}
		}
	}
}
