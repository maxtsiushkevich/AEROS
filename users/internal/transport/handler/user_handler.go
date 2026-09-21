package handler

import (
	"errors"
	"net/http"
	"pkg/helpers"
	domain_err "users/internal/domain/errors"
	"users/internal/domain/service"
	"users/internal/transport/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

func (h *UserHandler) Registration(c *gin.Context) {
	ctx := c.Request.Context()

	var req dto.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	access, refresh, err := h.service.RegisterUser(ctx, req.Name, req.Email, req.Birthday, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, domain_err.ErrInvalidName),
			errors.Is(err, domain_err.ErrInvalidEmail),
			errors.Is(err, domain_err.ErrBirthdayInvalid),
			errors.Is(err, domain_err.ErrNotAdult):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal error",
			})
		}
		return
	}

	c.SetCookie(
		"refresh_token",
		*refresh,
		60*60*24*30,
		"/",
		"",
		false,
		true,
	)

	c.JSON(http.StatusCreated, dto.TokensResponse{
		AccessToken: *access,
	})
}

func (h *UserHandler) Profile(c *gin.Context) {
	ctx := c.Request.Context()

	claims, ok := helpers.ClaimsFromContextGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "claims not found",
		})
		return
	}

	u, err := h.service.Profile(ctx, claims.Id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get profile",
		})
		return
	}

	c.JSON(http.StatusOK, dto.ToUserResponse(u))
}

func (h *UserHandler) Activate(c *gin.Context) {
	ctx := c.Request.Context()

	claims, ok := helpers.ClaimsFromContextGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "claims not found",
		})
		return
	}

	id := c.Param("userId")
	parsed, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Bad request",
		})
		return
	}

	if claims.Id != parsed {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not allowed to activate",
		})
		return
	}

	err = h.service.ActivateUser(ctx, parsed)
	if err != nil {
		if errors.Is(err, domain_err.ErrUserAlreadyActivated) {
			c.JSON(http.StatusConflict, gin.H{
				"error": "User is already activated",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to activate user",
		})
		return
	}
}
