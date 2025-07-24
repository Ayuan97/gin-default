package rediskey

const WhiteList = "yq_white_list" //白名单

// GetWhiteListKey 获取白名单key
func GetWhiteListKey() string {
	return WhiteList
}
