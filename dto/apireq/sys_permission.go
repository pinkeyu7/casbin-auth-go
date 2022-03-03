package apireq

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
