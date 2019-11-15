package security

//shared dtos

//PermissionDTO models a security permission
type PermissionDTO struct {
	Key         string `json:"key"`
	Description string `json:"description"`
}
