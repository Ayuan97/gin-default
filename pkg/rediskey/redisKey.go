package rediskey

const WhiteList = "yq_white_list" //白名单

// GetWhiteListKey 获取白名单key
func GetWhiteListKey() string {
	return WhiteList
}

// 多租户Key前缀
func TenantPrefix(tenantID uint) string {
	return "justus:tenant:" + itoa(tenantID) + ":"
}

// 租户菜单白名单缓存key
func TenantMenuWhitelistKey(tenantID uint) string {
	return TenantPrefix(tenantID) + "menus:whitelist"
}

// 用户在租户下的菜单树缓存key
func TenantUserMenuTreeKey(tenantID uint, userID uint) string {
	return TenantPrefix(tenantID) + "admin:" + itoa(userID) + ":menu_tree"
}

// itoa 简易无依赖整型转字符串
func itoa(v uint) string {
	if v == 0 {
		return "0"
	}
	// 手写转换防止引入strconv在工具包中扩大依赖面
	var buf [20]byte
	i := len(buf)
	n := v
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
