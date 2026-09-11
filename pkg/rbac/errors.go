package rbac

import (
	"errors"
)

var RoleExistsError = errors.New("role already exists")
var ActionExistsError = errors.New("action already exists")
var ResourceExistsError = errors.New("resource already exists")
var PermissionExistsError = errors.New("permission already exists")
var RolePermissionExistsError = errors.New("role permission already exists")

var RoleNotFound = errors.New("role not found")
var ActionNotFound = errors.New("action not found")
var ResourceNotFound = errors.New("resource not found")
var PermissionNotFound = errors.New("permission not found")
