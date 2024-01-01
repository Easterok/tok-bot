package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Run() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	e.GET("/", hello)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("API_PORT"))))
}

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
