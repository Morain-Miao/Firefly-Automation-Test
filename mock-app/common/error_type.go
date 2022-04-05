package common

const (
	ServerError     = 1000 // 系统错误
	NotFoundError   = 1001 // 401错误
	UnknownError    = 1002 // 未知错误
	ParameterError  = 1003 // 参数错误
	AuthError       = 1004 // 错误
	DataBaseError   = 1005 //数据库错误
	JsonDecodeError = 1006 //Json解码错误
)

// ServerException 500 错误处理
func ServerException(message string) *ApiException {
	return newApiException(ServerError, message, "系统错误")
}

// NotFoundException 404 错误
func NotFoundException(message string) *ApiException {
	return newApiException(NotFoundError, message, "404错误")
}

// UnknownException 未知错误
func UnknownException(message string) *ApiException {
	return newApiException(UnknownError, message, "未知错误")
}

// ParameterException 参数错误
func ParameterException(message string) *ApiException {
	return newApiException(ParameterError, message, "参数错误")
}

func DataBaseException(message string) *ApiException {
	return newApiException(DataBaseError, message, "数据库错误")
}

func JsonDecodeException(message string) *ApiException {
	return newApiException(JsonDecodeError, message, "Json解析错误！")
}
