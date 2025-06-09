package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/bypr1/gx-toko-kue/app/handler"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	e := echo.New()

	dsn := "root:123456@tcp(localhost:3306)/toko_kue?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to initialize database")
	}

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	hargaHandler := handler.NewHargaHandler(db)
	e.POST("/hitung-hpp", hargaHandler.HitungHPPPerKue)

	// Start server
	if err := e.Start(":8080"); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("Failed to start server")
	}

}
