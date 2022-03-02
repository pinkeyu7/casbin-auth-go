package apireq

type AddSystem struct {
	AccountId        int      `json:"account_id" validate:"required"`
	Name             string   `json:"name" validate:"required"`
	SystemType       string   `json:"system_type" validate:"required"`
	Tag              string   `json:"tag" validate:"required"`
	Email            string   `json:"email" validate:"required"`
	Address          string   `json:"address" validate:"omitempty,max=255"`
	Tel              string   `json:"tel" validate:"omitempty,max=20"`
	Uuid             string   `json:"uuid" validate:"required"`
	Quota            int      `json:"quota" validate:"omitempty"`
	IpAddress        []string `json:"ip_address" validate:"omitempty"`
	MacAddress       []string `json:"mac_address" validate:"omitempty"`
	Principal        string   `json:"principal" validate:"omitempty,max=64"`
	Salesman         string   `json:"salesman" validate:"omitempty,max=64"`
	SalesmanPhone    string   `json:"salesman_phone" validate:"omitempty,max=20"`
	CopyFromSystemId int      `json:"copy_from_system" validate:"omitempty"`
}

type EditSystem struct {
	AccountId  int      `json:"account_id" validate:"required"`
	Name       string   `json:"name" validate:"omitempty"`
	Address    string   `json:"address" validate:"omitempty,max=255"`
	Tel        string   `json:"tel" validate:"omitempty,max=32"`
	IsDisable  *bool    `json:"is_disable" validate:"required"`
	IpAddress  []string `json:"ip_address" validate:"omitempty"`
	MacAddress []string `json:"mac_address" validate:"omitempty"`
}
