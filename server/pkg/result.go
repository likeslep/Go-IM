package pkg

// 统一API返回格式
func Success(data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"code": 200,
		"msg":  "success",
		"data": data,
	}
}

func Fail(msg string) map[string]interface{} {
	return map[string]interface{}{
		"code": 500,
		"msg":  msg,
		"data": nil,
	}
}