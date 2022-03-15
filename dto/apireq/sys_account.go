package apireq

type ListSysAccount struct {
	AccountId int `form:"account_id" validate:"required"`
	Page      int `form:"page" validate:"required"`
	PerPage   int `form:"per_page" validate:"required"`
	SystemId  int `form:"system_id" validate:"omitempty"`
}

type AddSysAccount struct {
	AccountId int    `json:"account_id" validate:"required"`
	SystemId  int    `json:"system_id" validate:"required"`
	Account   string `json:"account" validate:"required"`
	Phone     string `json:"phone" validate:"omitempty"`
	Email     string `json:"email" validate:"required"`
	Name      string `json:"name" validate:"omitempty"`
	RoleId    int    `json:"role_id" validate:"required"`
}

type EditSysAccount struct {
	AccountId int    `json:"account_id" validate:"required"`
	Phone     string `json:"phone" validate:"omitempty"`
	Email     string `json:"email" validate:"omitempty"`
	Name      string `json:"name" validate:"omitempty"`
	IsDisable *bool  `json:"is_disable" validate:"omitempty"`
	RoleId    int    `json:"role_id" validate:"omitempty"`
}

type ForgotSysAccountPassword struct {
	AccountId int `json:"account_id" validate:"required"`
}

type ChangePassword struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,max=16,min=8,necsfield=OldPassword"`
}
