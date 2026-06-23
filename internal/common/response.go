package common

// APIResponse 标准 API 响应格式
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(data interface{}) APIResponse {
	return APIResponse{
		Success: true,
		Data:    data,
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(err error) APIResponse {
	return APIResponse{
		Success: false,
		Error:   err.Error(),
	}
}

// NewErrorResponseStr 创建错误响应（字符串）
func NewErrorResponseStr(errMsg string) APIResponse {
	return APIResponse{
		Success: false,
		Error:   errMsg,
	}
}
