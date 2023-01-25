package handler

import (
	"net/http"
	"strconv"

	"github.com/Anchousfish/rest/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createItem(ctx *gin.Context) {
	uid, err := getUserId(ctx)
	if err != nil {
		return
	}
	listId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid list id format")
		return
	}
	var inputItem models.TodoItem
	if err = ctx.BindJSON(&inputItem); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	idItem, err := h.s.TodoItem.CreateItem(uid, listId, inputItem)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id_Item": idItem,
	})
}

type allItemsResponse struct {
	Data []models.TodoItem `json:"data"`
}

func (h *Handler) getAllItems(ctx *gin.Context) {
	uid, err := getUserId(ctx)
	if err != nil {
		return
	}
	listId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid list id format")
		return
	}
	items, err := h.s.TodoItem.GetAllItems(uid, listId)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, allItemsResponse{
		Data: items,
	})

}
func (h *Handler) getItem(ctx *gin.Context) {
	uid, err := getUserId(ctx)
	if err != nil {
		return
	}
	itemId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid item id format")
		return
	}
	item, err := h.s.TodoItem.GetItem(uid, itemId)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, item)

}
func (h *Handler) updateItem(ctx *gin.Context) {
	uid, err := getUserId(ctx)
	if err != nil {
		return
	}
	listId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid id format")
		return
	}
	var input models.UpdateItemInput
	if err = ctx.BindJSON(&input); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	err = h.s.TodoItem.UpdateItem(listId, uid, input)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, statusResponse{
		Status: "OK",
	})

}
func (h *Handler) deleteItem(ctx *gin.Context) {
	uid, err := getUserId(ctx)
	if err != nil {
		return
	}
	itemId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid item id format")
		return
	}
	err = h.s.TodoItem.DeleteItem(uid, itemId)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, statusResponse{Status: "Deleted"})

}
