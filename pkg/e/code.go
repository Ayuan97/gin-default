package e

const (
	SUCCESS        = 200
	ERROR          = 202
	CONTENT_EMPTY  = 204
	INVALID_PARAMS = 400
	SIGN_ERROR     = 401

	ERROR_EXIST_TAG       = 10001
	ERROR_EXIST_TAG_FAIL  = 10002
	ERROR_NOT_EXIST_TAG   = 10003
	ERROR_GET_TAGS_FAIL   = 10004
	ERROR_COUNT_TAG_FAIL  = 10005
	ERROR_ADD_TAG_FAIL    = 10006
	ERROR_EDIT_TAG_FAIL   = 10007
	ERROR_DELETE_TAG_FAIL = 10008
	ERROR_EXPORT_TAG_FAIL = 10009
	ERROR_IMPORT_TAG_FAIL = 10010

	ERROR_NOT_EXIST_ARTICLE        = 10011
	ERROR_CHECK_EXIST_ARTICLE_FAIL = 10012
	ERROR_ADD_ARTICLE_FAIL         = 10013
	ERROR_DELETE_ARTICLE_FAIL      = 10014
	ERROR_EDIT_ARTICLE_FAIL        = 10015
	ERROR_COUNT_ARTICLE_FAIL       = 10016
	ERROR_GET_ARTICLES_FAIL        = 10017
	ERROR_GET_ARTICLE_FAIL         = 10018
	ERROR_GEN_ARTICLE_POSTER_FAIL  = 10019

	// 认证相关错误码
	ERROR_AUTH_CHECK_TOKEN_FAIL    = 20001
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT = 20002
	ERROR_AUTH_TOKEN               = 20003
	ERROR_AUTH                     = 20004

	// 文件上传相关错误码
	ERROR_UPLOAD_SAVE_IMAGE_FAIL    = 30001
	ERROR_UPLOAD_CHECK_IMAGE_FAIL   = 30002
	ERROR_UPLOAD_CHECK_IMAGE_FORMAT = 30003

	// 用户相关错误码
	ERROR_USER_NOT_FOUND      = 40001
	ERROR_USER_ALREADY_EXIST  = 40002
	ERROR_USER_CREATE_FAIL    = 40003
	ERROR_USER_UPDATE_FAIL    = 40004
	ERROR_USER_DELETE_FAIL    = 40005
	ERROR_USER_STATUS_INVALID = 40006

	// 权限和角色相关错误码
	ERROR_PERMISSION_DENIED       = 41001
	ERROR_ROLE_NOT_FOUND          = 41002
	ERROR_ROLE_ALREADY_EXIST      = 41003
	ERROR_ROLE_CREATE_FAIL        = 41004
	ERROR_ROLE_UPDATE_FAIL        = 41005
	ERROR_ROLE_DELETE_FAIL        = 41006
	ERROR_ROLE_IN_USE             = 41007
	ERROR_ADMIN_ROLE_PROTECT      = 41008
	ERROR_INSUFFICIENT_PERMISSION = 41009

	// 管理员相关错误码
	ERROR_ADMIN_NOT_FOUND      = 42001
	ERROR_ADMIN_CREATE_FAIL    = 42002
	ERROR_ADMIN_UPDATE_FAIL    = 42003
	ERROR_ADMIN_DELETE_FAIL    = 42004
	ERROR_ADMIN_SELF_OPERATION = 42005

	// 数据库相关错误码
	ERROR_DATABASE_CONNECTION = 50001
	ERROR_DATABASE_QUERY      = 50002
	ERROR_DATABASE_INSERT     = 50003
	ERROR_DATABASE_UPDATE     = 50004
	ERROR_DATABASE_DELETE     = 50005

	// 缓存相关错误码
	ERROR_CACHE_GET = 60001
	ERROR_CACHE_SET = 60002
	ERROR_CACHE_DEL = 60003

	// 文件相关错误码
	ERROR_FILE_NOT_FOUND  = 70001
	ERROR_FILE_READ_FAIL  = 70002
	ERROR_FILE_WRITE_FAIL = 70003

	// 网络相关错误码
	ERROR_NETWORK_TIMEOUT = 80001
	ERROR_NETWORK_REQUEST = 80002

	// 业务逻辑错误码
	ERROR_BUSINESS_LOGIC  = 90001
	ERROR_DATA_VALIDATION = 90002

	// 系统相关错误码
	ERROR_SYSTEM_MAINTENANCE = 91001
	ERROR_SYSTEM_OVERLOAD    = 91002
	ERROR_SYSTEM_CONFIG      = 91003

	// ZincSearch 相关错误码
	ERROR_ZINC_CONNECTION_FAILED = 95001
	ERROR_ZINC_INDEX_NOT_FOUND   = 95002
	ERROR_ZINC_INDEX_CREATE_FAIL = 95003
	ERROR_ZINC_DOCUMENT_FAIL     = 95004
	ERROR_ZINC_SEARCH_FAIL       = 95005
)
