package common

type ServiceName string

const (
	Order ServiceName = "order-service"
	User  ServiceName = "user-service"
)

func (sn ServiceName) String() string {
	return string(sn)
}

func (ServiceName) GetService(name string) ServiceName {
	switch name {
	case "order":
		return Order
	case "user":
		return User
	default:
		return ""
	}
}
