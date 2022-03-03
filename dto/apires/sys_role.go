package apires

import "casbin-auth-go/dto/model"

type SysRoleWithPermissionIds struct {
	*model.SysRole
	PermissionIds []int `json:"permission_ids"`
}
