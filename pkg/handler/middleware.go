package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader) //получение хэдера авторизации и его валидация
	if header == "" {
		newErrorResponse(c, http.StatusUnauthorized, "пустой заголовок авторизации")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "невалидный заголовок авторизации")
		return
	}

	// parse token
	userId, err := h.services.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	// запишем значение id в контекст (для того чтобы иметь доступ к id пользователя,
	//  который делает запрос в последующих обработчиках, которые вызываются после данной прослойки middleware)
	c.Set(userCtx, userId)
}

// функция приведения интерфейса id из контекста к инту
func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx) //возвращает интерфейс
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found in context")
		return 0, errors.New("user id not found in context")
	}

	idInt, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found in context")
		return 0, errors.New("user id not found in context")
	}
	return idInt, nil
}
