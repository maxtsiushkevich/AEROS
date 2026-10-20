package handler

import (
	"errors"
	"net/http"
	"pkg/helpers"
	domainerr "users/internal/domain/errors"
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
		case errors.Is(err, domainerr.ErrInvalidName),
			errors.Is(err, domainerr.ErrInvalidEmail),
			errors.Is(err, domainerr.ErrBirthdayInvalid),
			errors.Is(err, domainerr.ErrNotAdult):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, domainerr.ErrUserAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
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

func (h *UserHandler) ActivateUserByCode(c *gin.Context) {
	ctx := c.Request.Context()

	claims, ok := helpers.ClaimsFromContextGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "claims not found",
		})
		return
	}

	var req dto.ActivateUserByCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.service.ActivateUserByCode(ctx, claims.Id, req.Code)

	if err != nil {
		switch {
		case errors.Is(err, domainerr.ErrInvalidOTP):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid activation code",
			})
		case errors.Is(err, domainerr.ErrUserAlreadyActivated):
			c.JSON(http.StatusConflict, gin.H{
				"error": "user is already activated",
			})
		case errors.Is(err, domainerr.ErrUserActivationFailed):
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "activation code expired or unavailable",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to activate user",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "activated"})
}

func (h *UserHandler) RequestActivationCode(c *gin.Context) {
	ctx := c.Request.Context()

	claims, ok := helpers.ClaimsFromContextGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "claims not found",
		})
		return
	}

	err := h.service.RequestActivationCode(ctx, claims.Id, claims.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to send activation code",
		})
	}

	c.JSON(http.StatusAccepted, gin.H{"status": "activation code sent"})
}

func (h *UserHandler) AdminUserByAdmin(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("userId")
	parsed, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Bad request",
		})
		return
	}

	err = h.service.AdminActivateUser(ctx, parsed)
	if err != nil {
		if errors.Is(err, domainerr.ErrUserAlreadyActivated) {
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
