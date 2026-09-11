package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"pkg/httperr"
	"pkg/rbac"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type RbacHandler struct {
	service  rbac.AuthorizationService
	logger   *slog.Logger
	validate *validator.Validate
}

func NewRbacHandler(rbacService rbac.AuthorizationService, logger *slog.Logger) *RbacHandler {
	return &RbacHandler{
		service:  rbacService,
		logger:   logger,
		validate: validator.New(),
	}
}

func parseBody[T any](validate *validator.Validate, r *http.Request) (*T, error) {
	req := new(T)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return nil, err
	}

	if err := validate.Struct(req); err != nil {
		return nil, err
	}

	return req, nil
}

func (h *RbacHandler) CreateRole() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := parseBody[rbac.CreateRoleRequest](h.validate, r)
		if err != nil {
			httperr.Write(w, http.StatusBadRequest, err.Error())
			return
		}

		role, err := h.service.CreateRole(ctx, req.Name, req.Description)
		if err != nil {
			if err == rbac.RoleExistsError {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(role.ToResponse())
				return
			}
			httperr.Write(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(role.ToResponse())
	}
}

func (h *RbacHandler) DeleteRole() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roleName := r.PathValue("role")
		if roleName == "" {
			httperr.Write(w, http.StatusBadRequest, "role name is required")
			return
		}

		if err := h.service.DeleteRole(r.Context(), roleName); err != nil {
			if err == rbac.RoleNotFound {
				httperr.Write(w, http.StatusNotFound, err.Error())
				return
			}
			httperr.Write(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (h *RbacHandler) CreateAction() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := parseBody[rbac.CreateActionRequest](h.validate, r)
		if err != nil {
			httperr.Write(w, http.StatusBadRequest, err.Error())
			return
		}

		action, err := h.service.CreateAction(ctx, req.Name)
		if err != nil {
			if err == rbac.ActionExistsError {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(action.ToResponse())
				return
			}
			httperr.Write(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(action.ToResponse())
	}
}

func (h *RbacHandler) CreateResource() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := parseBody[rbac.CreateResourceRequest](h.validate, r)
		if err != nil {
			httperr.Write(w, http.StatusBadRequest, err.Error())
			return
		}

		resource, err := h.service.CreateResource(ctx, req.Name, req.Description)
		if err != nil {
			if err == rbac.ResourceExistsError {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(resource.ToResponse())
				return
			}
			httperr.Write(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resource.ToResponse())
	}
}

func (h *RbacHandler) CreatePermission() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req, err := parseBody[rbac.PermissionRequest](h.validate, r)
		if err != nil {
			httperr.Write(w, http.StatusBadRequest, err.Error())
			return
		}

		perm, err := h.service.CreatePermission(ctx, req.ResourceName, req.ActionName)
		if err != nil {
			if err == rbac.ResourceNotFound || err == rbac.ActionNotFound {
				httperr.Write(w, http.StatusNotFound, err.Error())
				return
			}
			if err == rbac.PermissionExistsError {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(perm.ToResponse())
				return
			}
			httperr.Write(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(perm.ToResponse())
	}
}

func (h *RbacHandler) GrantPermissionToRole() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		roleName := r.PathValue("role_name")

		req, err := parseBody[rbac.PermissionRequest](h.validate, r)
		if err != nil {
			httperr.Write(w, http.StatusBadRequest, err.Error())
			return
		}

		rp, err := h.service.GrantPermissionToRole(ctx, roleName, req.ResourceName, req.ActionName)
		if err != nil {
			if err == rbac.RoleNotFound || err == rbac.PermissionNotFound {
				httperr.Write(w, http.StatusNotFound, err.Error())
				return
			}
			if err == rbac.RolePermissionExistsError {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(rp.ToResponse())
				return
			}
			httperr.Write(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(rp.ToResponse())

	}
}

func (h *RbacHandler) RevokePermissionFromRole() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		roleName := r.PathValue("role_name")
		resource := r.URL.Query().Get("resource")
		action := r.URL.Query().Get("action")

		err := h.service.RevokePermissionFromRole(ctx, roleName, resource, action)
		if err != nil {
			if err == rbac.PermissionNotFound || err == rbac.RoleNotFound {
				httperr.Write(w, http.StatusNotFound, err.Error())
				return
			}
			httperr.Write(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
func (h *RbacHandler) AssignRoleToUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := uuid.Parse(r.PathValue("user_id"))
		if err != nil {
			httperr.Write(w, http.StatusBadRequest, "invalid user id")
			return
		}

		req, err := parseBody[rbac.AssignRoleRequest](h.validate, r)
		if err != nil {
			httperr.Write(w, http.StatusBadRequest, err.Error())
			return
		}

		err = h.service.AssignRoleToUser(r.Context(), userID, req.RoleName)
		if err != nil {
			if err == rbac.RoleNotFound {
				httperr.Write(w, http.StatusNotFound, err.Error())
				return
			}
			httperr.Write(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
func (h *RbacHandler) RemoveRoleFromUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := uuid.Parse(r.PathValue("user_id"))
		if err != nil {
			httperr.Write(w, http.StatusBadRequest, "invalid user id")
			return
		}

		roleName := r.PathValue("role_name")
		if roleName == "" {
			httperr.Write(w, http.StatusBadRequest, "role name is required")
			return
		}

		err = h.service.RemoveRoleFromUser(r.Context(), userID, roleName)
		if err != nil {
			httperr.Write(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
