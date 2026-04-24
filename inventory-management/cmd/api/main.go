package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"inventory-management/internal/handler"
	"inventory-management/internal/repository"
	"inventory-management/internal/service"
	"inventory-management/pkg/config"
	"inventory-management/pkg/database"
	"inventory-management/pkg/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("gagal membaca konfigurasi: %v", err)
	}

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("gagal terhubung ke database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("gagal menutup koneksi database: %v", err)
		}
	}()

	itemRepo := repository.NewItemRepository(db)
	itemService := service.NewItemService(itemRepo)
	itemHandler := handler.NewItemHandler(itemService)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(handler.JSONRecovery())

	router.GET("/health", func(c *gin.Context) {
		utils.Success(c, http.StatusOK, "service sehat", gin.H{"uptime": "ok"})
	})
	router.NoRoute(func(c *gin.Context) {
		utils.Error(c, http.StatusNotFound, "route tidak ditemukan", nil)
	})

	itemHandler.RegisterRoutes(router)

	srv := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("server berjalan di port %s", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("gagal menjalankan server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("signal shutdown diterima, mulai graceful shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown gagal: %v", err)
	}

	log.Println("server berhenti dengan aman")
}
