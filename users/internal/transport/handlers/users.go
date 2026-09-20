package handlers

import (
	"context"
	"fmt"
	"net/http"
	"pkg/helpers"
	auth "users/api/proto"
	"users/internal/domain/service"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type UserHandler struct {
	grpcConn *grpc.ClientConn
	service  service.UsersService
}

func NewUserHandler(grpcConn *grpc.ClientConn, service service.UsersService) *UserHandler {
	return &UserHandler{
		grpcConn: grpcConn,
		service:  service,
	}
}

func (h *UserHandler) Registration(c *gin.Context) {
	client := auth.NewAuthClient(h.grpcConn)

	resp, err := client.AddUser(context.Background(), &auth.AddUserRequest{
		Id:       "00000000-0000-0000-0000-100000000000",
		Password: "34mf9304mf3940fj43jf34iksdz",
		Email:    "max@gmail.com",
	})

	if err != nil {
		c.String(http.StatusInternalServerError, "Error calling gRPC service: %v", err)
		return
	}

	// В resp будут refresh и access токены, которые нужно записать в куки
	fmt.Println(resp)

	c.String(http.StatusOK, "pong")
}

func (h *UserHandler) Profile(c *gin.Context) {
	claims, ok := helpers.ClaimsFromContextGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "claims not found",
		})
		return
	}

	_, err := h.service.Profile(claims.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get profile",
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (h *UserHandler) Activate(c *gin.Context) {
	claims, ok := helpers.ClaimsFromContextGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "claims not found",
		})
		return
	}

	err := h.service.ActivateUser(claims.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "failed to activate user",
		})
		return
	}
}
