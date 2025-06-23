package main

import (
	"log"
	"net/http"
	"os"
	"rei-api"
	"rei-api/services"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	connection := os.Getenv("DATABASE_URL")
	db, err := sqlx.Open("postgres", connection)
	if err != nil {
		log.Fatal("Connection with postgres database failed: ", err)
	}
	defer db.Close()

	reader := rei.NewParser()

	fh := services.FeedHandler(db, reader)
	fh.RegisterFeedsRoutes(e)

	e.Logger.Fatal(e.Start(":1323"))
}
