package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ponomare0v/todo-go-app/pkg/models"
)

// @Summary Create todo list
// @Security ApiKeyAuth
// @Tags lists
// @Description create todo list
// @ID create-list
// @Accept json
// @Produce json
// @Param input body models.TodoList true "list info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists [post]
func (h *Handler) createList(c *gin.Context) { // хэндлеры для работы с эндпоинтами списков
	userId, err := getUserId(c) // получаем id пользователя из контекса после аутентификации
	if err != nil {
		return
	}

	var input models.TodoList
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	//// вот здесь передаем id из начала где получаем его из контекста, но получаем мы интерфейс,
	//  а передавать надо int, чтобы не приводить постоянно к int создадим функцию в middleware под названием getUserId
	id, err := h.services.TodoList.Create(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// . Для ответа (response) используем дополнительную структуру, в которой будет поле data типа слайса списков
type getAllListsResponse struct {
	Data []models.TodoList `json:"data"`
}

// @Summary Get all todo lists
// @Security ApiKeyAuth
// @Tags lists
// @Description get all todo lists of the authenticated user
// @ID get-all-lists
// @Accept json
// @Produce json
// @Success 200 {object} getAllListsResponse "list of todo lists"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/lists [get]
func (h *Handler) getAllLists(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	lists, err := h.services.TodoList.GetAll(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, getAllListsResponse{
		Data: lists,
	})
}

// @Summary Get todo list by ID
// @Security ApiKeyAuth
// @Tags lists
// @Description Get a todo list by its ID
// @ID get-list-by-id
// @Produce json
// @Param id path int true "List ID"
// @Success 200 {object} models.TodoList
// @Failure 400 {object} errorResponse "Invalid ID param"
// @Failure 404 {object} errorResponse "List not found"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /api/lists/{id} [get]
func (h *Handler) getListById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	//добавим тут также получение параметра id списка из пути запроса - для этого вызовем у контекста метод param,
	//  указав в качестве аргумента - имя параметра. Сразу обернем это в функцию Atoi из стандартной библиотеки strconv для перобразования строки в число.
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	list, err := h.services.TodoList.GetById(userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, list)
}

// @Summary Update todo list by ID
// @Security ApiKeyAuth
// @Tags lists
// @Description Update an existing todo list by its ID
// @ID update-list
// @Accept json
// @Produce json
// @Param id path int true "List ID"
// @Param input body models.UpdateListInput true "Updated list info"
// @Success 200 {object} statusResponse "List updated successfully"
// @Failure 400 {object} errorResponse "Invalid ID param or request body"
// @Failure 404 {object} errorResponse "List not found"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /api/lists/{id} [put]
func (h *Handler) updateList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	//добавим тут также получение параметра id списка из пути запроса - для этого вызовем у контекста метод param,
	//  указав в качестве аргумента - имя параметра. Сразу обернем это в функцию Atoi из стандартной библиотеки strconv для перобразования строки в число.
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input models.UpdateListInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.TodoList.Update(userId, id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

// @Summary Delete todo list by ID
// @Security ApiKeyAuth
// @Tags lists
// @Description Delete a todo list by its ID
// @ID delete-list
// @Accept json
// @Produce json
// @Param id path int true "List ID"
// @Success 200 {object} statusResponse "List deleted successfully"
// @Failure 400 {object} errorResponse "Invalid ID param"
// @Failure 404 {object} errorResponse "List not found"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /api/lists/{id} [delete]
func (h *Handler) deleteList(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	//добавим тут также получение параметра id списка из пути запроса - для этого вызовем у контекста метод param,
	//  указав в качестве аргумента - имя параметра. Сразу обернем это в функцию Atoi из стандартной библиотеки strconv для перобразования строки в число.
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	err = h.services.TodoList.Delete(userId, id)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{ // структура для ответа написана в handler/response.go
		Status: "ok",
	})
}
