package handler

import (
	"net/http"
	"strconv"

	"github.com/Anchousfish/rest/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createList(ctx *gin.Context) {
	uid, err := getUserId(ctx)
	if err != nil {
		return
	}
	var input models.TodoList
	if err := ctx.BindJSON(&input); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	idList, err := h.s.TodoList.CreateList(uid, input)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"id": idList,
	})
}

type allListsResponse struct {
	Data []models.TodoList `json:"data"`
}

func (h *Handler) getAllLists(ctx *gin.Context) {
	uid, err := getUserId(ctx)
	if err != nil {
		return
	}
	lists, err := h.s.TodoList.GetAllLists(uid)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, allListsResponse{
		Data: lists,
	})
}

func (h *Handler) getList(ctx *gin.Context) {
	uid, err := getUserId(ctx)
	if err != nil {
		return
	}
	listId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid id format")
		return
	}
	list, err := h.s.TodoList.GetList(listId, uid)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, list)

}

func (h *Handler) updateList(ctx *gin.Context) {
	uid, err := getUserId(ctx)
	if err != nil {
		return
	}
	listId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid id format")
		return
	}
	var input models.UpdateListInput
	if err = ctx.BindJSON(&input); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
	}

	err = h.s.TodoList.UpdateList(listId, uid, input)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, statusResponse{
		Status: "OK",
	})

}

func (h *Handler) deleteList(ctx *gin.Context) {
	uid, err := getUserId(ctx)
	if err != nil {
		return
	}
	listId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, "invalid id format")
		return
	}
	err = h.s.TodoList.DeleteList(listId, uid)
	if err != nil {
		NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, statusResponse{
		Status: "OK",
	})
}
