package handlers

import (
	"context"
	"fmt"
	"net/http"
	"pkg/helpers"
	auth "users/api/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type UserHandler struct {
	grpcConn *grpc.ClientConn
}

func NewUserHandler(grpcConn *grpc.ClientConn) *UserHandler {
	return &UserHandler{
		grpcConn: grpcConn,
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

func (h *UserHandler) Account(c *gin.Context) {
	_, ok := helpers.ClaimsFromContextGin(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "claims not found",
		})
		return
	}
}
