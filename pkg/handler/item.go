package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ponomare0v/todo-go-app/pkg/models"
)

// @Summary Create a new item
// @Security ApiKeyAuth
// @Tags items
// @Description Create a new item in a specific list
// @ID create-item
// @Accept json
// @Produce json
// @Param id path int true "List ID"
// @Param input body models.TodoItem true "Item data"
// @Success 200 {object} map[string]interface{} "Item created successfully"
// @Failure 400 {object} errorResponse "Invalid data or list ID"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /api/lists/{id}/items [post]
func (h *Handler) createItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	var input models.TodoItem
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.services.TodoItem.Create(userId, listId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// @Summary Get all items in a list
// @Security ApiKeyAuth
// @Tags items
// @Description Get all items from a specific list
// @ID get-all-items
// @Accept json
// @Produce json
// @Param id path int true "List ID"
// @Success 200 {array} models.TodoItem "List of items"
// @Failure 400 {object} errorResponse "Invalid list ID"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /api/lists/{id}/items [get]
func (h *Handler) getAllItems(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	listId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	items, err := h.services.TodoItem.GetAll(userId, listId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, items)

}

// @Summary Get an item by its ID
// @Security ApiKeyAuth
// @Tags items
// @Description Get a specific item by its ID
// @ID get-item-by-id
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Success 200 {object} models.TodoItem "Item details"
// @Failure 400 {object} errorResponse "Invalid item ID"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /api/items/{id} [get]
func (h *Handler) getItemById(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	item, err := h.services.TodoItem.GetById(userId, itemId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, item)
}

// @Summary Update an item by its ID
// @Security ApiKeyAuth
// @Tags items
// @Description Update the details of a specific item
// @ID update-item-by-id
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Param item body models.UpdateItemInput true "Item data to update"
// @Success 200 {object} statusResponse "Update successful"
// @Failure 400 {object} errorResponse "Invalid input or ID"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /api/items/{id} [put]
func (h *Handler) updateItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid id param")
		return
	}

	var input models.UpdateItemInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.TodoItem.Update(userId, id, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}

// @Summary Delete an item by its ID
// @Security ApiKeyAuth
// @Tags items
// @Description Delete a specific item by its ID
// @ID delete-item-by-id
// @Param id path int true "Item ID"
// @Success 200 {object} statusResponse "Delete successful"
// @Failure 400 {object} errorResponse "Invalid item ID"
// @Failure 500 {object} errorResponse "Internal server error"
// @Router /api/items/{id} [delete]
func (h *Handler) deleteItem(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		return
	}

	itemId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "invalid list id param")
		return
	}

	err = h.services.TodoItem.Delete(userId, itemId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
