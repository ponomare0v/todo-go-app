package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ponomare0v/todo-go-app/pkg/models"
)

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body models.User true "account info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *gin.Context) {
	var input models.User

	//парсим тело запроса и валидируем его
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	//передаем на слой ниже в сервис, из которого получаем id созданного юзера в бд
	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error()) //внутрення ошибка на сервере
		return
	}

	//если все окей без ошибок создалось и вернулось
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

//
//
//
//
//

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary SignIn
// @Tags auth
// @Description authenticate user and return token
// @ID sign-in
// @Accept json
// @Produce json
// @Param input body signInInput true "credentials"
// @Success 200 {object} map[string]string "token"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *gin.Context) {
	var input signInInput

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
	}

	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// возврат токена
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}
