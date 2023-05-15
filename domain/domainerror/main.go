package domainerror

type Interface interface {
	Code() string
	GetMessage() string
	GetHttpCode() int
}
