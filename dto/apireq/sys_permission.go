package apireq

type ListSysPermission struct {
	AccountId int `form:"account_id" validate:"required"`
	Page      int `form:"page" validate:"required"`
	PerPage   int `form:"per_page" validate:"required"`
	SystemId  int `form:"system_id" validate:"omitempty"`
}

type AddSysPermission struct {
	AccountId    int    `json:"account_id" validate:"required"`
	SystemId     int    `json:"system_id" validate:"required"`
	AllowApiPath string `json:"allow_api_path" validate:"omitempty"`
	Action       string `json:"action" validate:"omitempty"`
	Slug         string `json:"slug" validate:"required"`
	Description  string `json:"description" validate:"required"`
}

type EditSysPermission struct {
	AccountId   int    `json:"account_id" validate:"required"`
	Slug        string `json:"slug" validate:"required"`
	Description string `json:"description" validate:"required"`
}

type DeleteSysPermission struct {
	AccountId int `json:"account_id" validate:"required"`
}
