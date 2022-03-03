package apireq

type ListSysRole struct {
	AccountId int `form:"account_id" validate:"required"`
	Page      int `form:"page" validate:"required"`
	PerPage   int `form:"per_page" validate:"required"`
	SystemId  int `form:"system_id" validate:"omitempty"`
}

type AddSysRole struct {
	AccountId     int    `json:"account_id" validate:"required"`
	SystemId      int    `json:"system_id" validate:"required"`
	Name          string `json:"name" validate:"required"`
	DisplayName   string `json:"display_name" validate:"required"`
	PermissionIds []int  `json:"permission_ids" validate:"required"`
}

type EditSysRole struct {
	AccountId     int    `json:"account_id" validate:"required"`
	Name          string `json:"name" validate:"omitempty"`
	DisplayName   string `json:"display_name" validate:"omitempty"`
	IsDisable     *bool  `json:"is_disable" validate:"omitempty"`
	PermissionIds []int  `json:"permission_ids" validate:"omitempty"`
}
