package handler

func ErrorObj(errMsg string) map[string]interface{} {
	return map[string]interface{}{"error": errMsg}
}

func BadRequestObj() map[string]interface{} {
	return ErrorObj("bad request object")
}

func InternalErrorObj() map[string]interface{} {
	return ErrorObj("internal server error")
}
