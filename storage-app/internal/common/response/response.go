package response

type AppResponse[T any] struct {
	Success bool   `json:"sucess"`
	Message string `json:"message"`
	Data    T      `json:"data" swaggerignore:"true"`
}

func InitResponse[T any](success bool, message string, data T) AppResponse[T] {
	return AppResponse[T]{
		Success: success,
		Message: message,
		Data:    data,
	}
}
