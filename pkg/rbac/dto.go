package rbac

import (
	"time"

	"github.com/google/uuid"
)

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required" validate:"required,min=2,max=30,alpha"`
	Description string `json:"description" validate:"omitempty,max=100"`
}

type CreateActionRequest struct {
	Name string `json:"name" binding:"required" validate:"required,min=2,max=30,alpha"`
}

type CreateResourceRequest struct {
	Name        string `json:"name" binding:"required" validate:"required"`
	Description string `json:"description" validate:"omitempty,max=100"`
}

type PermissionRequest struct {
	ResourceName string `json:"resource_name" binding:"required" validate:"required,min=2"`
	ActionName   string `json:"action_name" binding:"required" validate:"required,min=2,max=30,alpha"`
}

type AssignRoleRequest struct {
	RoleName string `json:"role_name" binding:"required" validate:"required,min=2,max=30,alpha"`
}

type RemoveRoleRequest struct {
	UserID   uuid.UUID `json:"user_id" binding:"required" validate:"required,uuid"`
	RoleName string    `json:"role_name" binding:"required" validate:"required,min=2,max=30,alpha"`
}

type RoleResponse struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ActionResponse struct {
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ResourceResponse struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PermissionResponse struct {
	ResourceName string            `json:"resource_name"`
	ActionName   string            `json:"action_name"`
	Resource     *ResourceResponse `json:"resource,omitempty"`
	Action       *ActionResponse   `json:"action,omitempty"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

type RolePermissionResponse struct {
	RoleName     string              `json:"role_name"`
	ResourceName string              `json:"resource_name"`
	ActionName   string              `json:"action_name"`
	Role         *RoleResponse       `json:"role,omitempty"`
	Permission   *PermissionResponse `json:"permission,omitempty"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

type UserRoleResponse struct {
	UserID    uuid.UUID     `json:"user_id"`
	RoleName  string        `json:"role_name"`
	Role      *RoleResponse `json:"role,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// type RoleDetailResponse struct {
// 	Name        string                   `json:"name"`
// 	Description string                   `json:"description"`
// 	Permissions []RolePermissionResponse `json:"permissions,omitempty"`
// 	CreatedAt   time.Time                `json:"created_at"`
// 	UpdatedAt   time.Time                `json:"updated_at"`
// }

// type PermissionDetailResponse struct {
// 	ResourceName string            `json:"resource_name"`
// 	ActionName   string            `json:"action_name"`
// 	Resource     *ResourceResponse `json:"resource"`
// 	Action       *ActionResponse   `json:"action"`
// 	CreatedAt    time.Time         `json:"created_at"`
// 	UpdatedAt    time.Time         `json:"updated_at"`
// }

// type UserRoleDetailResponse struct {
// 	UserID    uuid.UUID           `json:"user_id"`
// 	RoleName  string              `json:"role_name"`
// 	Role      *RoleDetailResponse `json:"role"`
// 	CreatedAt time.Time           `json:"created_at"`
// 	UpdatedAt time.Time           `json:"updated_at"`
// }

func (r *Role) ToResponse() *RoleResponse {
	if r == nil {
		return nil
	}
	return &RoleResponse{
		Name:        r.Name,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func (a *Action) ToResponse() *ActionResponse {
	if a == nil {
		return nil
	}
	return &ActionResponse{
		Name:      a.Name,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func (res *Resource) ToResponse() *ResourceResponse {
	if res == nil {
		return nil
	}
	return &ResourceResponse{
		Name:        res.Name,
		Description: res.Description,
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
	}
}

func (p *Permission) ToResponse() *PermissionResponse {
	if p == nil {
		return nil
	}
	return &PermissionResponse{
		ResourceName: p.ResourceName,
		ActionName:   p.ActionName,
		Resource:     p.Resource.ToResponse(),
		Action:       p.Action.ToResponse(),
		CreatedAt:    p.CreatedAt,
		UpdatedAt:    p.UpdatedAt,
	}
}

func (rp *RolePermission) ToResponse() *RolePermissionResponse {
	if rp == nil {
		return nil
	}
	return &RolePermissionResponse{
		RoleName:     rp.RoleName,
		ResourceName: rp.ResourceName,
		ActionName:   rp.ActionName,
		Role:         rp.Role.ToResponse(),
		Permission:   rp.Permission.ToResponse(),
		CreatedAt:    rp.CreatedAt,
		UpdatedAt:    rp.UpdatedAt,
	}
}

func (ur *UserRole) ToResponse() *UserRoleResponse {
	if ur == nil {
		return nil
	}
	return &UserRoleResponse{
		UserID:    ur.UserID,
		Role:      ur.Role.ToResponse(),
		CreatedAt: ur.CreatedAt,
		UpdatedAt: ur.UpdatedAt,
	}
}
