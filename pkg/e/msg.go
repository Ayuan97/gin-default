package e

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	ERROR:          "fail",
	CONTENT_EMPTY:  "内容为空",
	INVALID_PARAMS: "请求参数错误",
	SIGN_ERROR:     "签名错误",

	ERROR_EXIST_TAG:       "已存在该标签名称",
	ERROR_EXIST_TAG_FAIL:  "获取已存在标签失败",
	ERROR_NOT_EXIST_TAG:   "该标签不存在",
	ERROR_GET_TAGS_FAIL:   "获取所有标签失败",
	ERROR_COUNT_TAG_FAIL:  "统计标签失败",
	ERROR_ADD_TAG_FAIL:    "新增标签失败",
	ERROR_EDIT_TAG_FAIL:   "修改标签失败",
	ERROR_DELETE_TAG_FAIL: "删除标签失败",
	ERROR_EXPORT_TAG_FAIL: "导出标签失败",
	ERROR_IMPORT_TAG_FAIL: "导入标签失败",

	ERROR_NOT_EXIST_ARTICLE:        "该文章不存在",
	ERROR_CHECK_EXIST_ARTICLE_FAIL: "检查文章是否存在失败",
	ERROR_ADD_ARTICLE_FAIL:         "新增文章失败",
	ERROR_DELETE_ARTICLE_FAIL:      "删除文章失败",
	ERROR_EDIT_ARTICLE_FAIL:        "修改文章失败",
	ERROR_COUNT_ARTICLE_FAIL:       "统计文章失败",
	ERROR_GET_ARTICLES_FAIL:        "获取多个文章失败",
	ERROR_GET_ARTICLE_FAIL:         "获取单个文章失败",
	ERROR_GEN_ARTICLE_POSTER_FAIL:  "生成文章海报失败",

	// 认证相关错误消息
	ERROR_AUTH_CHECK_TOKEN_FAIL:    "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: "Token已超时",
	ERROR_AUTH_TOKEN:               "Token生成失败",
	ERROR_AUTH:                     "Token错误",

	// 文件上传相关错误消息
	ERROR_UPLOAD_SAVE_IMAGE_FAIL:    "保存图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FAIL:   "检查图片失败",
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT: "校验图片错误，图片格式或大小有问题",

	// 用户相关错误消息
	ERROR_USER_NOT_FOUND:      "用户不存在",
	ERROR_USER_ALREADY_EXIST:  "用户已存在",
	ERROR_USER_CREATE_FAIL:    "创建用户失败",
	ERROR_USER_UPDATE_FAIL:    "更新用户失败",
	ERROR_USER_DELETE_FAIL:    "删除用户失败",
	ERROR_USER_STATUS_INVALID: "用户状态无效",

	// 权限和角色相关错误消息
	ERROR_PERMISSION_DENIED:       "权限不足",
	ERROR_ROLE_NOT_FOUND:          "角色不存在",
	ERROR_ROLE_ALREADY_EXIST:      "角色已存在",
	ERROR_ROLE_CREATE_FAIL:        "创建角色失败",
	ERROR_ROLE_UPDATE_FAIL:        "更新角色失败",
	ERROR_ROLE_DELETE_FAIL:        "删除角色失败",
	ERROR_ROLE_IN_USE:             "角色正在使用中，无法删除",
	ERROR_ADMIN_ROLE_PROTECT:      "管理员角色受保护，无法修改或删除",
	ERROR_INSUFFICIENT_PERMISSION: "权限不足，无法执行此操作",

	// 管理员相关错误消息
	ERROR_ADMIN_NOT_FOUND:      "管理员不存在",
	ERROR_ADMIN_CREATE_FAIL:    "创建管理员失败",
	ERROR_ADMIN_UPDATE_FAIL:    "更新管理员失败",
	ERROR_ADMIN_DELETE_FAIL:    "删除管理员失败",
	ERROR_ADMIN_SELF_OPERATION: "不能对自己执行此操作",

	// 数据库相关错误消息
	ERROR_DATABASE_CONNECTION: "数据库连接失败",
	ERROR_DATABASE_QUERY:      "数据库查询失败",
	ERROR_DATABASE_INSERT:     "数据库插入失败",
	ERROR_DATABASE_UPDATE:     "数据库更新失败",
	ERROR_DATABASE_DELETE:     "数据库删除失败",

	// 缓存相关错误消息
	ERROR_CACHE_GET: "缓存获取失败",
	ERROR_CACHE_SET: "缓存设置失败",
	ERROR_CACHE_DEL: "缓存删除失败",

	// 文件相关错误消息
	ERROR_FILE_NOT_FOUND:  "文件不存在",
	ERROR_FILE_READ_FAIL:  "文件读取失败",
	ERROR_FILE_WRITE_FAIL: "文件写入失败",

	// 网络相关错误消息
	ERROR_NETWORK_TIMEOUT: "网络请求超时",
	ERROR_NETWORK_REQUEST: "网络请求失败",

	// 业务逻辑错误消息
	ERROR_BUSINESS_LOGIC:  "业务逻辑错误",
	ERROR_DATA_VALIDATION: "数据验证失败",

	// 系统相关错误消息
	ERROR_SYSTEM_MAINTENANCE: "系统维护中",
	ERROR_SYSTEM_OVERLOAD:    "系统负载过高",
	ERROR_SYSTEM_CONFIG:      "系统配置错误",
}

// GetMsg get error information based on Code
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
