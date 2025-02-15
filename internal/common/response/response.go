package response

type AppResponse struct {
	Success bool        `json:"sucess"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func InitResponse(success bool, message string, data interface{}) AppResponse {
	return AppResponse{
		Success: success,
		Message: message,
		Data:    data,
	}
}
