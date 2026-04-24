package handler

import (
	"net/http"
	"strconv"

	"inventory-management/internal/models"
	"inventory-management/internal/service"
	"inventory-management/pkg/utils"

	"github.com/gin-gonic/gin"
)

// Layer handler berubah ketika kontrak HTTP (route, payload, status code) berubah,
// bukan saat aturan bisnis atau query data berubah.
type ItemHandler struct {
	service service.ItemService
}

func NewItemHandler(service service.ItemService) *ItemHandler {
	return &ItemHandler{service: service}
}

func (h *ItemHandler) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/api/v1/items")
	{
		group.POST("", h.Create)
		group.GET("", h.List)
		group.GET("/:id", h.GetByID)
		group.PUT("/:id", h.Update)
		group.DELETE("/:id", h.Delete)
	}
}

func (h *ItemHandler) Create(c *gin.Context) {
	var input models.CreateItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.HandleError(c, utils.NewInvalidInputError("payload tidak valid", err))
		return
	}

	item, err := h.service.Create(c.Request.Context(), input)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, http.StatusCreated, "item berhasil dibuat", item)
}

func (h *ItemHandler) GetByID(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	item, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, http.StatusOK, "item berhasil diambil", item)
}

func (h *ItemHandler) List(c *gin.Context) {
	items, err := h.service.List(c.Request.Context())
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, http.StatusOK, "daftar item berhasil diambil", items)
}

func (h *ItemHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	var input models.UpdateItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.HandleError(c, utils.NewInvalidInputError("payload tidak valid", err))
		return
	}

	item, err := h.service.Update(c.Request.Context(), id, input)
	if err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, http.StatusOK, "item berhasil diperbarui", item)
}

func (h *ItemHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		utils.HandleError(c, err)
		return
	}

	utils.Success(c, http.StatusOK, "item berhasil dihapus", nil)
}

func parseID(c *gin.Context) (int64, bool) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		utils.HandleError(c, utils.NewInvalidInputError("id harus berupa angka positif", err))
		return 0, false
	}
	return id, true
}
