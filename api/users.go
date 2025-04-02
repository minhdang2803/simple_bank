package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/minhdang2803/simple_bank/db/sqlc"
	"github.com/minhdang2803/simple_bank/utils"
)

type CreateUserParams struct {

	/// alphanum: contains ASCII alphanumeric characters only
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type CreateUserResponse struct {
	Username          string     `json:"username"`
	FullName          string     `json:"fullname"`
	Email             string     `json:"email"`
	PasswordChangedAt *time.Time `json:"password_changed_at"`
	CreateAt          *time.Time `json: "create_at"`
}

func (server *Server) CreateUser(ctx *gin.Context) {

	var request *CreateUserParams

	err := ctx.ShouldBindBodyWithJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := utils.HashedPassword(request.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	user, err := server.store.CreateUser(ctx, db.CreateUserParams{
		Username:       request.Username,
		HashedPassword: hashedPassword,
		FullName:       request.FullName,
		Email:          request.Email,
	})

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	response := CreateUserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: &(user.PasswordChangedAt),
		CreateAt:          &(user.CreatedAt),
	}

	ctx.JSON(http.StatusOK, response)
}

type GetUserParam struct {
	Username string `form:"username" binding:"required"`
}

func (server *Server) GetUser(ctx *gin.Context) {
	var request *GetUserParam

	err := ctx.ShouldBindQuery(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUser(ctx, request.Username)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (server *Server) UpdateUser(ctx *gin.Context) {
	var request *db.UpdateUserParams

	err := ctx.ShouldBindBodyWithJSON(&request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err = server.store.UpdateUser(ctx, *request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Succesfully update %s infomations", request.Username),
	})
}
