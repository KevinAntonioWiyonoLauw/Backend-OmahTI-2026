package handler

import (
	"net/http"

	"inventory-management/pkg/utils"

	"github.com/gin-gonic/gin"
)

// Middleware recovery dipisahkan agar strategi penanganan panic bisa berubah
// tanpa menyentuh business logic service/repository.
func JSONRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered any) {
		_ = recovered
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan internal", nil)
	})
}
