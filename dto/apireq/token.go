package apireq

type GetSysAccountToken struct {
	SystemId int    `json:"system_id" validate:"required"`
	Account  string `json:"account" validate:"required"`
	Password string `json:"password" validate:"required"`
}
