package service

type ServiceError struct {
	Code  int    `json:"code"`
	Short string `json:"short"`
	Long  string `json:"long"`
}

func (e *ServiceError) Error() string {
	return e.Long
}

const (
	GeneralErrorCode = 1000
)

var (
	InternalError = &ServiceError{GeneralErrorCode + 1, "internal", "Internal server error"}
)
