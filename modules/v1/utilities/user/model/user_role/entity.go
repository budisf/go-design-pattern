package model

type UserRolesRawQueryResult struct {
	ID           *uint   `json:"user_id"`
	AuthServerId *uint   `json:"auth_server_id"`
	Nip          *string `json:"nip"`
	Name         *string `json:"name"`
	RoleId       *uint   `json:"role_id"`
	RoleName     *string `json:"role_name"`
	RoleLabel    *string `json:"role_label"`
	ParentId     *uint   `json:"parent_id"`
	NameRoleHead *string `json:"name_role_head"`
	CreatedAt    *uint64 `json:"created_at"`
	UpdatedAt    *uint64 `json:"updated_at"`
}
